package store

import (
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

/**
 * Cart - 買い物かご
 */
type GetCartInput struct {
	SessionID string `validate:"required"`
}

type CalcCartInput struct {
	SessionID      string `validate:"required"`
	CoordinatorID  string `validate:"required"`
	BoxNumber      int64  `validate:"min=0"`
	PromotionCode  string `validate:"omitempty,len=8"`
	PrefectureCode int32  `validate:"min=0,max=47"`
	Pickup         bool   `validate:""`
}

type AddCartItemInput struct {
	SessionID string `validate:"required"`
	UserID    string `validate:""`
	UserAgent string `validate:""`
	ClientIP  string `validate:"omitempty,ip_addr"`
	ProductID string `validate:"required"`
	Quantity  int64  `validate:"min=1"`
}

type RemoveCartItemInput struct {
	SessionID string `validate:"required"`
	UserID    string `validate:""`
	UserAgent string `validate:""`
	ClientIP  string `validate:"omitempty,ip_addr"`
	BoxNumber int64  `validate:"min=0"`
	ProductID string `validate:"required"`
}

/**
 * Category - 商品カテゴリ
 */
type ListCategoriesOrderKey int32

const (
	ListCategoriesOrderByName ListCategoriesOrderKey = iota + 1
)

type ListCategoriesInput struct {
	Name   string                 `validate:"max=32"`
	Limit  int64                  `validate:"required,max=200"`
	Offset int64                  `validate:"min=0"`
	Orders []*ListCategoriesOrder `validate:"dive,required"`
}

type ListCategoriesOrder struct {
	Key        ListCategoriesOrderKey `validate:"required"`
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

/**
 * Checkout - 購入処理
 */
type GetCheckoutStateInput struct {
	UserID        string `validate:""`
	SessionID     string `validate:"required_without=UserID"`
	TransactionID string `validate:"required"`
}

type CheckoutCreditCardInput struct {
	CheckoutDetail
	Name              string `validate:"required"`
	Number            string `validate:"required,credit_card"`
	Month             int64  `validate:"min=1,max=12"`
	Year              int64  `validate:"min=2000,max=2100"`
	VerificationValue string `validate:"min=3,max=4,numeric"`
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

type CheckoutPaidyInput struct {
	CheckoutDetail
}

type CheckoutBankTransferInput struct {
	CheckoutDetail
}

type CheckoutPayEasyInput struct {
	CheckoutDetail
}

type CheckoutFreeInput struct {
	CheckoutDetail
}

type CheckoutDetail struct {
	CheckoutProductDetail    `validate:"-"`
	CheckoutExperienceDetail `validate:"-"`
	Type                     entity.OrderType `validate:"required"`
	UserID                   string           `validate:"required"`
	SessionID                string           `validate:"required"`
	RequestID                string           `validate:"required"`
	PromotionCode            string           `validate:"omitempty,len=8"`
	BillingAddressID         string           `validate:""`
	CallbackURL              string           `validate:"required,http_url"`
	Total                    int64            `validate:"min=0"`
	OrderRequest             string           `validate:"max=256"`
}

type CheckoutProductDetail struct {
	CoordinatorID     string    `validate:"required"`
	BoxNumber         int64     `validate:"min=0"`
	ShippingAddressID string    `validate:"required_without=Pickup"`
	Pickup            bool      `validate:""`
	PickupAt          time.Time `validate:"required_with=Pickup"`
	PickupLocation    string    `validate:"required_with=Pickup"`
}

type CheckoutExperienceDetail struct {
	ExperienceID          string `validate:"required"`
	AdultCount            int64  `validate:"min=0,max=99"`
	JuniorHighSchoolCount int64  `validate:"min=0,max=99"`
	ElementarySchoolCount int64  `validate:"min=0,max=99"`
	PreschoolCount        int64  `validate:"min=0,max=99"`
	SeniorCount           int64  `validate:"min=0,max=99"`
	Transportation        string `validate:"max=256"`
	RequestedDate         string `validate:"omitempty,date"`
	RequestedTime         string `validate:"omitempty,time"`
}

type NotifyPaymentAuthorizedInput struct {
	NotifyPaymentPayload
}

type NotifyPaymentCapturedInput struct {
	NotifyPaymentPayload
}

type NotifyPaymentFailedInput struct {
	NotifyPaymentPayload
}

type NotifyPaymentRefundedInput struct {
	NotifyPaymentPayload
	Type   entity.RefundType `validate:"required,oneof=1 2"`
	Reason string            `validate:"max=2000"`
	Total  int64             `validate:"min=0"`
}

type NotifyPaymentPayload struct {
	OrderID   string               `validate:"required"`
	PaymentID string               `validate:""`
	IssuedAt  time.Time            `validate:"required"`
	Status    entity.PaymentStatus `validate:"required"`
}

/**
 * Experience - 体験
 */
type ListExperiencesInput struct {
	Name            string `validate:"max=64"`
	PrefectureCode  int32  `validate:"min=0,max=47"`
	ShopID          string `validate:""`
	ProducerID      string `validate:""`
	OnlyPublished   bool   `validate:""`
	ExcludeFinished bool   `validate:""`
	ExcludeDeleted  bool   `validate:""`
	Limit           int64  `validate:"required_without=NoLimit,min=0,max=200"`
	Offset          int64  `validate:"min=0"`
	NoLimit         bool   `validate:""`
}

type ListExperiencesByGeolocationInput struct {
	ShopID          string  `validate:""`
	ProducerID      string  `validate:""`
	Latitude        float64 `validate:"min=-90,max=90"`
	Longitude       float64 `validate:"min=-180,max=180"`
	Radius          int64   `validate:"min=0"`
	OnlyPublished   bool    `validate:""`
	ExcludeFinished bool    `validate:""`
	ExcludeDeleted  bool    `validate:""`
}

type MultiGetExperiencesInput struct {
	ExperienceIDs []string `validate:"dive,required"`
}

type MultiGetExperiencesByRevisionInput struct {
	ExperienceRevisionIDs []int64 `validate:"dive,required"`
}

type GetExperienceInput struct {
	ExperienceID string `validate:"required"`
}

type CreateExperienceInput struct {
	ShopID                string                   `validate:"required"`
	CoordinatorID         string                   `validate:"required"`
	ProducerID            string                   `validate:"required"`
	TypeID                string                   `validate:"required"`
	Title                 string                   `validate:"required,max=128"`
	Description           string                   `validate:"required,max=20000"`
	Public                bool                     `validate:""`
	SoldOut               bool                     `validate:""`
	Media                 []*CreateExperienceMedia `validate:"max=8,unique=URL"`
	PriceAdult            int64                    `validate:"min=0"`
	PriceJuniorHighSchool int64                    `validate:"min=0"`
	PriceElementarySchool int64                    `validate:"min=0"`
	PricePreschool        int64                    `validate:"min=0"`
	PriceSenior           int64                    `validate:"min=0"`
	RecommendedPoints     []string                 `validate:"max=3,dive,max=128"`
	PromotionVideoURL     string                   `validate:"omitempty,url"`
	Duration              int64                    `validate:"min=0"`
	Direction             string                   `validate:"max=2000"`
	BusinessOpenTime      string                   `validate:"time"`
	BusinessCloseTime     string                   `validate:"time"`
	HostPostalCode        string                   `validate:"required,max=16,numeric"`
	HostPrefectureCode    int32                    `validate:"required"`
	HostCity              string                   `validate:"required,max=32"`
	HostAddressLine1      string                   `validate:"required,max=64"`
	HostAddressLine2      string                   `validate:"max=64"`
	StartAt               time.Time                `validate:"required"`
	EndAt                 time.Time                `validate:"required,gtfield=StartAt"`
}

type CreateExperienceMedia struct {
	URL         string `validate:"required,url"`
	IsThumbnail bool   `validate:""`
}

type UpdateExperienceInput struct {
	ExperienceID          string                   `validate:"required"`
	TypeID                string                   `validate:"required"`
	Title                 string                   `validate:"required,max=128"`
	Description           string                   `validate:"required,max=20000"`
	Public                bool                     `validate:""`
	SoldOut               bool                     `validate:""`
	Media                 []*UpdateExperienceMedia `validate:"max=8,unique=URL"`
	PriceAdult            int64                    `validate:"min=0"`
	PriceJuniorHighSchool int64                    `validate:"min=0"`
	PriceElementarySchool int64                    `validate:"min=0"`
	PricePreschool        int64                    `validate:"min=0"`
	PriceSenior           int64                    `validate:"min=0"`
	RecommendedPoints     []string                 `validate:"max=3,dive,max=128"`
	PromotionVideoURL     string                   `validate:"omitempty,url"`
	Duration              int64                    `validate:"min=0"`
	Direction             string                   `validate:"max=2000"`
	BusinessOpenTime      string                   `validate:"time"`
	BusinessCloseTime     string                   `validate:"time"`
	HostPostalCode        string                   `validate:"required,max=16,numeric"`
	HostPrefectureCode    int32                    `validate:"required"`
	HostCity              string                   `validate:"required,max=32"`
	HostAddressLine1      string                   `validate:"required,max=64"`
	HostAddressLine2      string                   `validate:"max=64"`
	StartAt               time.Time                `validate:"required"`
	EndAt                 time.Time                `validate:"required,gtfield=StartAt"`
}

type UpdateExperienceMedia struct {
	URL         string `validate:"required,url"`
	IsThumbnail bool   `validate:""`
}

type DeleteExperienceInput struct {
	ExperienceID string `validate:"required"`
}

/**
 * ExperienceReview - 体験レビュー
 */
type ListExperienceReviewsInput struct {
	ExperienceID string  `validate:"required"`
	UserID       string  `validate:""`
	Rates        []int64 `validate:"dive,min=1,max=5"`
	Limit        int64   `validate:"required_without=NoLimit,min=0,max=200"`
	NextToken    string  `validate:""`
	NoLimit      bool    `validate:""`
}

type GetExperienceReviewInput struct {
	ReviewID string `validate:"required"`
}

type CreateExperienceReviewInput struct {
	ExperienceID string `validate:"required"`
	UserID       string `validate:"required"`
	Rate         int64  `validate:"min=1,max=5"`
	Title        string `validate:"required,max=64"`
	Comment      string `validate:"required,max=2000"`
}

type UpdateExperienceReviewInput struct {
	ReviewID string `validate:"required"`
	Rate     int64  `validate:"min=1,max=5"`
	Title    string `validate:"required,max=64"`
	Comment  string `validate:"required,max=2000"`
}

type DeleteExperienceReviewInput struct {
	ReviewID string `validate:"required"`
}

type AggregateExperienceReviewsInput struct {
	ExperienceIDs []string `validate:"min=1,dive,required"`
}

/**
 * ExperienceReviewReaction - 体験レビューへのリアクション
 */
type UpsertExperienceReviewReactionInput struct {
	ReviewID     string                              `validate:"required"`
	UserID       string                              `validate:"required"`
	ReactionType entity.ExperienceReviewReactionType `validate:"required,oneof=1 2"`
}

type DeleteExperienceReviewReactionInput struct {
	ReviewID string `validate:"required"`
	UserID   string `validate:"required"`
}

type GetUserExperienceReviewReactionsInput struct {
	ExperienceID string `validate:"required"`
	UserID       string `validate:"required"`
}

/**
 * ExperienceType - 体験種別
 */
type ListExperienceTypesInput struct {
	Name   string `validate:"max=32"`
	Limit  int64  `validate:"required,max=200"`
	Offset int64  `validate:"min=0"`
}

type MultiGetExperienceTypesInput struct {
	ExperienceTypeIDs []string `validate:"dive,required"`
}

type GetExperienceTypeInput struct {
	ExperienceTypeID string `validate:"required"`
}

type CreateExperienceTypeInput struct {
	Name string `validate:"required,max=32"`
}

type UpdateExperienceTypeInput struct {
	ExperienceTypeID string `validate:"required"`
	Name             string `validate:"required,max=32"`
}

type DeleteExperienceTypeInput struct {
	ExperienceTypeID string `validate:"required"`
}

/**
 * Live - マルシェタイムテーブル
 */
type ListLivesInput struct {
	ScheduleIDs   []string `validate:"dive,required"`
	ProducerID    string   `validate:""`
	Limit         int64    `validate:"required_without=NoLimit,min=0,max=200"`
	Offset        int64    `validate:"min=0"`
	NoLimit       bool     `validate:""`
	OnlyPublished bool     `validate:""`
}

type GetLiveInput struct {
	LiveID string `validate:"required"`
}

type CreateLiveInput struct {
	ScheduleID string    `validate:"required"`
	ProducerID string    `validate:"required"`
	ProductIDs []string  `validate:"unique"`
	Comment    string    `validate:"required,max=2000"`
	StartAt    time.Time `validate:"required"`
	EndAt      time.Time `validate:"required,gtfield=StartAt"`
}

type UpdateLiveInput struct {
	LiveID     string    `validate:"required"`
	ProductIDs []string  `validate:"unique"`
	Comment    string    `validate:"required,max=2000"`
	StartAt    time.Time `validate:"required"`
	EndAt      time.Time `validate:"required,gtfield=StartAt"`
}

type DeleteLiveInput struct {
	LiveID string `validate:"required"`
}

/**
 * Order - 注文履歴
 */
type ListOrdersInput struct {
	ShopID   string               `validate:""`
	UserID   string               `validate:""`
	Types    []entity.OrderType   `validate:""`
	Statuses []entity.OrderStatus `validate:""`
	Limit    int64                `validate:"required,max=200"`
	Offset   int64                `validate:"min=0"`
}

type ListOrderUserIDsInput struct {
	ShopID string `validate:""`
	Limit  int64  `validate:"required,max=200"`
	Offset int64  `validate:"min=0"`
}

type GetOrderInput struct {
	OrderID string `validate:"required"`
}

type GetOrderByTransactionIDInput struct {
	UserID        string `validate:"required"`
	TransactionID string `validate:"required"`
}

type CaptureOrderInput struct {
	OrderID string `validate:"required"`
}

type DraftOrderInput struct {
	OrderID         string `validate:"required"`
	ShippingMessage string `validate:"max=2000"`
}

type CompleteProductOrderInput struct {
	OrderID         string `validate:"required"`
	ShippingMessage string `validate:"required,max=2000"`
}

type CompleteExperienceOrderInput struct {
	OrderID string `validate:"required"`
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
	ShopID       string    `validate:""`
	CreatedAtGte time.Time `validate:""`
	CreatedAtLt  time.Time `validate:""`
}

type AggregateOrdersByUserInput struct {
	ShopID  string   `validate:""`
	UserIDs []string `validate:"dive,required"`
}

type AggregateOrdersByPaymentMethodTypeInput struct {
	ShopID       string    `validate:""`
	CreatedAtGte time.Time `validate:""`
	CreatedAtLt  time.Time `validate:""`
}

type AggregateOrdersByPromotionInput struct {
	ShopID       string   `validate:""`
	PromotionIDs []string `validate:"dive,required"`
}

type AggregateOrdersByPeriodInput struct {
	ShopID       string                          `validate:""`
	PeriodType   entity.AggregateOrderPeriodType `validate:"required"`
	CreatedAtGte time.Time                       `validate:""`
	CreatedAtLt  time.Time                       `validate:""`
}

type ExportOrdersInput struct {
	ShopID          string                      `validate:""`
	ShippingCarrier entity.ShippingCarrier      `validate:"oneof=0 1 2"`
	EncodingType    codes.CharacterEncodingType `validate:"oneof=0 1"`
}

/**
 * PaymentSystem - 決済システム
 */
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

/**
 * PostalCode - 郵便番号
 */
type SearchPostalCodeInput struct {
	PostlCode string `validate:"required,numeric,len=7"`
}

/**
 * Product - 商品
 */
type ListProductsOrderKey int32

const (
	ListProductsOrderByName ListProductsOrderKey = iota + 1
	ListProductsOrderBySoldOut
	ListProductsOrderByPublic
	ListProductsOrderByInventory
	ListProductsOrderByOriginPrefecture
	ListProductsOrderByOriginCity
	ListProductsOrderByStartAt
	ListProductsOrderByCreatedAt
	ListProductsOrderByUpdatedAt
)

type ListProductsInput struct {
	Name             string                `validate:"max=128"`
	ShopID           string                `validate:""`
	ProducerID       string                `validate:""`
	ProducerIDs      []string              `validate:"dive,required"`
	Scopes           []entity.ProductScope `validate:""`
	ExcludeOutOfSale bool                  `validate:""`
	ExcludeDeleted   bool                  `validate:""`
	Limit            int64                 `validate:"required_without=NoLimit,min=0,max=200"`
	Offset           int64                 `validate:"min=0"`
	NoLimit          bool                  `validate:""`
	Orders           []*ListProductsOrder  `validate:"dive,required"`
}

type ListProductsOrder struct {
	Key        ListProductsOrderKey `validate:"required"`
	OrderByASC bool                 `validate:""`
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
	ShopID               string                   `validate:"required"`
	CoordinatorID        string                   `validate:"required"`
	ProducerID           string                   `validate:"required"`
	TypeID               string                   `validate:"required"`
	TagIDs               []string                 `validate:"max=8,dive,required"`
	Name                 string                   `validate:"required,max=128"`
	Description          string                   `validate:"required,max=20000"`
	Scope                entity.ProductScope      `validate:""`
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
	Box60Rate            int64                    `validate:"min=0,max=600"`
	Box80Rate            int64                    `validate:"min=0,max=250"`
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
	Scope                entity.ProductScope      `validate:""`
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
	Box60Rate            int64                    `validate:"min=0,max=600"`
	Box80Rate            int64                    `validate:"min=0,max=250"`
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

type DeleteProductInput struct {
	ProductID string `validate:"required"`
}

/**
 * ProductReview - 商品レビュー
 */
type ListProductReviewsInput struct {
	ProductID string  `validate:"required"`
	UserID    string  `validate:""`
	Rates     []int64 `validate:"dive,min=1,max=5"`
	Limit     int64   `validate:"required_without=NoLimit,min=0,max=200"`
	NextToken string  `validate:""`
	NoLimit   bool    `validate:""`
}

type GetProductReviewInput struct {
	ReviewID string `validate:"required"`
}

type CreateProductReviewInput struct {
	ProductID string `validate:"required"`
	UserID    string `validate:"required"`
	Rate      int64  `validate:"min=1,max=5"`
	Title     string `validate:"required,max=64"`
	Comment   string `validate:"required,max=2000"`
}

type UpdateProductReviewInput struct {
	ReviewID string `validate:"required"`
	Rate     int64  `validate:"min=1,max=5"`
	Title    string `validate:"required,max=64"`
	Comment  string `validate:"required,max=2000"`
}

type DeleteProductReviewInput struct {
	ReviewID string `validate:"required"`
}

type AggregateProductReviewsInput struct {
	ProductIDs []string `validate:"min=1,dive,required"`
}

/**
 * ProductReviewReaction - 商品レビューへのリアクション
 */
type UpsertProductReviewReactionInput struct {
	ReviewID     string                           `validate:"required"`
	UserID       string                           `validate:"required"`
	ReactionType entity.ProductReviewReactionType `validate:"required,oneof=1 2"`
}

type DeleteProductReviewReactionInput struct {
	ReviewID string `validate:"required"`
	UserID   string `validate:"required"`
}

type GetUserProductReviewReactionsInput struct {
	ProductID string `validate:"required"`
	UserID    string `validate:"required"`
}

/**
 * ProductTag - 商品タグ
 */
type ListProductTagsOrderKey int32

const (
	ListProductTagsOrderByName ListProductTagsOrderKey = iota + 1
)

type ListProductTagsInput struct {
	Name   string                  `validate:"max=32"`
	Limit  int64                   `validate:"required,max=200"`
	Offset int64                   `validate:"min=0"`
	Orders []*ListProductTagsOrder `validate:"dive,required"`
}

type ListProductTagsOrder struct {
	Key        ListProductTagsOrderKey `validate:"required"`
	OrderByASC bool                    `validate:""`
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

/**
 * ProductType - 品目
 */
type ListProductTypesOrderKey int32

const (
	ListProductTypesOrderByName ListProductTypesOrderKey = iota + 1
)

type ListProductTypesInput struct {
	Name       string                   `validate:"max=32"`
	CategoryID string                   `validate:""`
	Limit      int64                    `validate:"required,max=200"`
	Offset     int64                    `validate:"min=0"`
	Orders     []*ListProductTypesOrder `validate:"dive,required"`
}

type ListProductTypesOrder struct {
	Key        ListProductTypesOrderKey `validate:"required"`
	OrderByASC bool                     `validate:""`
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

type DeleteProductTypeInput struct {
	ProductTypeID string `validate:"required"`
}

/**
 * Promotion - プロモーション
 */
type ListPromotionsOrderKey int32

const (
	ListPromotionsOrderByTitle ListPromotionsOrderKey = iota + 1
	ListPromotionsOrderByPublic
	ListPromotionsOrderByStartAt
	ListPromotionsOrderByEndAt
	ListPromotionsOrderByCreatedAt
	ListPromotionsOrderByUpdatedAt
)

type ListPromotionsInput struct {
	ShopID        string                 `validate:""`
	Title         string                 `validate:"max=64"`
	Limit         int64                  `validate:"required,max=200"`
	Offset        int64                  `validate:"min=0"`
	Orders        []*ListPromotionsOrder `validate:"dive,required"`
	WithAllTarget bool                   `validate:""`
}

type ListPromotionsOrder struct {
	Key        ListPromotionsOrderKey `validate:"required"`
	OrderByASC bool                   `validate:""`
}

type MultiGetPromotionsInput struct {
	PromotionIDs []string `validate:"dive,required"`
}

type GetPromotionInput struct {
	PromotionID string `validate:"required"`
	ShopID      string `validate:""`
	OnlyEnabled bool   `validate:""`
}

type GetPromotionByCodeInput struct {
	PromotionCode string `validate:"required"`
	ShopID        string `validate:""`
	OnlyEnabled   bool   `validate:""`
}

type CreatePromotionInput struct {
	AdminID      string                   `validate:"required"`
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
	AdminID      string                   `validate:"required"`
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

/**
 * Schedule - マルシェ開催スケジュール
 */

type ListSchedulesInput struct {
	ShopID        string    `validate:""`
	ProducerID    string    `validate:""`
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
	ShopID          string    `validate:"required"`
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

type DeleteScheduleInput struct {
	ScheduleID string `validate:"required"`
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

/**
 * Shipping - 配送設定
 */
type ListShippingsByShopIDInput struct {
	ShopID string `validate:"required"`
	Limit  int64  `validate:"min=0,max=200"`
	Offset int64  `validate:"min=0"`
}

type ListShippingsByShopIDsInput struct {
	ShopIDs []string `validate:"dive,required"`
}

type ListShippingsByCoordinatorIDsInput struct {
	CoordinatorIDs []string `validate:"dive,required"`
}

type MultiGetShippingsByRevisionInput struct {
	ShippingRevisionIDs []int64 `validate:"dive,required"`
}

type GetShippingInput struct {
	ShippingID string `validate:"required"`
}

type GetDefaultShippingInput struct{}

type GetShippingByShopIDInput struct {
	ShopID string `validate:"required"`
}

type GetShippingByCoordinatorIDInput struct {
	CoordinatorID string `validate:"required"`
}

type CreateShippingInput struct {
	ShopID            string                `validate:"required"`
	CoordinatorID     string                `validate:"required"`
	Name              string                `validate:"required,max=64"`
	Box60Rates        []*CreateShippingRate `validate:"required,dive,required"`
	Box60Frozen       int64                 `validate:"min=0,lt=10000000000"`
	Box80Rates        []*CreateShippingRate `validate:"required,dive,required"`
	Box80Frozen       int64                 `validate:"min=0,lt=10000000000"`
	Box100Rates       []*CreateShippingRate `validate:"required,dive,required"`
	Box100Frozen      int64                 `validate:"min=0,lt=10000000000"`
	HasFreeShipping   bool                  `validate:""`
	FreeShippingRates int64                 `validate:"min=0,lt=10000000000"`
	InUse             bool                  `validate:""`
}

type CreateShippingRate struct {
	Name            string  `validate:"required"`
	Price           int64   `validate:"min=0,lt=10000000000"`
	PrefectureCodes []int32 `validate:"required"`
}

type UpdateShippingInput struct {
	ShippingID        string                `validate:"required"`
	Name              string                `validate:"required,max=64"`
	Box60Rates        []*UpdateShippingRate `validate:"required,dive,required"`
	Box60Frozen       int64                 `validate:"min=0,lt=10000000000"`
	Box80Rates        []*UpdateShippingRate `validate:"required,dive,required"`
	Box80Frozen       int64                 `validate:"min=0,lt=10000000000"`
	Box100Rates       []*UpdateShippingRate `validate:"required,dive,required"`
	Box100Frozen      int64                 `validate:"min=0,lt=10000000000"`
	HasFreeShipping   bool                  `validate:""`
	FreeShippingRates int64                 `validate:"min=0,lt=10000000000"`
}

type UpdateShippingRate struct {
	Name            string  `validate:"required"`
	Price           int64   `validate:"min=0,lt=10000000000"`
	PrefectureCodes []int32 `validate:"required"`
}

type UpdateShippingInUseInput struct {
	ShopID     string `validate:"required"`
	ShippingID string `validate:"required"`
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
	Price           int64   `validate:"min=0,lt=10000000000"`
	PrefectureCodes []int32 `validate:"required"`
}

type DeleteShippingInput struct {
	ShippingID string `validate:"required"`
}

/**
 * Deprecated: Shop - 店舗
 */
type ListShopsInput struct {
	CoordinatorIDs []string `validate:""`
	ProducerIDs    []string `validate:""`
	Limit          int64    `validate:"required_without=NoLimit,min=0,max=200"`
	Offset         int64    `validate:"min=0"`
	NoLimit        bool     `validate:""`
}

type ListShopProducersInput struct {
	ShopID  string `validate:""`
	Limit   int64  `validate:"required_without=NoLimit,min=0,max=200"`
	Offset  int64  `validate:"min=0"`
	NoLimit bool   `validate:""`
}

type MultiGetShopsInput struct {
	ShopIDs []string `validate:"dive,required"`
}

type GetShopInput struct {
	ShopID string `validate:"required"`
}

type GetShopByCoordinatorIDInput struct {
	CoordinatorID string `validate:"required"`
}

type CreateShopInput struct {
	CoordinatorID  string         `validate:"required"`
	Name           string         `validate:"required,max=64"`
	ProductTypeIDs []string       `validate:"dive,required"`
	BusinessDays   []time.Weekday `validate:"max=7,unique"`
}

type UpdateShopInput struct {
	ShopID         string         `validate:"required"`
	Name           string         `validate:"required,max=64"`
	ProductTypeIDs []string       `validate:"dive,required"`
	BusinessDays   []time.Weekday `validate:"max=7,unique"`
}

type DeleteShopInput struct {
	ShopID string `validate:"required"`
}

type RelateShopProducerInput struct {
	ShopID     string `validate:"required"`
	ProducerID string `validate:"required"`
}

type UnrelateShopProducerInput struct {
	ShopID     string `validate:"required"`
	ProducerID string `validate:"required"`
}

/**
 * Spot - スポット
 */
type ListSpotsInput struct {
	Name            string   `validate:"max=64"`
	TypeIDs         []string `validate:""`
	UserID          string   `validate:""`
	ExcludeApproved bool     `validate:""`
	ExcludeDisabled bool     `validate:""`
	Limit           int64    `validate:"required_without=NoLimit,min=0,max=200"`
	Offset          int64    `validate:"min=0"`
	NoLimit         bool     `validate:""`
}

type ListSpotsByGeolocationInput struct {
	TypeIDs         []string `validate:""`
	Latitude        float64  `validate:"min=-90,max=90"`
	Longitude       float64  `validate:"min=-180,max=180"`
	Radius          int64    `validate:"min=0"`
	ExcludeDisabled bool     `validate:""`
}

type GetSpotInput struct {
	SpotID string `validate:"required"`
}

type CreateSpotByUserInput struct {
	TypeID       string  `validate:"required"`
	UserID       string  `validate:"required"`
	Name         string  `validate:"required,max=64"`
	Description  string  `validate:"required,max=2000"`
	ThumbnailURL string  `validate:"omitempty,url"`
	Longitude    float64 `validate:"min=-180,max=180"`
	Latitude     float64 `validate:"min=-90,max=90"`
}

type CreateSpotByAdminInput struct {
	TypeID       string  `validate:"required"`
	AdminID      string  `validate:"required"`
	Name         string  `validate:"required,max=64"`
	Description  string  `validate:"required,max=2000"`
	ThumbnailURL string  `validate:"omitempty,url"`
	Longitude    float64 `validate:"min=-180,max=180"`
	Latitude     float64 `validate:"min=-90,max=90"`
}

type UpdateSpotInput struct {
	TypeID       string  `validate:"required"`
	SpotID       string  `validate:"required"`
	Name         string  `validate:"required,max=64"`
	Description  string  `validate:"required,max=2000"`
	ThumbnailURL string  `validate:"omitempty,url"`
	Longitude    float64 `validate:"min=-180,max=180"`
	Latitude     float64 `validate:"min=-90,max=90"`
}

type ApproveSpotInput struct {
	SpotID   string `validate:"required"`
	AdminID  string `validate:"required"`
	Approved bool   `validate:""`
}

type DeleteSpotInput struct {
	SpotID string `validate:"required"`
}

/**
 * SpotType - スポット種別
 */
type ListSpotTypesInput struct {
	Name   string `validate:"max=32"`
	Limit  int64  `validate:"required,max=200"`
	Offset int64  `validate:"min=0"`
}

type MultiGetSpotTypesInput struct {
	SpotTypeIDs []string `validate:"dive,required"`
}

type GetSpotTypeInput struct {
	SpotTypeID string `validate:"required"`
}

type CreateSpotTypeInput struct {
	Name string `validate:"required,max=32"`
}

type UpdateSpotTypeInput struct {
	SpotTypeID string `validate:"required"`
	Name       string `validate:"required,max=32"`
}

type DeleteSpotTypeInput struct {
	SpotTypeID string `validate:"required"`
}
