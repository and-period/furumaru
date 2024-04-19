package service

import (
	"sort"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type LiveProduct struct {
	response.LiveProduct
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
		LiveProduct: response.LiveProduct{
			ProductID:    product.ID,
			Name:         product.Name,
			Price:        product.Price,
			Inventory:    product.Inventory,
			ThumbnailURL: thumbnailURL,
		},
		isSale: isSale,
	}
}

func (p *LiveProduct) Response() *response.LiveProduct {
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

func (ps LiveProducts) Response() []*response.LiveProduct {
	res := make([]*response.LiveProduct, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
