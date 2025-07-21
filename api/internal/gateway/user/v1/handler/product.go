package handler

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products")

	r.GET("", h.ListProducts)
	r.GET("/:productId", h.GetProduct)
	r.GET("/merchant-feed", h.GetMerchantCenterFeed)
}

func (h *handler) ListProducts(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	var shopID string
	if coordinatorID := util.GetQuery(ctx, "coordinatorId", ""); coordinatorID != "" {
		coordinator, err := h.getCoordinator(ctx, coordinatorID)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		shopID = coordinator.ShopID
	}

	in := &store.ListProductsInput{
		Limit:            limit,
		Offset:           offset,
		OnlyPublished:    true,
		ExcludeOutOfSale: true,
		ExcludeDeleted:   true,
		ShopID:           shopID,
		Orders: []*store.ListProductsOrder{
			// 売り切れでないもの順 && 公開日時が新しいもの順
			{Key: store.ListProductsOrderBySoldOut, OrderByASC: true},
			{Key: store.ListProductsOrderByStartAt, OrderByASC: false},
		},
	}
	products, total, err := h.store.ListProducts(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(products) == 0 {
		res := &response.ProductsResponse{
			Products:     []*response.Product{},
			Coordinators: []*response.Coordinator{},
			Producers:    []*response.Producer{},
			Categories:   []*response.Category{},
			ProductTypes: []*response.ProductType{},
			ProductTags:  []*response.ProductTag{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		coordinators service.Coordinators
		producers    service.Producers
		categories   service.Categories
		productTypes service.ProductTypes
		productTags  service.ProductTags
		productRates service.ProductRates
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinators, err = h.multiGetCoordinators(ectx, products.CoordinatorIDs())
		return
	})
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, products.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, products.ProductTypeIDs())
		if err != nil || len(productTypes) == 0 {
			return
		}
		categories, err = h.multiGetCategories(ectx, productTypes.CategoryIDs())
		return
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, products.ProductTagIDs())
		return
	})
	eg.Go(func() (err error) {
		productRates, err = h.aggregateProductRates(ectx, products.IDs()...)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	details := &service.ProductDetailsParams{
		Categories:   categories.Map(),
		ProductTypes: productTypes.Map(),
		ProductRates: productRates.MapByProductID(),
	}
	sproducts := service.NewProducts(products, details)

	res := &response.ProductsResponse{
		Products:     sproducts.Response(),
		Coordinators: coordinators.Response(),
		Producers:    producers.Response(),
		Categories:   categories.Response(),
		ProductTypes: productTypes.Response(),
		ProductTags:  productTags.Response(),
		Total:        total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetProduct(ctx *gin.Context) {
	product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		coordinator *service.Coordinator
		producer    *service.Producer
		category    *service.Category
		productType *service.ProductType
		productTags service.ProductTags
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		coordinator, err = h.getCoordinator(ectx, product.CoordinatorID)
		return
	})
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, product.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		productType, err = h.getProductType(ectx, product.ProductTypeID)
		if err != nil {
			return
		}
		category, err = h.getCategory(ectx, productType.CategoryID)
		return
	})
	eg.Go(func() (err error) {
		productTags, err = h.multiGetProductTags(ectx, product.ProductTagIDs)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &response.ProductResponse{
		Product:     product.Response(),
		Coordinator: coordinator.Response(),
		Producer:    producer.Response(),
		Category:    category.Response(),
		ProductType: productType.Response(),
		ProductTags: productTags.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) listProducts(
	ctx context.Context,
	in *store.ListProductsInput,
) (service.Products, error) {
	products, _, err := h.store.ListProducts(ctx, in)
	if err != nil || len(products) == 0 {
		return service.Products{}, err
	}
	products = products.FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
}

func (h *handler) multiGetProducts(
	ctx context.Context,
	productIDs []string,
) (service.Products, error) {
	if len(productIDs) == 0 {
		return service.Products{}, nil
	}
	in := &store.MultiGetProductsInput{
		ProductIDs: productIDs,
	}
	products, err := h.store.MultiGetProducts(ctx, in)
	if err != nil || len(products) == 0 {
		return service.Products{}, err
	}
	products = products.FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
}

func (h *handler) multiGetProductsByRevision(
	ctx context.Context,
	revisionIDs []int64,
) (service.Products, error) {
	if len(revisionIDs) == 0 {
		return service.Products{}, nil
	}
	in := &store.MultiGetProductsByRevisionInput{
		ProductRevisionIDs: revisionIDs,
	}
	products, err := h.store.MultiGetProductsByRevision(ctx, in)
	if err != nil || len(products) == 0 {
		return service.Products{}, err
	}
	products = products.FilterByPublished()
	details, err := h.getProductDetails(ctx, products.IDs()...)
	if err != nil {
		return nil, err
	}
	return service.NewProducts(products, details), nil
}

func (h *handler) getProduct(ctx context.Context, productID string) (*service.Product, error) {
	in := &store.GetProductInput{
		ProductID: productID,
	}
	product, err := h.store.GetProduct(ctx, in)
	if err != nil {
		return nil, err
	}
	if !product.Public {
		// 非公開のものは利用者側に表示しない
		return nil, exception.ErrNotFound
	}
	details, err := h.getProductDetails(ctx, productID)
	if err != nil {
		return nil, err
	}
	category := details.Categories[product.TypeID]
	rate := details.ProductRates[productID]
	return service.NewProduct(product, category, rate), nil
}

func (h *handler) getProductDetails(
	ctx context.Context,
	productIDs ...string,
) (*service.ProductDetailsParams, error) {
	var (
		categories   service.Categories
		productTypes service.ProductTypes
		productRates service.ProductRates
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, productIDs)
		if err != nil {
			return
		}
		categories, err = h.multiGetCategories(ectx, productTypes.CategoryIDs())
		return
	})
	eg.Go(func() (err error) {
		productRates, err = h.aggregateProductRates(ectx, productIDs...)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	res := &service.ProductDetailsParams{
		Categories:   categories.Map(),
		ProductTypes: productTypes.Map(),
		ProductRates: productRates.MapByProductID(),
	}
	return res, nil
}

type MerchantCenterFeed struct {
	XMLName xml.Name               `xml:"rss"`
	Version string                 `xml:"version,attr"`
	Channel *MerchantCenterChannel `xml:"channel"`
}

type MerchantCenterChannel struct {
	Title       string                `xml:"title"`
	Link        string                `xml:"link"`
	Description string                `xml:"description"`
	Items       []*MerchantCenterItem `xml:"item"`
}

type MerchantCenterItem struct {
	ID                    string `xml:"g:id"`
	Title                 string `xml:"g:title"`
	Description           string `xml:"g:description"`
	Link                  string `xml:"g:link"`
	ImageLink             string `xml:"g:image_link"`
	Condition             string `xml:"g:condition"`
	Availability          string `xml:"g:availability"`
	Price                 string `xml:"g:price"`
	Brand                 string `xml:"g:brand"`
	GTIN                  string `xml:"g:gtin,omitempty"`
	MPN                   string `xml:"g:mpn,omitempty"`
	GoogleProductCategory string `xml:"g:google_product_category,omitempty"`
	ProductType           string `xml:"g:product_type,omitempty"`
	ItemGroupID           string `xml:"g:item_group_id,omitempty"`
	ShippingWeight        string `xml:"g:shipping_weight,omitempty"`
}

func (h *handler) GetMerchantCenterFeed(ctx *gin.Context) {
	const maxProducts = 10000

	in := &store.ListProductsInput{
		Limit:            maxProducts,
		Offset:           0,
		OnlyPublished:    true,
		ExcludeOutOfSale: true,
		ExcludeDeleted:   true,
		Orders: []*store.ListProductsOrder{
			{Key: store.ListProductsOrderByStartAt, OrderByASC: false},
		},
	}

	products, _, err := h.store.ListProducts(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	if len(products) == 0 {
		feed := &MerchantCenterFeed{
			Version: "2.0",
			Channel: &MerchantCenterChannel{
				Title:       "ふるマル - 全国ふるさとマルシェ",
				Link:        "https://furumaru.jp",
				Description: "地域・地方の特産品を扱うECマーケットプレイス",
				Items:       []*MerchantCenterItem{},
			},
		}
		ctx.Header("Content-Type", "application/xml; charset=utf-8")
		ctx.XML(http.StatusOK, feed)
		return
	}

	var (
		producers    service.Producers
		categories   service.Categories
		productTypes service.ProductTypes
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, products.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		productTypes, err = h.multiGetProductTypes(ectx, products.ProductTypeIDs())
		if err != nil || len(productTypes) == 0 {
			return
		}
		categories, err = h.multiGetCategories(ectx, productTypes.CategoryIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	producersMap := producers.Map()
	categoriesMap := categories.Map()
	productTypesMap := productTypes.Map()

	items := make([]*MerchantCenterItem, 0, len(products))
	for _, product := range products {
		if !product.Public {
			continue
		}

		producer := producersMap[product.ProducerID]
		productType := productTypesMap[product.TypeID]

		var category *service.Category
		if productType != nil {
			category = categoriesMap[productType.CategoryID]
		}

		description := product.Description
		if description == "" {
			description = product.Name
		}
		if len(description) > 5000 {
			description = description[:4997] + "..."
		}

		imageLink := ""
		if len(product.Media) > 0 {
			imageLink = product.Media[0].URL
		}

		brandName := ""
		if producer != nil {
			brandName = producer.Username
		}

		productTypeName := ""
		if productType != nil {
			productTypeName = productType.Name
		}

		availability := "in_stock"
		if product.Inventory <= 0 {
			availability = "out_of_stock"
		}

		price := fmt.Sprintf("%.0f JPY", float64(product.Price))

		item := &MerchantCenterItem{
			ID:             product.ID,
			Title:          product.Name,
			Description:    description,
			Link:           fmt.Sprintf("https://furumaru.jp/products/%s", product.ID),
			ImageLink:      imageLink,
			Condition:      "new",
			Availability:   availability,
			Price:          price,
			Brand:          brandName,
			ProductType:    productTypeName,
			ShippingWeight: fmt.Sprintf("%.0f g", product.Weight),
		}

		if category != nil {
			item.GoogleProductCategory = category.Name
		}

		items = append(items, item)
	}

	feed := &MerchantCenterFeed{
		Version: "2.0",
		Channel: &MerchantCenterChannel{
			Title:       "ふるマル - 全国ふるさとマルシェ",
			Link:        "https://furumaru.jp",
			Description: "地域・地方の特産品を扱うECマーケットプレイス",
			Items:       items,
		},
	}

	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	ctx.XML(http.StatusOK, feed)
}
