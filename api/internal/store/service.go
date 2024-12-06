//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package store

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Service interface {
	// カテゴリ
	ListCategories(ctx context.Context, in *ListCategoriesInput) (entity.Categories, int64, error)  // 一覧取得
	MultiGetCategories(ctx context.Context, in *MultiGetCategoriesInput) (entity.Categories, error) // 一覧取得(ID指定)
	GetCategory(ctx context.Context, in *GetCategoryInput) (*entity.Category, error)                // 1県取得
	CreateCategory(ctx context.Context, in *CreateCategoryInput) (*entity.Category, error)          // 登録
	UpdateCategory(ctx context.Context, in *UpdateCategoryInput) error                              // 更新
	DeleteCategory(ctx context.Context, in *DeleteCategoryInput) error                              // 削除
	// 品目
	ListProductTypes(ctx context.Context, in *ListProductTypesInput) (entity.ProductTypes, int64, error)  // 一覧取得
	MultiGetProductTypes(ctx context.Context, in *MultiGetProductTypesInput) (entity.ProductTypes, error) // 一覧取得(ID指定)
	GetProductType(ctx context.Context, in *GetProductTypeInput) (*entity.ProductType, error)             // １件取得
	CreateProductType(ctx context.Context, in *CreateProductTypeInput) (*entity.ProductType, error)       // 登録
	UpdateProductType(ctx context.Context, in *UpdateProductTypeInput) error                              // 更新
	DeleteProductType(ctx context.Context, in *DeleteProductTypeInput) error                              // 削除
	// 商品タグ
	ListProductTags(ctx context.Context, in *ListProductTagsInput) (entity.ProductTags, int64, error)  // 一覧取得
	MultiGetProductTags(ctx context.Context, in *MultiGetProductTagsInput) (entity.ProductTags, error) // 一覧取得(ID指定)
	GetProductTag(ctx context.Context, in *GetProductTagInput) (*entity.ProductTag, error)             // １件取得
	CreateProductTag(ctx context.Context, in *CreateProductTagInput) (*entity.ProductTag, error)       // 登録
	UpdateProductTag(ctx context.Context, in *UpdateProductTagInput) error                             // 更新
	DeleteProductTag(ctx context.Context, in *DeleteProductTagInput) error                             // 削除
	// 配送設定
	ListShippingsByCoordinatorIDs(ctx context.Context, in *ListShippingsByCoordinatorIDsInput) (entity.Shippings, error) // 一覧取得(コーディネータID指定)
	MultiGetShippingsByRevision(ctx context.Context, in *MultiGetShippingsByRevisionInput) (entity.Shippings, error)     // 一覧取得(変更履歴指定)
	GetDefaultShipping(ctx context.Context, in *GetDefaultShippingInput) (*entity.Shipping, error)                       // １件取得(デフォルト設定)
	GetShippingByCoordinatorID(ctx context.Context, in *GetShippingByCoordinatorIDInput) (*entity.Shipping, error)       // １件取得(コーディネータID設定)
	UpdateDefaultShipping(ctx context.Context, in *UpdateDefaultShippingInput) error                                     // 登録または更新(デフォルト設定)
	UpsertShipping(ctx context.Context, in *UpsertShippingInput) error                                                   // 登録または更新(コーディネータごとの設定)
	// 商品
	ListProducts(ctx context.Context, in *ListProductsInput) (entity.Products, int64, error)                      // 一覧取得
	MultiGetProducts(ctx context.Context, in *MultiGetProductsInput) (entity.Products, error)                     // 一覧取得(ID指定)
	MultiGetProductsByRevision(ctx context.Context, in *MultiGetProductsByRevisionInput) (entity.Products, error) // 一覧取得(変更履歴ID指定)
	GetProduct(ctx context.Context, in *GetProductInput) (*entity.Product, error)                                 // １件取得
	CreateProduct(ctx context.Context, in *CreateProductInput) (*entity.Product, error)                           // 登録
	UpdateProduct(ctx context.Context, in *UpdateProductInput) error                                              // 更新
	DeleteProduct(ctx context.Context, in *DeleteProductInput) error                                              // 削除
	// プロモーション
	ListPromotions(ctx context.Context, in *ListPromotionsInput) (entity.Promotions, int64, error)  // 一覧取得
	MultiGetPromotions(ctx context.Context, in *MultiGetPromotionsInput) (entity.Promotions, error) // 一覧取得(ID指定)
	GetPromotion(ctx context.Context, in *GetPromotionInput) (*entity.Promotion, error)             // 取得
	GetPromotionByCode(ctx context.Context, in *GetPromotionByCodeInput) (*entity.Promotion, error) // 取得(コード指定)
	CreatePromotion(ctx context.Context, in *CreatePromotionInput) (*entity.Promotion, error)       // 登録
	UpdatePromotion(ctx context.Context, in *UpdatePromotionInput) error                            // 更新
	DeletePromotion(ctx context.Context, in *DeletePromotionInput) error                            // 削除
	// マルシェ開催スケジュール
	ListSchedules(ctx context.Context, in *ListSchedulesInput) (entity.Schedules, int64, error)  // 一覧取得
	MultiGetSchedules(ctx context.Context, in *MultiGetSchedulesInput) (entity.Schedules, error) // 一覧取得(ID指定)
	GetSchedule(ctx context.Context, in *GetScheduleInput) (*entity.Schedule, error)             // １件取得
	CreateSchedule(ctx context.Context, in *CreateScheduleInput) (*entity.Schedule, error)       // 登録
	UpdateSchedule(ctx context.Context, in *UpdateScheduleInput) error                           // 更新
	DeleteSchedule(ctx context.Context, in *DeleteScheduleInput) error                           // 削除
	ApproveSchedule(ctx context.Context, in *ApproveScheduleInput) error                         // 承認
	PublishSchedule(ctx context.Context, in *PublishScheduleInput) error                         // 公開
	// マルシェタイムテーブル
	ListLives(ctx context.Context, in *ListLivesInput) (entity.Lives, int64, error) // 一覧取得
	GetLive(ctx context.Context, in *GetLiveInput) (*entity.Live, error)            // 取得
	CreateLive(ctx context.Context, in *CreateLiveInput) (*entity.Live, error)      // 登録
	UpdateLive(ctx context.Context, in *UpdateLiveInput) error                      // 更新
	DeleteLive(ctx context.Context, in *DeleteLiveInput) error                      // 削除
	// 注文履歴
	ListOrders(ctx context.Context, in *ListOrdersInput) (entity.Orders, int64, error)                                             // 一覧取得
	ListOrderUserIDs(ctx context.Context, in *ListOrderUserIDsInput) ([]string, int64, error)                                      // 注文したユーザーID一覧取得
	GetOrder(ctx context.Context, in *GetOrderInput) (*entity.Order, error)                                                        // １件取得
	GetOrderByTransactionID(ctx context.Context, in *GetOrderByTransactionIDInput) (*entity.Order, error)                          // １件取得(決済トランザクションID指定)
	CaptureOrder(ctx context.Context, in *CaptureOrderInput) error                                                                 // 注文確定
	DraftOrder(ctx context.Context, in *DraftOrderInput) error                                                                     // 注文の下書き保存
	CompleteOrder(ctx context.Context, in *CompleteOrderInput) error                                                               // 注文対応完了
	CancelOrder(ctx context.Context, in *CancelOrderInput) error                                                                   // 注文キャンセル
	RefundOrder(ctx context.Context, in *RefundOrderInput) error                                                                   // 注文返金依頼
	UpdateOrderFulfillment(ctx context.Context, in *UpdateOrderFulfillmentInput) error                                             // 注文配送情報更新
	AggregateOrders(ctx context.Context, in *AggregateOrdersInput) (entity.AggregatedOrders, error)                                // 集計結果一覧取得
	AggregateOrdersByPromotion(ctx context.Context, in *AggregateOrdersByPromotionInput) (entity.AggregatedOrderPromotions, error) // プロモーション利用履歴集計結果一覧取得
	ExportOrders(ctx context.Context, in *ExportOrdersInput) ([]byte, error)                                                       // 注文履歴一覧CSV出力
	// 買い物かご
	GetCart(ctx context.Context, in *GetCartInput) (*entity.Cart, error)                                // 取得
	CalcCart(ctx context.Context, in *CalcCartInput) (*entity.Cart, *entity.OrderPaymentSummary, error) // 購入前の支払い情報取得
	AddCartItem(ctx context.Context, in *AddCartItemInput) error                                        // 商品を追加
	RemoveCartItem(ctx context.Context, in *RemoveCartItemInput) error                                  // 商品を削除
	// 体験種別
	ListExperienceTypes(ctx context.Context, in *ListExperienceTypesInput) (entity.ExperienceTypes, int64, error)  // 一覧取得
	MultiGetExperienceTypes(ctx context.Context, in *MultiGetExperienceTypesInput) (entity.ExperienceTypes, error) // 一覧取得（ID指定）
	GetExperienceType(ctx context.Context, in *GetExperienceTypeInput) (*entity.ExperienceType, error)             // １件取得
	CreateExperienceType(ctx context.Context, in *CreateExperienceTypeInput) (*entity.ExperienceType, error)       // 登録
	UpdateExperienceType(ctx context.Context, in *UpdateExperienceTypeInput) error                                 // 更新
	DeleteExperienceType(ctx context.Context, in *DeleteExperienceTypeInput) error                                 // 削除
	// 体験
	ListExperiences(ctx context.Context, in *ListExperiencesInput) (entity.Experiences, int64, error)                      // 一覧取得
	MultiGetExperiences(ctx context.Context, in *MultiGetExperiencesInput) (entity.Experiences, error)                     // 一覧取得（ID指定）
	MultiGetExperiencesByRevision(ctx context.Context, in *MultiGetExperiencesByRevisionInput) (entity.Experiences, error) // 一覧取得(変更履歴ID指定)
	GetExperience(ctx context.Context, in *GetExperienceInput) (*entity.Experience, error)                                 // １件取得
	CreateExperience(ctx context.Context, in *CreateExperienceInput) (*entity.Experience, error)                           // 登録
	UpdateExperience(ctx context.Context, in *UpdateExperienceInput) error                                                 // 更新
	DeleteExperience(ctx context.Context, in *DeleteExperienceInput) error                                                 // 削除
	// 購入処理
	GetCheckoutState(ctx context.Context, in *GetCheckoutStateInput) (string, entity.PaymentStatus, error) // 支払い状態取得
	CheckoutCreditCard(ctx context.Context, in *CheckoutCreditCardInput) (string, error)                   // 支払い申請（クレジットカード）
	CheckoutPayPay(ctx context.Context, in *CheckoutPayPayInput) (string, error)                           // 支払い申請（PayPay）
	CheckoutLinePay(ctx context.Context, in *CheckoutLinePayInput) (string, error)                         // 支払い申請（LINE Pay）
	CheckoutMerpay(ctx context.Context, in *CheckoutMerpayInput) (string, error)                           // 支払い申請（メルペイ）
	CheckoutRakutenPay(ctx context.Context, in *CheckoutRakutenPayInput) (string, error)                   // 支払い申請（楽天ペイ）
	CheckoutAUPay(ctx context.Context, in *CheckoutAUPayInput) (string, error)                             // 支払い申請（au PAY）
	CheckoutFree(ctx context.Context, in *CheckoutFreeInput) (string, error)                               // 支払い申請（無料）
	NotifyPaymentCompleted(ctx context.Context, in *NotifyPaymentCompletedInput) error                     // 支払い通知
	NotifyPaymentRefunded(ctx context.Context, in *NotifyPaymentRefundedInput) error                       // 返金通知
	// 決済システム
	MultiGetPaymentSystems(ctx context.Context, in *MultiGetPaymentSystemsInput) (entity.PaymentSystems, error) // 一覧取得(種別指定)
	GetPaymentSystem(ctx context.Context, in *GetPaymentSystemInput) (*entity.PaymentSystem, error)             // １件取得
	UpdatePaymentSystem(ctx context.Context, in *UpdatePaymentStatusInput) error                                // 更新
	// 郵便番号
	SearchPostalCode(ctx context.Context, in *SearchPostalCodeInput) (*entity.PostalCode, error) // 検索
	// スポット種別
	ListSpotTypes(ctx context.Context, in *ListSpotTypesInput) (entity.SpotTypes, int64, error)  // 一覧取得
	MultiGetSpotTypes(ctx context.Context, in *MultiGetSpotTypesInput) (entity.SpotTypes, error) // 一覧取得(ID指定)
	GetSpotType(ctx context.Context, in *GetSpotTypeInput) (*entity.SpotType, error)             // １件取得
	CreateSpotType(ctx context.Context, in *CreateSpotTypeInput) (*entity.SpotType, error)       // 登録
	UpdateSpotType(ctx context.Context, in *UpdateSpotTypeInput) error                           // 更新
	DeleteSpotType(ctx context.Context, in *DeleteSpotTypeInput) error                           // 削除
	// スポット
	ListSpots(ctx context.Context, in *ListSpotsInput) (entity.Spots, int64, error)                    // 一覧取得
	ListSpotsByGeolocation(ctx context.Context, in *ListSpotsByGeolocationInput) (entity.Spots, error) // 一覧取得（座標指定）
	GetSpot(ctx context.Context, in *GetSpotInput) (*entity.Spot, error)                               // １件取得
	CreateSpotByUser(ctx context.Context, in *CreateSpotByUserInput) (*entity.Spot, error)             // 登録（購入者）
	CreateSpotByAdmin(ctx context.Context, in *CreateSpotByAdminInput) (*entity.Spot, error)           // 登録（管理者）
	UpdateSpot(ctx context.Context, in *UpdateSpotInput) error                                         // 更新
	DeleteSpot(ctx context.Context, in *DeleteSpotInput) error                                         // 削除
	ApproveSpot(ctx context.Context, in *ApproveSpotInput) error                                       // 承認
}
