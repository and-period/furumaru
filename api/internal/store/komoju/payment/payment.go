package payment

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/store/komoju"
	"go.uber.org/zap"
)

type client struct {
	client *komoju.APIClient
	logger *zap.Logger
	host   string
}

type Params struct {
	Logger       *zap.Logger
	Host         string // KOMOJU接続用URL
	ClientID     string // KOMOJU接続時のBasic認証ユーザー名
	ClientSecret string // KOMOJU接続時のBasic認証パスワード
}

func NewClient(cli *http.Client, params *Params, opts ...komoju.Option) komoju.Payment {
	return &client{
		client: komoju.NewAPIClient(cli, params.ClientID, params.ClientSecret, opts...),
		logger: params.Logger,
		host:   params.Host,
	}
}

func (c *client) Show(ctx context.Context, paymentID string) (*komoju.PaymentResponse, error) {
	const path = "/api/v1/payments/%s"
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodGet,
		Path:   path,
		Params: []interface{}{paymentID},
	}
	res := &komoju.PaymentResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type captureRequest struct {
	Amount int64 `json:"amount,omitempty"`
	Tax    int64 `json:"tax,omitempty"`
}

func (c *client) Capture(ctx context.Context, paymentID string) (*komoju.PaymentResponse, error) {
	const path = "/api/v1/payments/%s/capture"
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{paymentID},
		Body:   &captureRequest{},
	}
	res := &komoju.PaymentResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Cancel(ctx context.Context, paymentID string) (*komoju.PaymentResponse, error) {
	const path = "/api/v1/payments/%s/cancel"
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{paymentID},
	}
	res := &komoju.PaymentResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type refundRequest struct {
	Amount      int64  `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *client) Refund(ctx context.Context, params *komoju.RefundParams) (*komoju.PaymentResponse, error) {
	const path = "/api/v1/payments/%s/refund"
	body := &refundRequest{
		Amount:      params.Amount,
		Description: params.Description,
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.PaymentID},
		Body:   body,
	}
	res := &komoju.PaymentResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type refundRequestRequest struct {
	Amount                  int64  `json:"amount"`
	CustomerName            string `json:"customer_name"`
	BankName                string `json:"bank_name"`
	BankCode                string `json:"bank_code,omitempty"`
	BranchName              string `json:"branch_name,omitempty"`
	BranchNumber            string `json:"branch_number"`
	AccountType             string `json:"account_type"`
	AccountNumber           string `json:"account_number"`
	IncludePaymentMethodFee bool   `json:"include_payment_method_fee"`
	Description             string `json:"description"`
}

func (c *client) RefundRequest(ctx context.Context, params *komoju.RefundRequestParams) (*komoju.PaymentResponse, error) {
	const path = "/api/v1/payments/%s/refund_request"
	body := &refundRequestRequest{
		Amount:                  params.Amount,
		CustomerName:            params.CustomerName,
		BankName:                params.BankName,
		BankCode:                params.BankCode,
		BranchName:              params.BranchName,
		BranchNumber:            params.BranchNumber,
		AccountType:             string(params.AccountType),
		AccountNumber:           params.AccountNumber,
		IncludePaymentMethodFee: params.IncludePaymentMethodFee,
		Description:             params.Description,
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.PaymentID},
		Body:   body,
	}
	res := &komoju.PaymentResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}
