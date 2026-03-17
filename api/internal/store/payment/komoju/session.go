package komoju

import (
	"context"
	"net/http"
	"strconv"

	"github.com/and-period/furumaru/api/internal/store/payment"
)

const (
	sessionMode      = "payment"
	defaultExpiresIn = 7200 // 2時間
	defaultCurrency  = "JPY"
	defaultLocale    = "ja"
	defaultCountry   = "日本"
)

type createSessionRequest struct {
	ReturnURL          string                    `json:"return_url,omitempty"`
	Mode               string                    `json:"mode,omitempty"`
	Amount             int64                     `json:"amount,omitempty"`
	Currency           string                    `json:"currency,omitempty"`
	Email              string                    `json:"email,omitempty"`
	ExpiresInSeconds   int64                     `json:"expires_in_seconds,omitempty"`
	ExternalCustomerID string                    `json:"external_customer_id,omitempty"`
	Metadata           string                    `json:"string,omitempty"`
	PaymentTypes       []string                  `json:"payment_types,omitempty"`
	DefaultLocale      string                    `json:"default_locale,omitempty"`
	PaymentData        *createSessionPaymentData `json:"payment_data,omitempty"`
}

type createSessionPaymentData struct {
	Amount              int64                 `json:"amount,omitempty"`
	Currency            string                `json:"currency,omitempty"`
	ExternalOrderNumber string                `json:"external_order_num,omitempty"`
	Name                string                `json:"name,omitempty"`
	NameKana            string                `json:"name_kana,omitempty"`
	ShippingAddress     *createSessionAddress `json:"shipping_address,omitempty"`
	BillingAddress      *createSessionAddress `json:"billing_address,omitempty"`
	Capture             string                `json:"capture,omitempty"`
}

type createSessionAddress struct {
	ZipCode        string `json:"zipcode,omitempty"`
	Country        string `json:"country"`
	State          string `json:"state,omitempty"`
	City           string `json:"city,omitempty"`
	StreetAddress1 string `json:"street_address1,omitempty"`
	StreetAddress2 string `json:"street_address2,omitempty"`
}

func (p *provider) CreateSession(ctx context.Context, params *payment.CreateSessionParams) (*payment.CreateSessionResult, error) {
	const path = "/api/v1/sessions"
	types := paymentTypesFromMethodType(params.PaymentMethodType)
	typeStrs := make([]string, len(types))
	for i := range types {
		typeStrs[i] = string(types[i])
	}
	body := &createSessionRequest{
		ReturnURL:          params.CallbackURL,
		Mode:               sessionMode,
		Amount:             params.Amount,
		Currency:           defaultCurrency,
		Email:              params.Customer.Email,
		ExpiresInSeconds:   defaultExpiresIn,
		ExternalCustomerID: params.Customer.ID,
		PaymentTypes:       typeStrs,
		DefaultLocale:      defaultLocale,
		PaymentData: &createSessionPaymentData{
			Amount:              params.Amount,
			Currency:            defaultCurrency,
			ExternalOrderNumber: params.OrderID,
			Name:                params.Customer.Name,
			NameKana:            params.Customer.NameKana,
			Capture:             string(p.captureMode),
		},
	}
	if params.BillingAddress != nil {
		body.PaymentData.BillingAddress = &createSessionAddress{
			ZipCode:        params.BillingAddress.ZipCode,
			Country:        defaultCountry,
			State:          params.BillingAddress.Prefecture,
			City:           params.BillingAddress.City,
			StreetAddress1: params.BillingAddress.AddressLine1,
			StreetAddress2: params.BillingAddress.AddressLine2,
		}
	}
	if params.ShippingAddress != nil {
		body.PaymentData.ShippingAddress = &createSessionAddress{
			ZipCode:        params.ShippingAddress.ZipCode,
			Country:        defaultCountry,
			State:          params.ShippingAddress.Prefecture,
			City:           params.ShippingAddress.City,
			StreetAddress1: params.ShippingAddress.AddressLine1,
			StreetAddress2: params.ShippingAddress.AddressLine2,
		}
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Body:           body,
		IdempotencyKey: params.OrderID,
	}
	res := &sessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.CreateSessionResult{
		SessionID: res.ID,
		ReturnURL: res.ReturnURL,
	}, nil
}

