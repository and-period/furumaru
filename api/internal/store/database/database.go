//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"gorm.io/gorm"
)

type Params struct {
	Database *database.Client
	DynamoDB dynamodb.Client
}

type Database struct {
	Address     Address
	Category    Category
	Order       Order
	Product     Product
	ProductType ProductType
	Promotion   Promotion
	Rehearsal   Rehearsal
	Shipping    Shipping
	Schedule    Schedule
	Live        Live
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Address:     NewAddress(params.Database),
		Category:    NewCategory(params.Database),
		Live:        NewLive(params.Database),
		Order:       NewOrder(params.Database),
		Product:     NewProduct(params.Database),
		ProductType: NewProductType(params.Database),
		Promotion:   NewPromotion(params.Database),
		Rehearsal:   NewRehearsal(params.DynamoDB),
		Shipping:    NewShipping(params.Database),
		Schedule:    NewSchedule(params.Database),
	}
}

/**
 * interface
 */
type Address interface {
	MultiGet(ctx context.Context, addressIDs []string, fields ...string) (entity.Addresses, error)
}

type Category interface {
	List(ctx context.Context, params *ListCategoriesParams, fields ...string) (entity.Categories, error)
	Count(ctx context.Context, params *ListCategoriesParams) (int64, error)
	MultiGet(ctx context.Context, categoryIDs []string, fields ...string) (entity.Categories, error)
	Get(ctx context.Context, categoryID string, fields ...string) (*entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, categoryID, name string) error
	Delete(ctx context.Context, categoryID string) error
}

type Order interface {
	List(ctx context.Context, params *ListOrdersParams, fields ...string) (entity.Orders, error)
	Count(ctx context.Context, params *ListOrdersParams) (int64, error)
	Get(ctx context.Context, orderID string, fields ...string) (*entity.Order, error)
	Aggregate(ctx context.Context, userIDs []string) (entity.AggregatedOrders, error)
}

type Product interface {
	List(ctx context.Context, params *ListProductsParams, fields ...string) (entity.Products, error)
	Count(ctx context.Context, params *ListProductsParams) (int64, error)
	MultiGet(ctx context.Context, productIDs []string, fields ...string) (entity.Products, error)
	Get(ctx context.Context, productID string, fields ...string) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, productID string, params *UpdateProductParams) error
	UpdateMedia(ctx context.Context, productID string, set func(media entity.MultiProductMedia) bool) error
	Delete(ctx context.Context, productID string) error
}

type ProductType interface {
	List(ctx context.Context, params *ListProductTypesParams, fields ...string) (entity.ProductTypes, error)
	Count(ctx context.Context, params *ListProductTypesParams) (int64, error)
	MultiGet(ctx context.Context, productTypeIDs []string, fields ...string) (entity.ProductTypes, error)
	Get(ctx context.Context, productTypeID string, fields ...string) (*entity.ProductType, error)
	Create(ctx context.Context, productType *entity.ProductType) error
	Update(ctx context.Context, productTypeID, name, iconURL string) error
	UpdateIcons(ctx context.Context, productTypeID string, icons common.Images) error
	Delete(ctx context.Context, productTypeID string) error
}

type Promotion interface {
	List(ctx context.Context, params *ListPromotionsParams, fields ...string) (entity.Promotions, error)
	Count(ctx context.Context, params *ListPromotionsParams) (int64, error)
	Get(ctx context.Context, promotionID string, fields ...string) (*entity.Promotion, error)
	Create(ctx context.Context, promotion *entity.Promotion) error
	Update(ctx context.Context, promotionID string, params *UpdatePromotionParams) error
	Delete(ctx context.Context, promotionID string) error
}

type Rehearsal interface {
	Get(ctx context.Context, liveID string) (*entity.Rehearsal, error)
	Create(ctx context.Context, rehearsal *entity.Rehearsal) error
}

type Shipping interface {
	List(ctx context.Context, params *ListShippingsParams, fields ...string) (entity.Shippings, error)
	Count(ctx context.Context, params *ListShippingsParams) (int64, error)
	Get(ctx context.Context, shoppingID string, fields ...string) (*entity.Shipping, error)
	MultiGet(ctx context.Context, shippingIDs []string, fields ...string) (entity.Shippings, error)
	Create(ctx context.Context, shipping *entity.Shipping) error
	Update(ctx context.Context, shippingID string, params *UpdateShippingParams) error
	Delete(ctx context.Context, shippingID string) error
}

