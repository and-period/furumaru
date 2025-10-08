package entity

import (
	"net/url"
	"testing"

	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestTemplateBuilder(t *testing.T) {
	t.Parallel()
	u, err := url.Parse("http://example.com")
	assert.NoError(t, err)
	maker := NewUserURLMaker(u)
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
		OrderMetadata: sentity.OrderMetadata{
			OrderID:         "order-id",
			ShippingMessage: "ありがとうございます",
		},
		ID:            "order-id",
		UserID:        "user-id",
		CoordinatorID: "coordinator-id",
		PromotionID:   "promotion-id",
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
			Scope:         sentity.ProductScopePublic,
			ThumbnailURL:  "http://example.com/image.png",
		},
	}
	experience := &sentity.Experience{
		ID:            "experience-id",
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
		ThumbnailURL:       "http://example.com/thumbnail.png",
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
	tests := []struct {
		name    string
		execute func(builder *TemplateDataBuilder) *TemplateDataBuilder
		expect  map[string]interface{}
	}{
		{
			name: "data",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.Data(map[string]any{"key": "value"})
			},
			expect: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "year month",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.YearMonth(jst.Date(2022, 1, 2, 18, 30, 0, 0))
			},
			expect: map[string]interface{}{
				"年月": "2022年01月",
			},
		},
		{
			name: "name",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.Name("中村 広大")
			},
			expect: map[string]interface{}{
				"氏名": "中村 広大",
			},
		},
		{
			name: "email",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.Email("test@example.com")
			},
			expect: map[string]interface{}{
				"メールアドレス": "test@example.com",
			},
		},
		{
			name: "password",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.Password("!Qaz2wsx")
			},
			expect: map[string]interface{}{
				"パスワード": "!Qaz2wsx",
			},
		},
		{
			name: "web url",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.WebURL("http://example.com")
			},
			expect: map[string]interface{}{
				"サイトURL": "http://example.com",
			},
		},
		{
			name: "contact",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.Contact("件名", "本文")
			},
			expect: map[string]interface{}{
				"件名": "件名",
				"本文": "本文",
			},
		},
		{
			name: "live",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.Live("マルシェ", "濵田 海斗", now, now)
			},
			expect: map[string]interface{}{
				"タイトル":     "マルシェ",
				"コーディネータ名": "濵田 海斗",
				"開催日":      "2022-01-02",
				"開始時間":     "18:30",
				"終了時間":     "18:30",
			},
		},
		{
			name: "order payment",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderPayment(&order.OrderPayment)
			},
			expect: map[string]interface{}{
				"注文番号":  "order-id",
				"決済方法":  "クレジットカード決済",
				"商品金額":  "2000",
				"割引金額":  "500",
				"配送手数料": "500",
				"消費税":   "181",
				"合計金額":  "2000",
			},
		},
		{
			name: "order billing",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderBilling(addresses[1])
			},
			expect: map[string]interface{}{
				"郵便番号": "1000014",
				"住所":   "東京都 千代田区 永田町1-7-1",
			},
		},
		{
			name: "order fulfillment",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderFulfillment(order.OrderFulfillments, addresses)
			},
			expect: map[string]interface{}{
				"郵便番号": "1000014",
				"住所":   "東京都 千代田区 永田町1-7-1",
			},
		},
		{
			name: "order fulfillments is empty",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderFulfillment(sentity.OrderFulfillments{}, map[int64]*uentity.Address{})
			},
			expect: map[string]interface{}{},
		},
		{
			name: "addresss is empty",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderFulfillment(order.OrderFulfillments, map[int64]*uentity.Address{0: {}})
			},
			expect: map[string]interface{}{},
		},
		{
			name: "order items",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderItems(order.OrderItems, products)
			},
			expect: map[string]interface{}{
				"商品一覧": []map[string]string{{
					"商品名":      "おいしいじゃがいも",
					"サムネイルURL": "http://example.com/image.png",
					"購入数":      "2",
					"商品金額":     "1000",
					"合計金額":     "2000",
				}},
			},
		},
		{
			name: "order items is empty",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderItems(sentity.OrderItems{}, map[int64]*sentity.Product{})
			},
			expect: map[string]interface{}{
				"商品一覧": []map[string]string{},
			},
		},
		{
			name: "order experience",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.OrderExperience(&order.OrderExperience, experience)
			},
			expect: map[string]interface{}{
				"体験概要":     "じゃがいも収穫",
				"サムネイルURL": "http://example.com/thumbnail.png",
				"大人人数":     "2",
				"中学生人数":    "1",
				"小学生人数":    "0",
				"幼児人数":     "0",
				"シニア人数":    "0",
			},
		},
		{
			name: "shipped",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.Shipped(order.ShippingMessage)
			},
			expect: map[string]interface{}{
				"メッセージ": "ありがとうございます",
			},
		},
		{
			name: "review items",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.ReviewItems(order.OrderItems, products, maker)
			},
			expect: map[string]interface{}{
				"商品一覧": []map[string]string{{
					"商品名":      "おいしいじゃがいも",
					"サムネイルURL": "http://example.com/image.png",
					"レビューURL":  "http://example.com/reviews/products/product-id",
				}},
			},
		},
		{
			name: "review experience",
			execute: func(builder *TemplateDataBuilder) *TemplateDataBuilder {
				return builder.ReviewExperience(experience, maker)
			},
			expect: map[string]interface{}{
				"体験名":      "じゃがいも収穫",
				"サムネイルURL": "http://example.com/thumbnail.png",
				"レビューURL":  "http://example.com/reviews/experiences/experience-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			builder := NewTemplateDataBuilder()
			data := tt.execute(builder).Build()
			assert.Equal(t, tt.expect, data)
		})
	}
}
