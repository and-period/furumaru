package general

import (
	"strconv"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/exporter"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

var receiptHeaders = []string{
	"注文管理番号",
	"ユーザーID",
	"コーディネータID",
	"お届け希望日",
	"お届け希望時間帯",
	"お届け先名",
	"お届け先名（かな）",
	"お届け先電話番号",
	"お届け先郵便番号",
	"お届け先都道府県",
	"お届け先市区町村",
	"お届け先町名・番地",
	"お届け先ビル名・号室など",
	"ご依頼主名",
	"ご依頼主名（かな）",
	"ご依頼主電話番号",
	"ご依頼主郵便番号",
	"ご依頼主都道府県",
	"ご依頼主市区町村",
	"ご依頼主町名・番地",
	"ご依頼主ビル名・号室など",
	"商品コード1",
	"商品名1",
	"商品コード2",
	"商品名2",
	"決済手段",
	"商品金額",
	"割引金額",
	"配送手数料",
	"合計金額",
	"注文日時",
}

// Receipt - 送り状作成情報
type Receipt struct {
	OrderID                   string    // 注文管理番号
	UserID                    string    // ユーザーID
	CoordinatorID             string    // コーディネータID
	ExpectedDelveryDate       string    // お届け希望日
	ExpectedDeliveryTimeFrame string    // お届け希望時間帯
	DeliveryName              string    // お届け先名
	DeliveryNameKana          string    // お届け先名（かな）
	DeliveryPhoneNumber       string    // お届け先電話番号
	DeliveryPostalCode        string    // お届け先郵便番号
	DeliveryPrefecture        string    // お届け先都道府県
	DeliveryCity              string    // お届け先市区町村
	DeliveryAddressLine1      string    // お届け先町名・番地
	DeliveryAddressLine2      string    // お届け先ビル名・号室など
	ClientName                string    // お届け先名
	ClientNameKana            string    // お届け先名（かな）
	ClientPhoneNumber         string    // お届け先電話番号
	ClientPostalCode          string    // お届け先郵便番号
	ClientPrefecture          string    // お届け先都道府県
	ClientCity                string    // お届け先市区町村
	ClientAddressLine1        string    // お届け先町名・番地
	ClientAddressLine2        string    // お届け先ビル名・号室など
	Product1ID                string    // 商品1管理番号
	Product1Name              string    // 商品1品名
	Product1Quantity          int64     // 商品1数量
	Product2ID                string    // 商品2管理番号
	Product2Name              string    // 商品2品名
	Product2Quantity          int64     // 商品2数量
	Product3ID                string    // 商品3管理番号
	Product3Name              string    // 商品3品名
	Product3Quantity          int64     // 商品3数量
	PaymentMethod             string    // 決済手段
	Subtotal                  int64     // 商品金額
	Discount                  int64     // 割引金額
	ShippingFee               int64     // 配送手数料
	Total                     int64     // 合計金額
	OrderedAt                 time.Time // 注文日時
	ShippingType              string    // 配送方法
	ShippingSize              string    // 箱のサイズ
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

func NewReceipt(params *ReceiptParams) exporter.Receipt {
	receipt := &Receipt{}
	receipt.SetReceiptDetails(params.Order)
	receipt.SetDeliveryDetails(params.Addresses[params.Fulfillment.AddressRevisionID])
	receipt.SetClientDetails(params.Addresses[params.Order.OrderPayment.AddressRevisionID])
	receipt.SetPaymentDetails(&params.Order.OrderPayment)
	receipt.SetFulfillmentDetails(params.Fulfillment)
	receipt.SetProductDetails(params.Items, params.Products)
	return receipt
}

func (r *Receipt) SetReceiptDetails(order *entity.Order) {
	r.OrderID = order.ID
	r.UserID = order.UserID
	r.CoordinatorID = order.CoordinatorID
}

func (r *Receipt) SetDeliveryDetails(address *uentity.Address) {
	if address == nil {
		return
	}
	r.DeliveryName = address.Name()
	r.DeliveryNameKana = address.NameKana()
	r.DeliveryPhoneNumber = address.PhoneNumber
	r.DeliveryPostalCode = address.PostalCode
	r.DeliveryPrefecture = address.Prefecture
	r.DeliveryCity = address.City
	r.DeliveryAddressLine1 = address.AddressLine1
	r.DeliveryAddressLine2 = address.AddressLine2
}

func (r *Receipt) SetClientDetails(address *uentity.Address) {
	if address == nil {
		return
	}
	r.ClientName = address.Name()
	r.ClientNameKana = address.NameKana()
	r.ClientPhoneNumber = address.PhoneNumber
	r.ClientPostalCode = address.PostalCode
	r.ClientPrefecture = address.Prefecture
	r.ClientCity = address.City
	r.ClientAddressLine1 = address.AddressLine1
	r.ClientAddressLine2 = address.AddressLine2
}

func (r *Receipt) SetPaymentDetails(payment *entity.OrderPayment) {
	r.PaymentMethod = payment.MethodType.String()
	r.Subtotal = payment.Subtotal
	r.Discount = payment.Discount
	r.ShippingFee = payment.ShippingFee
	r.Total = payment.Total
	r.OrderedAt = payment.OrderedAt
}

func (r *Receipt) SetFulfillmentDetails(fulfillment *entity.OrderFulfillment) {
	r.ExpectedDelveryDate = ""       // TODO: 購入フローの改修時に対応
	r.ExpectedDeliveryTimeFrame = "" // TODO: 購入フローの改修時に対応
	r.ShippingType = fulfillment.ShippingType.String()
	r.ShippingSize = fulfillment.BoxSize.String()
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
	r.Product1Quantity = items[0].Quantity
	if len(items) < 2 {
		return
	}
	product, ok = products[items[1].ProductRevisionID]
	if !ok {
		return
	}
	r.Product2ID = product.ID
	r.Product2Name = product.Name
	r.Product2Quantity = items[1].Quantity
	if len(items) < 3 {
		return
	}
	product, ok = products[items[2].ProductRevisionID]
	if !ok {
		return
	}
	r.Product3ID = product.ID
	r.Product3Name = product.Name
	r.Product3Quantity = items[2].Quantity
}

func (r *Receipt) Header() []string {
	return receiptHeaders
}

func (r *Receipt) Record() []string {
	return []string{
		r.OrderID,
		r.UserID,
		r.CoordinatorID,
		r.ExpectedDelveryDate,
		r.ExpectedDeliveryTimeFrame,
		r.DeliveryName,
		r.DeliveryNameKana,
		r.DeliveryPhoneNumber,
		r.DeliveryPostalCode,
		r.DeliveryPrefecture,
		r.DeliveryCity,
		r.DeliveryAddressLine1,
		r.DeliveryAddressLine2,
		r.ClientName,
		r.ClientNameKana,
		r.ClientPhoneNumber,
		r.ClientPostalCode,
		r.ClientPrefecture,
		r.ClientCity,
		r.ClientAddressLine1,
		r.ClientAddressLine2,
		r.Product1ID,
		r.Product1Name,
		strconv.FormatInt(r.Product1Quantity, 10),
		r.Product2ID,
		r.Product2Name,
		strconv.FormatInt(r.Product2Quantity, 10),
		r.Product3ID,
		r.Product3Name,
		strconv.FormatInt(r.Product3Quantity, 10),
		r.PaymentMethod,
		strconv.FormatInt(r.Subtotal, 10),
		strconv.FormatInt(r.Discount, 10),
		strconv.FormatInt(r.ShippingFee, 10),
		strconv.FormatInt(r.Total, 10),
		r.OrderedAt.Format(time.DateTime),
		r.ShippingType,
		r.ShippingSize,
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
