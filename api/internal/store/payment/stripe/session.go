package stripe

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/payment"
	pkgstripe "github.com/and-period/furumaru/api/pkg/stripe"
	lib "github.com/stripe/stripe-go/v82"
)

func (p *provider) CreateSession(ctx context.Context, params *payment.CreateSessionParams) (*payment.CreateSessionResult, error) {
	in := &pkgstripe.GuestOrderParams{
		Email:             params.Customer.Email,
		PaymentMethodType: lib.PaymentMethodTypeCard,
		Amount:            params.Amount,
		Description:       params.OrderID,
		Metadata: map[string]string{
			"order_id": params.OrderID,
		},
	}
	pi, err := p.client.GuestOrder(ctx, in)
	if err != nil {
		return nil, err
	}
	return &payment.CreateSessionResult{
		SessionID: pi.ID,
		ReturnURL: pi.ClientSecret,
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
