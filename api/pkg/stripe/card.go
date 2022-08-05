package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/paymentintent"
	"github.com/stripe/stripe-go/v73/paymentmethod"
	"github.com/stripe/stripe-go/v73/setupintent"
	"go.uber.org/zap"
)

type OrderCardParams struct {
	CustomerID  string
	CardID      string
	Amount      int64
	Description string
	Metadata    map[string]string
}

// reference: https://stripe.com/docs/api/payment_methods/customer_list
func (c *client) ListCards(ctx context.Context, customerID string) ([]*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodListParams{
		ListParams: stripe.ListParams{Context: ctx},
		Customer:   stripe.String(customerID),
		Type:       stripe.String(string(stripe.PaymentMethodTypeCard)),
	}
	iter := paymentmethod.List(params)
	if err := iter.Err(); err != nil {
		return nil, err
	}
	res := make([]*stripe.PaymentMethod, 0)
	for iter.Next() {
		res = append(res, iter.PaymentMethod())
	}
	return res, nil
}

// reference: https://stripe.com/docs/api/payment_methods/retrieve
func (c *client) GetCard(ctx context.Context, customerID, cardID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodParams{
		Params:   stripe.Params{Context: ctx},
		Customer: stripe.String(customerID),
		Type:     stripe.String(string(stripe.PaymentMethodTypeCard)),
	}
	return paymentmethod.Get(cardID, params)
}

// reference: https://stripe.com/docs/api/setup_intents/create
func (c *client) SetupCard(ctx context.Context, customerID string) (*stripe.SetupIntent, error) {
	params := &stripe.SetupIntentParams{
		Params:             stripe.Params{Context: ctx},
		Customer:           stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{string(stripe.PaymentMethodTypeCard)}),
		Usage:              stripe.String(string(stripe.PaymentIntentSetupFutureUsageOnSession)),
	}
	return setupintent.New(params)
}

// reference: https://stripe.com/docs/api/payment_intents/create
func (c *client) OrderCard(ctx context.Context, in *OrderCardParams) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Context:  ctx,
			Metadata: in.Metadata,
		},
		Customer:           stripe.String(in.CustomerID),
		Description:        nullString(in.Description),
		Amount:             stripe.Int64(in.Amount),
		Currency:           stripe.String(string(stripe.CurrencyJPY)),
		PaymentMethod:      stripe.String(in.CardID),
		PaymentMethodTypes: stripe.StringSlice([]string{string(stripe.PaymentMethodTypeCard)}),
		CaptureMethod:      stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
	}
	var pi *stripe.PaymentIntent
	orderFn := func() (err error) {
		pi, err = paymentintent.New(params)
		return err
	}
	if err := c.do(ctx, orderFn); err != nil {
		c.logger.Error("Failed to order card",
			zap.String("customerId", in.CustomerID), zap.String("cardId", in.CardID), zap.Error(err))
		return nil, err
	}
	return pi, nil
}
