package session

import (
	"context"
	"net/http"
	"strconv"

	"github.com/and-period/furumaru/api/internal/store/komoju"
	"go.uber.org/zap"
)

const (
	sessionMode      = "payment"
	captureMode      = "manual"
	defaultExpiresIn = 7200 // 2時間
	defaultCurrency  = "JPY"
	defaultLocale    = "ja"
	defaultCountry   = "日本"
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

func NewClient(cli *http.Client, params *Params, opts ...komoju.Option) komoju.Session {
	return &client{
		client: komoju.NewAPIClient(cli, params.ClientID, params.ClientSecret, opts...),
		logger: params.Logger,
		host:   params.Host,
	}
}

func (c *client) Get(ctx context.Context, sessionID string) (*komoju.SessionResponse, error) {
	const path = "/api/v1/sessions/%s"
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodGet,
		Path:   path,
		Params: []interface{}{sessionID},
	}
	res := &komoju.SessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type createSessionRequest struct {
	ReturnURL          string                    `json:"return_url,omitempty"`           // 支払い後リダイレクトURL
	Mode               string                    `json:"mode,omitempty"`                 // 支払いAPI種別
	Amount             int64                     `json:"amount,omitempty"`               // 支払い金額
	Currency           string                    `json:"currency,omitempty"`             // 支払い時の通貨
	Email              string                    `json:"email,omitempty"`                // 顧客メールアドレス
	ExpiresInSeconds   int64                     `json:"expires_in_seconds,omitempty"`   // 支払い期限（sec）
	ExternalCustomerID string                    `json:"external_customer_id,omitempty"` // ふるマル顧客ID
	Metadata           string                    `json:"string,omitempty"`               // メタデータ
	PaymentTypes       []string                  `json:"payment_types,omitempty"`        // 支払い可能種別一覧
	DefaultLocale      string                    `json:"default_locale,omitempty"`       // ロケーション
	PaymentData        *createSessionPaymentData `json:"payment_data,omitempty"`         // 購入詳細情報
}

type createSessionPaymentData struct {
	Amount              int64                 `json:"amount,omitempty"`             // 支払い金額
	Currency            string                `json:"currency,omitempty"`           // 支払い時の通貨
	ExternalOrderNumber string                `json:"external_order_num,omitempty"` // ふるマル注文履歴ID
	Name                string                `json:"name,omitempty"`               // 顧客名
	NameKana            string                `json:"name_kana,omitempty"`          // 顧客名（かな）
	ShippingAddress     *createSessionAddress `json:"shipping_address,omitempty"`   // 配送先住所
	BillingAddress      *createSessionAddress `json:"billing_address,omitempty"`    // 請求先住所
	Capture             string                `json:"capture,omitempty"`            // 売上処理
}

type createSessionAddress struct {
	ZipCode        string `json:"zipcode,omitempty"`         // 郵便番号
	Country        string `json:"country"`                   // 国
	State          string `json:"state,omitempty"`           // 都道府県
	City           string `json:"city,omitempty"`            // 市区町村
	StreetAddress1 string `json:"street_address1,omitempty"` // 町名・番地
	StreetAddress2 string `json:"street_address2,omitempty"` // ビル名・号室など
}

func (c *client) Create(ctx context.Context, params *komoju.CreateSessionParams) (*komoju.SessionResponse, error) {
	const path = "/api/v1/sessions"
	types := make([]string, len(params.PaymentTypes))
	for i := range params.PaymentTypes {
		types[i] = string(params.PaymentTypes[i])
	}
	body := &createSessionRequest{
		ReturnURL:          params.CallbackURL,
		Mode:               sessionMode,
		Amount:             params.Amount,
		Currency:           defaultCurrency,
		Email:              params.Customer.Email,
		ExpiresInSeconds:   defaultExpiresIn,
		ExternalCustomerID: params.Customer.ID,
		PaymentTypes:       types,
		DefaultLocale:      defaultLocale,
		PaymentData: &createSessionPaymentData{
			Amount:              params.Amount,
			Currency:            defaultCurrency,
			ExternalOrderNumber: params.OrderID,
			Name:                params.Customer.Name,
			NameKana:            params.Customer.NameKana,
			ShippingAddress: &createSessionAddress{
				ZipCode:        params.ShippingAddress.ZipCode,
				Country:        defaultCountry,
				State:          params.ShippingAddress.Prefecture,
				City:           params.ShippingAddress.City,
				StreetAddress1: params.ShippingAddress.AddressLine1,
				StreetAddress2: params.ShippingAddress.AddressLine2,
			},
			BillingAddress: &createSessionAddress{
				ZipCode:        params.BillingAddress.ZipCode,
				Country:        defaultCountry,
				State:          params.BillingAddress.Prefecture,
				City:           params.BillingAddress.City,
				StreetAddress1: params.BillingAddress.AddressLine1,
				StreetAddress2: params.BillingAddress.AddressLine2,
			},
			Capture: captureMode,
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	}
	res := &komoju.SessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Cancel(ctx context.Context, sessionID string) (*komoju.SessionResponse, error) {
	const path = "/api/v1/sessions/%s/cancel"
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{sessionID},
	}
	res := &komoju.SessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderCreditCardRequest struct {
	Capture        string             `json:"capture"`         // 売上処理
	PaymentDetails *creditCardDetails `json:"payment_details"` // 決済詳細情報
}

type creditCardDetails struct {
	Type              string `json:"type"`                         // 決済種別
	Email             string `json:"email,omitempty"`              // 顧客メールアドレス
	Number            string `json:"number"`                       // クレジットカード番号
	Month             string `json:"month"`                        // 有効期限（月）
	Year              string `json:"year"`                         // 有効期限（年）
	VerificationValue string `json:"verification_value,omitempty"` // セキュリティコード
	Name              string `json:"name,omitempty"`               // 氏名
	FamilyName        string `json:"family_name,omitempty"`        // 氏名（姓）
	GivenName         string `json:"given_name,omitempty"`         // 氏名（名）
	ThreeDSecure      bool   `json:"three_d_secure,omitempty"`     // 3Dセキュアの有効化
}

func (c *client) OrderCreditCard(ctx context.Context, params *komoju.OrderCreditCardParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderCreditCardRequest{
		Capture: captureMode,
		PaymentDetails: &creditCardDetails{
			Type:              string(komoju.PaymentTypeCreditCard),
			Email:             params.Email,
			Number:            params.Number,
			Month:             strconv.FormatInt(params.Month, 10),
			Year:              strconv.FormatInt(params.Year, 10),
			VerificationValue: params.VerificationValue,
			Name:              params.Name,
			ThreeDSecure:      true,
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderBankTransferRequest struct {
	Capture        string               `json:"capture"`         // 売上処理
	PaymentDetails *bankTransferDetails `json:"payment_details"` // 決済詳細情報
}

type bankTransferDetails struct {
	Type           string `json:"type"`                       // 決済種別
	Email          string `json:"email"`                      // 顧客メールアドレス
	PhoneNumber    string `json:"phone"`                      // 顧客電話番号
	FamilyName     string `json:"family_name,omitempty"`      // 氏名（姓）
	GivenName      string `json:"given_name,omitempty"`       // 氏名（名）
	FamilyNameKana string `json:"family_name_kana,omitempty"` // 氏名（姓：かな）
	GivenNameKana  string `json:"given_name_kana,omitempty"`  // 氏名（名：かな）
}

func (c *client) OrderBankTransfer(ctx context.Context, params *komoju.OrderBankTransferParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderBankTransferRequest{
		Capture: captureMode,
		PaymentDetails: &bankTransferDetails{
			Type:           string(komoju.PaymentTypeBankTransfer),
			Email:          params.Email,
			PhoneNumber:    params.PhoneNumber,
			GivenName:      params.Firstname,
			FamilyName:     params.Lastname,
			GivenNameKana:  params.FirstnameKana,
			FamilyNameKana: params.LastnameKana,
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderKonbiniRequest struct {
	Capture        string          `json:"capture"`         // 売上処理
	PaymentDetails *konbiniDetails `json:"payment_details"` // 決済詳細情報
}

type konbiniDetails struct {
	Type        string `json:"type"`            // 決済種別
	Store       string `json:"store"`           // 店舗種別
	Email       string `json:"email"`           // 顧客メールアドレス
	PhoneNumber string `json:"phone,omitempty"` // 顧客電話番号
}

func (c *client) OrderKonbini(ctx context.Context, params *komoju.OrderKonbiniParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderKonbiniRequest{
		Capture: captureMode,
		PaymentDetails: &konbiniDetails{
			Type:  string(komoju.PaymentTypeKonbini),
			Store: string(params.Store),
			Email: params.Email,
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderPayPayRequest struct {
	Capture        string         `json:"capture"`         // 売上処理
	PaymentDetails *paypayDetails `json:"payment_details"` // 決済詳細情報
}

type paypayDetails struct {
	Type                string `json:"type"`                            // 決済種別
	UserAuthorizationID string `json:"user_authorization_id,omitempty"` // 認証ID
	ScannedCode         string `json:"scanned_code,omitempty"`          // バーコード
}

func (c *client) OrderPayPay(ctx context.Context, params *komoju.OrderPayPayParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderPayPayRequest{
		Capture: captureMode,
		PaymentDetails: &paypayDetails{
			Type: string(komoju.PaymentTypePayPay),
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderLinePayRequest struct {
	Capture        string          `json:"capture"`         // 売上処理
	PaymentDetails *linePayDetails `json:"payment_details"` // 決済詳細情報
}

type linePayDetails struct {
	Type        string `json:"type"`                   // 決済種別
	RegKey      string `json:"reg_key,omitempty"`      // 認証キー
	ScannedCode string `json:"scanned_code,omitempty"` // バーコード
}

func (c *client) OrderLinePay(ctx context.Context, params *komoju.OrderLinePayParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderLinePayRequest{
		Capture: captureMode,
		PaymentDetails: &linePayDetails{
			Type: string(komoju.PaymentTypeLinePay),
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderMerpayRequest struct {
	Capture        string         `json:"capture"`         // 売上処理
	PaymentDetails *merpayDetails `json:"payment_details"` // 決済詳細情報
}

type merpayDetails struct {
	Type string `json:"type"` // 決済種別
}

func (c *client) OrderMerpay(ctx context.Context, params *komoju.OrderMerpayParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderMerpayRequest{
		Capture: captureMode,
		PaymentDetails: &merpayDetails{
			Type: string(komoju.PaymentTypeMerpay),
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderRakutenPayRequest struct {
	Capture        string             `json:"capture"`         // 売上処理
	PaymentDetails *rakutenPayDetails `json:"payment_details"` // 決済詳細情報
}

type rakutenPayDetails struct {
	Type string `json:"type"` // 決済種別
}

func (c *client) OrderRakutenPay(ctx context.Context, params *komoju.OrderRakutenPayParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderMerpayRequest{
		Capture: captureMode,
		PaymentDetails: &merpayDetails{
			Type: string(komoju.PaymentTypeRakutenPay),
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

type orderAUPayRequest struct {
	Capture        string        `json:"capture"`         // 売上処理
	PaymentDetails *auPayDetails `json:"payment_details"` // 決済詳細情報
}

type auPayDetails struct {
	Type        string `json:"type"`                   // 決済種別
	ScannedCode string `json:"scanned_code,omitempty"` // バーコード
}

func (c *client) OrderAUPay(ctx context.Context, params *komoju.OrderAUPayParams) (*komoju.OrderSessionResponse, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderAUPayRequest{
		Capture: captureMode,
		PaymentDetails: &auPayDetails{
			Type: string(komoju.PaymentTypeAUPay),
		},
	}
	req := &komoju.APIParams{
		Host:   c.host,
		Method: http.MethodPost,
		Path:   path,
		Params: []interface{}{params.SessionID},
		Body:   body,
	}
	res := &komoju.OrderSessionResponse{}
	if err := c.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}
