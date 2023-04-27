package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListProducts)
	arg.POST("", h.CreateProduct)
	arg.GET("/:productId", h.filterAccessProduct, h.GetProduct)
	arg.PATCH("/:productId", h.filterAccessProduct, h.UpdateProduct)
	arg.DELETE("/:productId", h.filterAccessProduct, h.DeleteProduct)
}

func (h *handler) filterAccessProduct(ctx *gin.Context) {
	params := &filterAccessParams{
		coordinator: func(ctx *gin.Context) (bool, error) {
			producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
			if err != nil {
				return false, err
			}
			product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
			if err != nil {
				return false, err
			}
			return producers.Contains(product.ProducerID), nil
		},
		producer: func(ctx *gin.Context) (bool, error) {
			product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
			if err != nil {
				return false, err
			}
			return currentAdmin(ctx, product.ProducerID), nil
		},
	}
	if err := filterAccess(ctx, params); err != nil {
		httpError(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) ListProducts(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	orders, err := h.newProductOrders(ctx)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &store.ListProductsInput{
		Name:       util.GetQuery(ctx, "name", ""),
		ProducerID: util.GetQuery(ctx, "producerId", ""),
		Limit:      limit,
		Offset:     offset,
		Orders:     orders,
	}
	if getRole(ctx) == service.AdminRoleCoordinator {
		producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
		if err != nil {
			httpError(ctx, err)
			return
		}
		in.ProducerIDs = producers.IDs()
	}
	sproducts, total, err := h.store.ListProducts(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	products := service.NewProducts(sproducts)
	if len(products) == 0 {
		res := &response.ProductsResponse{
			Products: products.Response(),
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	if err := h.getProductsDetails(ctx, products...); err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductsResponse{
		Products: products.Response(),
		Total:    total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newProductOrders(ctx *gin.Context) ([]*store.ListProductsOrder, error) {
	products := map[string]sentity.ProductOrderBy{
		"name":             sentity.ProductOrderByName,
		"public":           sentity.ProductOrderByPublic,
		"inventory":        sentity.ProductOrderByInventory,
		"price":            sentity.ProductOrderByPrice,
		"originPrefecture": sentity.ProductOrderByOriginPrefecture,
		"originCity":       sentity.ProductOrderByOriginCity,
		"createdAt":        sentity.ProductOrderByCreatedAt,
		"updatedAt":        sentity.ProductOrderByUpdatedAt,
	}
	params := util.GetOrders(ctx)
	res := make([]*store.ListProductsOrder, len(params))
	for i, p := range params {
		key, ok := products[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderkey)
		}
		res[i] = &store.ListProductsOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

func (h *handler) GetProduct(ctx *gin.Context) {
	product, err := h.getProduct(ctx, util.GetParam(ctx, "productId"))
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ProductResponse{
		Product: product.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateProduct(ctx *gin.Context) {
	req := &request.CreateProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	if getRole(ctx).IsCoordinator() {
		producers, err := h.getProducersByCoordinatorID(ctx, getAdminID(ctx))
		if err != nil {
			httpError(ctx, err)
		}
		if !producers.Contains(req.ProducerID) {
			forbidden(ctx, errors.New("handler: not authorized this coordinator"))
			return
		}
	}

	var (
		producer    *service.Producer
		productType *service.ProductType
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producer, err = h.getProducer(ectx, req.ProducerID)
		return
	})
	eg.Go(func() (err error) {
		productType, err = h.getProductType(ctx, req.TypeID)
		return
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		badRequest(ctx, err)
		return
	}
	if err != nil {
		httpError(ctx, err)
		return
	}

	eg, ectx = errgroup.WithContext(ctx)
	productMedia := make([]*store.CreateProductMedia, len(req.Media))
	for i := range req.Media {
		i := i
		eg.Go(func() error {
			in := &media.UploadFileInput{
				URL: req.Media[i].URL,
			}
			url, err := h.media.UploadProductMedia(ectx, in)
			if err != nil {
				return err
			}
			productMedia[i] = &store.CreateProductMedia{
				URL:         url,
				IsThumbnail: req.Media[i].IsThumbnail,
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)
	in := &store.CreateProductInput{
		ProducerID:       req.ProducerID,
		TypeID:           req.TypeID,
		Name:             req.Name,
		Description:      req.Description,
		Public:           req.Public,
		Inventory:        req.Inventory,
		Weight:           weight,
		WeightUnit:       weightUnit,
		Item:             1, // 1固定
		ItemUnit:         req.ItemUnit,
		ItemDescription:  req.ItemDescription,
		Media:            productMedia,
		Price:            req.Price,
		DeliveryType:     service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:        req.Box60Rate,
		Box80Rate:        req.Box80Rate,
		Box100Rate:       req.Box100Rate,
		OriginPrefecture: req.OriginPrefecture,
		OriginCity:       req.OriginCity,
	}
	sproduct, err := h.store.CreateProduct(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	product := service.NewProduct(sproduct)

	product.Fill(productType, producer)

	res := &response.ProductResponse{
		Product: product.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateProduct(ctx *gin.Context) {
	req := &request.UpdateProductRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	eg, ectx := errgroup.WithContext(ctx)
	productMedia := make([]*store.UpdateProductMedia, len(req.Media))
	for i := range req.Media {
		i := i
		eg.Go(func() error {
			in := &media.UploadFileInput{
				URL: req.Media[i].URL,
			}
			url, err := h.media.UploadProductMedia(ectx, in)
			if err != nil {
				return err
			}
			productMedia[i] = &store.UpdateProductMedia{
				URL:         url,
				IsThumbnail: req.Media[i].IsThumbnail,
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	weight, weightUnit := service.NewProductWeightFromRequest(req.Weight)
	in := &store.UpdateProductInput{
		ProductID:        util.GetParam(ctx, "productId"),
		ProducerID:       req.ProducerID,
		TypeID:           req.TypeID,
		Name:             req.Name,
		Description:      req.Description,
		Public:           req.Public,
		Inventory:        req.Inventory,
		Weight:           weight,
		WeightUnit:       weightUnit,
		Item:             1, // 1固定
		ItemUnit:         req.ItemUnit,
		ItemDescription:  req.ItemDescription,
		Media:            productMedia,
		Price:            req.Price,
		DeliveryType:     service.DeliveryType(req.DeliveryType).StoreEntity(),
		Box60Rate:        req.Box60Rate,
		Box80Rate:        req.Box80Rate,
		Box100Rate:       req.Box100Rate,
		OriginPrefecture: req.OriginPrefecture,
		OriginCity:       req.OriginCity,
	}
	if err := h.store.UpdateProduct(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteProduct(ctx *gin.Context) {
	in := &store.DeleteProductInput{
		ProductID: util.GetParam(ctx, "productId"),
	}
	if err := h.store.DeleteProduct(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) multiGetProducts(ctx context.Context, productIDs []string) (service.Products, error) {
	in := &store.MultiGetProductsInput{
		ProductIDs: productIDs,
	}
	sproducts, err := h.store.MultiGetProducts(ctx, in)
	if err != nil {
		return nil, err
	}
	products := service.NewProducts(sproducts)
	if len(products) == 0 {
		return products, nil
	}
	if err := h.getProductsDetails(ctx, products...); err != nil {
		return nil, err
	}
	return products, nil
}

func (h *handler) getProduct(ctx context.Context, productID string) (*service.Product, error) {
	in := &store.GetProductInput{
		ProductID: productID,
	}
	sproduct, err := h.store.GetProduct(ctx, in)
	if err != nil {
		return nil, err
	}
	product := service.NewProduct(sproduct)
	if err := h.getProductsDetails(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (h *handler) getProductsDetails(ctx context.Context, products ...*service.Product) error {
	ps := service.Products(products)
	var (
		producers service.Producers
		types     service.ProductTypes
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		producers, err = h.multiGetProducers(ectx, ps.ProducerIDs())
		return
	})
	eg.Go(func() (err error) {
		types, err = h.multiGetProductTypes(ectx, ps.ProductTypeIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	ps.Fill(types.Map(), producers.Map())
	return nil
}
