package sagawa

import (
	"strconv"
	"strings"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/exporter"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

var receiptHeaders = []string{
	"お届け先コード取得区分",
	"お届け先コード",
	"お届け先電話番号",
	"お届け先郵便番号",
	"お届け先住所１",
	"お届け先住所２",
	"お届け先住所３",
	"お届け先名称１",
	"お届け先名称２",
	"お客様管理番号",
	"お客様コード",
	"部署ご担当者コード",
	"取得区分",
	"部署ご担当者コード",
	"部署ご担当者名称",
	"荷送人電話番号",
	"ご依頼主コード取得区分",
	"ご依頼主コード",
	"ご依頼主電話番号",
	"ご依頼主郵便番号",
	"ご依頼主住所１",
	"ご依頼主住所２",
	"ご依頼主名称１",
	"ご依頼主名称２",
	"荷姿",
	"品名１",
	"品名２",
	"品名３",
	"品名４",
	"品名５",
	"荷札荷姿",
	"荷札品名1",
	"荷札品名2",
	"荷札品名3",
	"荷札品名4",
	"荷札品名5",
	"荷札品名6",
	"荷札品名7",
	"荷札品名8",
	"荷札品名9",
	"荷札品名10",
	"荷札品名11",
	"出荷個数",
	"スピード指定",
	"クール便指定",
	"配達日",
	"配達指定時間帯",
	"配達指定時間（時分）",
	"代引金額",
	"消費税",
	"決済種別",
	"保険金額",
	"指定シール1",
	"指定シール2",
	"指定シール3",
	"営業所受取",
	"ＳＲＣ区分",
	"営業所受取営業所コード",
	"元着区分",
	"メールアドレス",
	"ご不在時連絡先",
	"出荷日",
	"お問い合せ送り状No.",
	"出荷場印字区分",
	"集約解除指定",
	"編集1",
	"編集2",
	"編集3",
	"編集4",
	"編集5",
	"編集6",
	"編集7",
	"編集8",
	"編集9",
	"編集10",
}

// ShippingType - クール便指定
type ShippingType string

const (
	ShippingTypeNormal        ShippingType = "001" // 宅急便
	ShippingTypeFreeze        ShippingType = "002" // 宅急便冷蔵（飛脚クール便（冷蔵））
	ShippingTypeRefrigeration ShippingType = "003" // 宅急便冷凍（飛脚クール便（冷凍））
)

// DeliveryTimeFrame - 	配達指定時間帯
type DeliveryTimeFrame string

const (
	DeliveryTimeFrameNone    DeliveryTimeFrame = ""   // 未指定
	DeliveryTimeFrameMorning DeliveryTimeFrame = "01" // 午前中
	DeliveryTimeFrame1214    DeliveryTimeFrame = "12" // 12：00～14：00
	DeliveryTimeFrame1416    DeliveryTimeFrame = "14" // 14：00～16：00
	DeliveryTimeFrame1618    DeliveryTimeFrame = "16" // 16：00～18：00
	DeliveryTimeFrame1820    DeliveryTimeFrame = "18" // 18：00～20：00
	DeliveryTimeFrame1821    DeliveryTimeFrame = "04" // 18：00～21：00
	DeliveryTimeFrame1921    DeliveryTimeFrame = "19" // 19：00～21：00
)

// Receipt - 佐川e飛伝Ⅲ
// See: https://contents.raku-uru.jp/documents/manual/data/ordercsv/ehiden3.html
type Receipt struct {
	DeliveryCodeType                  string            // お届け先コード取得区分
	DeliveryCode                      string            // お届け先コード
	DeliveryPhoneNumber               string            // お届け先電話番号
	DeliveryPostalCode                string            // お届け先郵便番号
	DeliveryAddress1                  string            // お届け先住所１
	DeliveryAddress2                  string            // お届け先住所２
	DeliveryAddress3                  string            // お届け先住所３
	DeliveryLastname                  string            // お届け先名称１
	DeliveryFirstname                 string            // お届け先名称２
	OrderID                           string            // お客様管理番号
	UserID                            string            // お客様コード
	DivisionClientCodeType            string            // 部署ご担当者コード取得区分
	DivisionClientCode                string            // 部署ご担当者コード
	DivixionClientName                string            // 部署ご担当者名称
	ShipperPhoneNumber                string            // 荷送人電話番号
	ClientCodeType                    string            // ご依頼主コード取得区分
	ClientCode                        string            // ご依頼主コード
	ClientPhoneNumber                 string            // ご依頼主電話番号
	ClientPostalCode                  string            // ご依頼主郵便番号
	ClientAddress1                    string            // ご依頼主住所１
	ClientAddress2                    string            // ご依頼主住所２
	ClientLastname                    string            // ご依頼主名称１
	ClientFirstname                   string            // ご依頼主名称２
	Packing                           string            // 荷姿
	Product1Name                      string            // 品名１
	Product2Name                      string            // 品名２
	Product3Name                      string            // 品名３
	Product4Name                      string            // 品名４
	Product5Name                      string            // 品名５
	Tagpackaging                      string            // 荷札荷姿
	TagProduct1Name                   string            // 荷札品名1
	TagProduct2Name                   string            // 荷札品名2
	TagProduct3Name                   string            // 荷札品名3
	TagProduct4Name                   string            // 荷札品名4
	TagProduct5Name                   string            // 荷札品名5
	TagProduct6Name                   string            // 荷札品名6
	TagProduct7Name                   string            // 荷札品名7
	TagProduct8Name                   string            // 荷札品名8
	TagProduct9Name                   string            // 荷札品名9
	TagProduct10Name                  string            // 荷札品名10
	TagProduct11Name                  string            // 荷札品名11
	ShipmentQuantity                  int64             // 出荷個数
	SpeedSpecification                string            // スピード指定
	ShippingType                      ShippingType      // クール便指定
	DeliveryDate                      string            // 配達日
	DeliveryTimeFrame                 DeliveryTimeFrame // 配達指定時間帯
	DeliveryTime                      string            // 配達指定時間（時分）
	DeliveryAmount                    int64             // 代引金額
	DeliveryTax                       int64             // 消費税
	PaymentMethod                     string            // 決済種別
	InsuranceAmount                   int64             // 保険金額
	DesignatedSticker1                string            // 指定シール1
	DesignatedSticker2                string            // 指定シール2
	DesignatedSticker3                string            // 指定シール3
	OfficePickUp                      string            // 営業所受取
	SrcType                           string            // ＳＲＣ区分
	OfficePickUpCode                  string            // 営業所受取営業所コード
	OriginalDestinationClassification string            // 元着区分
	Email                             string            // メールアドレス
	AbsenceContact                    string            // ご不在時連絡先
	ShippingDate                      string            // 出荷日
	ReceiptNumber                     string            // お問い合せ送り状No.
	ShippingFacilityPrintType         string            // 出荷場印字区分
	UnaggregatedSpec                  string            // 集約解除指定
	Reserve1                          string            // 編集1
	Reserve2                          string            // 編集2
	Reserve3                          string            // 編集3
	Reserve4                          string            // 編集4
	Reserve5                          string            // 編集5
	Reserve6                          string            // 編集6
	Reserve7                          string            // 編集7
	Reserve8                          string            // 編集8
	Reserve9                          string            // 編集9
	Reserve10                         string            // 編集10
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
	receipt.SetReceiptDetails(params.Order, params.Fulfillment)
	receipt.SetDeliveryDetails(params.Addresses[params.Fulfillment.AddressRevisionID])
	receipt.SetClientDetails(params.Addresses[params.Order.OrderPayment.AddressRevisionID])
	receipt.SetProductDetails(params.Items, params.Products)
	receipt.SetTagProductDetails()
	return receipt
}

func (r *Receipt) SetReceiptDetails(order *entity.Order, fulfillment *entity.OrderFulfillment) {
	r.OrderID = order.OrderID
	r.UserID = order.UserID
	r.DivisionClientCodeType = ""
	r.DivisionClientCode = ""
	r.DivixionClientName = ""
	r.ShipperPhoneNumber = ""
	r.ShipmentQuantity = 1 // 1固定
	r.SpeedSpecification = ""
	r.ShippingType = NewShippingType(fulfillment.ShippingType)
	r.DeliveryDate = ""      // TODO: 購入フローの改修時に対応
	r.DeliveryTimeFrame = "" // TODO: 購入フローの改修時に対応
	r.DeliveryTime = ""
	r.DeliveryAmount = 0
	r.DeliveryTax = 0
	r.PaymentMethod = ""
	r.InsuranceAmount = 0
	r.DesignatedSticker1 = ""
	r.DesignatedSticker2 = ""
	r.DesignatedSticker3 = ""
	r.SrcType = ""
	r.OfficePickUp = ""
	r.OfficePickUpCode = ""
	r.OriginalDestinationClassification = ""
	r.Email = ""
	r.AbsenceContact = ""
	r.ShippingDate = ""
	r.ReceiptNumber = ""
	r.ShippingFacilityPrintType = ""
	r.UnaggregatedSpec = ""
	r.Reserve1 = ""
	r.Reserve2 = ""
	r.Reserve3 = ""
	r.Reserve4 = ""
	r.Reserve5 = ""
	r.Reserve6 = ""
	r.Reserve7 = ""
	r.Reserve8 = ""
	r.Reserve9 = ""
	r.Reserve10 = ""
}

func (r *Receipt) SetDeliveryDetails(address *uentity.Address) {
	if address == nil {
		return
	}
	r.DeliveryCodeType = ""
	r.DeliveryCode = ""
	r.DeliveryPhoneNumber = address.PhoneNumber
	r.DeliveryPostalCode = address.PostalCode
	r.DeliveryAddress1 = address.Prefecture
	r.DeliveryAddress2 = strings.Join([]string{address.City, address.AddressLine1}, " ")
	r.DeliveryAddress3 = address.AddressLine2
	r.DeliveryLastname = address.Lastname
	r.DeliveryFirstname = address.Firstname
}

func (r *Receipt) SetClientDetails(address *uentity.Address) {
	if address == nil {
		return
	}
	r.ClientCodeType = ""
	r.ClientCode = ""
	r.ClientPhoneNumber = address.PhoneNumber
	r.ClientPostalCode = address.PostalCode
	r.ClientAddress1 = address.ShortPath()
	r.ClientAddress2 = address.AddressLine2
	r.ClientLastname = address.Lastname
	r.ClientFirstname = address.Firstname
}

func (r *Receipt) SetProductDetails(items entity.OrderItems, products map[int64]*entity.Product) {
	r.Packing = ""
	if len(items) < 1 {
		return
	}
	product, ok := products[items[0].ProductRevisionID]
	if !ok {
		return
	}
	r.Product1Name = product.Name
	if len(items) < 2 {
		return
	}
	product, ok = products[items[1].ProductRevisionID]
	if !ok {
		return
	}
	r.Product2Name = product.Name
	if len(items) < 3 {
		return
	}
	product, ok = products[items[2].ProductRevisionID]
	if !ok {
		return
	}
	r.Product3Name = product.Name
	if len(items) < 4 {
		return
	}
	product, ok = products[items[3].ProductRevisionID]
	if !ok {
		return
	}
	r.Product4Name = product.Name
	if len(items) < 5 {
		return
	}
	product, ok = products[items[4].ProductRevisionID]
	if !ok {
		return
	}
	r.Product5Name = product.Name
}

func (r *Receipt) SetTagProductDetails() {
	r.Tagpackaging = ""
	r.TagProduct1Name = ""
	r.TagProduct2Name = ""
	r.TagProduct3Name = ""
	r.TagProduct4Name = ""
	r.TagProduct5Name = ""
	r.TagProduct6Name = ""
	r.TagProduct7Name = ""
	r.TagProduct8Name = ""
	r.TagProduct9Name = ""
	r.TagProduct10Name = ""
	r.TagProduct11Name = ""
}

func (r *Receipt) Header() []string {
	return receiptHeaders
}

func (r *Receipt) Record() []string {
	return []string{
		r.DeliveryCodeType,
		r.DeliveryCode,
		r.DeliveryPhoneNumber,
		r.DeliveryPostalCode,
		r.DeliveryAddress1,
		r.DeliveryAddress2,
		r.DeliveryAddress3,
		r.DeliveryLastname,
		r.DeliveryFirstname,
		r.OrderID,
		r.UserID,
		r.DivisionClientCodeType,
		r.DivisionClientCode,
		r.DivixionClientName,
		r.ShipperPhoneNumber,
		r.ClientCodeType,
		r.ClientCode,
		r.ClientPhoneNumber,
		r.ClientPostalCode,
		r.ClientAddress1,
		r.ClientAddress2,
		r.ClientLastname,
		r.ClientFirstname,
		r.Packing,
		r.Product1Name,
		r.Product2Name,
		r.Product3Name,
		r.Product4Name,
		r.Product5Name,
		r.Tagpackaging,
		r.TagProduct1Name,
		r.TagProduct2Name,
		r.TagProduct3Name,
		r.TagProduct4Name,
		r.TagProduct5Name,
		r.TagProduct6Name,
		r.TagProduct7Name,
		r.TagProduct8Name,
		r.TagProduct9Name,
		r.TagProduct10Name,
		r.TagProduct11Name,
		strconv.FormatInt(r.ShipmentQuantity, 10),
		r.SpeedSpecification,
		string(r.ShippingType),
		r.DeliveryDate,
		string(r.DeliveryTimeFrame),
		r.DeliveryTime,
		strconv.FormatInt(r.DeliveryAmount, 10),
		strconv.FormatInt(r.DeliveryTax, 10),
		r.PaymentMethod,
		strconv.FormatInt(r.InsuranceAmount, 10),
		r.DesignatedSticker1,
		r.DesignatedSticker2,
		r.DesignatedSticker3,
		r.OfficePickUp,
		r.SrcType,
		r.OfficePickUpCode,
		r.OriginalDestinationClassification,
		r.Email,
		r.AbsenceContact,
		r.ShippingDate,
		r.ReceiptNumber,
		r.ShippingFacilityPrintType,
		r.UnaggregatedSpec,
		r.Reserve1,
		r.Reserve2,
		r.Reserve3,
		r.Reserve4,
		r.Reserve5,
		r.Reserve6,
		r.Reserve7,
		r.Reserve8,
		r.Reserve9,
		r.Reserve10,
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
