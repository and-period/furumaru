//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"
	"time"

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
	Category       Category
	Experience     Experience
	ExperienceType ExperienceType
	Live           Live
	Order          Order
	PaymentSystem  PaymentSystem
	Product        Product
	ProductTag     ProductTag
	ProductType    ProductType
	Promotion      Promotion
	Schedule       Schedule
	Shipping       Shipping
	Spot           Spot
}

/**
 * interface
 */
type Category interface {
	List(ctx context.Context, params *ListCategoriesParams, fields ...string) (entity.Categories, error)
	Count(ctx context.Context, params *ListCategoriesParams) (int64, error)
	MultiGet(ctx context.Context, categoryIDs []string, fields ...string) (entity.Categories, error)
	Get(ctx context.Context, categoryID string, fields ...string) (*entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, categoryID, name string) error
	Delete(ctx context.Context, categoryID string) error
}

type ListCategoriesOrderKey string

const (
	ListCategoriesOrderByName ListCategoriesOrderKey = "name"
)

type ListCategoriesParams struct {
	Name   string
	Limit  int
	Offset int
	Orders []*ListCategoriesOrder
}

type ListCategoriesOrder struct {
	Key        ListCategoriesOrderKey
	OrderByASC bool
}

type Experience interface {
	List(ctx context.Context, params *ListExperiencesParams, fields ...string) (entity.Experiences, error)
	Count(ctx context.Context, params *ListExperiencesParams) (int64, error)
	MultiGet(ctx context.Context, experienceIDs []string, fields ...string) (entity.Experiences, error)
	MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Experiences, error)
	Get(ctx context.Context, experienceID string, fields ...string) (*entity.Experience, error)
	Create(ctx context.Context, experience *entity.Experience) error
	Update(ctx context.Context, experienceID string, params *UpdateExperienceParams) error
	Delete(ctx context.Context, experienceID string) error
}

type ListExperiencesParams struct {
	Name           string
	HostPrefecture int32
	CoordinatorID  string
	ProducerID     string
	OnlyPublished  bool
	ExcludeDeleted bool
	EndAtGte       time.Time
	Limit          int
	Offset         int
}

type UpdateExperienceParams struct {
	TypeID                string
	Title                 string
	Description           string
	Public                bool
	SoldOut               bool
	Media                 entity.MultiExperienceMedia
	PriceAdult            int64
	PriceJuniorHighSchool int64
	PriceElementarySchool int64
	PricePreschool        int64
	PriceSenior           int64
	RecommendedPoints     []string
	PromotionVideoURL     string
	Duration              int64
	Direction             string
	BusinessOpenTime      string
	BusinessCloseTime     string
	HostPostalCode        string
	HostPrefectureCode    int32
	HostCity              string
	HostAddressLine1      string
	HostAddressLine2      string
	HostLongitude         float64
	HostLatitude          float64
	StartAt               time.Time
	EndAt                 time.Time
}

type ExperienceType interface {
	List(ctx context.Context, params *ListExperienceTypesParams, fields ...string) (entity.ExperienceTypes, error)
	Count(ctx context.Context, params *ListExperienceTypesParams) (int64, error)
	MultiGet(ctx context.Context, experienceTypeIDs []string, fields ...string) (entity.ExperienceTypes, error)
	Get(ctx context.Context, experienceTypeID string, fields ...string) (*entity.ExperienceType, error)
	Create(ctx context.Context, experienceType *entity.ExperienceType) error
	Update(ctx context.Context, experienceTypeID string, params *UpdateExperienceTypeParams) error
	Delete(ctx context.Context, experienceTypeID string) error
}

type ListExperienceTypesParams struct {
	Name   string
	Limit  int
	Offset int
}

type UpdateExperienceTypeParams struct {
	Name string
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
	ProducerID  string
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
	ListUserIDs(ctx context.Context, params *ListOrdersParams) ([]string, int64, error)
	Count(ctx context.Context, params *ListOrdersParams) (int64, error)
	Get(ctx context.Context, orderID string, fields ...string) (*entity.Order, error)
	GetByTransactionID(ctx context.Context, userID, transactionID string) (*entity.Order, error)
	GetByTransactionIDWithSessionID(ctx context.Context, sessionID, transactionID string) (*entity.Order, error)
	Create(ctx context.Context, order *entity.Order) error
	UpdatePayment(ctx context.Context, orderID string, params *UpdateOrderPaymentParams) error
	UpdateFulfillment(ctx context.Context, orderID, fulfillmentID string, params *UpdateOrderFulfillmentParams) error
	UpdateRefund(ctx context.Context, orderID string, params *UpdateOrderRefundParams) error
	Draft(ctx context.Context, orderID string, params *DraftOrderParams) error
	Complete(ctx context.Context, orderID string, params *CompleteOrderParams) error
	Aggregate(ctx context.Context, params *AggregateOrdersParams) (entity.AggregatedOrders, error)
	AggregateByPromotion(ctx context.Context, params *AggregateOrdersByPromotionParams) (entity.AggregatedOrderPromotions, error)
}

