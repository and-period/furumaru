package komoju

import "time"

// PaymentStatus represents a KOMOJU payment status.
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusCaptured   PaymentStatus = "captured"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusExpired    PaymentStatus = "expired"
	PaymentStatusFailed     PaymentStatus = "failed"
)

// SessionStatus represents a KOMOJU session status.
type SessionStatus string

const (
	SessionStatusPending   SessionStatus = "pending"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

// PaymentType represents a KOMOJU payment type.
type PaymentType string

const (
	PaymentTypeCreditCard   PaymentType = "credit_card"
	PaymentTypeBankTransfer PaymentType = "bank_transfer"
	PaymentTypeKonbini      PaymentType = "konbini"
	PaymentTypeLinePay      PaymentType = "linepay"
	PaymentTypeMerpay       PaymentType = "merpay"
	PaymentTypePayPay       PaymentType = "paypay"
	PaymentTypeRakutenPay   PaymentType = "rakutenpay"
	PaymentTypeAUPay        PaymentType = "aupay"
	PaymentTypePaidy        PaymentType = "paidy"
	PaymentTypePayEasy      PaymentType = "pay_easy"
)

// CaptureMode represents how payments are captured.
type CaptureMode string

const (
	CaptureModeManual CaptureMode = "manual"
	CaptureModeAuto   CaptureMode = "auto"
)

// paymentResponse is the KOMOJU API response for payment operations.
type paymentResponse struct {
	paymentInfo
}

type paymentInfo struct {
	ID                  string           `json:"id"`
	Resource            string           `json:"resource"`
	Status              PaymentStatus    `json:"status"`
	Amount              int64            `json:"amount"`
	Tax                 int64            `json:"tax"`
	Customer            string           `json:"customer,omitempty"`
	PaymentDeadline     time.Time        `json:"payment_deadline,omitempty"`
	PaymentDetails      *paymentDetails  `json:"payment_details,omitempty"`
	PaymentMethodFee    int64            `json:"payment_method_fee"`
	Total               int64            `json:"total"`
	Currency            string           `json:"currency"`
	Description         string           `json:"description,omitempty"`
	CapturedAt          time.Time        `json:"captured_at,omitempty"`
	ExternalOrderNumber string           `json:"external_order_num,omitempty"`
	CreatedAt           time.Time        `json:"created_at"`
	AmountRefunded      int64            `json:"amount_refunded"`
	Locale              string           `json:"locale"`
	Session             string           `json:"session"`
	CustomerFamilyName  string           `json:"customer_family_name,omitempty"`
	CustomerGivenName   string           `json:"customer_given_name,omitempty"`
	Refunds             []*refund        `json:"refunds"`
	RefundRequests      []*refundRequest `json:"refund_requests"`
}

type paymentDetails struct {
	Type               string    `json:"type"`
	Email              string    `json:"email,omitempty"`
	OrderID            string    `json:"order_id,omitempty"`
	Brand              string    `json:"brand,omitempty"`
	LastFourDigits     string    `json:"last_four_digits,omitempty"`
	Month              int64     `json:"month,omitempty"`
	Year               int64     `json:"year,omitempty"`
	Field55            string    `json:"field55,omitempty"`
	CardnetStatus      string    `json:"cardnet_status,omitempty"`
	ApprovalNumber     string    `json:"approval_number,omitempty"`
	Store              string    `json:"store,omitempty"`
	ConfirmationCode   string    `json:"confirmation_code,omitempty"`
	Receipt            string    `json:"receipt,omitempty"`
	BankID             string    `json:"bank_id,omitempty"`
	BankName           string    `json:"bank_name,omitempty"`
	CustomerID         string    `json:"customer_id,omitempty"`
	AccountBranchName  string    `json:"account_branch_name,omitempty"`
	AccountNumber      string    `json:"account_number,omitempty"`
	AccountType        string    `json:"account_type,omitempty"`
	AccountName        string    `json:"account_name,omitempty"`
	InstructionsURL    string    `json:"instructions_url,omitempty"`
	PaymentDeadline    time.Time `json:"payment_deadline,omitempty"`
	ChargeKey          string    `json:"charge_key,omitempty"`
	RedirectURL        string    `json:"redirect_url,omitempty"`
	ConfirmationID     string    `json:"confirmation_id,omitempty"`
	ExternalPaymentID  string    `json:"external_payment_id,omitempty"`
	TransactionKey     string    `json:"transaction_key,omitempty"`
	PaymentURL         string    `json:"payment_url,omitempty"`
	PaymentURLApp      string    `json:"payment_url_app,omitempty"`
	PaymentAccessToken string    `json:"payment_access_token,omitempty"`
	CVSCode            string    `json:"cvs_code,omitempty"`
}

type refund struct {
	ID          string    `json:"id"`
	Resource    string    `json:"resource"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Payment     string    `json:"payment"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Chargeback  bool      `json:"chargeback"`
}

type refundRequest struct {
	ID            string `json:"id"`
	Payment       string `json:"payment"`
	CustomerName  string `json:"customer_name"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	BranchName    string `json:"branch_name"`
	BranchNumber  string `json:"branch_number"`
	AccountNumber string `json:"account_number"`
	Status        string `json:"status"`
	Description   string `json:"description"`
}

// sessionResponse is the KOMOJU API response for session operations.
type sessionResponse struct {
	ID             string             `json:"id"`
	Resource       string             `json:"resource"`
	Mode           string             `json:"mode"`
	Amount         int64              `json:"amount"`
	Currency       string             `json:"currency"`
	SessionURL     string             `json:"session_url"`
	ReturnURL      string             `json:"return_url"`
	DefaultLocale  string             `json:"default_locale,omitempty"`
	PaymentMethods []*paymentMethod   `json:"payment_methods"`
	CreatedAt      time.Time          `json:"created_at"`
	CancelledAt    time.Time          `json:"cancelled_at,omitempty"`
	CompletedAt    time.Time          `json:"completed_at,omitempty"`
	Status         SessionStatus      `json:"status"`
	Expired        bool               `json:"expired"`
	Payment        *paymentInfo       `json:"payment,omitempty"`
	SecureToken    *secureToken       `json:"secure_token,omitempty"`
	PaymentData    *paymentData       `json:"payment_data"`
	CustomerID     string             `json:"customer_id,omitempty"`
	LineItems      []*sessionLineItem `json:"line_items,omitempty"`
}

type orderSessionResponse struct {
	RedirectURL string        `json:"redirect_url,omitempty"`
	Status      SessionStatus `json:"status"`
	Payment     *paymentInfo  `json:"payment"`
	AppURL      string        `json:"app_url,omitempty"`
}

type sessionLineItem struct {
	Amount      int64  `json:"amount"`
	Quantity    int64  `json:"quantity"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type paymentMethod struct {
	Type             string        `json:"type"`
	Brands           []interface{} `json:"brands,omitempty"`
	HashedGateway    string        `json:"hashed_gateway"`
	Offsite          bool          `json:"offsite,omitempty"`
	AdditionalFields []string      `json:"additional_fields,omitempty"`
	Amount           int64         `json:"amount,omitempty"`
	Currency         string        `json:"currency,omitempty"`
	ExchangeRate     float64       `json:"exchange_rate,omitempty"`
}

type paymentData struct {
	ExternalOrderNumber string `json:"external_order_num,omitempty"`
	Name                string `json:"name"`
	NameKana            string `json:"name_kana"`
	Capture             string `json:"capture"`
}

type secureToken struct {
	ID                 string    `json:"id,omitempty"`
	VerificationStatus string    `json:"verification_status,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
}
