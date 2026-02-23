//go:generate go tool mockgen -source=$GOFILE -package=mock_payment -destination=./../../../mock/store/payment/provider.go
package payment

import "context"

// Provider defines the interface for payment gateway operations.
// Implementations are provider-specific (e.g., KOMOJU, Stripe).
type Provider interface {
	// Session
	CreateSession(ctx context.Context, params *CreateSessionParams) (*CreateSessionResult, error)
	GetSession(ctx context.Context, sessionID string) (*GetSessionResult, error)

	// Order (per payment method)
	OrderCreditCard(ctx context.Context, params *OrderCreditCardParams) (*OrderResult, error)
	OrderBankTransfer(ctx context.Context, params *OrderBankTransferParams) (*OrderResult, error)
	OrderKonbini(ctx context.Context, params *OrderKonbiniParams) (*OrderResult, error)
	OrderPayPay(ctx context.Context, params *OrderPayPayParams) (*OrderResult, error)
	OrderLinePay(ctx context.Context, params *OrderLinePayParams) (*OrderResult, error)
	OrderMerpay(ctx context.Context, params *OrderMerpayParams) (*OrderResult, error)
	OrderRakutenPay(ctx context.Context, params *OrderRakutenPayParams) (*OrderResult, error)
	OrderAUPay(ctx context.Context, params *OrderAUPayParams) (*OrderResult, error)
	OrderPaidy(ctx context.Context, params *OrderPaidyParams) (*OrderResult, error)
	OrderPayEasy(ctx context.Context, params *OrderPayEasyParams) (*OrderResult, error)

	// Payment lifecycle
	ShowPayment(ctx context.Context, paymentID string) (*PaymentResult, error)
	CapturePayment(ctx context.Context, paymentID string) error
	CancelPayment(ctx context.Context, paymentID string) error
	RefundPayment(ctx context.Context, params *RefundParams) error

	// Error classification
	IsSessionFailed(err error) bool
}