type ListOrdersParams struct {
	CoordinatorID string
	UserID        string
	Types         []entity.OrderType
	Statuses      []entity.OrderStatus
	Limit         int
	Offset        int
}

type UpdateOrderPaymentParams struct {
	Status    entity.PaymentStatus
	PaymentID string
	IssuedAt  time.Time
}

type UpdateOrderFulfillmentParams struct {
	Status          entity.FulfillmentStatus
	ShippingCarrier entity.ShippingCarrier
	TrackingNumber  string
	ShippedAt       time.Time
}

type UpdateOrderRefundParams struct {
	Status       entity.PaymentStatus
	RefundType   entity.RefundType
	RefundTotal  int64
	RefundReason string
	IssuedAt     time.Time
}

type DraftOrderParams struct {
	ShippingMessage string
}

type CompleteOrderParams struct {
	ShippingMessage string
	CompletedAt     time.Time
}

type AggregateOrdersParams struct {
	CoordinatorID string
	UserIDs       []string
}

type AggregateOrdersByPromotionParams struct {
	CoordinatorID string
	PromotionIDs  []string
}

type PaymentSystem interface {
	MultiGet(ctx context.Context, methodTypes []entity.PaymentMethodType, fields ...string) (entity.PaymentSystems, error)
	Get(ctx context.Context, methodType entity.PaymentMethodType, fields ...string) (*entity.PaymentSystem, error)
	Update(ctx context.Context, methodType entity.PaymentMethodType, status entity.PaymentSystemStatus) error
}

type Product interface {
	List(ctx context.Context, params *ListProductsParams, fields ...string) (entity.Products, error)
	Count(ctx context.Context, params *ListProductsParams) (int64, error)
	MultiGet(ctx context.Context, productIDs []string, fields ...string) (entity.Products, error)
	MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Products, error)
	Get(ctx context.Context, productID string, fields ...string) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, productID string, params *UpdateProductParams) error
	DecreaseInventory(ctx context.Context, revisionID, quantity int64) error
	Delete(ctx context.Context, productID string) error
}

type ListProductsOrderKey string

const (
	ListProductsOrderByName             ListProductsOrderKey = "name"
	ListProductsOrderBySoldOut          ListProductsOrderKey = "CASE WHEN (inventory = 0) THEN 1 ELSE 0 END"
	ListProductsOrderByPublic           ListProductsOrderKey = "public"
	ListProductsOrderByInventory        ListProductsOrderKey = "inventory"
	ListProductsOrderByOriginPrefecture ListProductsOrderKey = "origin_prefecture"
	ListProductsOrderByOriginCity       ListProductsOrderKey = "origin_city"
	ListProductsOrderByStartAt          ListProductsOrderKey = "start_at"
	ListProductsOrderByCreatedAt        ListProductsOrderKey = "created_at"
	ListProductsOrderByUpdatedAt        ListProductsOrderKey = "updated_at"
)

type ListProductsParams struct {
	Name           string
	CoordinatorID  string
	ProducerID     string
	ProducerIDs    []string
	ProductTypeIDs []string
	ProductTagID   string
	OnlyPublished  bool
	ExcludeDeleted bool
	EndAtGte       time.Time
	Limit          int
	Offset         int
	Orders         []*ListProductsOrder
}

type ListProductsOrder struct {
	Key        ListProductsOrderKey
	OrderByASC bool
}

