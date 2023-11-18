package entity

import (
	"testing"

	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestTemplateBuilder(t *testing.T) {
	order := &sentity.Order{
		OrderPayment: sentity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			Status:            sentity.PaymentStatusAuthorized,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			MethodType:        sentity.PaymentMethodTypeCreditCard,
			Subtotal:          2000,
			Discount:          500,
			ShippingFee:       500,
			Tax:               200,
			Total:             2200,
			RefundTotal:       0,
			RefundType:        0,
			RefundReason:      "",
		},
		OrderFulfillments: []*sentity.OrderFulfillment{},
		OrderItems:        []*sentity.OrderItem{},
		ID:                "order-id",
		UserID:            "user-id",
		CoordinatorID:     "coordinator-id",
		PromotionID:       "promotion-id",
	}
	builder := NewTemplateDataBuilder().
		Data(map[string]string{"key": "value"}).
		YearMonth(jst.Date(2022, 1, 2, 18, 30, 0, 0)).
		Name("中村 広大").
		Email("test-user@and-period.jp").
		Password("!Qaz2wsx").
		WebURL("http://example.com").
		Contact("件名", "本文").
		Order(order)
	data := builder.Build()
	assert.Equal(t, "value", data["key"])
	assert.Equal(t, "2022年01月", data["年月"])
	assert.Equal(t, "中村 広大", data["氏名"])
	assert.Equal(t, "test-user@and-period.jp", data["メールアドレス"])
	assert.Equal(t, "!Qaz2wsx", data["パスワード"])
	assert.Equal(t, "http://example.com", data["サイトURL"])
	assert.Equal(t, "件名", data["件名"])
	assert.Equal(t, "本文", data["本文"])
	assert.Equal(t, "クレジットカード決済", data["支払い:決済方法"])
	assert.Equal(t, "2000", data["商品金額"])
	assert.Equal(t, "500", data["割引金額"])
	assert.Equal(t, "500", data["配送手数料"])
	assert.Equal(t, "200", data["消費税"])
	assert.Equal(t, "2200", data["合計金額"])
}
