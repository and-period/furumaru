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
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/japanese"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

var errDuplicateOrder = errors.New("service: duplicate order")

func (s *service) GetCheckoutState(
	ctx context.Context,
	in *store.GetCheckoutStateInput,
) (string, entity.PaymentStatus, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", entity.PaymentStatusUnknown, internalError(err)
	}
	var (
		order *entity.Order
		err   error
	)
	if in.UserID != "" {
		order, err = s.db.Order.GetByTransactionID(ctx, in.UserID, in.TransactionID)
	} else {
		order, err = s.db.Order.GetByTransactionIDWithSessionID(ctx, in.SessionID, in.TransactionID)
	}
	if err != nil {
		return "", entity.PaymentStatusUnknown, internalError(err)
	}
	if order.OrderPayment.Status != entity.PaymentStatusPending {
		return order.ID, order.OrderPayment.Status, nil
	}
	s.logger.Info("This checkout is pending", zap.String("transactionId", in.TransactionID))
	// 未払い状態の場合、KOMOJUから最新の状態を取得する
	res, err := s.komoju.Session.Get(ctx, in.TransactionID)
	if err != nil || res.Payment == nil {
		s.logger.Warn(
			"Failed to get session state",
			zap.String("transactionId", in.TransactionID),
			zap.Error(err),
		)
		return order.ID, entity.PaymentStatusUnknown, internalError(err)
	}
	return order.ID, entity.NewPaymentStatus(res.Payment.Status), nil
}