func (p *provider) GetSession(ctx context.Context, sessionID string) (*payment.GetSessionResult, error) {
	const path = "/api/v1/sessions/%s"
	req := &apiParams{
		Host:   p.host,
		Method: http.MethodGet,
		Path:   path,
		Params: []interface{}{sessionID},
	}
	res := &sessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	var status PaymentStatus
	if res.Payment != nil {
		status = res.Payment.Status
	}
	return &payment.GetSessionResult{
		PaymentStatus: convertPaymentStatus(status),
	}, nil
}

func (p *provider) IsSessionFailed(err error) bool {
	return isSessionFailed(err)
}

// Credit card order request types

type orderCreditCardRequest struct {
	Capture        string      `json:"capture"`
	PaymentDetails interface{} `json:"payment_details"`
}

type creditCardDetails struct {
	Type              string `json:"type"`
	Email             string `json:"email,omitempty"`
	Number            string `json:"number"`
	Month             string `json:"month"`
	Year              string `json:"year"`
	VerificationValue string `json:"verification_value,omitempty"`
	Name              string `json:"name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	ThreeDSecure      bool   `json:"three_d_secure,omitempty"`
}

type creditCardTokenDetails struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func (p *provider) OrderCreditCard(ctx context.Context, params *payment.OrderCreditCardParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	var details interface{}
	if params.Token != "" {
		details = &creditCardTokenDetails{
			Type:  string(PaymentTypeCreditCard),
			Token: params.Token,
		}
	} else {
		details = &creditCardDetails{
			Type:              string(PaymentTypeCreditCard),
			Email:             params.Email,
			Number:            params.Number,
			Month:             strconv.FormatInt(params.Month, 10),
			Year:              strconv.FormatInt(params.Year, 10),
			VerificationValue: params.VerificationValue,
			Name:              params.Name,
			ThreeDSecure:      true,
		}
	}
	body := &orderCreditCardRequest{
		Capture:        string(p.captureMode),
		PaymentDetails: details,
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// Bank transfer

type orderBankTransferRequest struct {
	Capture        string               `json:"capture"`
	PaymentDetails *bankTransferDetails `json:"payment_details"`
}

type bankTransferDetails struct {
	Type           string `json:"type"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone"`
	FamilyName     string `json:"family_name,omitempty"`
	GivenName      string `json:"given_name,omitempty"`
	FamilyNameKana string `json:"family_name_kana,omitempty"`
	GivenNameKana  string `json:"given_name_kana,omitempty"`
}

