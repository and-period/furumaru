package stripe

import (
	"context"
	"log/slog"

	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/stripe/stripe-go/v82"
)

type OrderParams struct {
	CustomerID        string
	PaymentMethodID   string
	PaymentMethodType stripe.PaymentMethodType
	Amount            int64
	Description       string
	Metadata          map[string]string
}

type GuestOrderParams struct {
	Email             string
	PaymentMethodType stripe.PaymentMethodType
	Amount            int64
	Description       string
	Metadata          map[string]string
}

// reference: https://stripe.com/docs/api/payment_methods/attach
func (c *client) AttachPayment(ctx context.Context, customerID, paymentID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Params:   stripe.Params{Context: ctx},
		Customer: stripe.String(customerID),
	}
	var pm *stripe.PaymentMethod
	attachFn := func() (err error) {
		pm, err = c.paymentmethod.Attach(paymentID, params)
		return err
	}
	if err := c.do(ctx, attachFn); err != nil {
		slog.ErrorContext(ctx, "Failed to attach payment",
			slog.String("customerId", customerID), slog.String("paymentMethodId", paymentID), log.Error(err))
		return nil, err
	}
	return pm, nil
}

// reference: https://stripe.com/docs/api/payment_methods/detach
func (c *client) DetachPayment(ctx context.Context, customerID, paymentID string) error {
	params := &stripe.PaymentMethodDetachParams{
		Params: stripe.Params{Context: ctx},
	}
	if _, err := c.paymentmethod.Detach(paymentID, params); err != nil {
		slog.ErrorContext(ctx, "Failed to detach payment",
			slog.String("customerId", customerID), slog.String("paymentMethodId", paymentID), log.Error(err))
		return err
	}
	return nil
}

// reference: https://stripe.com/docs/api/customers/update
func (c *client) UpdateDefaultPayment(ctx context.Context, customerID, paymentID string) error {
	params := &stripe.CustomerParams{
		Params: stripe.Params{Context: ctx},
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentID),
		},
	}
	if _, err := c.customer.Update(customerID, params); err != nil {
		slog.ErrorContext(ctx, "Failed to update default payment method",
			slog.String("customerId", customerID), slog.String("paymentMethodId", paymentID), log.Error(err))
		return err
	}
	return nil
}

// reference: https://stripe.com/docs/api/payment_intents/create
func (c *client) Order(ctx context.Context, in *OrderParams) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Context:  ctx,
			Metadata: in.Metadata,
		},
		Customer:           stripe.String(in.CustomerID),
		Description:        nullString(in.Description),
		Amount:             stripe.Int64(in.Amount),
		Currency:           stripe.String(string(stripe.CurrencyJPY)),
		PaymentMethod:      stripe.String(in.PaymentMethodID),
		PaymentMethodTypes: stripe.StringSlice([]string{string(in.PaymentMethodType)}),
		CaptureMethod:      stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
		UseStripeSDK:       stripe.Bool(true),
	}
	var pi *stripe.PaymentIntent
	orderFn := func() (err error) {
		pi, err = c.paymentintent.New(params)
		return err
	}
	if err := c.do(ctx, orderFn); err != nil {
		slog.ErrorContext(ctx, "Failed to order",
			slog.String("customerId", in.CustomerID),
			slog.String("paymentMethodId", in.PaymentMethodID),
			slog.String("paymentMethodType", string(in.PaymentMethodType)),
			log.Error(err))
		return nil, err
	}
	return pi, nil
}

// reference: https://stripe.com/docs/api/payment_intents/create
func (c *client) GuestOrder(ctx context.Context, in *GuestOrderParams) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Context:  ctx,
			Metadata: in.Metadata,
		},
		ReceiptEmail:       stripe.String(in.Email),
		Description:        nullString(in.Description),
		Amount:             stripe.Int64(in.Amount),
		Currency:           stripe.String(string(stripe.CurrencyJPY)),
		PaymentMethodTypes: stripe.StringSlice([]string{string(in.PaymentMethodType)}),
		CaptureMethod:      stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
		UseStripeSDK:       stripe.Bool(true),
	}
	var pi *stripe.PaymentIntent
	orderFn := func() (err error) {
		pi, err = c.paymentintent.New(params)
		return err
	}
	if err := c.do(ctx, orderFn); err != nil {
		slog.ErrorContext(ctx, "Failed to guest order",
			slog.String("paymentMethodType", string(in.PaymentMethodType)),
			log.Error(err))
		return nil, err
	}
	return pi, nil
}

// reference: https://stripe.com/docs/api/payment_intents/capture
func (c *client) Capture(ctx context.Context, transactionID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentCaptureParams{
		Params: stripe.Params{Context: ctx},
	}
	var pi *stripe.PaymentIntent
	captureFn := func() (err error) {
		pi, err = c.paymentintent.Capture(transactionID, params)
		return err
	}
	if err := c.do(ctx, captureFn); err != nil {
		slog.ErrorContext(ctx, "Failed to capture", slog.String("transactionId", transactionID), log.Error(err))
		return nil, err
	}
	return pi, nil
}

// reference: https://stripe.com/docs/api/payment_intents/retrieve
func (c *client) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{Context: ctx},
	}
	var pi *stripe.PaymentIntent
	getFn := func() (err error) {
		pi, err = c.paymentintent.Get(paymentIntentID, params)
		return err
	}
	if err := c.do(ctx, getFn); err != nil {
		slog.ErrorContext(ctx, "Failed to get payment intent",
			slog.String("paymentIntentId", paymentIntentID), log.Error(err))
		return nil, err
	}
	return pi, nil
}

// reference: https://stripe.com/docs/api/payment_intents/cancel
func (c *client) Cancel(
	ctx context.Context, transactionID string, reason stripe.PaymentIntentCancellationReason,
) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentCancelParams{
		Params:             stripe.Params{Context: ctx},
		CancellationReason: nullString(string(reason)),
	}
	var pi *stripe.PaymentIntent
	cancelFn := func() (err error) {
		pi, err = c.paymentintent.Cancel(transactionID, params)
		return err
	}
	if err := c.do(ctx, cancelFn); err != nil {
		slog.ErrorContext(ctx, "Failed to cancel", slog.String("transactionId", transactionID), log.Error(err))
		return nil, err
	}
	return pi, nil
}