func (s *service) CheckoutCreditCard(
	ctx context.Context,
	in *store.CheckoutCreditCardInput,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	cardParams := &entity.NewCreditCardParams{
		Name:   in.Name,
		Number: in.Number,
		Month:  in.Month,
		Year:   in.Year,
		CVV:    in.VerificationValue,
	}
	card := entity.NewCreditCard(cardParams)
	if err := card.Validate(s.now()); err != nil {
		return "", fmt.Errorf(
			"service: invalid credit card: %s: %w",
			err.Error(),
			exception.ErrInvalidArgument,
		)
	}
	payFn := func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderCreditCardParams{
			SessionID:         sessionID,
			Name:              card.Name,
			Number:            card.Number,
			Month:             card.Month,
			Year:              card.Year,
			VerificationValue: card.CVV,
			Email:             params.customer.Email(),
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

func (s *service) CheckoutPayPay(
	ctx context.Context,
	in *store.CheckoutPayPayInput,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, _ *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
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

func (s *service) CheckoutLinePay(
	ctx context.Context,
	in *store.CheckoutLinePayInput,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, _ *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderLinePayParams{
			SessionID: sessionID,
		}
		return s.komoju.Session.OrderLinePay(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypeLinePay,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutMerpay(
	ctx context.Context,
	in *store.CheckoutMerpayInput,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, _ *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderMerpayParams{
			SessionID: sessionID,
		}
		return s.komoju.Session.OrderMerpay(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypeMerpay,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutRakutenPay(
	ctx context.Context,
	in *store.CheckoutRakutenPayInput,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, _ *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderRakutenPayParams{
			SessionID: sessionID,
		}
		return s.komoju.Session.OrderRakutenPay(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypeRakutenPay,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutAUPay(ctx context.Context, in *store.CheckoutAUPayInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, _ *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderAUPayParams{
			SessionID: sessionID,
		}
		return s.komoju.Session.OrderAUPay(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypeAUPay,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutPaidy(ctx context.Context, in *store.CheckoutPaidyInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderPaidyParams{
			SessionID: sessionID,
			Email:     params.customer.Email(),
			Name:      params.customer.Name(),
		}
		return s.komoju.Session.OrderPaidy(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypePaidy,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutBankTransfer(
	ctx context.Context,
	in *store.CheckoutBankTransferInput,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderBankTransferParams{
			SessionID:     sessionID,
			Email:         params.customer.Email(),
			PhoneNumber:   params.billingAddress.PhoneNumber,
			Lastname:      params.billingAddress.Lastname,
			Firstname:     params.billingAddress.Firstname,
			LastnameKana:  japanese.HiraganaToKatakana(params.billingAddress.LastnameKana),
			FirstnameKana: japanese.HiraganaToKatakana(params.billingAddress.FirstnameKana),
		}
		return s.komoju.Session.OrderBankTransfer(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypeBankTransfer,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutPayEasy(
	ctx context.Context,
	in *store.CheckoutPayEasyInput,
) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	payFn := func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		in := &komoju.OrderPayEasyParams{
			SessionID:     sessionID,
			Email:         params.customer.Email(),
			PhoneNumber:   params.billingAddress.PhoneNumber,
			Lastname:      params.billingAddress.Lastname,
			Firstname:     params.billingAddress.Firstname,
			LastnameKana:  japanese.HiraganaToKatakana(params.billingAddress.LastnameKana),
			FirstnameKana: japanese.HiraganaToKatakana(params.billingAddress.FirstnameKana),
		}
		return s.komoju.Session.OrderPayEasy(ctx, in)
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypePayEasy,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

func (s *service) CheckoutFree(ctx context.Context, in *store.CheckoutFreeInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	if in.Total != 0 {
		return "", fmt.Errorf("service: invalid total: %w", exception.ErrInvalidArgument)
	}
	payFn := func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
		return &komoju.OrderSessionResponse{}, nil
	}
	params := &checkoutParams{
		payload:           &in.CheckoutDetail,
		paymentMethodType: entity.PaymentMethodTypeNone,
		payFn:             payFn,
	}
	return s.checkout(ctx, params)
}

type checkoutParams struct {
	payload           *store.CheckoutDetail
	paymentMethodType entity.PaymentMethodType
	customer          *uentity.User
	billingAddress    *uentity.Address
	shippingAddress   *uentity.Address
	payFn             func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error)
}

type checkoutDetailParams struct {
	customer       *uentity.User
	billingAddress *uentity.Address
}

func (s *service) checkout(ctx context.Context, params *checkoutParams) (string, error) {
	var checkoutFn func(context.Context, *checkoutParams) (string, error)
	switch params.payload.Type {
	case entity.OrderTypeProduct:
		if err := s.validator.Struct(params.payload.CheckoutProductDetail); err != nil {
			return "", internalError(err)
		}
		checkoutFn = s.checkoutProduct
	case entity.OrderTypeExperience:
		if err := s.validator.Struct(params.payload.CheckoutExperienceDetail); err != nil {
			return "", internalError(err)
		}
		checkoutFn = s.checkoutExperience
	default:
		return "", fmt.Errorf("service: invalid order type: %w", exception.ErrInvalidArgument)
	}

	eg, ectx := errgroup.WithContext(ctx)
	// ユーザーの取得
	eg.Go(func() (err error) {
		in := &user.GetUserInput{
			UserID: params.payload.UserID,
		}
		params.customer, err = s.user.GetUser(ectx, in)
		return
	})
	// 請求先住所の取得
	eg.Go(func() (err error) {
		in := &user.GetAddressInput{
			UserID:    params.payload.UserID,
			AddressID: params.payload.BillingAddressID,
		}
		params.billingAddress, err = s.user.GetAddress(ectx, in)
		return
	})
	// 配送先住所の取得
	eg.Go(func() (err error) {
		if params.payload.ShippingAddressID == "" {
			return
		}
		in := &user.GetAddressInput{
			UserID:    params.payload.UserID,
			AddressID: params.payload.ShippingAddressID,
		}
		params.shippingAddress, err = s.user.GetAddress(ectx, in)
		return
	})
	if err := eg.Wait(); err != nil {
		return "", internalError(err)
	}

	return checkoutFn(ctx, params)
}

func (s *service) checkoutProduct(ctx context.Context, params *checkoutParams) (string, error) {
	// 店舗の有効性検証
	shop, err := s.db.Shop.GetByCoordinatorID(ctx, params.payload.CoordinatorID)
	if err != nil {
		return "", internalError(err)
	}
	if !shop.Enabled() {
		return "", fmt.Errorf("service: shop is disabled: %w", exception.ErrForbidden)
	}
	var (
		shipping  *entity.Shipping
		cart      *entity.Cart
		promotion *entity.Promotion
	)
	eg, ectx := errgroup.WithContext(ctx)
	// カートの取得
	eg.Go(func() (err error) {
		cart, err = s.getCart(ctx, params.payload.SessionID)
		return
	})
	// 配送設定の取得
	eg.Go(func() (err error) {
		shipping, err = s.getShippingByCoordinatorID(ectx, params.payload.CoordinatorID)
		return
	})
	// プロモーションの取得
	eg.Go(func() (err error) {
		if params.payload.PromotionCode == "" {
			return
		}
		promotion, err = s.db.Promotion.GetByCode(ectx, params.payload.PromotionCode)
		return
	})
	err = eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return "", fmt.Errorf("service: not found: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return "", internalError(err)
	}
	// プロモーションの有効性検証
	if params.payload.PromotionCode != "" && !promotion.IsEnabled(shop.ID) {
		s.logger.Warn(
			"Failed to disable promotion",
			zap.String(
				"userId",
				params.payload.UserID,
			),
			zap.String("code", params.payload.PromotionCode),
		)
		return "", fmt.Errorf("service: disable promotion: %w", exception.ErrFailedPrecondition)
	}
	// 購入する買い物かごのみ取得
	baskets := cart.Baskets.FilterByCoordinatorID(params.payload.CoordinatorID).
		FilterByBoxNumber(params.payload.BoxNumber)
	if len(baskets) == 0 {
		return "", fmt.Errorf("service: no target baskets: %w", exception.ErrInvalidArgument)
	}
	// 在庫一覧を取得
	productIDs := baskets.ProductIDs()
	products, err := s.db.Product.MultiGet(ctx, productIDs)
	if err != nil {
		return "", internalError(err)
	}
	products = products.FilterBySales()
	// 商品がすべて販売中かの確認
	if len(products) > 0 && len(productIDs) != len(products) {
		s.logger.Warn(
			"Failed because there are products outside the sales period",
			zap.String(
				"userId",
				params.payload.UserID,
			),
			zap.String("sessionId", params.payload.SessionID),
			zap.String(
				"coordinatorId",
				params.payload.CoordinatorID,
			),
			zap.Int64("boxNumber", params.payload.BoxNumber),
		)
		return "", fmt.Errorf(
			"service: there are products outside the sales period: %w",
			exception.ErrFailedPrecondition,
		)
	}
	// 在庫の過不足確認
	if err := baskets.VerifyQuantities(products.Map()); err != nil {
		s.logger.Warn(
			"Failed to verify quantities in baskets",
			zap.String(
				"userId",
				params.payload.UserID,
			),
			zap.String("sessionId", params.payload.SessionID),
			zap.String(
				"coordinatorId",
				params.payload.CoordinatorID,
			),
			zap.Int64("boxNumber", params.payload.BoxNumber),
		)
		return "", fmt.Errorf(
			"service: insufficient stock: %w: %s",
			exception.ErrFailedPrecondition,
			err.Error(),
		)
	}
	// 注文インスタンスの生成
	oparams := &entity.NewProductOrderParams{
		OrderID:           params.payload.RequestID,
		SessionID:         params.payload.SessionID,
		ShopID:            shop.ID,
		CoordinatorID:     params.payload.CoordinatorID,
		Customer:          params.customer,
		BillingAddress:    params.billingAddress,
		ShippingAddress:   params.shippingAddress,
		Shipping:          shipping,
		Baskets:           baskets,
		Products:          products,
		PaymentMethodType: params.paymentMethodType,
		Promotion:         promotion,
	}
	order, err := entity.NewProductOrder(oparams)
	if err != nil {
		return "", internalError(err)
	}
	// チェックサム
	if params.payload.Total != order.Total {
		s.logger.Warn(
			"Failed to checksum before checkout",
			zap.String(
				"userId",
				params.payload.UserID,
			),
			zap.String("sessionId", params.payload.SessionID),
			zap.String(
				"coordinatorId",
				params.payload.CoordinatorID,
			),
			zap.Int64("boxNumber", params.payload.BoxNumber),
			zap.Int64(
				"payload.total",
				params.payload.Total,
			),
			zap.Any("payment", order.OrderPayment),
		)
		return "", fmt.Errorf("service: unmatch total: %w", exception.ErrInvalidArgument)
	}
	// 支払い処理
	var (
		redirectURL string
		afterFn     func(context.Context)
	)
	if order.Total > 0 {
		redirectURL, afterFn, err = s.executePaymentOrder(ctx, order, params)
	} else {
		redirectURL, afterFn, err = s.executeFreeOrder(ctx, order, params)
	}
	if errors.Is(err, errDuplicateOrder) {
		return redirectURL, nil
	}
	if err != nil {
		return "", err
	}
	s.waitGroup.Add(3)
	// 支払い完了後の処理
	go func() {
		defer s.waitGroup.Done()
		afterFn(context.Background())
	}()
	// 買い物かごの削除（即時決済の場合は仮売上通知のタイミングで削除するため、ここでは削除しない）
	go func() {
		defer s.waitGroup.Done()
		if order.IsImmediatePayment() {
			return
		}
		cart.RemoveBaskets(baskets.BoxNumbers()...)
		if err := s.refreshCart(context.Background(), cart); err != nil {
			s.logger.Error("Failed to refresh cart after checkout",
				zap.Any("payload", params.payload), zap.Any("order", order),
				zap.Int32("methodType", int32(params.paymentMethodType)), zap.Error(err))
		}
	}()
	// 商品の在庫を減算
	go func() {
		defer s.waitGroup.Done()
		s.decreaseProductInventories(context.Background(), order.OrderItems)
	}()
	return redirectURL, nil
}

func (s *service) checkoutExperience(ctx context.Context, params *checkoutParams) (string, error) {
	var (
		experience *entity.Experience
		promotion  *entity.Promotion
	)
	eg, ectx := errgroup.WithContext(ctx)
	// 体験の取得
	eg.Go(func() (err error) {
		experience, err = s.db.Experience.Get(ectx, params.payload.ExperienceID)
		return
	})
	// プロモーションの取得
	eg.Go(func() (err error) {
		if params.payload.PromotionCode == "" {
			return
		}
		promotion, err = s.db.Promotion.GetByCode(ectx, params.payload.PromotionCode)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return "", fmt.Errorf("service: not found: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return "", internalError(err)
	}
	// 店舗の有効性検証
	shop, err := s.db.Shop.Get(ctx, experience.ShopID)
	if err != nil {
		return "", internalError(err)
	}
	if !shop.Enabled() {
		return "", fmt.Errorf("service: shop is disabled: %w", exception.ErrForbidden)
	}
	// プロモーションの有効性検証
	if params.payload.PromotionCode != "" && !promotion.IsEnabled(experience.ShopID) {
		s.logger.Warn(
			"Failed to disable promotion",
			zap.String(
				"userId",
				params.payload.UserID,
			),
			zap.String("code", params.payload.PromotionCode),
		)
		return "", fmt.Errorf("service: disable promotion: %w", exception.ErrFailedPrecondition)
	}
	// 体験が販売中かの確認
	if experience.Status != entity.ExperienceStatusAccepting {
		s.logger.Warn(
			"Failed to checkout because the experience is not accepting",
			zap.String(
				"userId",
				params.payload.UserID,
			),
			zap.String("experienceId", params.payload.ExperienceID),
		)
		return "", fmt.Errorf(
			"service: the experience is not accepting: %w",
			exception.ErrFailedPrecondition,
		)
	}
	// 注文インスタンスの生成
	oparams := &entity.NewExperienceOrderParams{
		OrderID:               params.payload.RequestID,
		SessionID:             params.payload.SessionID,
		ShopID:                shop.ID,
		CoordinatorID:         experience.CoordinatorID,
		Customer:              params.customer,
		BillingAddress:        params.billingAddress,
		Experience:            experience,
		PaymentMethodType:     params.paymentMethodType,
		Promotion:             promotion,
		AdultCount:            params.payload.AdultCount,
		JuniorHighSchoolCount: params.payload.JuniorHighSchoolCount,
		ElementarySchoolCount: params.payload.ElementarySchoolCount,
		PreschoolCount:        params.payload.PreschoolCount,
		SeniorCount:           params.payload.SeniorCount,
		Transportation:        params.payload.Transportation,
		RequetsedDate:         params.payload.RequestedDate,
		RequetsedTime:         params.payload.RequestedTime,
	}
	order, err := entity.NewExperienceOrder(oparams)
	if err != nil {
		return "", internalError(err)
	}
	// チェックサム
	if params.payload.Total != order.Total {
		s.logger.Warn(
			"Failed to checksum before checkout",
			zap.String(
				"userId",
				params.payload.UserID,
			),
			zap.String("sessionId", params.payload.SessionID),
			zap.String(
				"coordinatorId",
				params.payload.CoordinatorID,
			),
			zap.String("experienceId", params.payload.ExperienceID),
			zap.Int64(
				"payload.total",
				params.payload.Total,
			),
			zap.Any("payment", order.OrderPayment),
		)
		return "", fmt.Errorf("service: unmatch total: %w", exception.ErrInvalidArgument)
	}
	var (
		redirectURL string
		afterFn     func(context.Context)
	)
	// 支払い処理
	if order.Total == 0 {
		redirectURL, afterFn, err = s.executeFreeOrder(ctx, order, params)
	} else {
		redirectURL, afterFn, err = s.executePaymentOrder(ctx, order, params)
	}
	if errors.Is(err, errDuplicateOrder) {
		return redirectURL, nil
	}
	if err != nil {
		return "", err
	}
	s.waitGroup.Add(1)
	// 支払い完了後の処理
	go func() {
		defer s.waitGroup.Done()
		afterFn(context.Background())
	}()
	return redirectURL, nil
}

func (s *service) executePaymentOrder(
	ctx context.Context, order *entity.Order, params *checkoutParams,
) (string, func(context.Context), error) {
	sparams := &komoju.CreateSessionParams{
		OrderID:      order.ID,
		Amount:       order.Total,
		CallbackURL:  params.payload.CallbackURL,
		PaymentTypes: entity.NewKomojuPaymentTypes(params.paymentMethodType),
		Customer: &komoju.CreateSessionCustomer{
			ID:    params.customer.ID,
			Name:  params.billingAddress.Name(),
			Email: params.customer.Email(),
		},
		BillingAddress: &komoju.CreateSessionAddress{
			ZipCode:      params.billingAddress.PostalCode,
			Prefecture:   params.billingAddress.Prefecture,
			City:         params.billingAddress.City,
			AddressLine1: params.billingAddress.AddressLine1,
			AddressLine2: params.billingAddress.AddressLine2,
		},
	}
	if params.shippingAddress != nil {
		sparams.ShippingAddress = &komoju.CreateSessionAddress{
			ZipCode:      params.shippingAddress.PostalCode,
			Prefecture:   params.shippingAddress.Prefecture,
			City:         params.shippingAddress.City,
			AddressLine1: params.shippingAddress.AddressLine1,
			AddressLine2: params.shippingAddress.AddressLine2,
		}
	}
	// 決済トランザクションの発行
	session, err := s.komoju.Session.Create(ctx, sparams)
	if err != nil {
		return "", nil, internalError(err)
	}
	order.SetTransaction(session.ID, s.now())
	// 注文履歴レコードの登録
	if err := s.db.Order.Create(ctx, order); err != nil {
		return "", nil, internalError(err)
	}
	// 決済依頼処理
	cparams := &checkoutDetailParams{
		customer:       params.customer,
		billingAddress: params.billingAddress,
	}
	pay, err := params.payFn(ctx, session.ID, cparams)
	if komoju.IsSessionFailed(err) && session.ReturnURL != "" {
		// 支払い状態取得エンドポイントから状態取得ができるよう、session_idを付与する
		return fmt.Sprintf(
			"%s?session_id=%s",
			session.ReturnURL,
			session.ID,
		), func(ctx context.Context) {}, errDuplicateOrder
	}
	if err != nil {
		return "", nil, internalError(err)
	}
	return pay.RedirectURL, func(context.Context) {}, nil
}

func (s *service) executeFreeOrder(
	ctx context.Context, order *entity.Order, params *checkoutParams,
) (string, func(context.Context), error) {
	// 注文履歴レコードの登録
	order.SetTransaction("", s.now())
	if err := s.db.Order.Create(ctx, order); err != nil {
		return "", nil, internalError(err)
	}
	// 支払い完了通知
	afterFn := func(ctx context.Context) {
		if err := s.notifyPaymentCompleted(ctx, order); err != nil {
			s.logger.Error("Failed to notify payment completed", zap.Error(err))
		}
	}
	redirectURL := fmt.Sprintf("%s?session_id=%s", params.payload.CallbackURL, order.ID)
	return redirectURL, afterFn, nil
}

func (s *service) getShippingByCoordinatorID(
	ctx context.Context,
	coordinatorID string,
) (*entity.Shipping, error) {
	shipping, err := s.db.Shipping.GetByCoordinatorID(ctx, coordinatorID)
	if errors.Is(err, database.ErrNotFound) {
		// 配送設定が完了していないコーディネータの場合、デフォルト設定を使用
		return s.db.Shipping.GetDefault(ctx)
	}
	return shipping, err
}

func (s *service) decreaseProductInventories(ctx context.Context, items entity.OrderItems) {
	sem := semaphore.NewWeighted(3)
	for _, item := range items {
		if err := sem.Acquire(ctx, 1); err != nil {
			s.logger.Error("Failed to semaphore acuire", zap.Any("item", item), zap.Error(err))
			return
		}
		s.waitGroup.Add(1)
		go func(item *entity.OrderItem) {
			defer func() {
				sem.Release(1)
				s.waitGroup.Done()
			}()
			err := s.decreaseProductInventory(ctx, item.ProductRevisionID, item.Quantity)
			if err != nil {
				s.logger.Error(
					"Failed to decrease product inventory",
					zap.Any("item", item),
					zap.Error(err),
				)
			}
		}(item)
	}
}

func (s *service) decreaseProductInventory(ctx context.Context, revisionID, quantity int64) error {
	exec := func() error {
		return s.db.Product.DecreaseInventory(ctx, revisionID, quantity)
	}
	const maxRetries = 5
	retry := backoff.NewExponentialBackoff(maxRetries)
	return backoff.Retry(ctx, retry, exec, backoff.WithRetryablel(s.isRetryable))
}
