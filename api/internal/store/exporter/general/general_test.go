package general

import (
	"bytes"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/exporter"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceipts(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &ReceiptsParams{
		Orders: entity.Orders{
			{
				ID:            "order-id",
				UserID:        "user-id",
				PromotionID:   "",
				CoordinatorID: "coordinator-id",
				CreatedAt:     now,
				UpdatedAt:     now,
				OrderPayment: entity.OrderPayment{
					OrderID:           "order-id",
					AddressRevisionID: 1,
					TransactionID:     "transaction-id",
					Status:            entity.PaymentStatusPending,
					MethodType:        entity.PaymentMethodTypeCreditCard,
					Subtotal:          1100,
					Discount:          0,
					ShippingFee:       500,
					Tax:               145,
					Total:             1600,
					CreatedAt:         now,
					UpdatedAt:         now,
					OrderedAt:         now,
				},
				OrderFulfillments: entity.OrderFulfillments{
					{
						ID:                "fulfillment-id",
						OrderID:           "order-id",
						AddressRevisionID: 1,
						Status:            entity.FulfillmentStatusUnfulfilled,
						TrackingNumber:    "",
						ShippingCarrier:   entity.ShippingCarrierUnknown,
						ShippingType:      entity.ShippingTypeNormal,
						BoxNumber:         1,
						BoxSize:           entity.ShippingSize60,
						CreatedAt:         now,
						UpdatedAt:         now,
					},
				},
				OrderItems: entity.OrderItems{
					{
						FulfillmentID:     "fulfillment-id",
						ProductRevisionID: 1,
						OrderID:           "order-id",
						Quantity:          1,
						CreatedAt:         now,
						UpdatedAt:         now,
					},
				},
			},
		},
		Addresses: map[int64]*uentity.Address{
			1: {
				AddressRevision: uentity.AddressRevision{
					ID:             1,
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-5678",
				},
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Products: map[int64]*entity.Product{
			1: {
				ID:              "product-id",
				TypeID:          "type-id",
				TagIDs:          []string{"tag-id"},
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: entity.MultiProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				ExpirationDate:    7,
				StorageMethodType: entity.StorageMethodTypeNormal,
				DeliveryType:      entity.DeliveryTypeNormal,
				Box60Rate:         50,
				Box80Rate:         40,
				Box100Rate:        30,
				OriginPrefecture:  "滋賀県",
				OriginCity:        "彦根市",
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     400,
					Cost:      300,
					CreatedAt: now,
					UpdatedAt: now,
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
	expect := []exporter.Receipt{
		&Receipt{
			OrderID:                   "order-id",
			UserID:                    "user-id",
			CoordinatorID:             "coordinator-id",
			ExpectedDelveryDate:       "",
			ExpectedDeliveryTimeFrame: "",
			DeliveryName:              "&. 購入者",
			DeliveryNameKana:          "あんどどっと こうにゅうしゃ",
			DeliveryPhoneNumber:       "090-1234-5678",
			DeliveryPostalCode:        "1000014",
			DeliveryPrefecture:        "東京都",
			DeliveryCity:              "千代田区",
			DeliveryAddressLine1:      "永田町1-7-1",
			DeliveryAddressLine2:      "",
			ClientName:                "&. 購入者",
			ClientNameKana:            "あんどどっと こうにゅうしゃ",
			ClientPhoneNumber:         "090-1234-5678",
			ClientPostalCode:          "1000014",
			ClientPrefecture:          "東京都",
			ClientCity:                "千代田区",
			ClientAddressLine1:        "永田町1-7-1",
			ClientAddressLine2:        "",
			Product1ID:                "product-id",
			Product1Name:              "新鮮なじゃがいも",
			Product1Quantity:          1,
			Product2ID:                "",
			Product2Name:              "",
			Product2Quantity:          0,
			Product3ID:                "",
			Product3Name:              "",
			Product3Quantity:          0,
			PaymentMethod:             "クレジットカード決済",
			Subtotal:                  1100,
			Discount:                  0,
			ShippingFee:               500,
			Total:                     1600,
			OrderedAt:                 now,
			ShippingType:              "通常配送",
			ShippingSize:              "60",
		},
	}
	actual := NewReceipts(params)
	assert.Equal(t, expect, actual)
}

func TestReceipt_Write(t *testing.T) {
	t.Parallel()
	now := jst.Date(2024, 1, 23, 18, 30, 0, 0)
	receipt := &Receipt{
		OrderID:                   "order-id",
		UserID:                    "user-id",
		CoordinatorID:             "coordinator-id",
		ExpectedDelveryDate:       "",
		ExpectedDeliveryTimeFrame: "",
		DeliveryName:              "&. 購入者",
		DeliveryNameKana:          "あんどどっと こうにゅうしゃ",
		DeliveryPhoneNumber:       "090-1234-5678",
		DeliveryPostalCode:        "1000014",
		DeliveryPrefecture:        "東京都",
		DeliveryCity:              "千代田区",
		DeliveryAddressLine1:      "永田町1-7-1",
		DeliveryAddressLine2:      "",
		ClientName:                "&. 購入者",
		ClientNameKana:            "あんどどっと こうにゅうしゃ",
		ClientPhoneNumber:         "090-1234-5678",
		ClientPostalCode:          "1000014",
		ClientPrefecture:          "東京都",
		ClientCity:                "千代田区",
		ClientAddressLine1:        "永田町1-7-1",
		ClientAddressLine2:        "",
		Product1ID:                "product-id",
		Product1Name:              "新鮮なじゃがいも",
		Product1Quantity:          1,
		Product2ID:                "",
		Product2Name:              "",
		Product2Quantity:          0,
		Product3ID:                "",
		Product3Name:              "",
		Product3Quantity:          0,
		PaymentMethod:             "クレジットカード決済",
		Subtotal:                  1100,
		Discount:                  0,
		ShippingFee:               500,
		Total:                     1600,
		OrderedAt:                 now,
		ShippingType:              "常温・冷蔵便",
		ShippingSize:              "60",
	}
	tests := []struct {
		name         string
		encodingType codes.CharacterEncodingType
		receipt      *Receipt
		expect       string
	}{
		{
			name:         "success utf-8",
			encodingType: codes.CharacterEncodingTypeUTF8,
			receipt:      receipt,
			expect: "注文管理番号,ユーザーID,コーディネータID,お届け希望日,お届け希望時間帯,お届け先名,お届け先名（かな）,お届け先電話番号,お届け先郵便番号,お届け先都道府県,お届け先市区町村,お届け先町名・番地,お届け先ビル名・号室など,ご依頼主名,ご依頼主名（かな）,ご依頼主電話番号,ご依頼主郵便番号,ご依頼主都道府県,ご依頼主市区町村,ご依頼主町名・番地,ご依頼主ビル名・号室など,商品コード1,商品名1,商品1数量,商品コード2,商品名2,商品2数量,商品コード3,商品名3,商品3数量,決済手段,商品金額,割引金額,配送手数料,合計金額,注文日時,配送方法,箱のサイズ\n" +
				"order-id,user-id,coordinator-id,,,&. 購入者,あんどどっと こうにゅうしゃ,090-1234-5678,1000014,東京都,千代田区,永田町1-7-1,,&. 購入者,あんどどっと こうにゅうしゃ,090-1234-5678,1000014,東京都,千代田区,永田町1-7-1,,product-id,新鮮なじゃがいも,1,,,0,,,0,クレジットカード決済,1100,0,500,1600,2024-01-23 18:30:00,常温・冷蔵便,60\n",
		},
		{
			name:         "success shift-jis",
			encodingType: codes.CharacterEncodingTypeShiftJIS,
			receipt:      receipt,
			expect: "\x92\x8d\x95\xb6\x8aǗ\x9d\x94ԍ\x86,\x83\x86\x81[\x83U\x81[ID,\x83R\x81[\x83f\x83B\x83l\x81[\x83^ID,\x82\xa8\x93͂\xaf\x8a\xf3\x96]\x93\xfa,\x82\xa8\x93͂\xaf\x8a\xf3\x96]\x8e\x9e\x8aԑ\xd1,\x82\xa8\x93͂\xaf\x90於,\x82\xa8\x93͂\xaf\x90於\x81i\x82\xa9\x82ȁj,\x82\xa8\x93͂\xaf\x90\xe6\x93d\x98b\x94ԍ\x86,\x82\xa8\x93͂\xaf\x90\xe6\x97X\x95֔ԍ\x86,\x82\xa8\x93͂\xaf\x90\xe6\x93s\x93\xb9\x95{\x8c\xa7,\x82\xa8\x93͂\xaf\x90\xe6\x8es\x8b撬\x91\xba,\x82\xa8\x93͂\xaf\x90撬\x96\xbc\x81E\x94Ԓn,\x82\xa8\x93͂\xaf\x90\xe6\x83r\x83\x8b\x96\xbc\x81E\x8d\x86\x8e\xba\x82Ȃ\xc7,\x82\xb2\x88˗\x8a\x8e喼,\x82\xb2\x88˗\x8a\x8e喼\x81i\x82\xa9\x82ȁj,\x82\xb2\x88˗\x8a\x8e\xe5\x93d\x98b\x94ԍ\x86,\x82\xb2\x88˗\x8a\x8e\xe5\x97X\x95֔ԍ\x86,\x82\xb2\x88˗\x8a\x8e\xe5\x93s\x93\xb9\x95{\x8c\xa7,\x82\xb2\x88˗\x8a\x8e\xe5\x8es\x8b撬\x91\xba,\x82\xb2\x88˗\x8a\x8e咬\x96\xbc\x81E\x94Ԓn,\x82\xb2\x88˗\x8a\x8e\xe5\x83r\x83\x8b\x96\xbc\x81E\x8d\x86\x8e\xba\x82Ȃ\xc7,\x8f\xa4\x95i\x83R\x81[\x83h1,\x8f\xa4\x95i\x96\xbc1,\x8f\xa4\x95i1\x90\x94\x97\xca,\x8f\xa4\x95i\x83R\x81[\x83h2,\x8f\xa4\x95i\x96\xbc2,\x8f\xa4\x95i2\x90\x94\x97\xca,\x8f\xa4\x95i\x83R\x81[\x83h3,\x8f\xa4\x95i\x96\xbc3,\x8f\xa4\x95i3\x90\x94\x97\xca,\x8c\x88\x8dώ\xe8\x92i,\x8f\xa4\x95i\x8b\xe0\x8az,\x8a\x84\x88\xf8\x8b\xe0\x8az,\x94z\x91\x97\x8e萔\x97\xbf,\x8d\x87\x8cv\x8b\xe0\x8az,\x92\x8d\x95\xb6\x93\xfa\x8e\x9e,\x94z\x91\x97\x95\xfb\x96@,\x94\xa0\x82̃T\x83C\x83Y\n" +
				"order-id,user-id,coordinator-id,,,&. \x8dw\x93\xfc\x8e\xd2,\x82\xa0\x82\xf1\x82ǂǂ\xc1\x82\xc6 \x82\xb1\x82\xa4\x82ɂイ\x82\xb5\x82\xe1,090-1234-5678,1000014,\x93\x8c\x8b\x9e\x93s,\x90\xe7\x91\xe3\x93c\x8b\xe6,\x89i\x93c\x92\xac1-7-1,,&. \x8dw\x93\xfc\x8e\xd2,\x82\xa0\x82\xf1\x82ǂǂ\xc1\x82\xc6 \x82\xb1\x82\xa4\x82ɂイ\x82\xb5\x82\xe1,090-1234-5678,1000014,\x93\x8c\x8b\x9e\x93s,\x90\xe7\x91\xe3\x93c\x8b\xe6,\x89i\x93c\x92\xac1-7-1,,product-id,\x90V\x91N\x82Ȃ\xb6\x82Ⴊ\x82\xa2\x82\xe0,1,,,0,,,0,\x83N\x83\x8c\x83W\x83b\x83g\x83J\x81[\x83h\x8c\x88\x8d\xcf,1100,0,500,1600,2024-01-23 18:30:00,\x8f퉷\x81E\x97①\x95\xd6,60\n",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			writer := exporter.NewExporter(buf, tt.encodingType)
			err := writer.WriteHeader(&Receipt{})
			require.NoError(t, err)
			err = writer.WriteBody(tt.receipt)
			require.NoError(t, err)
			err = writer.Flush()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, buf.String())
		})
	}
}
