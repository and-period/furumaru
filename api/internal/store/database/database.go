//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

var (
	ErrInvalidArgument    = &Error{err: errors.New("database: invalid argument")}
	ErrNotFound           = &Error{err: errors.New("database: not found")}
	ErrAlreadyExists      = &Error{err: errors.New("database: already exists")}
	ErrFailedPrecondition = &Error{err: errors.New("database: failed precondition")}
	ErrCanceled           = &Error{err: errors.New("database: canceled")}
	ErrDeadlineExceeded   = &Error{err: errors.New("database: deadline exceeded")}
	ErrInternal           = &Error{err: errors.New("database: internal error")}
	ErrUnknown            = &Error{err: errors.New("database: unknown")}
)

type Database struct {
	Address     Address
	Category    Category
	Order       Order
	Product     Product
	ProductTag  ProductTag
	ProductType ProductType
	Promotion   Promotion
	Shipping    Shipping
	Schedule    Schedule
	Live        Live
}

/**
 * interface
 */
type Address interface {
	List(ctx context.Context, params *ListAddressesParams, fields ...string) (entity.Addresses, error)
	Count(ctx context.Context, params *ListAddressesParams) (int64, error)
	MultiGet(ctx context.Context, addressIDs []string, fields ...string) (entity.Addresses, error)
	Get(ctx context.Context, addressID string, fields ...string) (*entity.Address, error)
	Create(ctx context.Context, address *entity.Address) error
	Update(ctx context.Context, addressID, userID string, params *UpdateAddressParams) error
	Delete(ctx context.Context, addressID string) error
}

type ListAddressesParams struct {
	UserID string
	Limit  int
	Offset int
}

type UpdateAddressParams struct {
	Lastname     string
	Firstname    string
	PostalCode   string
	Prefecture   int64
	City         string
	AddressLine1 string
	AddressLine2 string
	PhoneNumber  string
	IsDefault    bool
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

type Live interface {
	List(ctx context.Context, params *ListLivesParams, fields ...string) (entity.Lives, error)
	Count(ctx context.Context, params *ListLivesParams) (int64, error)
	Get(ctx context.Context, liveID string, fields ...string) (*entity.Live, error)
	Create(ctx context.Context, live *entity.Live) error
	Update(ctx context.Context, liveID string, params *UpdateLiveParams) error
	Delete(ctx context.Context, liveID string) error
}

type ListLivesParams struct {
	ScheduleIDs []string
	Limit       int
	Offset      int
}

type UpdateLiveParams struct {
	ProductIDs []string
	Comment    string
	StartAt    time.Time
	EndAt      time.Time
}

type Order interface {
	List(ctx context.Context, params *ListOrdersParams, fields ...string) (entity.Orders, error)
	Count(ctx context.Context, params *ListOrdersParams) (int64, error)
	Get(ctx context.Context, orderID string, fields ...string) (*entity.Order, error)
	Aggregate(ctx context.Context, userIDs []string) (entity.AggregatedOrders, error)
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

type ListProductsParams struct {
	Name          string
	CoordinatorID string
	ProducerID    string
	ProducerIDs   []string
	OnlyPublished bool
	Limit         int
	Offset        int
	Orders        []*ListProductsOrder
}

type ListProductsOrder struct {
	Key        entity.ProductOrderBy
	OrderByASC bool
}

type UpdateProductParams struct {
	TypeID            string
	TagIDs            []string
	Name              string
	Description       string
	Public            bool
	Inventory         int64
	Weight            int64
	WeightUnit        entity.WeightUnit
	Item              int64
	ItemUnit          string
	ItemDescription   string
	Media             entity.MultiProductMedia
	Price             int64
	Cost              int64
	ExpirationDate    int64
	RecommendedPoints []string
	StorageMethodType entity.StorageMethodType
	DeliveryType      entity.DeliveryType
	Box60Rate         int64
	Box80Rate         int64
	Box100Rate        int64
	OriginPrefecture  int64
	OriginCity        string
	StartAt           time.Time
	EndAt             time.Time
}

type ProductTag interface {
	List(ctx context.Context, params *ListProductTagsParams, fields ...string) (entity.ProductTags, error)
	Count(ctx context.Context, params *ListProductTagsParams) (int64, error)
	MultiGet(ctx context.Context, productTagIDs []string, fields ...string) (entity.ProductTags, error)
	Get(ctx context.Context, productTagID string, fields ...string) (*entity.ProductTag, error)
	Create(ctx context.Context, category *entity.ProductTag) error
	Update(ctx context.Context, productTagID, name string) error
	Delete(ctx context.Context, productTagID string) error
}

type ListProductTagsParams struct {
	Name   string
	Limit  int
	Offset int
	Orders []*ListProductTagsOrder
}

type ListProductTagsOrder struct {
	Key        entity.ProductTagOrderBy
	OrderByASC bool
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

type Promotion interface {
	List(ctx context.Context, params *ListPromotionsParams, fields ...string) (entity.Promotions, error)
	Count(ctx context.Context, params *ListPromotionsParams) (int64, error)
	MultiGet(ctx context.Context, promotionIDs []string, fields ...string) (entity.Promotions, error)
	Get(ctx context.Context, promotionID string, fields ...string) (*entity.Promotion, error)
	Create(ctx context.Context, promotion *entity.Promotion) error
	Update(ctx context.Context, promotionID string, params *UpdatePromotionParams) error
	Delete(ctx context.Context, promotionID string) error
}

type ListPromotionsParams struct {
	Title  string
	Limit  int
	Offset int
	Orders []*ListPromotionsOrder
}

type ListPromotionsOrder struct {
	Key        entity.PromotionOrderBy
	OrderByASC bool
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

type Schedule interface {
	List(ctx context.Context, params *ListSchedulesParams, fields ...string) (entity.Schedules, error)
	Count(ctx context.Context, params *ListSchedulesParams) (int64, error)
	MultiGet(ctx context.Context, scheduleIDs []string, fields ...string) (entity.Schedules, error)
	Get(ctx context.Context, scheduleID string, fields ...string) (*entity.Schedule, error)
	Create(ctx context.Context, schedule *entity.Schedule) error
	Update(ctx context.Context, scheduleID string, params *UpdateScheduleParams) error
	UpdateThumbnails(ctx context.Context, scheduleID string, thumbnails common.Images) error
	Approve(ctx context.Context, scheduleID string, params *ApproveScheduleParams) error
}

type ListSchedulesParams struct {
	StartAtGte    time.Time
	StartAtLt     time.Time
	EndAtGte      time.Time
	EndAtLt       time.Time
	OnlyPublished bool
	Limit         int
	Offset        int
}

type UpdateScheduleParams struct {
	Title           string
	Description     string
	ThumbnailURL    string
	ImageURL        string
	OpeningVideoURL string
	Public          bool
	StartAt         time.Time
	EndAt           time.Time
}

type ApproveScheduleParams struct {
	Approved        bool
	ApprovedAdminID string
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

type ListShippingsParams struct {
	CoordinatorID string
	Name          string
	Limit         int
	Offset        int
	Orders        []*ListShippingsOrder
}

type ListShippingsOrder struct {
	Key        entity.ShippingOrderBy
	OrderByASC bool
}

type UpdateShippingParams struct {
	Name               string
	IsDefault          bool
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

type Error struct {
	err error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}
