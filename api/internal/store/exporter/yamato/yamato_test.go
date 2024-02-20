package yamato

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
			OrderID:                                "order-id",
			ServiceType:                            ServiceTypePrepayment,
			ShippingType:                           ShippingTypeNormal,
			YamatoOrderID:                          "",
			ExpectedShippingDate:                   "",
			ExpectedDeliveryDate:                   "",
			ExpectedDeliveryTimeFrame:              DeliveryTimeFrameNone,
			DeliveryCode:                           "1",
			DeliveryPhoneNumber:                    "090-1234-5678",
			DeliveryPhoneNumberExtension:           "",
			DeliveryPostalCode:                     "1000014",
			DeliveryAddress1:                       "東京都 千代田区 永田町1-7-1",
			DeliveryAddress2:                       "",
			DeliveryCompany1:                       "",
			DeliveryCompany2:                       "",
			DeliveryName:                           "&. 購入者",
			DeliveryNameKana:                       "あんどどっと こうにゅうしゃ",
			DeliveryNameHonorific:                  "様",
			ClientCode:                             "1",
			ClientPhoneNumber:                      "090-1234-5678",
			ClientPhoneNumberExtension:             "",
			ClientPostalCode:                       "1000014",
			ClientAddress1:                         "東京都 千代田区 永田町1-7-1",
			ClientAddress2:                         "",
			ClientName:                             "&. 購入者",
			ClientNameKana:                         "あんどどっと こうにゅうしゃ",
			Product1ID:                             "product-id",
			Product1Name:                           "新鮮なじゃがいも",
			Product2ID:                             "",
			Product2Name:                           "",
			Handling1:                              "",
			Handling2:                              "",
			Note:                                   "",
			WebCollectDeliveryAmount:               0,
			WebCollectDeliveryTax:                  0,
			BranchHoldType:                         1,
			BranchCode:                             "",
			IssuedQuantity:                         1,
			PrintingQuantity:                       "",
			BillingCode:                            "",
			BillingGroupCode:                       "",
			ShippingCostControlNumber:              "",
			WebCollectRegistration:                 0,
			WebCollectStoreNumber:                  "",
			WebCollectAcceptanceNumber1:            "",
			WebCollectAcceptanceNumber2:            "",
			WebCollectAcceptanceNumber3:            "",
			DeliveryNotificationEmailUsage:         0,
			DeliveryNotificationEmail:              "",
			InputMethod:                            "",
			DeliveryNotificationEmailMessage:       "",
			DeliveryCompletionEmailUsage:           0,
			DeliveryCompletionEmail:                "",
			DeliveryCompletionEmailMessage:         "",
			CollectionAgencyUsage:                  0,
			Reserve1:                               "",
			CollectionAgencyBillingAmount:          0,
			CollectionAgencyBillingTax:             0,
			CollectionAgencyBillingPostalCode:      "",
			CollectionAgencyBillingAddress1:        "",
			CollectionAgencyBillingAddress2:        "",
			CollectionAgencyBillingCompany1:        "",
			CollectionAgencyBillingCompany2:        "",
			CollectionAgencyBillingName:            "",
			CollectionAgencyBillingNameKana:        "",
			CollectionAgencyInquiryPostalCode:      "",
			CollectionAgencyInquiryAddress1:        "",
			CollectionAgencyInquiryAddress2:        "",
			CollectionAgencyInquiryPhoneNumber:     "",
			CollectionAgencyManagementNumber:       "",
			CollectionAgencyItemName:               "",
			CollectionAgencyRemarks:                "",
			MultiplePackagesGroupingKey:            "",
			SearchKey1Title:                        "ふるマルユーザーID",
			SearchKey1Value:                        "user-id",
			SearchKey2Title:                        "ふるマル注文履歴ID",
			SearchKey2Value:                        "order-id",
			SearchKey3Title:                        "",
			SearchKey3Value:                        "",
			SearchKey4Title:                        "",
			SearchKey4Value:                        "",
			SearchKey5Title:                        "",
			SearchKey5Value:                        "",
			Reserve2:                               "",
			Reserve3:                               "",
			MailDropNotificationEmailUsage:         0,
			MailDropNotificationEmail:              "",
			MailDropNotificationEmailMessage:       "",
			MailDropCompletionDeliveryEmailUsage:   0,
			MailDropCompletionDeliveryEmail:        "",
			MailDropCompletionDeliveryEmailMessage: "",
			MailDropCompletionClientEmailusage:     "",
			MailDropCompletionClientEmail:          "",
			MailDropCompletionClientEmailMessage:   "",
		},
	}
	actual := NewReceipts(params)
	assert.Equal(t, expect, actual)
}

