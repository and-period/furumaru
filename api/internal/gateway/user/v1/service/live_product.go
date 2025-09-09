package service

import (
	"sort"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type LiveProduct struct {
	types.LiveProduct
	isSale bool
}

type LiveProducts []*LiveProduct

func NewLiveProduct(product *entity.Product) *LiveProduct {
	var (
		thumbnailURL string
		isSale       bool
	)
	for _, media := range product.Media {
		if !media.IsThumbnail {
			continue
		}
		thumbnailURL = media.URL
		break
	}
	if product.Status == entity.ProductStatusForSale && product.Inventory > 0 {
		isSale = true
	}
	return &LiveProduct{
		LiveProduct: types.LiveProduct{
			ProductID:    product.ID,
			Name:         product.Name,
			Price:        product.Price,
			Inventory:    product.Inventory,
			ThumbnailURL: thumbnailURL,
		},
		isSale: isSale,
	}
}

func (p *LiveProduct) Response() *types.LiveProduct {
	return &p.LiveProduct
}

func NewLiveProducts(products entity.Products) LiveProducts {
	res := make(LiveProducts, len(products))
	for i := range products {
		res[i] = NewLiveProduct(products[i])
	}
	return res.SortByIsSale()
}

func (ps LiveProducts) SortByIsSale() LiveProducts {
	sort.SliceStable(ps, func(i, j int) bool {
		if !ps[j].isSale {
			return ps[i].isSale
		}
		return ps[i].ProductID <= ps[j].ProductID
	})
	return ps
}

func (ps LiveProducts) Response() []*types.LiveProduct {
	res := make([]*types.LiveProduct, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
