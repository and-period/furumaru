package types

// OrderType - 注文種別
type OrderType int32

const (
	OrderTypeUnknown    OrderType = 0
	OrderTypeProduct    OrderType = 1 // 商品
	OrderTypeExperience OrderType = 2 // 体験
)

// OrderStatus - 注文ステータス
type OrderStatus int32

const (
	OrderStatusUnknown   OrderStatus = 0
	OrderStatusUnpaid    OrderStatus = 1 // 支払い待ち
	OrderStatusWaiting   OrderStatus = 2 // 受注待ち
	OrderStatusPreparing OrderStatus = 3 // 発送準備中
	OrderStatusShipped   OrderStatus = 4 // 発送完了
	OrderStatusCompleted OrderStatus = 5 // 完了
	OrderStatusCanceled  OrderStatus = 6 // キャンセル
	OrderStatusRefunded  OrderStatus = 7 // 返金
	OrderStatusFailed    OrderStatus = 8 // 失敗
)

// OrderShippingType - 発送方法
type OrderShippingType int32

const (
	OrderShippingTypeUnknown  OrderShippingType = 0
	OrderShippingTypeNone     OrderShippingType = 1 // 発送なし
	OrderShippingTypeStandard OrderShippingType = 2 // 通常配送
	OrderShippingTypePickup   OrderShippingType = 3 // 店舗受取
)

// Order - 注文履歴情報
type Order struct {
	ID              string              `json:"id"`              // 注文履歴ID
	UserID          string              `json:"userId"`          // ユーザーID
	CoordinatorID   string              `json:"coordinatorId"`   // コーディネータID
	PromotionID     string              `json:"promotionId"`     // プロモーションID
	ManagementID    int64               `json:"managementId"`    // 注文管理用ID
	ShippingMessage string              `json:"shippingMessage"` // 発送連絡時のメッセージ
	Type            OrderType           `json:"type"`            // 注文種別
	Status          OrderStatus         `json:"status"`          // 注文ステータス
	Metadata        *OrderMetadata      `json:"metadata"`        // 注文付加情報
	Payment         *OrderPayment       `json:"payment"`         // 支払い情報
	Refund          *OrderRefund        `json:"refund"`          // 注文キャンセル情報
	Fulfillments    []*OrderFulfillment `json:"fulfillments"`    // 配送情報一覧
	Items           []*OrderItem        `json:"items"`           // 注文商品一覧
	Experience      *OrderExperience    `json:"experience"`      // 注文体験情報
	CreatedAt       int64               `json:"createdAt"`       // 登録日時
	UpdatedAt       int64               `json:"updatedAt"`       // 更新日時
	CompletedAt     int64               `json:"completedAt"`     // 対応完了日時
}

// OrderMetadata - 注文付加情報
type OrderMetadata struct {
	OrderRequest   string `json:"orderRequest"`   // 要望・質問など自由入力
	PickupAt       int64  `json:"pickupAt"`       // 受け取り日時
	PickupLocation string `json:"pickupLocation"` // 受け取り場所
}

// OrderItem - 注文商品情報
type OrderItem struct {
	FulfillmentID string `json:"fulfillmentId"` // 配送情報ID
	ProductID     string `json:"productId"`     // 商品ID
	Price         int64  `json:"price"`         // 購入価格(税込)
	Quantity      int64  `json:"quantity"`      // 購入数量
}

// OrderExperience - 注文体験情報
type OrderExperience struct {
	ExperienceID          string                  `json:"experienceId"`          // 体験ID
	AdultCount            int64                   `json:"adultCount"`            // 大人人数
	AdultPrice            int64                   `json:"adultPrice"`            // 大人価格
	JuniorHighSchoolCount int64                   `json:"juniorHighSchoolCount"` // 中学生人数
	JuniorHighSchoolPrice int64                   `json:"juniorHighSchoolPrice"` // 中学生価格
	ElementarySchoolCount int64                   `json:"elementarySchoolCount"` // 小学生人数
	ElementarySchoolPrice int64                   `json:"elementarySchoolPrice"` // 小学生価格
	PreschoolCount        int64                   `json:"preschoolCount"`        // 幼児人数
	PreschoolPrice        int64                   `json:"preschoolPrice"`        // 幼児価格
	SeniorCount           int64                   `json:"seniorCount"`           // シニア人数
	SeniorPrice           int64                   `json:"seniorPrice"`           // シニア価格
	Remarks               *OrderExperienceRemarks `json:"remarks"`               // 備考
}

// OrderExperienceRemarks - 体験希望情報
type OrderExperienceRemarks struct {
	Transportation string `json:"transportation"` // 交通手段
	RequestedDate  string `json:"requestedDate"`  // 体験希望日
	RequestedTime  string `json:"requestedTime"`  // 体験希望時間
}

type DraftOrderRequest struct {
	ShippingMessage string `json:"shippingMessage" validate:"omitempty,max=2000"` // 発送連絡時のメッセージ
}

type CompleteOrderRequest struct {
	ShippingMessage string `json:"shippingMessage" validate:"omitempty,max=2000"` // 発送連絡時のメッセージ
}

type RefundOrderRequest struct {
	Description string `json:"description" validate:"required,max=2000"` // 返金理由
}

type UpdateOrderFulfillmentRequest struct {
	ShippingCarrier int32  `json:"shippingCarrier" validate:"required"` // 配送会社
	TrackingNumber  string `json:"trackingNumber" validate:"required"`  // 伝票番号
}

type ExportOrdersRequest struct {
	ShippingCarrier       int32 `json:"shippingCarrier" validate:"required"` // 配送会社
	CharacterEncodingType int32 `json:"characterEncodingType" validate:""`   // 文字コード種別
}

type OrderResponse struct {
	Order       *Order       `json:"order"`       // 注文履歴情報
	User        *User        `json:"user"`        // 購入者情報
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Promotion   *Promotion   `json:"promotion"`   // プロモーション情報
	Products    []*Product   `json:"products"`    // 商品一覧
	Experience  *Experience  `json:"experience"`  // 体験情報
}

type OrdersResponse struct {
	Orders       []*Order       `json:"orders"`       // 注文履歴一覧
	Users        []*User        `json:"users"`        // 購入者一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Promotions   []*Promotion   `json:"promotions"`   // プロモーション一覧
	Total        int64          `json:"total"`        // 注文履歴合計数
}
