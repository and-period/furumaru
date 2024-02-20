package yamato

import (
	"strconv"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/exporter"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/text/width"
)

var receiptHeaders = []string{
	"お客様管理番号",
	"送り状種類",
	"クール区分",
	"伝票番号",
	"出荷予定日",
	"お届け予定日",
	"配送時間帯",
	"お届け先コード",
	"お届け先電話番号",
	"お届け先電話番号枝番",
	"お届け先郵便番号",
	"お届け先住所",
	"お届け先アパートマンション名",
	"お届け先会社・部門１",
	"お届け先会社・部門２",
	"お届け先名",
	"お届け先名(ｶﾅ)",
	"敬称",
	"ご依頼主コード",
	"ご依頼主電話番号",
	"ご依頼主電話番号枝番",
	"ご依頼主郵便番号",
	"ご依頼主住所",
	"ご依頼主アパートマンション名",
	"ご依頼主名",
	"ご依頼主名(ｶﾅ)",
	"品名コード１",
	"品名１",
	"品名コード２",
	"品名２",
	"荷扱い１",
	"荷扱い２",
	"記事",
	"ｺﾚｸﾄ代金引換額（税込)",
	"内消費税額等",
	"止置き",
	"営業所コード",
	"発行枚数",
	"個数口枠の印字",
	"請求先顧客コード",
	"請求先分類コード",
	"運賃管理番号",
	"クロネコwebコレクトデータ登録",
	"クロネコwebコレクト加盟店番号",
	"クロネコwebコレクト申込受付番号１",
	"クロネコwebコレクト申込受付番号２",
	"クロネコwebコレクト申込受付番号３",
	"お届け予定ｅメール利用区分",
	"お届け予定ｅメールe-mailアドレス",
	"入力機種",
	"お届け予定ｅメールメッセージ",
	"お届け完了ｅメール利用区分",
	"お届け完了ｅメールe-mailアドレス",
	"お届け完了ｅメールメッセージ",
	"クロネコ収納代行利用区分",
	"予備",
	"収納代行請求金額(税込)",
	"収納代行内消費税額等",
	"収納代行請求先郵便番号",
	"収納代行請求先住所",
	"収納代行請求先住所（アパートマンション名）",
	"収納代行請求先会社・部門名１",
	"収納代行請求先会社・部門名２",
	"収納代行請求先名(漢字)",
	"収納代行請求先名(カナ)",
	"収納代行問合せ先郵便番号",
	"収納代行問合せ先住所",
	"収納代行問合せ先住所（アパートマンション名）",
	"収納代行問合せ先電話番号",
	"収納代行管理番号",
	"収納代行品名",
	"収納代行備考",
	"複数口くくりキー",
	"検索キータイトル1",
	"検索キー1",
	"検索キータイトル2",
	"検索キー2",
	"検索キータイトル3",
	"検索キー3",
	"検索キータイトル4",
	"検索キー4",
	"検索キータイトル5",
	"検索キー5",
	"予備",
	"予備",
	"投函予定メール利用区分",
	"投函予定メールe-mailアドレス",
	"投函予定メールメッセージ",
	"投函完了メール（お届け先宛）利用区分",
	"投函予定メール（お届け先宛）e-mailアドレス",
	"投函予定メール（お届け先宛）メッセージ",
	"投函完了メール（ご依頼主宛）利用区分",
	"投函予定メール（ご依頼主宛）e-mailアドレス",
	"投函予定メール（ご依頼主宛）メッセージ",
}

// ServiceType - ヤマト運輸送り状種別
type ServiceType int32

const (
	ServiceTypePrepayment              ServiceType = 0 // 発払い
	ServiceTypeCollect                 ServiceType = 2 // コレクト
	ServiceTypeKuronekoYuMail          ServiceType = 3 // クロネコゆうメール
	ServiceTypeTime                    ServiceType = 4 // タイム
	ServiceTypeCashOnDelivery          ServiceType = 5 // 着払い
	ServiceTypeMultiplePrepayment      ServiceType = 6 // 発払い（複数口）
	ServiceTypeNekoPos                 ServiceType = 7 // ネコポス・クロネコゆうパケット
	ServiceTypeTakkyubinCompact        ServiceType = 8 // 宅急便コンパクト
	ServiceTypeTakkyubinCompactCollect ServiceType = 9 // 宅急便コンパクトコレクト
)

