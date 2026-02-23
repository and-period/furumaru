package komoju

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/store/payment"
)

func (p *provider) ShowPayment(ctx context.Context, paymentID string) (*payment.PaymentResult, error) {
	const path = "/api/v1/payments/%s"
	req := &apiParams{
		Host:   p.host,
		Method: http.MethodGet,
		Path:   path,
		Params: []interface{}{paymentID},
	}
	res := &paymentResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.PaymentResult{
		Status: convertPaymentStatus(res.Status),
	}, nil
}

type captureRequest struct {
	Amount int64 `json:"amount,omitempty"`
	Tax    int64 `json:"tax,omitempty"`
}

func (p *provider) CapturePayment(ctx context.Context, paymentID string) error {
	const path = "/api/v1/payments/%s/capture"
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{paymentID},
		Body:           &captureRequest{},
		IdempotencyKey: paymentID,
	}
	res := &paymentResponse{}
	return p.client.do(ctx, req, res)
}

func (p *provider) CancelPayment(ctx context.Context, paymentID string) error {
	const path = "/api/v1/payments/%s/cancel"
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{paymentID},
		IdempotencyKey: paymentID,
	}
	res := &paymentResponse{}
	return p.client.do(ctx, req, res)
}

type refundRequestBody struct {
	Amount      int64  `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
}

func (p *provider) RefundPayment(ctx context.Context, params *payment.RefundParams) error {
	const path = "/api/v1/payments/%s/refund"
	body := &refundRequestBody{
		Amount:      params.Amount,
		Description: params.Description,
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.PaymentID},
		Body:           body,
		IdempotencyKey: params.PaymentID,
	}
	res := &paymentResponse{}
	return p.client.do(ctx, req, res)
}
