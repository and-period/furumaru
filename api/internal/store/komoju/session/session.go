package session

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/store/komoju"
	"go.uber.org/zap"
)

type sessionClient struct {
	client *komoju.APIClient
	logger *zap.Logger
	host   string
}

type SessionParams struct {
	Logger       *zap.Logger
	Host         string // KOMOJU接続用URL
	ClientID     string // KOMOJU接続時のBasic認証ユーザー名
	ClientSecret string // KOMOJU接続時のBasic認証パスワード
}

func NewSessionClient(client *http.Client, params *SessionParams, opts ...komoju.Option) komoju.Session {
	return &sessionClient{
		client: komoju.NewAPIClient(client, params.ClientID, params.ClientSecret, opts...),
		logger: params.Logger,
		host:   params.Host,
	}
}

func (c *sessionClient) Show(_ context.Context, sessionID string) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received show session event", zap.String("sessionId", sessionID))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) Create(_ context.Context, params *komoju.CreateSessionParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received create session event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) Cancel(_ context.Context, sessionID string) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received cancel session event", zap.String("sessionId", sessionID))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecuteCreditCard(_ context.Context, params *komoju.ExecuteCreditCardParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute credit card event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecuteBankTransfer(_ context.Context, params *komoju.ExecuteBankTransferParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute bank transfer event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecuteKonbini(_ context.Context, params *komoju.ExecuteKonbiniParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute konbini event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecutePayPay(_ context.Context, params *komoju.ExecutePayPayParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute paypay event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecuteLinePay(_ context.Context, params *komoju.ExecuteLinePayParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute line pay event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecuteMerpay(_ context.Context, params *komoju.ExecuteMerpayParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute merpay event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecuteRakutenPay(_ context.Context, params *komoju.ExecuteRakutenPayParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute rakuten pay event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}

func (c *sessionClient) ExecuteAUPay(_ context.Context, params *komoju.ExecuteAUPayParams) (*komoju.SessionResponse, error) {
	// TODO: 詳細の実装
	c.logger.Debug("Received execute au pay event", zap.Any("params", params))
	return nil, komoju.ErrNotImplemented
}
