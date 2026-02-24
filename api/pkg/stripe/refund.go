package stripe

import (
	"context"
	"log/slog"

	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/stripe/stripe-go/v82"
)

// reference: https://stripe.com/docs/api/refunds/create
func (c *client) Refund(ctx context.Context, paymentIntentID string, amount int64, reason string) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		Params:        stripe.Params{Context: ctx},
		PaymentIntent: stripe.String(paymentIntentID),
		Amount:        stripe.Int64(amount),
		Reason:        nullString(reason),
	}
	var r *stripe.Refund
	refundFn := func() (err error) {
		r, err = c.refund.New(params)
		return err
	}
	if err := c.do(ctx, refundFn); err != nil {
		slog.ErrorContext(ctx, "Failed to refund",
			slog.String("paymentIntentId", paymentIntentID),
			slog.Int64("amount", amount),
			log.Error(err))
		return nil, err
	}
	return r, nil
}