func (p *provider) OrderBankTransfer(ctx context.Context, params *payment.OrderBankTransferParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderBankTransferRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &bankTransferDetails{
			Type:           string(PaymentTypeBankTransfer),
			Email:          params.Email,
			PhoneNumber:    params.PhoneNumber,
			GivenName:      params.Firstname,
			FamilyName:     params.Lastname,
			GivenNameKana:  params.FirstnameKana,
			FamilyNameKana: params.LastnameKana,
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// Konbini

type orderKonbiniRequest struct {
	Capture        string          `json:"capture"`
	PaymentDetails *konbiniDetails `json:"payment_details"`
}

type konbiniDetails struct {
	Type        string `json:"type"`
	Store       string `json:"store"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone,omitempty"`
}

func (p *provider) OrderKonbini(ctx context.Context, params *payment.OrderKonbiniParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderKonbiniRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &konbiniDetails{
			Type:  string(PaymentTypeKonbini),
			Store: string(params.Store),
			Email: params.Email,
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// PayPay

type orderPayPayRequest struct {
	Capture        string         `json:"capture"`
	PaymentDetails *paypayDetails `json:"payment_details"`
}

type paypayDetails struct {
	Type                string `json:"type"`
	UserAuthorizationID string `json:"user_authorization_id,omitempty"`
	ScannedCode         string `json:"scanned_code,omitempty"`
}

func (p *provider) OrderPayPay(ctx context.Context, params *payment.OrderPayPayParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderPayPayRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &paypayDetails{
			Type: string(PaymentTypePayPay),
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// LINE Pay

type orderLinePayRequest struct {
	Capture        string          `json:"capture"`
	PaymentDetails *linePayDetails `json:"payment_details"`
}

type linePayDetails struct {
	Type        string `json:"type"`
	RegKey      string `json:"reg_key,omitempty"`
	ScannedCode string `json:"scanned_code,omitempty"`
}

func (p *provider) OrderLinePay(ctx context.Context, params *payment.OrderLinePayParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderLinePayRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &linePayDetails{
			Type: string(PaymentTypeLinePay),
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// Merpay

type orderMerpayRequest struct {
	Capture        string         `json:"capture"`
	PaymentDetails *merpayDetails `json:"payment_details"`
}

type merpayDetails struct {
	Type string `json:"type"`
}

func (p *provider) OrderMerpay(ctx context.Context, params *payment.OrderMerpayParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderMerpayRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &merpayDetails{
			Type: string(PaymentTypeMerpay),
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// Rakuten Pay

type orderRakutenPayRequest struct {
	Capture        string             `json:"capture"`
	PaymentDetails *rakutenPayDetails `json:"payment_details"`
}

type rakutenPayDetails struct {
	Type string `json:"type"`
}

func (p *provider) OrderRakutenPay(ctx context.Context, params *payment.OrderRakutenPayParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderRakutenPayRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &rakutenPayDetails{
			Type: string(PaymentTypeRakutenPay),
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// AU Pay

type orderAUPayRequest struct {
	Capture        string        `json:"capture"`
	PaymentDetails *auPayDetails `json:"payment_details"`
}

type auPayDetails struct {
	Type        string `json:"type"`
	ScannedCode string `json:"scanned_code,omitempty"`
}

func (p *provider) OrderAUPay(ctx context.Context, params *payment.OrderAUPayParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderAUPayRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &auPayDetails{
			Type: string(PaymentTypeAUPay),
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// Paidy

type orderPaidyRequest struct {
	Capture        string        `json:"capture"`
	PaymentDetails *paidyDetails `json:"payment_details"`
}

type paidyDetails struct {
	Type         string `json:"type"`
	CustomerName string `json:"customer_name"`
	Email        string `json:"email,omitempty"`
}

func (p *provider) OrderPaidy(ctx context.Context, params *payment.OrderPaidyParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderPaidyRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &paidyDetails{
			Type:         string(PaymentTypePaidy),
			CustomerName: params.Name,
			Email:        params.Email,
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}

// Pay-easy

type orderPayEasyRequest struct {
	Capture        string          `json:"capture"`
	PaymentDetails *payEasyDetails `json:"payment_details"`
}

type payEasyDetails struct {
	Type           string `json:"type"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone"`
	FamilyName     string `json:"family_name,omitempty"`
	GivenName      string `json:"given_name,omitempty"`
	FamilyNameKana string `json:"family_name_kana,omitempty"`
	GivenNameKana  string `json:"given_name_kana,omitempty"`
}

func (p *provider) OrderPayEasy(ctx context.Context, params *payment.OrderPayEasyParams) (*payment.OrderResult, error) {
	const path = "/api/v1/sessions/%s/pay"
	body := &orderPayEasyRequest{
		Capture: string(p.captureMode),
		PaymentDetails: &payEasyDetails{
			Type:           string(PaymentTypePayEasy),
			Email:          params.Email,
			PhoneNumber:    params.PhoneNumber,
			GivenName:      params.Firstname,
			FamilyName:     params.Lastname,
			GivenNameKana:  params.FirstnameKana,
			FamilyNameKana: params.LastnameKana,
		},
	}
	req := &apiParams{
		Host:           p.host,
		Method:         http.MethodPost,
		Path:           path,
		Params:         []interface{}{params.SessionID},
		Body:           body,
		IdempotencyKey: params.SessionID,
	}
	res := &orderSessionResponse{}
	if err := p.client.do(ctx, req, res); err != nil {
		return nil, err
	}
	return &payment.OrderResult{
		RedirectURL: res.RedirectURL,
	}, nil
}
