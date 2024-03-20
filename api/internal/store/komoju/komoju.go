//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package komoju

import (
	"context"
	"time"
)

type Payment interface {
	Capture(ctx context.Context, paymentID string) (*PaymentResponse, error)                  // 売上確定処理
	Cancel(ctx context.Context, paymentID string) (*PaymentResponse, error)                   // キャンセル
	Refund(ctx context.Context, params *RefundParams) (*PaymentResponse, error)               // 返金処理
	RefundRequest(ctx context.Context, params *RefundRequestParams) (*PaymentResponse, error) // 返金要求処理
}

type RefundParams struct {
	PaymentID   string // ペイメントID
	Amount      int64  // 返金金額
	Description string // 返金理由
}

type RefundRequestParams struct {
	PaymentID               string          // ペイメントID
	Amount                  int64           // 返金金額
	CustomerName            string          // 顧客名（半角カナ）
	BankName                string          // 金融機関名
	BankCode                string          // 金融機関コード
	BranchName              string          // 支店名
	BranchNumber            string          // 支店コード
	AccountType             BankAccountType // 預金種別
	AccountNumber           string          // 口座番号
	IncludePaymentMethodFee bool            // 手数料を含むか
	Description             string          // 返金理由
}

type PaymentResponse struct {
	*PaymentInfo
}

type PaymentInfo struct {
	ID                  string           `json:"id"`
	Resource            string           `json:"resource"`
	Status              PaymentStatus    `json:"status"`
	Amount              int64            `json:"amount"`
	Tax                 int64            `json:"tax"`
	Cstomer             string           `json:"customer,omitempty"`
	PaymentDeadline     time.Time        `json:"payment_deadline,omitempty"`
	PaymentDetails      *PaymentDetails  `json:"payment_details,omitempty"`
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
	Refunds             []*Refund        `json:"refunds"`
	RefundRequests      []*RefundRequest `json:"refund_requests"`
}

type PaymentDetails struct {
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
	BankName           string    `json:"bank_name,omitempty"`
	AccountBranchName  string    `json:"account_branch_name,omitempty"`
	AccountNumber      string    `json:"account_number,omitempty"`
	AccountType        string    `json:"account_type,omitempty"`
	AccountName        string    `json:"account_name,omitempty"`
	InstructionsURL    string    `json:"instructions_url,omitempty"`
	PaymentDeadline    time.Time `json:"payment_deadline,omitempty"`
	ChargeKey          string    `json:"charge_key,omitempty"`
	RedirectURL        string    `json:"redirect_url,omitempty"`
	ExternalPaymentID  string    `json:"external_payment_id,omitempty"`
	TransactionKey     string    `json:"transaction_key,omitempty"`
	PaymentURLApp      string    `json:"payment_url_app,omitempty"`
	PaymentAccessToken string    `json:"payment_access_token,omitempty"`
}

type Refund struct {
	ID          string    `json:"id"`
	Resource    string    `json:"resource"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Payment     string    `json:"payment"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Chargeback  bool      `json:"chargeback"`
}