// ShippingType - ヤマト運輸クール区分
type ShippingType int32

const (
	ShippingTypeNormal        ShippingType = 0 // 通常
	ShippingTypeFreeze        ShippingType = 1 // クール冷凍
	ShippingTypeRefrigeration ShippingType = 2 // クール冷蔵
)

// DeliveryTimeFrame - ヤマト運輸配送時間帯
type DeliveryTimeFrame string

const (
	DeliveryTimeFrameNone    DeliveryTimeFrame = ""     // 指定なし
	DeliveryTimeFrameMorning DeliveryTimeFrame = "0812" // 午前中
	DeliveryTimeFrame1416    DeliveryTimeFrame = "1416" // 14〜16時
	DeliveryTimeFrame1618    DeliveryTimeFrame = "1618" // 16〜18時
	DeliveryTimeFrame1820    DeliveryTimeFrame = "1820" // 18〜20時
	DeliveryTimeFrame1921    DeliveryTimeFrame = "1921" // 19〜21時
	DeliveryTimeFrameUntil10 DeliveryTimeFrame = "0010" // 午前10時まで
	DeliveryTimeFrameUntil17 DeliveryTimeFrame = "0017" // 午後5時まで
)

// Receipt - ヤマト運輸送り状情報
// See: https://bmypage.kuroneko.co.jp/bmypage/pdf/new_exchange1.pdf
type Receipt struct {
	OrderID                                string            // お客様管理番号
	ServiceType                            ServiceType       // 送り状種類
	ShippingType                           ShippingType      // クール区分
	YamatoOrderID                          string            // 伝票番号
	ExpectedShippingDate                   string            // 出荷予定日
	ExpectedDeliveryDate                   string            // お届け予定日
	ExpectedDeliveryTimeFrame              DeliveryTimeFrame // 配送時間帯
	DeliveryCode                           string            // お届け先コード
	DeliveryPhoneNumber                    string            // お届け先電話番号
	DeliveryPhoneNumberExtension           string            // お届け先電話番号枝番
	DeliveryPostalCode                     string            // お届け先郵便番号
	DeliveryAddress1                       string            // お届け先住所
	DeliveryAddress2                       string            // お届け先アパートマンション名
	DeliveryCompany1                       string            // お届け先会社・部門１
	DeliveryCompany2                       string            // お届け先会社・部門２
	DeliveryName                           string            // お届け先名
	DeliveryNameKana                       string            // お届け先名(ｶﾅ)
	DeliveryNameHonorific                  string            // 敬称
	ClientCode                             string            // ご依頼主コード
	ClientPhoneNumber                      string            // ご依頼主電話番号
	ClientPhoneNumberExtension             string            // ご依頼主電話番号枝番
	ClientPostalCode                       string            // ご依頼主郵便番号
	ClientAddress1                         string            // ご依頼主住所
	ClientAddress2                         string            // ご依頼主アパートマンション名
	ClientName                             string            // ご依頼主名
	ClientNameKana                         string            // ご依頼主名(ｶﾅ)
	Product1ID                             string            // 品名コード１
	Product1Name                           string            // 品名１
	Product2ID                             string            // 品名コード２
	Product2Name                           string            // 品名２
	Handling1                              string            // 荷扱い１
	Handling2                              string            // 荷扱い２
	Note                                   string            // 記事
	WebCollectDeliveryAmount               int64             // ｺﾚｸﾄ代金引換額（税込)
	WebCollectDeliveryTax                  int64             // 内消費税額等
	BranchHoldType                         int64             // 止置き
	BranchCode                             string            // 営業所コード
	IssuedQuantity                         int64             // 発行枚数
	PrintingQuantity                       string            // 個数口枠の印字
	BillingCode                            string            // 請求先顧客コード
	BillingGroupCode                       string            // 請求先分類コード
	ShippingCostControlNumber              string            // 運賃管理番号
	WebCollectRegistration                 int64             // クロネコwebコレクトデータ登録
	WebCollectStoreNumber                  string            // クロネコwebコレクト加盟店番号
	WebCollectAcceptanceNumber1            string            // クロネコwebコレクト申込受付番号１
	WebCollectAcceptanceNumber2            string            // クロネコwebコレクト申込受付番号２
	WebCollectAcceptanceNumber3            string            // クロネコwebコレクト申込受付番号３
	DeliveryNotificationEmailUsage         int64             // お届け予定ｅメール利用区分
	DeliveryNotificationEmail              string            // お届け予定ｅメールe-mailアドレス
	InputMethod                            string            // 入力機種
	DeliveryNotificationEmailMessage       string            // お届け予定ｅメールメッセージ
	DeliveryCompletionEmailUsage           int64             // お届け完了ｅメール利用区分
	DeliveryCompletionEmail                string            // お届け完了ｅメールe-mailアドレス
	DeliveryCompletionEmailMessage         string            // お届け完了ｅメールメッセージ
	CollectionAgencyUsage                  int64             // クロネコ収納代行利用区分
	Reserve1                               string            // 予備
	CollectionAgencyBillingAmount          int64             // 収納代行請求金額(税込)
	CollectionAgencyBillingTax             int64             // 収納代行内消費税額等
	CollectionAgencyBillingPostalCode      string            // 収納代行請求先郵便番号
	CollectionAgencyBillingAddress1        string            // 収納代行請求先住所
	CollectionAgencyBillingAddress2        string            // 収納代行請求先住所（アパートマンション名）
	CollectionAgencyBillingCompany1        string            // 収納代行請求先会社・部門名１
	CollectionAgencyBillingCompany2        string            // 収納代行請求先会社・部門名２
	CollectionAgencyBillingName            string            // 収納代行請求先名(漢字)
	CollectionAgencyBillingNameKana        string            // 収納代行請求先名(カナ)
	CollectionAgencyInquiryPostalCode      string            // 収納代行問合せ先郵便番号
	CollectionAgencyInquiryAddress1        string            // 収納代行問合せ先住所
	CollectionAgencyInquiryAddress2        string            // 収納代行問合せ先住所（アパートマンション名）
	CollectionAgencyInquiryPhoneNumber     string            // 収納代行問合せ先電話番号
	CollectionAgencyManagementNumber       string            // 収納代行管理番号
	CollectionAgencyItemName               string            // 収納代行品名
	CollectionAgencyRemarks                string            // 収納代行備考
	MultiplePackagesGroupingKey            string            // 複数口くくりキー
	SearchKey1Title                        string            // 検索キータイトル1
	SearchKey1Value                        string            // 検索キー1
	SearchKey2Title                        string            // 検索キータイトル2
	SearchKey2Value                        string            // 検索キー2
	SearchKey3Title                        string            // 検索キータイトル3
	SearchKey3Value                        string            // 検索キー3
	SearchKey4Title                        string            // 検索キータイトル4
	SearchKey4Value                        string            // 検索キー4
	SearchKey5Title                        string            // 検索キータイトル5
	SearchKey5Value                        string            // 検索キー5
	Reserve2                               string            // 予備
	Reserve3                               string            // 予備
	MailDropNotificationEmailUsage         int64             // 投函予定メール利用区分
	MailDropNotificationEmail              string            // 投函予定メールe-mailアドレス
	MailDropNotificationEmailMessage       string            // 投函予定メールメッセージ
	MailDropCompletionDeliveryEmailUsage   int64             // 投函完了メール（お届け先宛）利用区分
	MailDropCompletionDeliveryEmail        string            // 投函予定メール（お届け先宛）e-mailアドレス
	MailDropCompletionDeliveryEmailMessage string            // 投函予定メール（お届け先宛）メッセージ
	MailDropCompletionClientEmailusage     string            // 投函完了メール（ご依頼主宛）利用区分
	MailDropCompletionClientEmail          string            // 投函予定メール（ご依頼主宛）e-mailアドレス
	MailDropCompletionClientEmailMessage   string            // 投函予定メール（ご依頼主宛）メッセージ
}

