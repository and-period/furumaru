package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type LiveProduct struct {
	response.LiveProduct
}

type LiveProducts []*LiveProduct

func NewLiveProduct(product *entity.Product) *LiveProduct {
	var (
		thumbnailURL string
		thumbnails   Images
	)
	for _, media := range product.Media {
		if !media.IsThumbnail {
			continue
		}
		thumbnailURL = media.URL
		thumbnails = NewImages(media.Images)
		break
	}
	return &LiveProduct{
		LiveProduct: response.LiveProduct{
			ProductID:    product.ID,
			Name:         product.Name,
			Price:        product.Price,
			Inventory:    product.Inventory,
			ThumbnailURL: thumbnailURL,
			Thumbnails:   thumbnails.Response(),
		},
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
	return res
}

func (ps LiveProducts) Response() []*response.LiveProduct {
	res := make([]*response.LiveProduct, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
