package stripe

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/payment"
	lib "github.com/stripe/stripe-go/v82"
)

func (p *provider) ShowPayment(ctx context.Context, paymentID string) (*payment.PaymentResult, error) {
	pi, err := p.client.GetPaymentIntent(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	return &payment.PaymentResult{
		Status: convertPaymentIntentStatus(pi.Status),
	}, nil
}

func (p *provider) CapturePayment(ctx context.Context, paymentID string) error {
	_, err := p.client.Capture(ctx, paymentID)
	return err
}

func (p *provider) CancelPayment(ctx context.Context, paymentID string) error {
	_, err := p.client.Cancel(ctx, paymentID, lib.PaymentIntentCancellationReasonRequestedByCustomer)
	return err
}

func (p *provider) RefundPayment(ctx context.Context, params *payment.RefundParams) error {
	_, err := p.client.Refund(ctx, params.PaymentID, params.Amount, params.Description)
	return err
}
