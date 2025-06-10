package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/exporter"
	"github.com/and-period/furumaru/api/internal/store/exporter/general"
	"github.com/and-period/furumaru/api/internal/store/exporter/sagawa"
	"github.com/and-period/furumaru/api/internal/store/exporter/yamato"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListOrders(ctx context.Context, in *store.ListOrdersInput) (entity.Orders, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListOrdersParams{
		ShopID:   in.ShopID,
		UserID:   in.UserID,
		Types:    in.Types,
		Statuses: in.Statuses,
		Limit:    int(in.Limit),
		Offset:   int(in.Offset),
	}
	var (
		orders entity.Orders
		total  int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		orders, err = s.db.Order.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Order.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return orders, total, nil
}

func (s *service) ListOrderUserIDs(ctx context.Context, in *store.ListOrderUserIDsInput) ([]string, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListOrdersParams{
		ShopID: in.ShopID,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	userIDs, total, err := s.db.Order.ListUserIDs(ctx, params)
	return userIDs, total, internalError(err)
}

func (s *service) GetOrder(ctx context.Context, in *store.GetOrderInput) (*entity.Order, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	return order, internalError(err)
}

func (s *service) GetOrderByTransactionID(ctx context.Context, in *store.GetOrderByTransactionIDInput) (*entity.Order, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	order, err := s.db.Order.GetByTransactionID(ctx, in.UserID, in.TransactionID)
	return order, internalError(err)
}

func (s *service) CaptureOrder(ctx context.Context, in *store.CaptureOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Capturable() {
		return fmt.Errorf("service: this order cannot be capture: %w", exception.ErrFailedPrecondition)
	}
	_, err = s.komoju.Payment.Capture(ctx, order.PaymentID)
	return internalError(err)
}

func (s *service) DraftOrder(ctx context.Context, in *store.DraftOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Preservable() {
		return fmt.Errorf("service: this order cannot be save: %w", exception.ErrFailedPrecondition)
	}
	params := &database.DraftOrderParams{
		ShippingMessage: in.ShippingMessage,
	}
	err = s.db.Order.Draft(ctx, in.OrderID, params)
	return internalError(err)
}

func (s *service) CompleteProductOrder(ctx context.Context, in *store.CompleteProductOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Completable() {
		return fmt.Errorf("service: this order cannot be complete: %w", exception.ErrFailedPrecondition)
	}
	params := &database.CompleteOrderParams{
		ShippingMessage: in.ShippingMessage,
		CompletedAt:     s.now(),
	}
	if err := s.db.Order.Complete(ctx, in.OrderID, params); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		if order.Type == entity.OrderTypeExperience {
			return // 配送しないため通知不要
		}
		in := &messenger.NotifyOrderShippedInput{
			OrderID: order.ID,
		}
		if err := s.messenger.NotifyOrderShipped(context.Background(), in); err != nil {
			s.logger.Error("Failed to notify order shipped", zap.String("orderId", order.ID), zap.Error(err))
		}
	}()
	return nil
}

func (s *service) CompleteExperienceOrder(ctx context.Context, in *store.CompleteExperienceOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Completable() {
		return fmt.Errorf("service: this order cannot be complete: %w", exception.ErrFailedPrecondition)
	}
	params := &database.CompleteOrderParams{
		CompletedAt: s.now(),
	}
	if err := s.db.Order.Complete(ctx, in.OrderID, params); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		in := &messenger.NotifyReviewRequestInput{
			OrderID: order.ID,
		}
		if err := s.messenger.NotifyReviewRequest(context.Background(), in); err != nil {
			s.logger.Error("Failed to notify review request", zap.String("orderId", order.ID), zap.Error(err))
		}
	}()
	return nil
}

func (s *service) CancelOrder(ctx context.Context, in *store.CancelOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Cancelable() {
		return fmt.Errorf("service: this order cannot be canceled: %w", exception.ErrFailedPrecondition)
	}
	_, err = s.komoju.Payment.Cancel(ctx, order.PaymentID)
	return internalError(err)
}

func (s *service) RefundOrder(ctx context.Context, in *store.RefundOrderInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	order, err := s.db.Order.Get(ctx, in.OrderID)
	if err != nil {
		return internalError(err)
	}
	if !order.Refundable() {
		return fmt.Errorf("service: this order cannot be refund: %w", exception.ErrFailedPrecondition)
	}
	params := &komoju.RefundParams{
		PaymentID:   order.PaymentID,
		Amount:      order.Total,
		Description: in.Description,
	}
	_, err = s.komoju.Payment.Refund(ctx, params)
	return internalError(err)
}

func (s *service) UpdateOrderFulfillment(ctx context.Context, in *store.UpdateOrderFulfillmentInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateOrderFulfillmentParams{
		Status:          entity.FulfillmentStatusFulfilled,
		ShippingCarrier: in.ShippingCarrier,
		TrackingNumber:  in.TrackingNumber,
		ShippedAt:       s.now(),
	}
	err := s.db.Order.UpdateFulfillment(ctx, in.OrderID, in.FulfillmentID, params)
	return internalError(err)
}

func (s *service) AggregateOrders(ctx context.Context, in *store.AggregateOrdersInput) (*entity.AggregatedOrder, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateOrdersParams{
		ShopID:       in.ShopID,
		CreatedAtGte: in.CreatedAtGte,
		CreatedAtLt:  in.CreatedAtLt,
	}
	order, err := s.db.Order.Aggregate(ctx, params)
	return order, internalError(err)
}

