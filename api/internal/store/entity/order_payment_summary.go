package entity

import (
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/shopspring/decimal"
)

const taxRate int64 = 10 // 税率（10%)

var (
	percent    = decimal.NewFromInt(100)
	taxPercent = decimal.NewFromInt(taxRate).Div(percent)
)

// チェックアウト前の支払い情報
type OrderPaymentSummary struct {
	Subtotal    int64 // 購入金額
	Discount    int64 // 割引金額
	ShippingFee int64 // 配送手数料
	Tax         int64 // 消費税
	TaxRate     int64 // 消費税率(%)
	Total       int64 // 合計金額
}

type NewOrderPaymentSummaryParams struct {
	Address   *entity.Address
	Baskets   CartBaskets
	Products  Products
	Shipping  *Shipping
	Promotion *Promotion
}

func NewOrderPaymentSummary(params *NewOrderPaymentSummaryParams) (*OrderPaymentSummary, error) {
	var shippingFee int64
	// 商品購入価格の算出
	subtotal, err := params.Baskets.TotalPrice(params.Products.Map())
	if err != nil {
		return nil, err
	}
	// 商品配送料金の算出
	for _, basket := range params.Baskets {
		if params.Address == nil {
			break // 配送先情報がない場合、配送料金は算出しない
		}
		fee, err := params.Shipping.CalcShippingFee(basket.BoxSize, basket.BoxType, subtotal, params.Address.PrefectureCode)
		if err != nil {
			return nil, err
		}
		shippingFee += fee
	}
	// 割引金額の算出
	discount := params.Promotion.CalcDiscount(subtotal, shippingFee)
	// 支払い金額の算出
	dsubtotal := decimal.NewFromInt(subtotal).Add(decimal.NewFromInt(shippingFee))
	ddiscount := decimal.NewFromInt(discount)
	dtax := dsubtotal.Sub(ddiscount).Mul(taxPercent)
	dtotal := dsubtotal.Sub(ddiscount).Add(dtax)
	return &OrderPaymentSummary{
		Subtotal:    subtotal,
		Discount:    discount,
		ShippingFee: shippingFee,
		Tax:         dtax.IntPart(),
		TaxRate:     taxRate,
		Total:       dtotal.IntPart(),
	}, nil
}
