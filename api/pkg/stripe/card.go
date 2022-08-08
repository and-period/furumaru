package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/paymentmethod"
	"github.com/stripe/stripe-go/v73/setupintent"
)

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