type ReceiptParams struct {
	Order       *entity.Order
	Fulfillment *entity.OrderFulfillment
	Items       entity.OrderItems
	Addresses   map[int64]*uentity.Address
	Products    map[int64]*entity.Product
}

type ReceiptsParams struct {
	Orders    entity.Orders
	Addresses map[int64]*uentity.Address
	Products  map[int64]*entity.Product
}

func NewShippingType(typ entity.ShippingType) ShippingType {
	switch typ {
	case entity.ShippingTypeNormal:
		return ShippingTypeNormal
	case entity.ShippingTypeFrozen:
		return ShippingTypeFreeze
	default:
		return ShippingTypeNormal
	}
}

func NewReceipt(params *ReceiptParams) exporter.Receipt {
	receipt := &Receipt{}
	receipt.SetReceiptDetails(params.Fulfillment)
	receipt.SetDeliveryDetails(params.Addresses[params.Fulfillment.AddressRevisionID])
	receipt.SetClientDetails(params.Addresses[params.Order.OrderPayment.AddressRevisionID])
	receipt.SetProductDetails(params.Items, params.Products)
	receipt.SetWebCollectDetails()
	receipt.SetCollectionAgencyDetails()
	receipt.SetDeliveryNotificationDetails()
	receipt.SetMailDropNotificationDetails()
	receipt.SetSearchDetails(params.Order)
	return receipt
}