type RefundRequest struct {
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

type Session interface {
	Get(ctx context.Context, sessionID string) (*SessionResponse, error)                                   // 決済情報の照会
	Create(ctx context.Context, params *CreateSessionParams) (*SessionResponse, error)                     // 決済トランザクションの作成
	Cancel(ctx context.Context, sessionID string) (*SessionResponse, error)                                // 決済キャンセル
	OrderCreditCard(ctx context.Context, params *OrderCreditCardParams) (*OrderSessionResponse, error)     // クレジット決済
	OrderBankTransfer(ctx context.Context, params *OrderBankTransferParams) (*OrderSessionResponse, error) // 銀行振込決済
	OrderKonbini(ctx context.Context, params *OrderKonbiniParams) (*OrderSessionResponse, error)           // コンビニ決済
	OrderPayPay(ctx context.Context, params *OrderPayPayParams) (*OrderSessionResponse, error)             // PayPay決済
	OrderLinePay(ctx context.Context, params *OrderLinePayParams) (*OrderSessionResponse, error)           // LINE Pay決済
	OrderMerpay(ctx context.Context, params *OrderMerpayParams) (*OrderSessionResponse, error)             // メルペイ決済
	OrderRakutenPay(ctx context.Context, params *OrderRakutenPayParams) (*OrderSessionResponse, error)     // 楽天ペイ決済
	OrderAUPay(ctx context.Context, params *OrderAUPayParams) (*OrderSessionResponse, error)               // au PAY決済
}

type CreateSessionParams struct {
	OrderID         string                 // 支払いID（ふるマル）
	Amount          int64                  // 支払い金額
	CallbackURL     string                 // 支払い後リダイレクトURL
	PaymentTypes    []PaymentType          // 決済種別一覧
	Customer        *CreateSessionCustomer // 顧客情報
	BillingAddress  *CreateSessionAddress  // 請求先住所
	ShippingAddress *CreateSessionAddress  // 配送先住所
}

type CreateSessionProduct struct {
	Amount      int64  // 商品単価
	Description string // 商品詳細
	Quantity    int64  // 商品個数
}

type CreateSessionCustomer struct {
	ID       string // 顧客ID
	Name     string // 顧客名
	NameKana string // 顧客名（かな）
	Email    string // メールアドレス
}

type CreateSessionAddress struct {
	ZipCode      string // 郵便番号
	Prefecture   string // 都道府県
	City         string // 市区町村
	AddressLine1 string // 町名・番地
	AddressLine2 string // ビル名・号室など
}

type OrderCreditCardParams struct {
	SessionID         string // セッションID
	Number            string // カード番号
	Month             int64  // 有効期限（月）
	Year              int64  // 有効期限（年）
	VerificationValue string // セキュリティコード
	Email             string // メールアドレス
	Name              string // 氏名（カード名義）
}

type OrderBankTransferParams struct {
	SessionID     string // セッションID
	Email         string // メールアドレス
	PhoneNumber   string // 電話番号
	Lastname      string // 氏名（姓）
	Firstname     string // 氏名（名）
	LastnameKana  string // 氏名（姓：かな）
	FirstnameKana string // 氏名（名：かな）
}

type OrderKonbiniParams struct {
	SessionID string      // セッションID
	Store     KonbiniType // 店舗種別
	Email     string      // メールアドレス
}

type OrderPayPayParams struct {
	SessionID string // セッションID
}

type OrderLinePayParams struct {
	SessionID string // セッションID
}

type OrderMerpayParams struct {
	SessionID string // セッションID
}

type OrderRakutenPayParams struct {
	SessionID string // セッションID
}

type OrderAUPayParams struct {
	SessionID string // セッションID
}

type SessionResponse struct {
	ID             string             `json:"id"`
	Resource       string             `json:"resource"`
	Mode           string             `json:"mode"`
	Amount         int64              `json:"amount"`
	Currency       string             `json:"currency"`
	SessionURL     string             `json:"session_url"`
	ReturnURL      string             `json:"return_url"`
	DefaultLocale  string             `json:"default_locale,omitempty"`
	PaymentMethods []*PaymentMethod   `json:"payment_methods"`
	CreatedAt      time.Time          `json:"created_at"`
	CancelledAt    time.Time          `json:"cancelled_at,omitempty"`
	CompletedAt    time.Time          `json:"completed_at,omitempty"`
	Status         SessionStatus      `json:"status"`
	Expired        bool               `json:"expired"`
	Payment        *PaymentInfo       `json:"payment,omitempty"`
	SecureToken    *SecureToken       `json:"secure_token,omitempty"`
	PaymentData    *PaymentData       `json:"payment_data"`
	CustomerID     string             `json:"customer_id,omitempty"`
	LineItems      []*SessionLineItem `json:"line_items,omitempty"`
}

type OrderSessionResponse struct {
	RedirectURL string        `json:"redirect_url,omitempty"`
	Status      SessionStatus `json:"status"`
	Payment     *PaymentInfo  `json:"payment"`
	AppURL      string        `json:"app_url,omitempty"`
}

type SessionLineItem struct {
	Amount      int64  `json:"amount"`
	Quantity    int64  `json:"quantity"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type PaymentMethod struct {
	Type             string        `json:"type"`
	Brands           []interface{} `json:"brands,omitempty"`
	HashedGateway    string        `json:"hashed_gateway"`
	Offsite          bool          `json:"offsite,omitempty"`
	AdditionalFields []string      `json:"additional_fields,omitempty"`
	Amount           int64         `json:"amount,omitempty"`
	Currency         string        `json:"currency,omitempty"`
	ExchangeRate     float64       `json:"exchange_rate,omitempty"`
}

type PaymentData struct {
	ExternalOrderNumber string `json:"external_order_num,omitempty"`
	Name                string `json:"name"`
	NameKana            string `json:"name_kana"`
	Capture             string `json:"capture"`
}

type SecureToken struct {
	ID                 string    `json:"id,omitempty"`
	VerificationStatus string    `json:"verification_status,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
}

type Komoju struct {
	Payment Payment
	Session Session
}

type Params struct {
	Payment Payment
	Session Session
}

func NewKomoju(params *Params) *Komoju {
	return &Komoju{
		Payment: params.Payment,
		Session: params.Session,
	}
}
