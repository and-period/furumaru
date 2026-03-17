package stripe

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/payment"
	pkgstripe "github.com/and-period/furumaru/api/pkg/stripe"
)

func (p *provider) CreateSession(ctx context.Context, params *payment.CreateSessionParams) (*payment.CreateSessionResult, error) {
	successURL := params.CallbackURL + "?session_id={CHECKOUT_SESSION_ID}"
	cancelURL := params.CallbackURL + "?canceled=true"
	in := &pkgstripe.CreateCheckoutSessionParams{
		Amount:      params.Amount,
		Currency:    "jpy",
		Description: params.OrderID,
		SuccessURL:  successURL,
		CancelURL:   cancelURL,
		Email:       params.Customer.Email,
		Metadata: map[string]string{
			"order_id": params.OrderID,
		},
	}
	cs, err := p.client.CreateCheckoutSession(ctx, in)
	if err != nil {
		return nil, err
	}
	return &payment.CreateSessionResult{
		SessionID:  cs.PaymentIntent.ID,
		SessionURL: cs.URL,
	}, nil
}

func (p *provider) GetSession(ctx context.Context, sessionID string) (*payment.GetSessionResult, error) {
	pi, err := p.client.GetPaymentIntent(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return &payment.GetSessionResult{
		PaymentStatus: convertPaymentIntentStatus(pi.Status),
	}, nil
}

func (p *provider) IsSessionFailed(err error) bool {
	return isSessionFailed(err)
}
