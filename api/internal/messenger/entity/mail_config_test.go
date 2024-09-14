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
			Tax:               181,
			Total:             2000,
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
		OrderExperience: sentity.OrderExperience{
			OrderID:               "order-id",
			ExperienceRevisionID:  1,
			AdultCount:            2,
			JuniorHighSchoolCount: 1,
			ElementarySchoolCount: 0,
			PreschoolCount:        0,
			SeniorCount:           0,
			Remarks: sentity.OrderExperienceRemarks{
				Transportation: "電車",
				RequestedDate:  now,
				RequestedTime:  now,
			},
		},
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
	experience := &sentity.Experience{
		CoordinatorID: "coordinator-id",
		ProducerID:    "producer-id",
		TypeID:        "experience-type-id",
		Title:         "じゃがいも収穫",
		Description:   "じゃがいもを収穫する体験",
		Public:        true,
		SoldOut:       false,
		Status:        sentity.ExperienceStatusAccepting,
		Media: sentity.MultiExperienceMedia{{
			URL:         "http://example.com/thumbnail.png",
			IsThumbnail: true,
		}},
		RecommendedPoints: []string{
			"じゃがいもを収穫する",
			"じゃがいもを食べる",
			"じゃがいもを持ち帰る",
		},
		PromotionVideoURL:  "http://example.com/promotion.mp4",
		Duration:           60,
		Direction:          "彦根駅から徒歩10分",
		BusinessOpenTime:   "1000",
		BusinessCloseTime:  "1800",
		HostPostalCode:     "5220061",
		HostPrefecture:     "滋賀県",
		HostPrefectureCode: 25,
		HostCity:           "彦根市",
		HostAddressLine1:   "金亀町１−１",
		HostAddressLine2:   "",
		HostLongitude:      136.251739,
		HostLatitude:       35.276833,
		StartAt:            now.AddDate(0, -1, 0),
		EndAt:              now.AddDate(0, 1, 0),
		ExperienceRevision: sentity.ExperienceRevision{
			ID:                    0,
			ExperienceID:          "", // ignore
			PriceAdult:            1000,
			PriceJuniorHighSchool: 800,
			PriceElementarySchool: 600,
			PricePreschool:        400,
			PriceSenior:           200,
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
				PhoneNumber:    "090-1234-1234",
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
			OrderBilling(addresses[1]).
			OrderFulfillment(order.OrderFulfillments, addresses).
			OrderItems(order.OrderItems, products).
			OrderExperience(&order.OrderExperience, experience).
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
		assert.Equal(t, "181", data["消費税"])
		assert.Equal(t, "2000", data["合計金額"])
		assert.Equal(t, "1000014", data["郵便番号"])
		assert.Equal(t, "東京都 千代田区 永田町1-7-1", data["住所"])
		assert.Equal(t, "ありがとうございます", data["メッセージ"])
		assert.Equal(t, "じゃがいも収穫", data["体験概要"])
		assert.Equal(t, "2", data["大人人数"])
		assert.Equal(t, "1", data["中学生人数"])
		assert.Equal(t, "0", data["小学生人数"])
		assert.Equal(t, "0", data["幼児人数"])
		assert.Equal(t, "0", data["シニア人数"])
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
