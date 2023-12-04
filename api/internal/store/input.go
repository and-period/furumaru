package store

import (
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ListCategoriesInput struct {
	Name   string                 `validate:"max=32"`
	Limit  int64                  `validate:"required,max=200"`
	Offset int64                  `validate:"min=0"`
	Orders []*ListCategoriesOrder `validate:"dive,required"`
}

type ListCategoriesOrder struct {
	Key        entity.CategoryOrderBy `validate:"required"`
	OrderByASC bool                   `validate:""`
}

type MultiGetCategoriesInput struct {
	CategoryIDs []string `validate:"dive,required"`
}

type GetCategoryInput struct {
	CategoryID string `validate:"required"`
}

type CreateCategoryInput struct {
	Name string `validate:"required,max=32"`
}

type UpdateCategoryInput struct {
	CategoryID string `validate:"required"`
	Name       string `validate:"required,max=32"`
}

type DeleteCategoryInput struct {
	CategoryID string `validate:"required"`
}

type ListProductTypesInput struct {
	Name       string                   `validate:"max=32"`
	CategoryID string                   `validate:""`
	Limit      int64                    `validate:"required,max=200"`
	Offset     int64                    `validate:"min=0"`
	Orders     []*ListProductTypesOrder `validate:"dive,required"`
}

type ListProductTypesOrder struct {
	Key        entity.ProductTypeOrderBy `validate:"required"`
	OrderByASC bool                      `validate:""`
}

type MultiGetProductTypesInput struct {
	ProductTypeIDs []string `validate:"dive,required"`
}

type GetProductTypeInput struct {
	ProductTypeID string `validate:"required"`
}

type CreateProductTypeInput struct {
	Name       string `validate:"required,max=32"`
	IconURL    string `validate:"required"`
	CategoryID string `validate:"required"`
}

type UpdateProductTypeInput struct {
	ProductTypeID string `validate:"required"`
	Name          string `validate:"required,max=32"`
	IconURL       string `validate:"required"`
}

type UpdateProductTypeIconsInput struct {
	ProductTypeID string        `validate:"required"`
	Icons         common.Images `validate:""`
}

type DeleteProductTypeInput struct {
	ProductTypeID string `validate:"required"`
}

type ListProductTagsInput struct {
	Name   string                  `validate:"max=32"`
	Limit  int64                   `validate:"required,max=200"`
	Offset int64                   `validate:"min=0"`
	Orders []*ListProductTagsOrder `validate:"dive,required"`
}

type ListProductTagsOrder struct {
	Key        entity.ProductTagOrderBy `validate:"required"`
	OrderByASC bool                     `validate:""`
}

type MultiGetProductTagsInput struct {
	ProductTagIDs []string `validate:"dive,required"`
}

type GetProductTagInput struct {
	ProductTagID string `validate:"required"`
}

type CreateProductTagInput struct {
	Name string `validate:"required,max=32"`
}

type UpdateProductTagInput struct {
	ProductTagID string `validate:"required"`
	Name         string `validate:"required,max=32"`
}

type DeleteProductTagInput struct {
	ProductTagID string `validate:"required"`
}

type ListShippingsByCoordinatorIDsInput struct {
	CoordinatorIDs []string `validate:"dive,required"`
}

type MultiGetShippingsByRevisionInput struct {
	ShippingRevisionIDs []int64 `validate:"dive,required"`
}

type GetDefaultShippingInput struct{}

type GetShippingByCoordinatorIDInput struct {
	CoordinatorID string `validate:"required"`
}

type UpsertShippingInput struct {
	CoordinatorID     string                `validate:"required"`
	Box60Rates        []*UpsertShippingRate `validate:"required,dive,required"`
	Box60Frozen       int64                 `validate:"min=0,lt=10000000000"`
	Box80Rates        []*UpsertShippingRate `validate:"required,dive,required"`
	Box80Frozen       int64                 `validate:"min=0,lt=10000000000"`
	Box100Rates       []*UpsertShippingRate `validate:"required,dive,required"`
	Box100Frozen      int64                 `validate:"min=0,lt=10000000000"`
	HasFreeShipping   bool                  `validate:""`
	FreeShippingRates int64                 `validate:"min=0,lt=10000000000"`
}

type UpsertShippingRate struct {
	Name            string  `validate:"required"`
	Price           int64   `validate:"required,lt=10000000000"`
	PrefectureCodes []int32 `validate:"required"`
}

type UpdateDefaultShippingInput struct {
	Box60Rates        []*UpdateDefaultShippingRate `validate:"required,dive,required"`
	Box60Frozen       int64                        `validate:"min=0,lt=10000000000"`
	Box80Rates        []*UpdateDefaultShippingRate `validate:"required,dive,required"`
	Box80Frozen       int64                        `validate:"min=0,lt=10000000000"`
	Box100Rates       []*UpdateDefaultShippingRate `validate:"required,dive,required"`
	Box100Frozen      int64                        `validate:"min=0,lt=10000000000"`
	HasFreeShipping   bool                         `validate:""`
	FreeShippingRates int64                        `validate:"min=0,lt=10000000000"`
}

type UpdateDefaultShippingRate struct {
	Name            string  `validate:"required"`
	Price           int64   `validate:"required,lt=10000000000"`
	PrefectureCodes []int32 `validate:"required"`
}

type ListProductsInput struct {
	Name          string               `validate:"max=128"`
	CoordinatorID string               `validate:""`
	ProducerID    string               `validate:""`
	ProducerIDs   []string             `validate:"dive,required"`
	OnlyPublished bool                 `validate:""`
	Limit         int64                `validate:"required,max=200"`
	Offset        int64                `validate:"min=0"`
	Orders        []*ListProductsOrder `validate:"dive,required"`
}

type ListProductsOrder struct {
	Key        entity.ProductOrderBy `validate:"required"`
	OrderByASC bool                  `validate:""`
}

type MultiGetProductsInput struct {
	ProductIDs []string `validate:"dive,required"`
}

type MultiGetProductsByRevisionInput struct {
	ProductRevisionIDs []int64 `validate:"dive,required"`
}

type GetProductInput struct {
	ProductID string `validate:"required"`
}

type CreateProductInput struct {
	CoordinatorID        string                   `validate:"required"`
	ProducerID           string                   `validate:"required"`
	TypeID               string                   `validate:"required"`
	TagIDs               []string                 `validate:"max=8,dive,required"`
	Name                 string                   `validate:"required,max=128"`
	Description          string                   `validate:"required,max=20000"`
	Public               bool                     `validate:""`
	Inventory            int64                    `validate:"min=0"`
	Weight               int64                    `validate:"min=0"`
	WeightUnit           entity.WeightUnit        `validate:"required,oneof=1 2"`
	Item                 int64                    `validate:"min=1"`
	ItemUnit             string                   `validate:"required,max=16"`
	ItemDescription      string                   `validate:"required,max=64"`
	Media                []*CreateProductMedia    `validate:"max=8,unique=URL"`
	Price                int64                    `validate:"min=0"`
	Cost                 int64                    `validate:"min=0"`
	ExpirationDate       int64                    `validate:"min=0"`
	RecommendedPoints    []string                 `validate:"max=3,dive,max=128"`
	StorageMethodType    entity.StorageMethodType `validate:"required,oneof=1 2 3 4"`
	DeliveryType         entity.DeliveryType      `validate:"required,oneof=1 2 3"`
	Box60Rate            int64                    `validate:"min=0,max=100"`
	Box80Rate            int64                    `validate:"min=0,max=100"`
	Box100Rate           int64                    `validate:"min=0,max=100"`
	OriginPrefectureCode int32                    `validate:"required"`
	OriginCity           string                   `validate:"max=32"`
	StartAt              time.Time                `validate:"required"`
	EndAt                time.Time                `validate:"required,gtfield=StartAt"`
}

type CreateProductMedia struct {
	URL         string `validate:"required,url"`
	IsThumbnail bool   `validate:""`
}

type UpdateProductInput struct {
	ProductID            string                   `validate:"required"`
	TypeID               string                   `validate:"required"`
	TagIDs               []string                 `validate:"max=8,dive,required"`
	Name                 string                   `validate:"required,max=128"`
	Description          string                   `validate:"required,max=20000"`
	Public               bool                     `validate:""`
	Inventory            int64                    `validate:"min=0"`
	Weight               int64                    `validate:"min=0"`
	WeightUnit           entity.WeightUnit        `validate:"required,oneof=1 2"`
	Item                 int64                    `validate:"min=1"`
	ItemUnit             string                   `validate:"required,max=16"`
	ItemDescription      string                   `validate:"required,max=64"`
	Media                []*UpdateProductMedia    `validate:"max=8,unique=URL"`
	Price                int64                    `validate:"min=0"`
	Cost                 int64                    `validate:"min=0"`
	ExpirationDate       int64                    `validate:"min=0"`
	RecommendedPoints    []string                 `validate:"max=3,dive,max=128"`
	StorageMethodType    entity.StorageMethodType `validate:"required,oneof=1 2 3 4"`
	DeliveryType         entity.DeliveryType      `validate:"required,oneof=1 2 3"`
	Box60Rate            int64                    `validate:"min=0,max=100"`
	Box80Rate            int64                    `validate:"min=0,max=100"`
	Box100Rate           int64                    `validate:"min=0,max=100"`
	OriginPrefectureCode int32                    `validate:"required"`
	OriginCity           string                   `validate:"max=32"`
	StartAt              time.Time                `validate:"required"`
	EndAt                time.Time                `validate:"required,gtfield=StartAt"`
}

type UpdateProductMedia struct {
	URL         string `validate:"required,url"`
	IsThumbnail bool   `validate:""`
}

type UpdateProductMediaInput struct {
	ProductID string                     `validate:"required"`
	Images    []*UpdateProductMediaImage `validate:"unique=OriginURL,dive,required"`
}

type UpdateProductMediaImage struct {
	OriginURL string        `validate:"required"`
	Images    common.Images `validate:""`
}

type DeleteProductInput struct {
	ProductID string `validate:"required"`
}

type ListPromotionsInput struct {
	Title  string                 `validate:"max=64"`
	Limit  int64                  `validate:"required,max=200"`
	Offset int64                  `validate:"min=0"`
	Orders []*ListPromotionsOrder `validate:"dive,required"`
}

type ListPromotionsOrder struct {
	Key        entity.PromotionOrderBy `validate:"required"`
	OrderByASC bool                    `validate:""`
}

type MultiGetPromotionsInput struct {
	PromotionIDs []string `validate:"dive,required"`
}

type GetPromotionInput struct {
	PromotionID string `validate:"required"`
}

type CreatePromotionInput struct {
	Title        string                   `validate:"required,max=64"`
	Description  string                   `validate:"required,max=2000"`
	Public       bool                     `validate:""`
	DiscountType entity.DiscountType      `validate:"required,oneof=1 2 3"`
	DiscountRate int64                    `validate:"min=0"`
	Code         string                   `validate:"len=8"`
	CodeType     entity.PromotionCodeType `validate:"required,oneof=1 2"`
	StartAt      time.Time                `validate:"required"`
	EndAt        time.Time                `validate:"required,gtfield=StartAt"`
}

type UpdatePromotionInput struct {
	PromotionID  string                   `validate:"required"`
	Title        string                   `validate:"required,max=64"`
	Description  string                   `validate:"required,max=2000"`
	Public       bool                     `validate:""`
	DiscountType entity.DiscountType      `validate:"required,oneof=1 2 3"`
	DiscountRate int64                    `validate:"min=0"`
	Code         string                   `validate:"len=8"`
	CodeType     entity.PromotionCodeType `validate:"required,oneof=1 2"`
	StartAt      time.Time                `validate:"required"`
	EndAt        time.Time                `validate:"required,gtfield=StartAt"`
}

type DeletePromotionInput struct {
	PromotionID string `validate:"required"`
}

type ListSchedulesInput struct {
	CoordinatorID string    `validate:""`
	StartAtGte    time.Time `validate:""`
	StartAtLt     time.Time `validate:""`
	EndAtGte      time.Time `validate:""`
	EndAtLt       time.Time `validate:""`
	OnlyPublished bool      `validate:""`
	Limit         int64     `validate:"required_without=NoLimit,min=0,max=200"`
	Offset        int64     `validate:"min=0"`
	NoLimit       bool      `validate:""`
}

type MultiGetSchedulesInput struct {
	ScheduleIDs []string `validate:"dive,required"`
}

type GetScheduleInput struct {
	ScheduleID string `validate:"required"`
}

type CreateScheduleInput struct {
	CoordinatorID   string    `validate:"required"`
	Title           string    `validate:"required,max=64"`
	Description     string    `validate:"required,max=2000"`
	ThumbnailURL    string    `validate:"url"`
	ImageURL        string    `validate:"url"`
	OpeningVideoURL string    `validate:"url"`
	Public          bool      `validate:""`
	StartAt         time.Time `validate:"required"`
	EndAt           time.Time `validate:"required,gtfield=StartAt"`
}

type UpdateScheduleInput struct {
	ScheduleID      string    `validate:"required"`
	Title           string    `validate:"required,max=64"`
	Description     string    `validate:"required,max=2000"`
	ThumbnailURL    string    `validate:"url"`
	ImageURL        string    `validate:"url"`
	OpeningVideoURL string    `validate:"url"`
	StartAt         time.Time `validate:"required"`
	EndAt           time.Time `validate:"required,gtfield=StartAt"`
}

type UpdateScheduleThumbnailsInput struct {
	ScheduleID string        `validate:"required"`
	Thumbnails common.Images `validate:""`
}

type ApproveScheduleInput struct {
	ScheduleID string `validate:"required"`
	AdminID    string `validate:"required"`
	Approved   bool   `validate:""`
}

type PublishScheduleInput struct {
	ScheduleID string `validate:"required"`
	Public     bool   `validate:""`
}

type GetLiveInput struct {
	LiveID string `validate:"required"`
}

type ListLivesInput struct {
	ScheduleIDs []string `validate:"dive,required"`
	Limit       int64    `validate:"required_without=NoLimit,min=0,max=200"`
	Offset      int64    `validate:"min=0"`
	NoLimit     bool     `validate:""`
}

type CreateLiveInput struct {
	ScheduleID string    `validate:"required"`
	ProducerID string    `validate:"required"`
	ProductIDs []string  `validate:"min=0,max=8,unique"`
	Comment    string    `validate:"required,max=2000"`
	StartAt    time.Time `validate:"required"`
	EndAt      time.Time `validate:"required,gtfield=StartAt"`
}

type UpdateLiveInput struct {
	LiveID     string    `validate:"required"`
	ProductIDs []string  `validate:"min=0,max=8,unique"`
	Comment    string    `validate:"required,max=2000"`
	StartAt    time.Time `validate:"required"`
	EndAt      time.Time `validate:"required,gtfield=StartAt"`
}

type DeleteLiveInput struct {
	LiveID string `validate:"required"`
}

type ListOrdersInput struct {
	CoordinatorID string `validate:""`
	UserID        string `validate:""`
	Limit         int64  `validate:"required,max=200"`
	Offset        int64  `validate:"min=0"`
}

type GetOrderInput struct {
	OrderID string `validate:"required"`
}

type CaptureOrderInput struct {
	OrderID string `validate:"required"`
}

type DraftOrderInput struct {
	OrderID         string `validate:"required"`
	ShippingMessage string `validate:"max=2000"`
}

type CompleteOrderInput struct {
	OrderID         string `validate:"required"`
	ShippingMessage string `validate:"required,max=2000"`
}

type CancelOrderInput struct {
	OrderID string `validate:"required"`
}

type RefundOrderInput struct {
	OrderID     string `validate:"required"`
	Description string `validate:"required"`
}

type UpdateOrderFulfillmentInput struct {
	OrderID         string                 `validate:"required"`
	FulfillmentID   string                 `validate:"required"`
	ShippingCarrier entity.ShippingCarrier `validate:"required,oneof=1 2"`
	TrackingNumber  string                 `validate:"required"`
}

type AggregateOrdersInput struct {
	CoordinatorID string   `validate:""`
	UserIDs       []string `validate:"dive,required"`
}

type GetCartInput struct {
	SessionID string `validate:"required"`
}

type AddCartItemInput struct {
	SessionID string `validate:"required"`
	ProductID string `validate:"required"`
	Quantity  int64  `validate:"min=1"`
}

type RemoveCartItemInput struct {
	SessionID string `validate:"required"`
	BoxNumber int64  `validate:"required"`
	ProductID string `validate:"required"`
}

type CheckoutCreditCardInput struct {
	CheckoutDetail
	Number            string `validate:"required,numeric"`
	Month             int64  `validate:"min=1,max=12"`
	Year              int64  `validate:"min=2000,max=2100"`
	VerificationValue string `validate:"required,numeric,max=4"`
}

type CheckoutPayPayInput struct {
	CheckoutDetail
}

type CheckoutLinePayInput struct {
	CheckoutDetail
}

type CheckoutMerpayInput struct {
	CheckoutDetail
}

type CheckoutRakutenPayInput struct {
	CheckoutDetail
}

type CheckoutAUPayInput struct {
	CheckoutDetail
}

type CheckoutDetail struct {
	UserID            string `validate:"required"`
	SessionID         string `validate:"required"`
	CoordinatorID     string `validate:"required"`
	BoxNumber         int64  `validate:"min=0"`
	PromotionID       string `validate:""`
	BillingAddressID  string `validate:"required"`
	ShippingAddressID string `validate:"required"`
	CallbackURL       string `validate:"required,http_url"`
	Total             int64  `validate:"required"`
}

type MultiGetPaymentSystemsInput struct {
	MethodTypes []entity.PaymentMethodType `validate:"dive,required"`
}

type GetPaymentSystemInput struct {
	MethodType entity.PaymentMethodType `validate:"required"`
}

type UpdatePaymentStatusInput struct {
	MethodType entity.PaymentMethodType   `validate:"required"`
	Status     entity.PaymentSystemStatus `validate:"required"`
}

type SearchPostalCodeInput struct {
	PostlCode string `validate:"required,numeric,len=7"`
}

type NotifyPaymentCompletedInput struct {
	OrderID   string               `validate:"required"`
	PaymentID string               `validate:"required"`
	Status    entity.PaymentStatus `validate:"required"`
	IssuedAt  time.Time            `validate:"required"`
}

type NotifyPaymentRefundedInput struct {
	OrderID  string               `validate:"required"`
	Status   entity.PaymentStatus `validate:"required"`
	Type     entity.RefundType    `validate:"required,oneof=1 2"`
	Reason   string               `validate:"max=2000"`
	Total    int64                `validate:"min=0"`
	IssuedAt time.Time            `validate:"required"`
}
