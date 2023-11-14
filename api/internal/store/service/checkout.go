package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) CheckoutCreditCard(ctx context.Context, in *store.CheckoutCreditCardInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderCreditCardParams{
			SessionID:         sessionID,
			Number:            in.Number,
			Month:             in.Month,
			Year:              in.Year,
			VerificationValue: in.VerificationValue,
			Email:             params.Customer.Email(),
			Lastname:          params.BillingAddress.Lastname,
			Firstname:         params.BillingAddress.Firstname,
		}
		return s.komoju.Session.OrderCreditCard(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypeCreditCard,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutPayPay(ctx context.Context, in *store.CheckoutPayPayInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderPayPayParams{
			SessionID: sessionID,
		}
		return s.komoju.Session.OrderPayPay(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypePayPay,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) NotifyPaymentCompleted(ctx context.Context, in *store.NotifyPaymentCompletedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateOrderPaymentParams{
		Status:    entity.NewPaymentStatus(komoju.PaymentStatus(in.Status)),
		PaymentID: in.PaymentID,
		IssuedAt:  in.IssuedAt,
	}
	err := s.db.Order.UpdatePaymentStatus(ctx, in.OrderID, params)
	if errors.Is(err, database.ErrFailedPrecondition) {
		s.logger.Warn("Order is not updatable",
			zap.String("orderId", in.OrderID), zap.String("status", in.Status), zap.Time("issuedAt", in.IssuedAt))
		return nil
	}
	return internalError(err)
}

type checkoutParams struct {
	payload           *store.CheckoutDetail
	paymentMethodType entity.PaymentMethodType
	payFn             func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error)
}

func (s *service) checkout(ctx context.Context, params *checkoutParams) (string, error) {
	var (
		customer        *uentity.User
		billingAddress  *uentity.Address
		shippingAddress *uentity.Address
		shipping        *entity.Shipping
		cart            *entity.Cart
		promotion       *entity.Promotion
	)
	eg, ectx := errgroup.WithContext(ctx)
	// ユーザーの取得
	eg.Go(func() (err error) {
		in := &user.GetUserInput{
			UserID: params.payload.UserID,
		}
		customer, err = s.user.GetUser(ectx, in)
		return
	})
	// 請求先/配送先住所の取得
	eg.Go(func() (err error) {
		in := &user.GetAddressInput{
			UserID:    params.payload.UserID,
			AddressID: params.payload.BillingAddressID,
		}
		billingAddress, err = s.user.GetAddress(ectx, in)
		return
	})
	eg.Go(func() (err error) {
		in := &user.GetAddressInput{
			UserID:    params.payload.UserID,
			AddressID: params.payload.ShippingAddressID,
		}
		shippingAddress, err = s.user.GetAddress(ctx, in)
		return
	})
	// カートの取得
	eg.Go(func() (err error) {
		cart, err = s.getCart(ctx, params.payload.SessionID)
		return
	})
	// 配送設定の取得
	eg.Go(func() (err error) {
		shipping, err = s.getShippingByCoordinatorID(ctx, params.payload.CoordinatorID)
		return
	})
	// プロモーションの取得
	eg.Go(func() (err error) {
		if params.payload.PromotionID == "" {
			return
		}
		promotion, err = s.db.Promotion.Get(ctx, params.payload.PromotionID)
		return
	})
	if err := eg.Wait(); err != nil {
		return "", internalError(err)
	}
	// プロモーションの有効性検証
	if params.payload.PromotionID != "" && !promotion.IsEnabled(s.now()) {
		s.logger.Warn("Failed to user disable promotion",
			zap.String("userId", params.payload.UserID), zap.String("promotionId", params.payload.PromotionID))
		return "", fmt.Errorf("service: disable promotion: %w", exception.ErrFailedPrecondition)
	}
	// 購入する買い物かごのみ取得
	baskets := cart.Baskets.FilterByCoordinatorID(params.payload.CoordinatorID).FilterByBoxNumber(params.payload.BoxNumber)
	if len(baskets) == 0 {
		return "", fmt.Errorf("service: no target baskets: %w", exception.ErrInvalidArgument)
	}
	// 在庫一覧を取得
	products, err := s.db.Product.MultiGet(ctx, baskets.ProductIDs())
	if err != nil {
		return "", internalError(err)
	}
	// 在庫の過不足確認
	if err := baskets.VerifyQuantities(products.Map()); err != nil {
		s.logger.Warn("Failed to verify quantities in baskets",
			zap.String("userId", params.payload.UserID), zap.String("sessionId", params.payload.SessionID),
			zap.String("coordinatorId", params.payload.CoordinatorID), zap.Int64("boxNumber", params.payload.BoxNumber))
		return "", fmt.Errorf("service: insufficient stock: %w: %s", exception.ErrFailedPrecondition, err.Error())
	}
	// 注文インスタンスの生成
	oparams := &entity.NewOrderParams{
		CoordinatorID:     params.payload.CoordinatorID,
		Customer:          customer,
		BillingAddress:    billingAddress,
		ShippingAddress:   shippingAddress,
		Shipping:          shipping,
		Baskets:           cart.Baskets,
		Products:          products,
		PaymentMethodType: params.paymentMethodType,
		Promotion:         promotion,
	}
	order, err := entity.NewOrder(oparams)
	if err != nil {
		return "", internalError(err)
	}
	// チェックサム
	if params.payload.Total != order.OrderPayment.Total {
		s.logger.Warn("Failed to checksum before checkout",
			zap.String("userId", params.payload.UserID), zap.String("sessionId", params.payload.SessionID),
			zap.String("coordinatorId", params.payload.CoordinatorID), zap.Int64("boxNumber", params.payload.BoxNumber),
			zap.Int64("payload.total", params.payload.Total), zap.Any("payment", order.OrderPayment))
		return "", fmt.Errorf("service: unmatch total: %w", exception.ErrInvalidArgument)
	}
	// 決済トランザクションの発行
	sparams := &komoju.CreateSessionParams{
		OrderID:      order.ID,
		Amount:       order.OrderPayment.Total,
		CallbackURL:  params.payload.CallbackURL,
		PaymentTypes: entity.NewKomojuPaymentTypes(params.paymentMethodType),
		Customer: &komoju.CreateSessionCustomer{
			ID:    customer.ID,
			Name:  billingAddress.Name(),
			Email: customer.Email(),
		},
		BillingAddress: &komoju.CreateSessionAddress{
			ZipCode:      billingAddress.PostalCode,
			Prefecture:   billingAddress.Prefecture,
			City:         billingAddress.City,
			AddressLine1: billingAddress.AddressLine1,
			AddressLine2: billingAddress.AddressLine2,
		},
		ShippingAddress: &komoju.CreateSessionAddress{
			ZipCode:      shippingAddress.PostalCode,
			Prefecture:   shippingAddress.Prefecture,
			City:         shippingAddress.City,
			AddressLine1: shippingAddress.AddressLine1,
			AddressLine2: shippingAddress.AddressLine2,
		},
	}
	session, err := s.komoju.Session.Create(ctx, sparams)
	if err != nil {
		return "", internalError(err)
	}
	order.OrderPayment.SetTransactionID(session.ID, s.now())
	// 注文履歴レコードの登録
	if err := s.db.Order.Create(ctx, order); err != nil {
		return "", internalError(err)
	}
	// 決済依頼処理
	pay, err := params.payFn(ctx, session.ID, oparams)
	if err != nil {
		return "", internalError(err)
	}
	// 買い物かごの削除
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		cart.RemoveBaskets(baskets.BoxNumbers()...)
		if err := s.refreshCart(context.Background(), cart); err != nil {
			s.logger.Error("Failed to refresh cart after checkout",
				zap.Any("payload", params.payload), zap.Any("order", order),
				zap.Int32("methodType", int32(params.paymentMethodType)), zap.Error(err))
		}
	}()
	return pay.RedirectURL, nil
}

func (s *service) getShippingByCoordinatorID(ctx context.Context, coordinatorID string) (*entity.Shipping, error) {
	shipping, err := s.db.Shipping.GetByCoordinatorID(ctx, coordinatorID)
	if errors.Is(err, exception.ErrNotFound) {
		// 配送設定が完了していないコーディネータの場合、デフォルト設定を使用
		return s.db.Shipping.GetDefault(ctx)
	}
	return shipping, err
}
