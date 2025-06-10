//go:generate go tool mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
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
	ErrPermissionDenied   = &Error{err: errors.New("database: permission denied")}
	ErrCanceled           = &Error{err: errors.New("database: canceled")}
	ErrDeadlineExceeded   = &Error{err: errors.New("database: deadline exceeded")}
	ErrInternal           = &Error{err: errors.New("database: internal error")}
	ErrUnknown            = &Error{err: errors.New("database: unknown")}
)

type Database struct {
	CartActionLog            CartActionLog
	Category                 Category
	Experience               Experience
	ExperienceReview         ExperienceReview
	ExperienceReviewReaction ExperienceReviewReaction
	ExperienceType           ExperienceType
	Live                     Live
	Order                    Order
	PaymentSystem            PaymentSystem
	Product                  Product
	ProductReview            ProductReview
	ProductReviewReaction    ProductReviewReaction
	ProductTag               ProductTag
	ProductType              ProductType
	Promotion                Promotion
	Schedule                 Schedule
	Shipping                 Shipping
	Shop                     Shop
	Spot                     Spot
	SpotType                 SpotType
}

/**
 * interface
 */
type CartActionLog interface {
	Create(ctx context.Context, log *entity.CartActionLog) error
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
	ListByGeolocation(ctx context.Context, params *ListExperiencesByGeolocationParams, fields ...string) (entity.Experiences, error)
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
	ShopID         string
	CoordinatorID  string // Deprecated
	ProducerID     string
	OnlyPublished  bool
	ExcludeDeleted bool
	EndAtGte       time.Time
	Limit          int
	Offset         int
}

