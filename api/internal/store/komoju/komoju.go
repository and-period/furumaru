//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package komoju

import "context"

type Payment interface {
	Capture(ctx context.Context, paymentID string) (*PaymentResponse, error) // 売上確定処理
	Refund(ctx context.Context, paymentID string) (*PaymentResponse, error)  // 返金処理
}

type PaymentResponse struct{}

type Session interface {
	Show(ctx context.Context, sessionID string) (*SessionResponse, error)                                 // 決済情報の照会
	Create(ctx context.Context, params *CreateSessionParams) (*SessionResponse, error)                    // 決済トランザクションの作成
	Cancel(ctx context.Context, sessionID string) (*SessionResponse, error)                               // 決済キャンセル
	ExecuteCreditCard(ctx context.Context, params *ExecuteCreditCardParams) (*SessionResponse, error)     // クレジット決済依頼
	ExecuteBankTransfer(ctx context.Context, params *ExecuteBankTransferParams) (*SessionResponse, error) // 銀行振込決済依頼
	ExecuteKonbini(ctx context.Context, params *ExecuteKonbiniParams) (*SessionResponse, error)           // コンビニ決済依頼
	ExecutePayPay(ctx context.Context, params *ExecutePayPayParams) (*SessionResponse, error)             // PayPay決済依頼
	ExecuteLinePay(ctx context.Context, params *ExecuteLinePayParams) (*SessionResponse, error)           // LINE Pay決済依頼
	ExecuteMerpay(ctx context.Context, params *ExecuteMerpayParams) (*SessionResponse, error)             // メルペイ決済依頼
	ExecuteRakutenPay(ctx context.Context, params *ExecuteRakutenPayParams) (*SessionResponse, error)     // 楽天ペイ決済依頼
	ExecuteAUPay(ctx context.Context, params *ExecuteAUPayParams) (*SessionResponse, error)               // au PAY決済依頼
}

type CreateSessionParams struct {
	OrderID         string                  // 支払いID（ふるマル）
	Amount          int64                   // 支払い金額
	CallbackURL     string                  // 支払い後リダイレクトURL
	Products        []*CreateSessionProduct // 商品情報
	Customer        *CreateSessionCustomer  // 顧客情報
	BillingAddress  *CreateSessionAddress   // 請求先住所
	ShippingAddress *CreateSessionAddress   // 配送先住所
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

type ExecuteCreditCardParams struct{}

type ExecuteBankTransferParams struct{}

type ExecuteKonbiniParams struct{}

type ExecutePayPayParams struct{}

type ExecuteLinePayParams struct{}

type ExecuteMerpayParams struct{}

type ExecuteRakutenPayParams struct{}

type ExecuteAUPayParams struct{}

type SessionResponse struct{}

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
