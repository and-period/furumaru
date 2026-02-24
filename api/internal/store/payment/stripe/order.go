package stripe

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/payment"
)

func (p *provider) OrderCreditCard(ctx context.Context, params *payment.OrderCreditCardParams) (*payment.OrderResult, error) {
	// Stripe ではフロントエンドで PaymentIntent を confirm するため、
	// バックエンド側での追加操作は不要
	return &payment.OrderResult{}, nil
}

func (p *provider) OrderBankTransfer(_ context.Context, _ *payment.OrderBankTransferParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderKonbini(_ context.Context, _ *payment.OrderKonbiniParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderPayPay(_ context.Context, _ *payment.OrderPayPayParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderLinePay(_ context.Context, _ *payment.OrderLinePayParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderMerpay(_ context.Context, _ *payment.OrderMerpayParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderRakutenPay(_ context.Context, _ *payment.OrderRakutenPayParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderAUPay(_ context.Context, _ *payment.OrderAUPayParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderPaidy(_ context.Context, _ *payment.OrderPaidyParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}

func (p *provider) OrderPayEasy(_ context.Context, _ *payment.OrderPayEasyParams) (*payment.OrderResult, error) {
	return nil, ErrNotSupported
}
