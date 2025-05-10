package sagawa

import (
	"bytes"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/exporter"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
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
			DeliveryCodeType:                  "",
			DeliveryCode:                      "",
			DeliveryPhoneNumber:               "090-1234-5678",
			DeliveryPostalCode:                "1000014",
			DeliveryAddress1:                  "東京都",
			DeliveryAddress2:                  "千代田区 永田町1-7-1",
			DeliveryAddress3:                  "",
			DeliveryLastname:                  "&.",
			DeliveryFirstname:                 "購入者",
			OrderID:                           "order-id",
			UserID:                            "user-id",
			DivisionClientCodeType:            "",
			DivisionClientCode:                "",
			DivixionClientName:                "",
			ShipperPhoneNumber:                "",
			ClientCodeType:                    "",
			ClientCode:                        "",
			ClientPhoneNumber:                 "090-1234-5678",
			ClientPostalCode:                  "1000014",
			ClientAddress1:                    "東京都 千代田区 永田町1-7-1",
			ClientAddress2:                    "",
			ClientLastname:                    "&.",
			ClientFirstname:                   "購入者",
			Packing:                           "",
			Product1Name:                      "新鮮なじゃがいも",
			Product2Name:                      "",
			Product3Name:                      "",
			Product4Name:                      "",
			Product5Name:                      "",
			Tagpackaging:                      "",
			TagProduct1Name:                   "",
			TagProduct2Name:                   "",
			TagProduct3Name:                   "",
			TagProduct4Name:                   "",
			TagProduct5Name:                   "",
			TagProduct6Name:                   "",
			TagProduct7Name:                   "",
			TagProduct8Name:                   "",
			TagProduct9Name:                   "",
			TagProduct10Name:                  "",
			TagProduct11Name:                  "",
			ShipmentQuantity:                  1,
			SpeedSpecification:                "",
			ShippingType:                      "001",
			DeliveryDate:                      "",
			DeliveryTimeFrame:                 "",
			DeliveryTime:                      "",
			DeliveryAmount:                    0,
			DeliveryTax:                       0,
			PaymentMethod:                     "",
			InsuranceAmount:                   0,
			DesignatedSticker1:                "",
			DesignatedSticker2:                "",
			DesignatedSticker3:                "",
			OfficePickUp:                      "",
			SrcType:                           "",
			OfficePickUpCode:                  "",
			OriginalDestinationClassification: "",
			Email:                             "",
			AbsenceContact:                    "",
			ShippingDate:                      "",
			ReceiptNumber:                     "",
			ShippingFacilityPrintType:         "",
			UnaggregatedSpec:                  "",
			Reserve1:                          "",
			Reserve2:                          "",
			Reserve3:                          "",
			Reserve4:                          "",
			Reserve5:                          "",
			Reserve6:                          "",
			Reserve7:                          "",
			Reserve8:                          "",
			Reserve9:                          "",
			Reserve10:                         "",
		},
	}
	actual := NewReceipts(params)
	assert.Equal(t, expect, actual)
}