func (r *Receipt) SetReceiptDetails(fulfillment *entity.OrderFulfillment) {
	r.OrderID = fulfillment.OrderID
	r.ServiceType = ServiceTypePrepayment // 0：発払い（支払い時に送料も含めているため）
	r.ShippingType = NewShippingType(fulfillment.ShippingType)
	r.YamatoOrderID = ""                                // データ入力用は空白を指定
	r.ExpectedShippingDate = ""                         // TODO: 購入フローの改修時に対応
	r.ExpectedDeliveryDate = ""                         // TODO: 購入フローの改修時に対応
	r.ExpectedDeliveryTimeFrame = DeliveryTimeFrameNone // TODO: 購入フローの改修時に対応
	r.Handling1 = ""
	r.Handling2 = ""
	r.Note = ""
	r.BranchHoldType = 1 // 1：利用する
	r.BranchCode = ""
	r.IssuedQuantity = 1
	r.PrintingQuantity = ""
	r.BillingCode = ""
	r.BillingGroupCode = ""
	r.ShippingCostControlNumber = ""
	r.MultiplePackagesGroupingKey = ""
	r.InputMethod = ""
	r.Reserve1 = ""
	r.Reserve2 = ""
	r.Reserve3 = ""
}

func (r *Receipt) SetDeliveryDetails(address *uentity.Address) {
	if address == nil {
		return
	}
	r.DeliveryCode = strconv.FormatInt(address.AddressRevision.ID, 10)
	r.DeliveryPhoneNumber = address.PhoneNumber
	r.DeliveryPhoneNumberExtension = "" // 未使用のため空文字固定
	r.DeliveryPostalCode = address.PostalCode
	r.DeliveryAddress1 = address.ShortPath()
	r.DeliveryAddress2 = address.AddressLine2
	r.DeliveryCompany1 = "" // 未使用のため空文字固定
	r.DeliveryCompany2 = "" // 未使用のため空文字固定
	r.DeliveryName = address.Name()
	r.DeliveryNameKana = width.Narrow.String(address.NameKana())
	r.DeliveryNameHonorific = "様" // 敬称は「様」固定
}

func (r *Receipt) SetClientDetails(address *uentity.Address) {
	if address == nil {
		return
	}
	r.ClientCode = strconv.FormatInt(address.AddressRevision.ID, 10)
	r.ClientPhoneNumber = address.PhoneNumber
	r.ClientPhoneNumberExtension = "" // 未使用のため空文字固定
	r.ClientPostalCode = address.PostalCode
	r.ClientAddress1 = address.ShortPath()
	r.ClientAddress2 = address.AddressLine2
	r.ClientName = address.Name()
	r.ClientNameKana = width.Narrow.String(address.NameKana())
}