type UpdateProductParams struct {
	TypeID               string
	TagIDs               []string
	Name                 string
	Description          string
	Public               bool
	Inventory            int64
	Weight               int64
	WeightUnit           entity.WeightUnit
	Item                 int64
	ItemUnit             string
	ItemDescription      string
	Media                entity.MultiProductMedia
	Price                int64
	Cost                 int64
	ExpirationDate       int64
	RecommendedPoints    []string
	StorageMethodType    entity.StorageMethodType
	DeliveryType         entity.DeliveryType
	Box60Rate            int64
	Box80Rate            int64
	Box100Rate           int64
	OriginPrefectureCode int32
	OriginCity           string
	StartAt              time.Time
	EndAt                time.Time
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

type ListProductTagsOrderKey string

const (
	ListProductTagsOrderByName ListProductTagsOrderKey = "name"
)

type ListProductTagsParams struct {
	Name   string
	Limit  int
	Offset int
	Orders []*ListProductTagsOrder
}

type ListProductTagsOrder struct {
	Key        ListProductTagsOrderKey
	OrderByASC bool
}

type ProductType interface {
	List(ctx context.Context, params *ListProductTypesParams, fields ...string) (entity.ProductTypes, error)
	Count(ctx context.Context, params *ListProductTypesParams) (int64, error)
	MultiGet(ctx context.Context, productTypeIDs []string, fields ...string) (entity.ProductTypes, error)
	Get(ctx context.Context, productTypeID string, fields ...string) (*entity.ProductType, error)
	Create(ctx context.Context, productType *entity.ProductType) error
	Update(ctx context.Context, productTypeID, name, iconURL string) error
	Delete(ctx context.Context, productTypeID string) error
}

type ListProductTypesOrderKey string

const (
	ListProductTypesOrderByName ListProductTypesOrderKey = "name"
)

type ListProductTypesParams struct {
	Name       string
	CategoryID string
	Limit      int
	Offset     int
	Orders     []*ListProductTypesOrder
}

type ListProductTypesOrder struct {
	Key        ListProductTypesOrderKey
	OrderByASC bool
}

type Promotion interface {
	List(ctx context.Context, params *ListPromotionsParams, fields ...string) (entity.Promotions, error)
	Count(ctx context.Context, params *ListPromotionsParams) (int64, error)
	MultiGet(ctx context.Context, promotionIDs []string, fields ...string) (entity.Promotions, error)
	Get(ctx context.Context, promotionID string, fields ...string) (*entity.Promotion, error)
	GetByCode(ctx context.Context, promotionID string, fields ...string) (*entity.Promotion, error)
	Create(ctx context.Context, promotion *entity.Promotion) error
	Update(ctx context.Context, promotionID string, params *UpdatePromotionParams) error
	Delete(ctx context.Context, promotionID string) error
}

type ListPromotionsOrderKey string

const (
	ListPromotionsOrderByTitle       ListPromotionsOrderKey = "title"
	ListPromotionsOrderByPublic      ListPromotionsOrderKey = "public"
	ListPromotionsOrderByPublishedAt ListPromotionsOrderKey = "published_at"
	ListPromotionsOrderByStartAt     ListPromotionsOrderKey = "start_at"
	ListPromotionsOrderByEndAt       ListPromotionsOrderKey = "end_at"
	ListPromotionsOrderByCreatedAt   ListPromotionsOrderKey = "created_at"
	ListPromotionsOrderByUpdatedAt   ListPromotionsOrderKey = "updated_at"
)

type ListPromotionsParams struct {
	Title  string
	Limit  int
	Offset int
	Orders []*ListPromotionsOrder
}

type ListPromotionsOrder struct {
	Key        ListPromotionsOrderKey
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
	Delete(ctx context.Context, scheduleID string) error
	Approve(ctx context.Context, scheduleID string, params *ApproveScheduleParams) error
	Publish(ctx context.Context, scheduleID string, public bool) error
}

type ListSchedulesParams struct {
	CoordinatorID string
	ProducerID    string
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
	StartAt         time.Time
	EndAt           time.Time
}

type ApproveScheduleParams struct {
	Approved        bool
	ApprovedAdminID string
}

type Shipping interface {
	ListByCoordinatorIDs(ctx context.Context, coordinatorIDs []string, fields ...string) (entity.Shippings, error)
	MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Shippings, error)
	GetDefault(ctx context.Context, fields ...string) (*entity.Shipping, error)
	GetByCoordinatorID(ctx context.Context, coordinatorID string, fields ...string) (*entity.Shipping, error)
	Create(ctx context.Context, shipping *entity.Shipping) error
	Update(ctx context.Context, shippingID string, params *UpdateShippingParams) error
}

type UpdateShippingParams struct {
	Box60Rates        entity.ShippingRates
	Box60Frozen       int64
	Box80Rates        entity.ShippingRates
	Box80Frozen       int64
	Box100Rates       entity.ShippingRates
	Box100Frozen      int64
	HasFreeShipping   bool
	FreeShippingRates int64
}

type Spot interface {
	List(ctx context.Context, params *ListSpotsParams, fields ...string) (entity.Spots, error)
	ListByGeolocation(ctx context.Context, params *ListSpotsByGeolocationParams, fields ...string) (entity.Spots, error)
	Count(ctx context.Context, params *ListSpotsParams) (int64, error)
	Get(ctx context.Context, spotID string, fields ...string) (*entity.Spot, error)
	Create(ctx context.Context, spot *entity.Spot) error
	Update(ctx context.Context, spotID string, params *UpdateSpotParams) error
	Delete(ctx context.Context, spotID string) error
	Approve(ctx context.Context, spotID string, params *ApproveSpotParams) error
}

type ListSpotsParams struct {
	Name            string
	UserID          string
	ExcludeApproved bool
	ExcludeDisabled bool
	Limit           int
	Offset          int
}

type ListSpotsByGeolocationParams struct {
	Longitude float64
	Latitude  float64
	Radius    int
}

type UpdateSpotParams struct {
	Name         string
	Description  string
	ThumbnailURL string
	Longitude    float64
	Latitude     float64
}

type ApproveSpotParams struct {
	Approved        bool
	ApprovedAdminID string
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
