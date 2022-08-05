package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/stripe/stripe-go/v73/paymentintent"
	"github.com/stripe/stripe-go/v73/paymentmethod"
	"go.uber.org/zap"
)

// reference: https://stripe.com/docs/api/payment_methods/attach
func (c *client) AttachPayment(ctx context.Context, customerID, paymentID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Params:   stripe.Params{Context: ctx},
		Customer: stripe.String(customerID),
	}
	var pm *stripe.PaymentMethod
	attachFn := func() (err error) {
		pm, err = paymentmethod.Attach(paymentID, params)
		return err
	}
	if err := c.do(ctx, attachFn); err != nil {
		c.logger.Error("Failed to attach payment",
			zap.String("customerId", customerID), zap.String("paymentMethodId", paymentID), zap.Error(err))
	}
	return pm, nil
}

// reference: https://stripe.com/docs/api/payment_methods/detach
func (c *client) DetachPayment(ctx context.Context, customerID, paymentID string) error {
	params := &stripe.PaymentMethodDetachParams{
		Params: stripe.Params{Context: ctx},
	}
	if _, err := paymentmethod.Detach(paymentID, params); err != nil {
		c.logger.Error("Failed to detach payment",
			zap.String("customerId", customerID), zap.String("paymentMethodId", paymentID), zap.Error(err))
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
	if _, err := customer.Update(customerID, params); err != nil {
		c.logger.Error("Failed to update default payment method",
			zap.String("customerId", customerID), zap.String("paymentMethodId", paymentID), zap.Error(err))
		return err
	}
	return nil
}

// reference: https://stripe.com/docs/api/payment_intents/capture
func (c *client) Capture(ctx context.Context, transactionID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentCaptureParams{
		Params: stripe.Params{Context: ctx},
	}
	var pi *stripe.PaymentIntent
	captureFn := func() (err error) {
		pi, err = paymentintent.Capture(transactionID, params)
		return err
	}
	if err := c.do(ctx, captureFn); err != nil {
		c.logger.Error("Failed to capture", zap.String("transactionId", transactionID), zap.Error(err))
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
		pi, err = paymentintent.Cancel(transactionID, params)
		return err
	}
	if err := c.do(ctx, cancelFn); err != nil {
		c.logger.Error("Failed to cancel", zap.String("transactionId", transactionID), zap.Error(err))
		return nil, err
	}
	return pi, nil
}