func TestReceipt_Write(t *testing.T) {
	t.Parallel()
	receipt := &Receipt{
		DeliveryPhoneNumber: "090-1234-5678",
		DeliveryPostalCode:  "1000014",
		DeliveryAddress1:    "東京都",
		DeliveryAddress2:    "千代田区 永田町1-7-1",
		DeliveryAddress3:    "",
		DeliveryLastname:    "&.",
		DeliveryFirstname:   "購入者",
		OrderID:             "order-id",
		UserID:              "user-id",
		ClientPhoneNumber:   "090-1234-5678",
		ClientPostalCode:    "1000014",
		ClientAddress1:      "東京都 千代田区 永田町1-7-1",
		ClientAddress2:      "",
		ClientLastname:      "&.",
		ClientFirstname:     "購入者",
		Packing:             "",
		Product1Name:        "新鮮なじゃがいも",
		Product2Name:        "",
		Product3Name:        "",
		Product4Name:        "",
		Product5Name:        "",
		ShipmentQuantity:    1,
		ShippingType:        "001",
		DeliveryDate:        "",
		DeliveryTimeFrame:   "",
		DeliveryTime:        "",
		DeliveryAmount:      0,
		DeliveryTax:         0,
		Email:               "",
		AbsenceContact:      "",
		ShippingDate:        "",
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
			expect: "お届け先コード取得区分,お届け先コード,お届け先電話番号,お届け先郵便番号,お届け先住所１,お届け先住所２,お届け先住所３,お届け先名称１,お届け先名称２,お客様管理番号,お客様コード,部署ご担当者コード,取得区分,部署ご担当者コード,部署ご担当者名称,荷送人電話番号,ご依頼主コード取得区分,ご依頼主コード,ご依頼主電話番号,ご依頼主郵便番号,ご依頼主住所１,ご依頼主住所２,ご依頼主名称１,ご依頼主名称２,荷姿,品名１,品名２,品名３,品名４,品名５,荷札荷姿,荷札品名1,荷札品名2,荷札品名3,荷札品名4,荷札品名5,荷札品名6,荷札品名7,荷札品名8,荷札品名9,荷札品名10,荷札品名11,出荷個数,スピード指定,クール便指定,配達日,配達指定時間帯,配達指定時間（時分）,代引金額,消費税,決済種別,保険金額,指定シール1,指定シール2,指定シール3,営業所受取,ＳＲＣ区分,営業所受取営業所コード,元着区分,メールアドレス,ご不在時連絡先,出荷日,お問い合せ送り状No.,出荷場印字区分,集約解除指定,編集1,編集2,編集3,編集4,編集5,編集6,編集7,編集8,編集9,編集10\n" +
				",,090-1234-5678,1000014,東京都,千代田区 永田町1-7-1,,&.,購入者,order-id,user-id,,,,,,,090-1234-5678,1000014,東京都 千代田区 永田町1-7-1,,&.,購入者,,新鮮なじゃがいも,,,,,,,,,,,,,,,,,1,,001,,,,0,0,,0,,,,,,,,,,,,,,,,,,,,,,,\n",
		},
		{
			name:         "success shift-jis",
			encodingType: codes.CharacterEncodingTypeShiftJIS,
			receipt:      receipt,
			expect: "\x82\xa8\x93͂\xaf\x90\xe6\x83R\x81[\x83h\x8e擾\x8b敪,\x82\xa8\x93͂\xaf\x90\xe6\x83R\x81[\x83h,\x82\xa8\x93͂\xaf\x90\xe6\x93d\x98b\x94ԍ\x86,\x82\xa8\x93͂\xaf\x90\xe6\x97X\x95֔ԍ\x86,\x82\xa8\x93͂\xaf\x90\xe6\x8fZ\x8f\x8a\x82P,\x82\xa8\x93͂\xaf\x90\xe6\x8fZ\x8f\x8a\x82Q,\x82\xa8\x93͂\xaf\x90\xe6\x8fZ\x8f\x8a\x82R,\x82\xa8\x93͂\xaf\x90於\x8f̂P,\x82\xa8\x93͂\xaf\x90於\x8f̂Q,\x82\xa8\x8bq\x97l\x8aǗ\x9d\x94ԍ\x86,\x82\xa8\x8bq\x97l\x83R\x81[\x83h,\x95\x94\x8f\x90\x82\xb2\x92S\x93\x96\x8e҃R\x81[\x83h,\x8e擾\x8b敪,\x95\x94\x8f\x90\x82\xb2\x92S\x93\x96\x8e҃R\x81[\x83h,\x95\x94\x8f\x90\x82\xb2\x92S\x93\x96\x8eҖ\xbc\x8f\xcc,\x89ב\x97\x90l\x93d\x98b\x94ԍ\x86,\x82\xb2\x88˗\x8a\x8e\xe5\x83R\x81[\x83h\x8e擾\x8b敪,\x82\xb2\x88˗\x8a\x8e\xe5\x83R\x81[\x83h,\x82\xb2\x88˗\x8a\x8e\xe5\x93d\x98b\x94ԍ\x86,\x82\xb2\x88˗\x8a\x8e\xe5\x97X\x95֔ԍ\x86,\x82\xb2\x88˗\x8a\x8e\xe5\x8fZ\x8f\x8a\x82P,\x82\xb2\x88˗\x8a\x8e\xe5\x8fZ\x8f\x8a\x82Q,\x82\xb2\x88˗\x8a\x8e喼\x8f̂P,\x82\xb2\x88˗\x8a\x8e喼\x8f̂Q,\x89\u05cep,\x95i\x96\xbc\x82P,\x95i\x96\xbc\x82Q,\x95i\x96\xbc\x82R,\x95i\x96\xbc\x82S,\x95i\x96\xbc\x82T,\x89\u05ceD\x89\u05cep,\x89\u05ceD\x95i\x96\xbc1,\x89\u05ceD\x95i\x96\xbc2,\x89\u05ceD\x95i\x96\xbc3,\x89\u05ceD\x95i\x96\xbc4,\x89\u05ceD\x95i\x96\xbc5,\x89\u05ceD\x95i\x96\xbc6,\x89\u05ceD\x95i\x96\xbc7,\x89\u05ceD\x95i\x96\xbc8,\x89\u05ceD\x95i\x96\xbc9,\x89\u05ceD\x95i\x96\xbc10,\x89\u05ceD\x95i\x96\xbc11,\x8fo\x89\u05cc\u0090\x94,\x83X\x83s\x81[\x83h\x8ew\x92\xe8,\x83N\x81[\x83\x8b\x95֎w\x92\xe8,\x94z\x92B\x93\xfa,\x94z\x92B\x8ew\x92莞\x8aԑ\xd1,\x94z\x92B\x8ew\x92莞\x8aԁi\x8e\x9e\x95\xaa\x81j,\x91\xe3\x88\xf8\x8b\xe0\x8az,\x8f\xc1\x94\xef\x90\xc5,\x8c\x88\x8dώ\xed\x95\xca,\x95ی\xaf\x8b\xe0\x8az,\x8ew\x92\xe8\x83V\x81[\x83\x8b1,\x8ew\x92\xe8\x83V\x81[\x83\x8b2,\x8ew\x92\xe8\x83V\x81[\x83\x8b3,\x89c\x8bƏ\x8a\x8e\xf3\x8e\xe6,\x82r\x82q\x82b\x8b敪,\x89c\x8bƏ\x8a\x8e\xf3\x8e\xe6\x89c\x8bƏ\x8a\x83R\x81[\x83h,\x8c\xb3\x92\x85\x8b敪,\x83\x81\x81[\x83\x8b\x83A\x83h\x83\x8c\x83X,\x82\xb2\x95s\x8dݎ\x9e\x98A\x97\x8d\x90\xe6,\x8fo\x89ד\xfa,\x82\xa8\x96₢\x8d\x87\x82\xb9\x91\x97\x82\xe8\x8f\xf3No.,\x8fo\x89\u05cf\xea\x88\U000ce68b敪,\x8fW\x96\xf1\x89\xf0\x8f\x9c\x8ew\x92\xe8,\x95ҏW1,\x95ҏW2,\x95ҏW3,\x95ҏW4,\x95ҏW5,\x95ҏW6,\x95ҏW7,\x95ҏW8,\x95ҏW9,\x95ҏW10\n" +
				",,090-1234-5678,1000014,\x93\x8c\x8b\x9e\x93s,\x90\xe7\x91\xe3\x93c\x8b\xe6 \x89i\x93c\x92\xac1-7-1,,&.,\x8dw\x93\xfc\x8e\xd2,order-id,user-id,,,,,,,090-1234-5678,1000014,\x93\x8c\x8b\x9e\x93s \x90\xe7\x91\xe3\x93c\x8b\xe6 \x89i\x93c\x92\xac1-7-1,,&.,\x8dw\x93\xfc\x8e\xd2,,\x90V\x91N\x82Ȃ\xb6\x82Ⴊ\x82\xa2\x82\xe0,,,,,,,,,,,,,,,,,1,,001,,,,0,0,,0,,,,,,,,,,,,,,,,,,,,,,,\n",
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
