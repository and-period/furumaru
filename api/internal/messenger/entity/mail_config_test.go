package entity

import (
	"testing"

	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestTemplateBuilder(t *testing.T) {
	now := jst.Date(2022, 1, 2, 18, 30, 0, 0)
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
		OrderFulfillments: []*sentity.OrderFulfillment{{
			ID:                "fulfillment-id",
			OrderID:           "order-id",
			AddressRevisionID: 1,
			Status:            sentity.FulfillmentStatusFulfilled,
			TrackingNumber:    "tracking-number",
			ShippingCarrier:   sentity.ShippingCarrierYamato,
			ShippingType:      sentity.ShippingTypeNormal,
			BoxNumber:         1,
			BoxSize:           sentity.ShippingSize60,
			BoxRate:           100,
		}},
		OrderItems: []*sentity.OrderItem{{
			FulfillmentID:     "fulfillment-id",
			ProductRevisionID: 1,
			OrderID:           "order-id",
			Quantity:          2,
		}},
		ID:              "order-id",
		UserID:          "user-id",
		CoordinatorID:   "coordinator-id",
		PromotionID:     "promotion-id",
		ShippingMessage: "ありがとうございます",
	}
	products := map[int64]*sentity.Product{
		1: {
			ProductRevision: sentity.ProductRevision{
				ID:        1,
				ProductID: "product-id",
				Price:     1000,
				Cost:      500,
			},
			ID:            "product-id",
			CoordinatorID: "coordinator-id",
			ProducerID:    "producer-id",
			TypeID:        "product-type-id",
			TagIDs:        []string{"tag-id"},
			Name:          "おいしいじゃがいも",
			Public:        true,
			ThumbnailURL:  "http://example.com/image.png",
		},
	}
	addresses := map[int64]*uentity.Address{
		1: {
			AddressRevision: uentity.AddressRevision{
				ID:             1,
				AddressID:      "address-id",
				Lastname:       "&.",
				Firstname:      "太郎",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "たろう",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
			ID:        "address-id",
			UserID:    "user-id",
			IsDefault: true,
		},
	}
	t.Run("success", func(t *testing.T) {
		builder := NewTemplateDataBuilder().
			Data(map[string]any{"key": "value"}).
			YearMonth(jst.Date(2022, 1, 2, 18, 30, 0, 0)).
			Name("中村 広大").
			Email("test-user@and-period.jp").
			Password("!Qaz2wsx").
			WebURL("http://example.com").
			Contact("件名", "本文").
			Live("マルシェ", "濵田 海斗", now, now).
			OrderPayment(&order.OrderPayment).
			OrderFulfillment(order.OrderFulfillments, addresses).
			OrderItems(order.OrderItems, products).
			Shipped(order.ShippingMessage)
		data := builder.Build()
		assert.Equal(t, "value", data["key"])
		assert.Equal(t, "2022年01月", data["年月"])
		assert.Equal(t, "中村 広大", data["氏名"])
		assert.Equal(t, "test-user@and-period.jp", data["メールアドレス"])
		assert.Equal(t, "!Qaz2wsx", data["パスワード"])
		assert.Equal(t, "http://example.com", data["サイトURL"])
		assert.Equal(t, "件名", data["件名"])
		assert.Equal(t, "本文", data["本文"])
		assert.Equal(t, "マルシェ", data["タイトル"])
		assert.Equal(t, "濵田 海斗", data["コーディネータ名"])
		assert.Equal(t, "2022-01-02", data["開催日"])
		assert.Equal(t, "18:30", data["開始時間"])
		assert.Equal(t, "18:30", data["終了時間"])
		assert.Equal(t, "order-id", data["注文番号"])
		assert.Equal(t, "クレジットカード決済", data["決済方法"])
		assert.Equal(t, "2000", data["商品金額"])
		assert.Equal(t, "500", data["割引金額"])
		assert.Equal(t, "500", data["配送手数料"])
		assert.Equal(t, "200", data["消費税"])
		assert.Equal(t, "2200", data["合計金額"])
		assert.Equal(t, "1000014", data["郵便番号"])
		assert.Equal(t, "東京都 千代田区 永田町1-7-1", data["住所"])
		assert.Equal(t, "ありがとうございます", data["メッセージ"])
		assert.Equal(t, []map[string]string{{
			"商品名":      "おいしいじゃがいも",
			"サムネイルURL": "http://example.com/image.png",
			"購入数":      "2",
			"商品金額":     "1000",
			"合計金額":     "2000",
		}}, data["商品一覧"])
	})
	t.Run("order fulfillments is empty", func(t *testing.T) {
		builder := NewTemplateDataBuilder().
			OrderFulfillment(sentity.OrderFulfillments{}, map[int64]*uentity.Address{})
		data := builder.Build()
		assert.Nil(t, data["郵便番号"])
		assert.Nil(t, data["住所"])
	})
	t.Run("addresss is empty", func(t *testing.T) {
		builder := NewTemplateDataBuilder().
			OrderFulfillment(order.OrderFulfillments, map[int64]*uentity.Address{0: {}})
		data := builder.Build()
		assert.Nil(t, data["郵便番号"])
		assert.Nil(t, data["住所"])
	})
	t.Run("products is empty", func(t *testing.T) {
		builder := NewTemplateDataBuilder().
			OrderItems(order.OrderItems, map[int64]*sentity.Product{})
		data := builder.Build()
		assert.Equal(t, []map[string]string{{
			"商品名":      "",
			"サムネイルURL": "",
			"購入数":      "2",
			"商品金額":     "0",
			"合計金額":     "0",
		}}, data["商品一覧"])
	})
}