func (r *Receipt) SetProductDetails(items entity.OrderItems, products map[int64]*entity.Product) {
	if len(items) < 1 {
		return
	}
	product, ok := products[items[0].ProductRevisionID]
	if !ok {
		return
	}
	r.Product1ID = product.ID
	r.Product1Name = product.Name
	if len(items) < 2 {
		return
	}
	product, ok = products[items[1].ProductRevisionID]
	if !ok {
		return
	}
	r.Product2ID = product.ID
	r.Product2Name = product.Name
}

func (r *Receipt) SetWebCollectDetails() {
	r.WebCollectDeliveryAmount = 0
	r.WebCollectDeliveryTax = 0
	r.WebCollectRegistration = 0
	r.WebCollectStoreNumber = ""
	r.WebCollectAcceptanceNumber1 = ""
	r.WebCollectAcceptanceNumber2 = ""
	r.WebCollectAcceptanceNumber3 = ""
}

func (r *Receipt) SetCollectionAgencyDetails() {
	r.CollectionAgencyUsage = 0 // 0：利用しない
	r.CollectionAgencyBillingAmount = 0
	r.CollectionAgencyBillingTax = 0
	r.CollectionAgencyBillingPostalCode = ""
	r.CollectionAgencyBillingAddress1 = ""
	r.CollectionAgencyBillingAddress2 = ""
	r.CollectionAgencyBillingCompany1 = ""
	r.CollectionAgencyBillingCompany2 = ""
	r.CollectionAgencyBillingName = ""
	r.CollectionAgencyBillingNameKana = ""
	r.CollectionAgencyInquiryPostalCode = ""
	r.CollectionAgencyInquiryAddress1 = ""
	r.CollectionAgencyInquiryAddress2 = ""
	r.CollectionAgencyInquiryPhoneNumber = ""
	r.CollectionAgencyManagementNumber = ""
	r.CollectionAgencyItemName = ""
	r.CollectionAgencyRemarks = ""
}

func (r *Receipt) SetDeliveryNotificationDetails() {
	r.DeliveryNotificationEmailUsage = 0 // 0：利用しない
	r.DeliveryNotificationEmail = ""
	r.DeliveryNotificationEmailMessage = ""
	r.DeliveryCompletionEmailUsage = 0 // 0：利用しない
	r.DeliveryCompletionEmail = ""
	r.DeliveryCompletionEmailMessage = ""
}

func (r *Receipt) SetMailDropNotificationDetails() {
	r.MailDropNotificationEmailUsage = 0 // 0：利用しない
	r.MailDropNotificationEmail = ""
	r.MailDropNotificationEmailMessage = ""
	r.MailDropCompletionDeliveryEmailUsage = 0 // 0：利用しない
	r.MailDropCompletionDeliveryEmail = ""
	r.MailDropCompletionDeliveryEmailMessage = ""
	r.MailDropCompletionClientEmailusage = ""
	r.MailDropCompletionClientEmail = ""
	r.MailDropCompletionClientEmailMessage = ""
}

func (r *Receipt) SetSearchDetails(order *entity.Order) {
	r.SearchKey1Title = "ふるマルユーザーID"
	r.SearchKey1Value = order.UserID
	r.SearchKey2Title = "ふるマル注文履歴ID"
	r.SearchKey2Value = order.ID
	r.SearchKey3Title = ""
	r.SearchKey3Value = ""
	r.SearchKey4Title = ""
	r.SearchKey4Value = ""
	r.SearchKey5Title = ""
	r.SearchKey5Value = ""
}

func (r *Receipt) Header() []string {
	return receiptHeaders
}