func (s *service) AggregateOrdersByUser(ctx context.Context, in *store.AggregateOrdersByUserInput) (entity.AggregatedUserOrders, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateOrdersByUserParams{
		ShopID:  in.ShopID,
		UserIDs: in.UserIDs,
	}
	orders, err := s.db.Order.AggregateByUser(ctx, params)
	return orders, internalError(err)
}

func (s *service) AggregateOrdersByPaymentMethodType(
	ctx context.Context,
	in *store.AggregateOrdersByPaymentMethodTypeInput,
) (entity.AggregatedOrderPayments, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateOrdersByPaymentMethodTypeParams{
		ShopID:             in.ShopID,
		PaymentMethodTypes: entity.AllPaymentMethodTypes,
		CreatedAtGte:       in.CreatedAtGte,
		CreatedAtLt:        in.CreatedAtLt,
	}
	orders, err := s.db.Order.AggregateByPaymentMethodType(ctx, params)
	return orders, internalError(err)
}

func (s *service) AggregateOrdersByPromotion(
	ctx context.Context,
	in *store.AggregateOrdersByPromotionInput,
) (entity.AggregatedOrderPromotions, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateOrdersByPromotionParams{
		ShopID:       in.ShopID,
		PromotionIDs: in.PromotionIDs,
	}
	orders, err := s.db.Order.AggregateByPromotion(ctx, params)
	return orders, internalError(err)
}

func (s *service) AggregateOrdersByPeriod(
	ctx context.Context,
	in *store.AggregateOrdersByPeriodInput,
) (entity.AggregatedPeriodOrders, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.AggregateOrdersByPeriodParams{
		ShopID:       in.ShopID,
		PeriodType:   in.PeriodType,
		CreatedAtGte: in.CreatedAtGte,
		CreatedAtLt:  in.CreatedAtLt,
	}
	orders, err := s.db.Order.AggregateByPeriod(ctx, params)
	return orders, internalError(err)
}

func (s *service) ExportOrders(ctx context.Context, in *store.ExportOrdersInput) ([]byte, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.ListOrdersParams{
		ShopID:   in.ShopID,
		Statuses: []entity.OrderStatus{entity.OrderStatusPreparing},
	}
	orders, err := s.db.Order.List(ctx, params)
	if err != nil {
		return nil, internalError(err)
	}
	var (
		addresses uentity.Addresses
		products  entity.Products
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if len(orders) == 0 {
			return
		}
		params := &user.MultiGetAddressesByRevisionInput{
			AddressRevisionIDs: orders.AddressRevisionIDs(),
		}
		addresses, err = s.user.MultiGetAddressesByRevision(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		if len(orders) == 0 {
			return
		}
		products, err = s.db.Product.MultiGetByRevision(ectx, orders.ProductRevisionIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, internalError(err)
	}
	payload := &exportOrdersParams{
		orders:    orders,
		addresses: addresses.MapByRevision(),
		products:  products.MapByRevision(),
	}
	buf := &bytes.Buffer{}
	switch in.ShippingCarrier {
	case entity.ShippingCarrierYamato:
		err = s.exportYamatoOrders(buf, payload)
	case entity.ShippingCarrierSagawa:
		err = s.exportSagawaOrders(buf, payload)
	default:
		err = s.exportGeneralOrders(buf, payload)
	}
	if err != nil {
		return nil, internalError(err)
	}
	return buf.Bytes(), nil
}

type exportOrdersParams struct {
	encodingType codes.CharacterEncodingType
	orders       entity.Orders
	addresses    map[int64]*uentity.Address
	products     map[int64]*entity.Product
}

func (s *service) exportGeneralOrders(writer io.Writer, payload *exportOrdersParams) error {
	client := exporter.NewExporter(writer, payload.encodingType)
	if err := client.WriteHeader(&general.Receipt{}); err != nil {
		return err
	}
	params := &general.ReceiptsParams{
		Orders:    payload.orders,
		Addresses: payload.addresses,
		Products:  payload.products,
	}
	receipts := general.NewReceipts(params)
	for i := range receipts {
		if err := client.WriteBody(receipts[i]); err != nil {
			return err
		}
	}
	return client.Flush()
}

func (s *service) exportYamatoOrders(writer io.Writer, payload *exportOrdersParams) error {
	client := exporter.NewExporter(writer, payload.encodingType)
	if err := client.WriteHeader(&yamato.Receipt{}); err != nil {
		return err
	}
	params := &yamato.ReceiptsParams{
		Orders:    payload.orders,
		Addresses: payload.addresses,
		Products:  payload.products,
	}
	receipts := yamato.NewReceipts(params)
	for i := range receipts {
		if err := client.WriteBody(receipts[i]); err != nil {
			return err
		}
	}
	return client.Flush()
}

func (s *service) exportSagawaOrders(writer io.Writer, payload *exportOrdersParams) error {
	client := exporter.NewExporter(writer, payload.encodingType)
	if err := client.WriteHeader(&sagawa.Receipt{}); err != nil {
		return err
	}
	params := &sagawa.ReceiptsParams{
		Orders:    payload.orders,
		Addresses: payload.addresses,
		Products:  payload.products,
	}
	receipts := sagawa.NewReceipts(params)
	for i := range receipts {
		if err := client.WriteBody(receipts[i]); err != nil {
			return err
		}
	}
	return client.Flush()
}