type ListExperiencesByGeolocationParams struct {
	ShopID         string
	CoordinatorID  string // Deprecated
	ProducerID     string
	Longitude      float64
	Latitude       float64
	Radius         int64
	OnlyPublished  bool
	ExcludeDeleted bool
	EndAtGte       time.Time
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

type ExperienceReview interface {
	List(ctx context.Context, params *ListExperienceReviewsParams, fields ...string) (entity.ExperienceReviews, string, error)
	Get(ctx context.Context, experienceReviewID string, fields ...string) (*entity.ExperienceReview, error)
	Create(ctx context.Context, experienceReview *entity.ExperienceReview) error
	Update(ctx context.Context, experienceReviewID string, params *UpdateExperienceReviewParams) error
	Delete(ctx context.Context, experienceReviewID string) error
	Aggregate(ctx context.Context, params *AggregateExperienceReviewsParams) (entity.AggregatedExperienceReviews, error)
}

type ListExperienceReviewsParams struct {
	ExperienceID string
	UserID       string
	Rates        []int64
	Limit        int64
	NextToken    string
}

type UpdateExperienceReviewParams struct {
	Rate    int64
	Title   string
	Comment string
}

type AggregateExperienceReviewsParams struct {
	ExperienceIDs []string
}

type ExperienceReviewReaction interface {
	Upsert(ctx context.Context, reaction *entity.ExperienceReviewReaction) error
	Delete(ctx context.Context, experienceReviewID, userID string) error
	GetUserReactions(ctx context.Context, experienceID, userID string) (entity.ExperienceReviewReactions, error)
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
	UpdateAuthorized(ctx context.Context, orderID string, params *UpdateOrderAuthorizedParams) error
	UpdateCaptured(ctx context.Context, orderID string, params *UpdateOrderCapturedParams) error
	UpdateFailed(ctx context.Context, orderID string, params *UpdateOrderFailedParams) error
	UpdateRefunded(ctx context.Context, orderID string, params *UpdateOrderRefundedParams) error
	UpdateFulfillment(ctx context.Context, orderID, fulfillmentID string, params *UpdateOrderFulfillmentParams) error
	Draft(ctx context.Context, orderID string, params *DraftOrderParams) error
	Complete(ctx context.Context, orderID string, params *CompleteOrderParams) error
	Aggregate(ctx context.Context, params *AggregateOrdersParams) (*entity.AggregatedOrder, error)
	AggregateByUser(ctx context.Context, params *AggregateOrdersByUserParams) (entity.AggregatedUserOrders, error)
	AggregateByPaymentMethodType(ctx context.Context, params *AggregateOrdersByPaymentMethodTypeParams) (entity.AggregatedOrderPayments, error)
	AggregateByPromotion(ctx context.Context, params *AggregateOrdersByPromotionParams) (entity.AggregatedOrderPromotions, error)
	AggregateByPeriod(ctx context.Context, params *AggregateOrdersByPeriodParams) (entity.AggregatedPeriodOrders, error)
}

type ListOrdersParams struct {
	ShopID   string
	UserID   string
	Types    []entity.OrderType
	Statuses []entity.OrderStatus
	Limit    int
	Offset   int
}

type UpdateOrderAuthorizedParams struct {
	PaymentID string
	IssuedAt  time.Time
}

type UpdateOrderCapturedParams struct {
	PaymentID string
	IssuedAt  time.Time
}

type UpdateOrderFailedParams struct {
	Status    entity.PaymentStatus
	PaymentID string
	IssuedAt  time.Time
}

type UpdateOrderRefundedParams struct {
	Status       entity.PaymentStatus
	RefundType   entity.RefundType
	RefundTotal  int64
	RefundReason string
	IssuedAt     time.Time
}

type UpdateOrderFulfillmentParams struct {
	Status          entity.FulfillmentStatus
	ShippingCarrier entity.ShippingCarrier
	TrackingNumber  string
	ShippedAt       time.Time
}

type DraftOrderParams struct {
	ShippingMessage string
}

type CompleteOrderParams struct {
	ShippingMessage string
	CompletedAt     time.Time
}

type AggregateOrdersParams struct {
	ShopID       string
	CreatedAtGte time.Time
	CreatedAtLt  time.Time
}

type AggregateOrdersByUserParams struct {
	ShopID  string
	UserIDs []string
}

type AggregateOrdersByPaymentMethodTypeParams struct {
	ShopID             string
	PaymentMethodTypes []entity.PaymentMethodType
	CreatedAtGte       time.Time
	CreatedAtLt        time.Time
}

type AggregateOrdersByPromotionParams struct {
	ShopID       string
	PromotionIDs []string
}

type AggregateOrdersByPeriodParams struct {
	ShopID       string
	PeriodType   entity.AggregateOrderPeriodType
	CreatedAtGte time.Time
	CreatedAtLt  time.Time
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
	ShopID         string
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

type ProductReview interface {
	List(ctx context.Context, params *ListProductReviewsParams, fields ...string) (entity.ProductReviews, string, error)
	Get(ctx context.Context, productReviewID string, fields ...string) (*entity.ProductReview, error)
	Create(ctx context.Context, productReview *entity.ProductReview) error
	Update(ctx context.Context, productReviewID string, params *UpdateProductReviewParams) error
	Delete(ctx context.Context, productReviewID string) error
	Aggregate(ctx context.Context, params *AggregateProductReviewsParams) (entity.AggregatedProductReviews, error)
}

type ListProductReviewsParams struct {
	ProductID string
	UserID    string
	Rates     []int64
	Limit     int64
	NextToken string
}

type UpdateProductReviewParams struct {
	Rate    int64
	Title   string
	Comment string
}

type AggregateProductReviewsParams struct {
	ProductIDs []string
}

type ProductReviewReaction interface {
	Upsert(ctx context.Context, reaction *entity.ProductReviewReaction) error
	Delete(ctx context.Context, productReviewID, userID string) error
	GetUserReactions(ctx context.Context, productID, userID string) (entity.ProductReviewReactions, error)
}

type ProductTag interface {
	List(ctx context.Context, params *ListProductTagsParams, fields ...string) (entity.ProductTags, error)
	Count(ctx context.Context, params *ListProductTagsParams) (int64, error)
	MultiGet(ctx context.Context, productTagIDs []string, fields ...string) (entity.ProductTags, error)
	Get(ctx context.Context, productTagID string, fields ...string) (*entity.ProductTag, error)
	Create(ctx context.Context, productTag *entity.ProductTag) error
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
	ListPromotionsOrderByTitle     ListPromotionsOrderKey = "title"
	ListPromotionsOrderByPublic    ListPromotionsOrderKey = "public"
	ListPromotionsOrderByStartAt   ListPromotionsOrderKey = "start_at"
	ListPromotionsOrderByEndAt     ListPromotionsOrderKey = "end_at"
	ListPromotionsOrderByCreatedAt ListPromotionsOrderKey = "created_at"
	ListPromotionsOrderByUpdatedAt ListPromotionsOrderKey = "updated_at"
)

type ListPromotionsParams struct {
	ShopID        string
	Title         string
	Limit         int
	Offset        int
	Orders        []*ListPromotionsOrder
	WithAllTarget bool
}

type ListPromotionsOrder struct {
	Key        ListPromotionsOrderKey
	OrderByASC bool
}

type UpdatePromotionParams struct {
	Title        string
	Description  string
	Public       bool
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
	ShopID        string
	CoordinatorID string // Deprecated
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
	List(ctx context.Context, params *ListShippingsParams, fields ...string) (entity.Shippings, error)
	ListByShopIDs(ctx context.Context, shopIDs []string, fields ...string) (entity.Shippings, error)
	Count(ctx context.Context, params *ListShippingsParams) (int64, error)
	MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Shippings, error)
	Get(ctx context.Context, shippingID string, fields ...string) (*entity.Shipping, error)
	GetDefault(ctx context.Context, fields ...string) (*entity.Shipping, error)
	GetByShopID(ctx context.Context, shopID string, fields ...string) (*entity.Shipping, error)
	GetByCoordinatorID(ctx context.Context, coordinatorID string, fields ...string) (*entity.Shipping, error) // Depcecated
	Create(ctx context.Context, shipping *entity.Shipping) error
	Update(ctx context.Context, shippingID string, params *UpdateShippingParams) error
	UpdateInUse(ctx context.Context, shopID, shippingID string) error
	Delete(ctx context.Context, shippingID string) error
}

type ListShippingsParams struct {
	ShopID    string
	ShopIDs   []string
	OnlyInUse bool
	Limit     int
	Offset    int
}

type UpdateShippingParams struct {
	Name              string
	Box60Rates        entity.ShippingRates
	Box60Frozen       int64
	Box80Rates        entity.ShippingRates
	Box80Frozen       int64
	Box100Rates       entity.ShippingRates
	Box100Frozen      int64
	HasFreeShipping   bool
	FreeShippingRates int64
}

type Shop interface {
	List(ctx context.Context, params *ListShopsParams, fields ...string) (entity.Shops, error)
	MultiGet(ctx context.Context, shopIDs []string, fields ...string) (entity.Shops, error)
	Count(ctx context.Context, params *ListShopsParams) (int64, error)
	Get(ctx context.Context, shopID string, fields ...string) (*entity.Shop, error)
	GetByCoordinatorID(ctx context.Context, coordinatorID string, fields ...string) (*entity.Shop, error)
	Create(ctx context.Context, shop *entity.Shop) error
	Update(ctx context.Context, shopID string, params *UpdateShopParams) error
	Delete(ctx context.Context, shopID string) error
	RemoveProductType(ctx context.Context, productTypeID string) error
	ListProducers(ctx context.Context, params *ListShopProducersParams) ([]string, error)
	RelateProducer(ctx context.Context, shopID, producerID string) error
	UnrelateProducer(ctx context.Context, shopID, producerID string) error
}

type ListShopsParams struct {
	CoordinatorIDs []string
	ProducerIDs    []string
	Limit          int
	Offset         int
}

type UpdateShopParams struct {
	Name           string
	ProductTypeIDs []string
	BusinessDays   []time.Weekday
}

type ListShopProducersParams struct {
	ShopID string
	Limit  int
	Offset int
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
	SpotTypeIDs     []string
	UserID          string
	ExcludeApproved bool
	ExcludeDisabled bool
	Limit           int
	Offset          int
}

type ListSpotsByGeolocationParams struct {
	SpotTypeIDs     []string
	Longitude       float64
	Latitude        float64
	Radius          int64
	ExcludeDisabled bool
}

type UpdateSpotParams struct {
	SpotTypeID     string
	Name           string
	Description    string
	ThumbnailURL   string
	Longitude      float64
	Latitude       float64
	PostalCode     string
	PrefectureCode int32
	City           string
	AddressLine1   string
	AddressLine2   string
}

type ApproveSpotParams struct {
	Approved        bool
	ApprovedAdminID string
}

type SpotType interface {
	List(ctx context.Context, params *ListSpotTypesParams, fields ...string) (entity.SpotTypes, error)
	Count(ctx context.Context, params *ListSpotTypesParams) (int64, error)
	MultiGet(ctx context.Context, spotTypeIDs []string, fields ...string) (entity.SpotTypes, error)
	Get(ctx context.Context, spotTypeID string, fields ...string) (*entity.SpotType, error)
	Create(ctx context.Context, spotType *entity.SpotType) error
	Update(ctx context.Context, spotTypeID string, params *UpdateSpotTypeParams) error
	Delete(ctx context.Context, spotTypeID string) error
}

type ListSpotTypesParams struct {
	Name   string
	Limit  int
	Offset int
}

type UpdateSpotTypeParams struct {
	Name string
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
