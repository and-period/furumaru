//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package stripe

import (
	"context"

	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/stripe/stripe-go/v73/paymentintent"
	"github.com/stripe/stripe-go/v73/paymentmethod"
	"github.com/stripe/stripe-go/v73/setupintent"
)

type Client interface {
	// #############################################
	// 顧客関連
	// #############################################
	// 顧客取得
	GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error)
	// 顧客登録
	CreateCustomer(ctx context.Context, in *CreateCustomerParams) (*stripe.Customer, error)
	// 顧客削除
	DeleteCustomer(ctx context.Context, customerID string) error
	// #############################################
	// 決済関連
	// #############################################
	// 決済要求
	Order(ctx context.Context, in *OrderParams) (*stripe.PaymentIntent, error)
	// 決済要求(ゲストユーザー)
	GuestOrder(ctx context.Context, in *GuestOrderParams) (*stripe.PaymentIntent, error)
	// 決済確定
	Capture(ctx context.Context, transactionID string) (*stripe.PaymentIntent, error)
	// 決済キャンセル
	Cancel(ctx context.Context, transactionID string, reason stripe.PaymentIntentCancellationReason) (*stripe.PaymentIntent, error)
	// #############################################
	// 決済方法 (共通)
	// #############################################
	// 顧客と決済手段の関連付け
	AttachPayment(ctx context.Context, customerID, paymentID string) (*stripe.PaymentMethod, error)
	// 顧客と決済手段の関連付けを解除
	DetachPayment(ctx context.Context, customerID, paymentID string) error
	// 顧客のデフォルト決済手段の更新
	UpdateDefaultPayment(ctx context.Context, customerID, paymentID string) error
	// #############################################
	// 決済方法 (クレジットカード)
	// #############################################
	// クレジットカード一覧取得
	ListCards(ctx context.Context, customerID string) ([]*stripe.PaymentMethod, error)
	// クレジットカード取得
	GetCard(ctx context.Context, customerID, paymentID string) (*stripe.PaymentMethod, error)
	// クレジットカード登録用の一時トークンを発行
	SetupCard(ctx context.Context, customerID string) (*stripe.SetupIntent, error)
}

type Receiver interface {
	// 受信イベントの検証
	Receive(payload []byte, signature string) (*stripe.Event, error)
}

type Params struct {
	SecretKey  string
	WebhookKey string
}

type client struct {
	maxRetries    int64
	customer      customer.Client
	paymentintent paymentintent.Client
	paymentmethod paymentmethod.Client
	setupintent   setupintent.Client
}

type options struct {
	maxRetries int64
}

type Option func(opts *options)

func WithMaxRetries(maxRetries int64) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func NewClient(params *Params, opts ...Option) Client {
	dopts := &options{
		maxRetries: 3,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &client{
		maxRetries: dopts.maxRetries,
		customer: customer.Client{
			B:   stripe.GetBackend(stripe.APIBackend),
			Key: params.SecretKey,
		},
		paymentintent: paymentintent.Client{
			B:   stripe.GetBackend(stripe.APIBackend),
			Key: params.SecretKey,
		},
		paymentmethod: paymentmethod.Client{
			B:   stripe.GetBackend(stripe.APIBackend),
			Key: params.SecretKey,
		},
		setupintent: setupintent.Client{
			B:   stripe.GetBackend(stripe.APIBackend),
			Key: params.SecretKey,
		},
	}
}

func (c *client) do(ctx context.Context, fn func() error) error {
	retry := backoff.NewExponentialBackoff(c.maxRetries)
	return backoff.Retry(ctx, retry, fn, backoff.WithRetryablel(isRetryable))
}

func isRetryable(_ error) bool {
	// TODO: 設定する
	return false
}

type receiver struct {
	webhookKey string
}

func NewReceiver(params *Params, opts ...Option) Receiver {
	dopts := &options{}
	for i := range opts {
		opts[i](dopts)
	}
	return &receiver{
		webhookKey: params.WebhookKey,
	}
}

func nullString(val string) *string {
	if val == "" {
		return nil
	}
	return stripe.String(val)
}
