package entity

import (
	"strconv"
	"time"

	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// EmailTemplateID - メールテンプレートID
type EmailTemplateID string

const (
	EmailTemplateIDAdminRegister       EmailTemplateID = "admin-register"        // 管理者登録
	EmailTemplateIDAdminResetPassword  EmailTemplateID = "admin-reset-password"  // 管理者パスワードリセット
	EmailTemplateIDUserReceivedContact EmailTemplateID = "user-received-contact" // お問い合わせ受領
	EmailTemplateIDUserOrderAuthorized EmailTemplateID = "user-order-authorized" // 支払い完了
	EmailTemplateIDUserOrderShipped    EmailTemplateID = "user-order-shipped"    // 発送完了
	EmailTemplateIDUserStartLive       EmailTemplateID = "user-start-live"       // ライブ配信開始
)

// MailConfig - メール送信設定
type MailConfig struct {
	TemplateID    EmailTemplateID        `json:"templateId"`    // メールテンプレートID
	Substitutions map[string]interface{} `json:"substitutions"` // メール動的内容
}

type TemplateDataBuilder struct {
	data map[string]interface{}
}

func NewTemplateDataBuilder() *TemplateDataBuilder {
	return &TemplateDataBuilder{
		data: map[string]interface{}{},
	}
}

func (b *TemplateDataBuilder) Build() map[string]any {
	return b.data
}

func (b *TemplateDataBuilder) Data(data map[string]any) *TemplateDataBuilder {
	if data != nil {
		b.data = data
	}
	return b
}

func (b *TemplateDataBuilder) YearMonth(yearMonth time.Time) *TemplateDataBuilder {
	b.data["年月"] = jst.Format(yearMonth, "2006年01月")
	return b
}

func (b *TemplateDataBuilder) Name(name string) *TemplateDataBuilder {
	b.data["氏名"] = name
	return b
}

func (b *TemplateDataBuilder) Email(email string) *TemplateDataBuilder {
	b.data["メールアドレス"] = email
	return b
}

func (b *TemplateDataBuilder) Password(password string) *TemplateDataBuilder {
	b.data["パスワード"] = password
	return b
}

func (b *TemplateDataBuilder) WebURL(url string) *TemplateDataBuilder {
	b.data["サイトURL"] = url
	return b
}

func (b *TemplateDataBuilder) Contact(title, body string) *TemplateDataBuilder {
	b.data["件名"] = title
	b.data["本文"] = body
	return b
}

func (b *TemplateDataBuilder) Live(title, coordinator string, startAt, endAt time.Time) *TemplateDataBuilder {
	b.data["タイトル"] = title
	b.data["コーディネータ名"] = coordinator
	b.data["開催日"] = startAt.Format(time.DateOnly)
	b.data["開始時間"] = startAt.Format("15:04")
	b.data["終了時間"] = endAt.Format("15:04")
	return b
}

func (b *TemplateDataBuilder) Order(order *sentity.Order) *TemplateDataBuilder {
	b.data["注文番号"] = order.ID
	b.data["決済方法"] = newPaymentMethodName(order.OrderPayment.MethodType)
	b.data["商品金額"] = strconv.FormatInt(order.OrderPayment.Subtotal, 10)
	b.data["割引金額"] = strconv.FormatInt(order.OrderPayment.Discount, 10)
	b.data["配送手数料"] = strconv.FormatInt(order.OrderPayment.ShippingFee, 10)
	b.data["消費税"] = strconv.FormatInt(order.OrderPayment.Tax, 10)
	b.data["合計金額"] = strconv.FormatInt(order.OrderPayment.Total, 10)
	return b
}

func (b *TemplateDataBuilder) OrderItems(items sentity.OrderItems, products map[int64]*sentity.Product) *TemplateDataBuilder {
	data := make([]map[string]string, 0, len(items))
	for _, item := range items {
		product, ok := products[item.ProductRevisionID]
		if !ok {
			product = &sentity.Product{}
		}
		data = append(data, newOrderItem(item, product))
	}
	b.data["商品一覧"] = data
	return b
}

func (b *TemplateDataBuilder) Shipped(message string) *TemplateDataBuilder {
	b.data["メッセージ"] = message
	return b
}

/**
 * private
 */
func newPaymentMethodName(typ sentity.PaymentMethodType) string {
	switch typ {
	case sentity.PaymentMethodTypeCash:
		return "代金支払い"
	case sentity.PaymentMethodTypeCreditCard:
		return "クレジットカード決済"
	case sentity.PaymentMethodTypeKonbini:
		return "コンビニ決済"
	case sentity.PaymentMethodTypeBankTranser:
		return "銀行振込決済"
	case sentity.PaymentMethodTypePayPay:
		return "QR決済（PayPay）"
	case sentity.PaymentMethodTypeLinePay:
		return "QR決済（LINE Pay）"
	case sentity.PaymentMethodTypeMerpay:
		return "QR決済（メルペイ）"
	case sentity.PaymentMethodTypeRakutenPay:
		return "QR決済（楽天ペイ）"
	case sentity.PaymentMethodTypeAUPay:
		return "QR決済（au PAY）"
	default:
		return ""
	}
}

func newOrderItem(item *sentity.OrderItem, product *sentity.Product) map[string]string {
	return map[string]string{
		"商品名":      product.Name,
		"サムネイルURL": product.ThumbnailURL,
		"購入数":      strconv.FormatInt(item.Quantity, 10),
		"商品金額":     strconv.FormatInt(product.Price, 10),
		"合計金額":     strconv.FormatInt(product.Price*item.Quantity, 10),
	}
}