func TestReceipt_Write(t *testing.T) {
	t.Parallel()
	receipt := &Receipt{
		OrderID:          "order-id",
		ServiceType:      ServiceTypePrepayment,
		ShippingType:     ShippingTypeNormal,
		DeliveryName:     "&. 利用者",
		DeliveryNameKana: "あんどどっと りようしゃ",
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
			expect: "お客様管理番号,送り状種類,クール区分,伝票番号,出荷予定日,お届け予定日,配送時間帯,お届け先コード,お届け先電話番号,お届け先電話番号枝番,お届け先郵便番号,お届け先住所,お届け先アパートマンション名,お届け先会社・部門１,お届け先会社・部門２,お届け先名,お届け先名(ｶﾅ),敬称,ご依頼主コード,ご依頼主電話番号,ご依頼主電話番号枝番,ご依頼主郵便番号,ご依頼主住所,ご依頼主アパートマンション名,ご依頼主名,ご依頼主名(ｶﾅ),品名コード１,品名１,品名コード２,品名２,荷扱い１,荷扱い２,記事,ｺﾚｸﾄ代金引換額（税込),内消費税額等,止置き,営業所コード,発行枚数,個数口枠の印字,請求先顧客コード,請求先分類コード,運賃管理番号,クロネコwebコレクトデータ登録,クロネコwebコレクト加盟店番号,クロネコwebコレクト申込受付番号１,クロネコwebコレクト申込受付番号２,クロネコwebコレクト申込受付番号３,お届け予定ｅメール利用区分,お届け予定ｅメールe-mailアドレス,入力機種,お届け予定ｅメールメッセージ,お届け完了ｅメール利用区分,お届け完了ｅメールe-mailアドレス,お届け完了ｅメールメッセージ,クロネコ収納代行利用区分,予備,収納代行請求金額(税込),収納代行内消費税額等,収納代行請求先郵便番号,収納代行請求先住所,収納代行請求先住所（アパートマンション名）,収納代行請求先会社・部門名１,収納代行請求先会社・部門名２,収納代行請求先名(漢字),収納代行請求先名(カナ),収納代行問合せ先郵便番号,収納代行問合せ先住所,収納代行問合せ先住所（アパートマンション名）,収納代行問合せ先電話番号,収納代行管理番号,収納代行品名,収納代行備考,複数口くくりキー,検索キータイトル1,検索キー1,検索キータイトル2,検索キー2,検索キータイトル3,検索キー3,検索キータイトル4,検索キー4,検索キータイトル5,検索キー5,予備,予備,投函予定メール利用区分,投函予定メールe-mailアドレス,投函予定メールメッセージ,投函完了メール（お届け先宛）利用区分,投函予定メール（お届け先宛）e-mailアドレス,投函予定メール（お届け先宛）メッセージ,投函完了メール（ご依頼主宛）利用区分,投函予定メール（ご依頼主宛）e-mailアドレス,投函予定メール（ご依頼主宛）メッセージ\n" +
				"order-id,0,0,,,,,,,,,,,,,&. 利用者,あんどどっと りようしゃ,,,,,,,,,,,,,,,,,0,0,0,,0,,,,,0,,,,,0,,,,0,,,0,,0,0,,,,,,,,,,,,,,,,,,,,,,,,,,,,0,,,0,,,,,\n",
		},
		{
			name:         "success shift-jis",
			encodingType: codes.CharacterEncodingTypeShiftJIS,
			receipt:      receipt,
			expect: "\x82\xa8\x8bq\x97l\x8aǗ\x9d\x94ԍ\x86,\x91\x97\x82\xe8\x8f\xf3\x8e\xed\x97\xde,\x83N\x81[\x83\x8b\x8b敪,\x93`\x95[\x94ԍ\x86,\x8fo\x89ח\\\x92\xe8\x93\xfa,\x82\xa8\x93͂\xaf\x97\\\x92\xe8\x93\xfa,\x94z\x91\x97\x8e\x9e\x8aԑ\xd1,\x82\xa8\x93͂\xaf\x90\xe6\x83R\x81[\x83h,\x82\xa8\x93͂\xaf\x90\xe6\x93d\x98b\x94ԍ\x86,\x82\xa8\x93͂\xaf\x90\xe6\x93d\x98b\x94ԍ\x86\x8e}\x94\xd4,\x82\xa8\x93͂\xaf\x90\xe6\x97X\x95֔ԍ\x86,\x82\xa8\x93͂\xaf\x90\xe6\x8fZ\x8f\x8a,\x82\xa8\x93͂\xaf\x90\xe6\x83A\x83p\x81[\x83g\x83}\x83\x93\x83V\x83\x87\x83\x93\x96\xbc,\x82\xa8\x93͂\xaf\x90\xe6\x89\xef\x8eЁE\x95\x94\x96\xe5\x82P,\x82\xa8\x93͂\xaf\x90\xe6\x89\xef\x8eЁE\x95\x94\x96\xe5\x82Q,\x82\xa8\x93͂\xaf\x90於,\x82\xa8\x93͂\xaf\x90於(\xb6\xc5),\x8ch\x8f\xcc,\x82\xb2\x88˗\x8a\x8e\xe5\x83R\x81[\x83h,\x82\xb2\x88˗\x8a\x8e\xe5\x93d\x98b\x94ԍ\x86,\x82\xb2\x88˗\x8a\x8e\xe5\x93d\x98b\x94ԍ\x86\x8e}\x94\xd4,\x82\xb2\x88˗\x8a\x8e\xe5\x97X\x95֔ԍ\x86,\x82\xb2\x88˗\x8a\x8e\xe5\x8fZ\x8f\x8a,\x82\xb2\x88˗\x8a\x8e\xe5\x83A\x83p\x81[\x83g\x83}\x83\x93\x83V\x83\x87\x83\x93\x96\xbc,\x82\xb2\x88˗\x8a\x8e喼,\x82\xb2\x88˗\x8a\x8e喼(\xb6\xc5),\x95i\x96\xbc\x83R\x81[\x83h\x82P,\x95i\x96\xbc\x82P,\x95i\x96\xbc\x83R\x81[\x83h\x82Q,\x95i\x96\xbc\x82Q,\x89\u05c8\xb5\x82\xa2\x82P,\x89\u05c8\xb5\x82\xa2\x82Q,\x8bL\x8e\x96,\xbaڸđ\xe3\x8b\xe0\x88\xf8\x8a\xb7\x8az\x81i\x90ō\x9e),\x93\xe0\x8f\xc1\x94\xef\x90Ŋz\x93\x99,\x8e~\x92u\x82\xab,\x89c\x8bƏ\x8a\x83R\x81[\x83h,\x94\xad\x8ds\x96\x87\x90\x94,\x8c\u0090\x94\x8c\xfb\x98g\x82̈\xf3\x8e\x9a,\x90\xbf\x8b\x81\x90\xe6\x8cڋq\x83R\x81[\x83h,\x90\xbf\x8b\x81\x90敪\x97ރR\x81[\x83h,\x89^\x92\xc0\x8aǗ\x9d\x94ԍ\x86,\x83N\x83\x8d\x83l\x83Rweb\x83R\x83\x8c\x83N\x83g\x83f\x81[\x83^\x93o\x98^,\x83N\x83\x8d\x83l\x83Rweb\x83R\x83\x8c\x83N\x83g\x89\xc1\x96\xbf\x93X\x94ԍ\x86,\x83N\x83\x8d\x83l\x83Rweb\x83R\x83\x8c\x83N\x83g\x90\\\x8d\x9e\x8e\xf3\x95t\x94ԍ\x86\x82P,\x83N\x83\x8d\x83l\x83Rweb\x83R\x83\x8c\x83N\x83g\x90\\\x8d\x9e\x8e\xf3\x95t\x94ԍ\x86\x82Q,\x83N\x83\x8d\x83l\x83Rweb\x83R\x83\x8c\x83N\x83g\x90\\\x8d\x9e\x8e\xf3\x95t\x94ԍ\x86\x82R,\x82\xa8\x93͂\xaf\x97\\\x92肅\x83\x81\x81[\x83\x8b\x97\x98\x97p\x8b敪,\x82\xa8\x93͂\xaf\x97\\\x92肅\x83\x81\x81[\x83\x8be-mail\x83A\x83h\x83\x8c\x83X,\x93\xfc\x97͋@\x8e\xed,\x82\xa8\x93͂\xaf\x97\\\x92肅\x83\x81\x81[\x83\x8b\x83\x81\x83b\x83Z\x81[\x83W,\x82\xa8\x93͂\xaf\x8a\xae\x97\xb9\x82\x85\x83\x81\x81[\x83\x8b\x97\x98\x97p\x8b敪,\x82\xa8\x93͂\xaf\x8a\xae\x97\xb9\x82\x85\x83\x81\x81[\x83\x8be-mail\x83A\x83h\x83\x8c\x83X,\x82\xa8\x93͂\xaf\x8a\xae\x97\xb9\x82\x85\x83\x81\x81[\x83\x8b\x83\x81\x83b\x83Z\x81[\x83W,\x83N\x83\x8d\x83l\x83R\x8e\xfb\x94[\x91\xe3\x8ds\x97\x98\x97p\x8b敪,\x97\\\x94\xf5,\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x8b\xe0\x8az(\x90ō\x9e),\x8e\xfb\x94[\x91\xe3\x8ds\x93\xe0\x8f\xc1\x94\xef\x90Ŋz\x93\x99,\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x90\xe6\x97X\x95֔ԍ\x86,\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x90\xe6\x8fZ\x8f\x8a,\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x90\xe6\x8fZ\x8f\x8a\x81i\x83A\x83p\x81[\x83g\x83}\x83\x93\x83V\x83\x87\x83\x93\x96\xbc\x81j,\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x90\xe6\x89\xef\x8eЁE\x95\x94\x96喼\x82P,\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x90\xe6\x89\xef\x8eЁE\x95\x94\x96喼\x82Q,\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x90於(\x8a\xbf\x8e\x9a),\x8e\xfb\x94[\x91\xe3\x8ds\x90\xbf\x8b\x81\x90於(\x83J\x83i),\x8e\xfb\x94[\x91\xe3\x8ds\x96⍇\x82\xb9\x90\xe6\x97X\x95֔ԍ\x86,\x8e\xfb\x94[\x91\xe3\x8ds\x96⍇\x82\xb9\x90\xe6\x8fZ\x8f\x8a,\x8e\xfb\x94[\x91\xe3\x8ds\x96⍇\x82\xb9\x90\xe6\x8fZ\x8f\x8a\x81i\x83A\x83p\x81[\x83g\x83}\x83\x93\x83V\x83\x87\x83\x93\x96\xbc\x81j,\x8e\xfb\x94[\x91\xe3\x8ds\x96⍇\x82\xb9\x90\xe6\x93d\x98b\x94ԍ\x86,\x8e\xfb\x94[\x91\xe3\x8ds\x8aǗ\x9d\x94ԍ\x86,\x8e\xfb\x94[\x91\xe3\x8ds\x95i\x96\xbc,\x8e\xfb\x94[\x91\xe3\x8ds\x94\xf5\x8dl,\x95\xa1\x90\x94\x8c\xfb\x82\xad\x82\xad\x82\xe8\x83L\x81[,\x8c\x9f\x8d\xf5\x83L\x81[\x83^\x83C\x83g\x83\x8b1,\x8c\x9f\x8d\xf5\x83L\x81[1,\x8c\x9f\x8d\xf5\x83L\x81[\x83^\x83C\x83g\x83\x8b2,\x8c\x9f\x8d\xf5\x83L\x81[2,\x8c\x9f\x8d\xf5\x83L\x81[\x83^\x83C\x83g\x83\x8b3,\x8c\x9f\x8d\xf5\x83L\x81[3,\x8c\x9f\x8d\xf5\x83L\x81[\x83^\x83C\x83g\x83\x8b4,\x8c\x9f\x8d\xf5\x83L\x81[4,\x8c\x9f\x8d\xf5\x83L\x81[\x83^\x83C\x83g\x83\x8b5,\x8c\x9f\x8d\xf5\x83L\x81[5,\x97\\\x94\xf5,\x97\\\x94\xf5,\x93\x8a\x94\x9f\x97\\\x92胁\x81[\x83\x8b\x97\x98\x97p\x8b敪,\x93\x8a\x94\x9f\x97\\\x92胁\x81[\x83\x8be-mail\x83A\x83h\x83\x8c\x83X,\x93\x8a\x94\x9f\x97\\\x92胁\x81[\x83\x8b\x83\x81\x83b\x83Z\x81[\x83W,\x93\x8a\x94\x9f\x8a\xae\x97\xb9\x83\x81\x81[\x83\x8b\x81i\x82\xa8\x93͂\xaf\x90戶\x81j\x97\x98\x97p\x8b敪,\x93\x8a\x94\x9f\x97\\\x92胁\x81[\x83\x8b\x81i\x82\xa8\x93͂\xaf\x90戶\x81je-mail\x83A\x83h\x83\x8c\x83X,\x93\x8a\x94\x9f\x97\\\x92胁\x81[\x83\x8b\x81i\x82\xa8\x93͂\xaf\x90戶\x81j\x83\x81\x83b\x83Z\x81[\x83W,\x93\x8a\x94\x9f\x8a\xae\x97\xb9\x83\x81\x81[\x83\x8b\x81i\x82\xb2\x88˗\x8a\x8e制\x81j\x97\x98\x97p\x8b敪,\x93\x8a\x94\x9f\x97\\\x92胁\x81[\x83\x8b\x81i\x82\xb2\x88˗\x8a\x8e制\x81je-mail\x83A\x83h\x83\x8c\x83X,\x93\x8a\x94\x9f\x97\\\x92胁\x81[\x83\x8b\x81i\x82\xb2\x88˗\x8a\x8e制\x81j\x83\x81\x83b\x83Z\x81[\x83W\n" +
				"order-id,0,0,,,,,,,,,,,,,&. \x97\x98\x97p\x8e\xd2,\x82\xa0\x82\xf1\x82ǂǂ\xc1\x82\xc6 \x82\xe8\x82悤\x82\xb5\x82\xe1,,,,,,,,,,,,,,,,,0,0,0,,0,,,,,0,,,,,0,,,,0,,,0,,0,0,,,,,,,,,,,,,,,,,,,,,,,,,,,,0,,,0,,,,,\n",
		},
	}
	for _, tt := range tests {
		tt := tt
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