func (r *Receipt) Record() []string {
	return []string{
		r.OrderID,
		strconv.FormatInt(int64(r.ServiceType), 10),
		strconv.FormatInt(int64(r.ShippingType), 10),
		r.YamatoOrderID,
		r.ExpectedShippingDate,
		r.ExpectedDeliveryDate,
		string(r.ExpectedDeliveryTimeFrame),
		r.DeliveryCode,
		r.DeliveryPhoneNumber,
		r.DeliveryPhoneNumberExtension,
		r.DeliveryPostalCode,
		r.DeliveryAddress1,
		r.DeliveryAddress2,
		r.DeliveryCompany1,
		r.DeliveryCompany2,
		r.DeliveryName,
		r.DeliveryNameKana,
		r.DeliveryNameHonorific,
		r.ClientCode,
		r.ClientPhoneNumber,
		r.ClientPhoneNumberExtension,
		r.ClientPostalCode,
		r.ClientAddress1,
		r.ClientAddress2,
		r.ClientName,
		r.ClientNameKana,
		r.Product1ID,
		r.Product1Name,
		r.Product2ID,
		r.Product2Name,
		r.Handling1,
		r.Handling2,
		r.Note,
		strconv.FormatInt(r.WebCollectDeliveryAmount, 10),
		strconv.FormatInt(r.WebCollectDeliveryTax, 10),
		strconv.FormatInt(r.BranchHoldType, 10),
		r.BranchCode,
		strconv.FormatInt(r.IssuedQuantity, 10),
		r.PrintingQuantity,
		r.BillingCode,
		r.BillingGroupCode,
		r.ShippingCostControlNumber,
		strconv.FormatInt(r.WebCollectRegistration, 10),
		r.WebCollectStoreNumber,
		r.WebCollectAcceptanceNumber1,
		r.WebCollectAcceptanceNumber2,
		r.WebCollectAcceptanceNumber3,
		strconv.FormatInt(r.DeliveryNotificationEmailUsage, 10),
		r.DeliveryNotificationEmail,
		r.InputMethod,
		r.DeliveryNotificationEmailMessage,
		strconv.FormatInt(r.DeliveryCompletionEmailUsage, 10),
		r.DeliveryCompletionEmail,
		r.DeliveryCompletionEmailMessage,
		strconv.FormatInt(r.CollectionAgencyUsage, 10),
		r.Reserve1,
		strconv.FormatInt(r.CollectionAgencyBillingAmount, 10),
		strconv.FormatInt(r.CollectionAgencyBillingTax, 10),
		r.CollectionAgencyBillingPostalCode,
		r.CollectionAgencyBillingAddress1,
		r.CollectionAgencyBillingAddress2,
		r.CollectionAgencyBillingCompany1,
		r.CollectionAgencyBillingCompany2,
		r.CollectionAgencyBillingName,
		r.CollectionAgencyBillingNameKana,
		r.CollectionAgencyInquiryPostalCode,
		r.CollectionAgencyInquiryAddress1,
		r.CollectionAgencyInquiryAddress2,
		r.CollectionAgencyInquiryPhoneNumber,
		r.CollectionAgencyManagementNumber,
		r.CollectionAgencyItemName,
		r.CollectionAgencyRemarks,
		r.MultiplePackagesGroupingKey,
		r.SearchKey1Title,
		r.SearchKey1Value,
		r.SearchKey2Title,
		r.SearchKey2Value,
		r.SearchKey3Title,
		r.SearchKey3Value,
		r.SearchKey4Title,
		r.SearchKey4Value,
		r.SearchKey5Title,
		r.SearchKey5Value,
		r.Reserve2,
		r.Reserve3,
		strconv.FormatInt(r.MailDropNotificationEmailUsage, 10),
		r.MailDropNotificationEmail,
		r.MailDropNotificationEmailMessage,
		strconv.FormatInt(r.MailDropCompletionDeliveryEmailUsage, 10),
		r.MailDropCompletionDeliveryEmail,
		r.MailDropCompletionDeliveryEmailMessage,
		r.MailDropCompletionClientEmailusage,
		r.MailDropCompletionClientEmail,
		r.MailDropCompletionClientEmailMessage,
	}
}

func NewReceipts(params *ReceiptsParams) []exporter.Receipt {
	res := make([]exporter.Receipt, 0, len(params.Orders))
	for _, order := range params.Orders {
		itemsMap := order.OrderItems.GroupByFulfillmentID()
		for _, fulfillment := range order.OrderFulfillments {
			in := &ReceiptParams{
				Order:       order,
				Fulfillment: fulfillment,
				Items:       itemsMap[fulfillment.ID],
				Addresses:   params.Addresses,
				Products:    params.Products,
			}
			res = append(res, NewReceipt(in))
		}
	}
	return res
}