type Schedule interface {
	Create(ctx context.Context, schedule *entity.Schedule, lives entity.Lives, products entity.LiveProducts) error
}

type Live interface {
	Get(ctx context.Context, liveID string, fields ...string) (*entity.Live, error)
	Update(ctx context.Context, liveID string, params *UpdateLiveParams) error
	UpdateLivePublic(ctx context.Context, liveID string, params *UpdateLivePublicParams) error
}

/**
 * params
 */
type ListCategoriesParams struct {
	Name   string
	Limit  int
	Offset int
	Orders []*ListCategoriesOrder
}

type ListCategoriesOrder struct {
	Key        entity.CategoryOrderBy
	OrderByASC bool
}

func (p *ListCategoriesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

type UpdateLiveParams struct {
	LiveProducts entity.LiveProducts
	ProducerID   string
	Title        string
	Description  string
	StartAt      time.Time
	EndAt        time.Time
}

type UpdateLivePublicParams struct {
	Published    bool
	Canceled     bool
	ChannelArn   string
	StreamKeyArn string
}

type ListOrdersParams struct {
	CoordinatorID string
	Limit         int
	Offset        int
	Orders        []*ListOrdersOrder
}

type ListOrdersOrder struct {
	Key        entity.OrderOrderBy
	OrderByASC bool
}

func (p *ListOrdersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

type ListProductsParams struct {
	Name        string
	ProducerID  string
	ProducerIDs []string
	Limit       int
	Offset      int
	Orders      []*ListProductsOrder
}

type ListProductsOrder struct {
	Key        entity.ProductOrderBy
	OrderByASC bool
}

func (p *ListProductsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.ProducerID != "" {
		stmt = stmt.Where("producer_id = ?", p.ProducerID)
	}
	if len(p.ProducerIDs) > 0 {
		stmt = stmt.Where("producer_id IN (?)", p.ProducerIDs)
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

type UpdateProductParams struct {
	ProducerID       string
	TypeID           string
	Name             string
	Description      string
	Public           bool
	Inventory        int64
	Weight           int64
	WeightUnit       entity.WeightUnit
	Item             int64
	ItemUnit         string
	ItemDescription  string
	Media            entity.MultiProductMedia
	Price            int64
	DeliveryType     entity.DeliveryType
	Box60Rate        int64
	Box80Rate        int64
	Box100Rate       int64
	OriginPrefecture string
	OriginCity       string
}

type ListProductTypesParams struct {
	Name       string
	CategoryID string
	Limit      int
	Offset     int
	Orders     []*ListProductTypesOrder
}

type ListProductTypesOrder struct {
	Key        entity.ProductTypeOrderBy
	OrderByASC bool
}

func (p *ListProductTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.CategoryID != "" {
		stmt = stmt.Where("category_id = ?", p.CategoryID)
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

type ListPromotionsParams struct {
	Limit  int
	Offset int
	Orders []*ListPromotionsOrder
}

type ListPromotionsOrder struct {
	Key        entity.PromotionOrderBy
	OrderByASC bool
}

func (p *ListPromotionsParams) stmt(stmt *gorm.DB) *gorm.DB {
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

type UpdatePromotionParams struct {
	Title        string
	Description  string
	Public       bool
	PublishedAt  time.Time
	DiscountType entity.DiscountType
	DiscountRate int64
	Code         string
	CodeType     entity.PromotionCodeType
	StartAt      time.Time
	EndAt        time.Time
}

type ListShippingsParams struct {
	Limit  int
	Offset int
	Orders []*ListShippingsOrder
}

type ListShippingsOrder struct {
	Key        entity.ShippingOrderBy
	OrderByASC bool
}

func (p *ListShippingsParams) stmt(stmt *gorm.DB) *gorm.DB {
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

type UpdateShippingParams struct {
	Name               string
	Box60Rates         entity.ShippingRates
	Box60Refrigerated  int64
	Box60Frozen        int64
	Box80Rates         entity.ShippingRates
	Box80Refrigerated  int64
	Box80Frozen        int64
	Box100Rates        entity.ShippingRates
	Box100Refrigerated int64
	Box100Frozen       int64
	HasFreeShipping    bool
	FreeShippingRates  int64
}
