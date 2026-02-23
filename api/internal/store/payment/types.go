package payment

import "github.com/and-period/furumaru/api/internal/store/entity"

// CreateSessionParams contains parameters for creating a payment session.
type CreateSessionParams struct {
	OrderID         string                 // 注文ID（冪等キーとして使用）
	Amount          int64                  // 支払い金額
	CallbackURL     string                 // 支払い後リダイレクトURL
	PaymentMethodType entity.PaymentMethodType // 決済手段
	Customer        *CreateSessionCustomer // 顧客情報
	BillingAddress  *CreateSessionAddress  // 請求先住所
	ShippingAddress *CreateSessionAddress  // 配送先住所
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

// CreateSessionResult is the result of creating a payment session.
type CreateSessionResult struct {
	SessionID string // セッションID
	ReturnURL string // リダイレクトURL
}

// GetSessionResult is the result of getting a payment session.
type GetSessionResult struct {
	PaymentStatus entity.PaymentStatus // 決済ステータス（provider内部で変換済み）
}

// OrderCreditCardParams contains parameters for credit card payment.
type OrderCreditCardParams struct {
	SessionID         string // セッションID
	Token             string // カードトークン（PCI DSS 4.0 準拠）
	Number            string // カード番号（後方互換）
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
	LastnameKana  string // 氏名（姓：カナ）
	FirstnameKana string // 氏名（名：カナ）
}

// KonbiniType represents a convenience store type.
type KonbiniType string

const (
	KonbiniTypeDailyYamazaki KonbiniType = "daily-yamazaki"
	KonbiniTypeFamilyMart    KonbiniType = "family-mart"
	KonbiniTypeLawson        KonbiniType = "lawson"
	KonbiniTypeMinistop      KonbiniType = "ministop"
	KonbiniTypeSeicomart     KonbiniType = "seicomart"
	KonbiniTypeSevenEleven   KonbiniType = "seven-eleven"
)

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

type OrderPaidyParams struct {
	SessionID string // セッションID
	Email     string // メールアドレス
	Name      string // 氏名
}

type OrderPayEasyParams struct {
	SessionID     string // セッションID
	Email         string // メールアドレス
	PhoneNumber   string // 電話番号
	Lastname      string // 氏名（姓）
	Firstname     string // 氏名（名）
	LastnameKana  string // 氏名（姓：カナ）
	FirstnameKana string // 氏名（名：カナ）
}

// OrderResult is the result of an order operation.
type OrderResult struct {
	RedirectURL string // リダイレクトURL（外部決済の場合）
}

// PaymentResult is the result of a payment query.
type PaymentResult struct {
	Status entity.PaymentStatus // 決済ステータス
}

// RefundParams contains parameters for refunding a payment.
type RefundParams struct {
	PaymentID   string // ペイメントID
	Amount      int64  // 返金金額
	Description string // 返金理由
}
