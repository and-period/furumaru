package entity

import "time"

type CartActionLogType int32

const (
	CartActionLogTypeUnknown        CartActionLogType = 0
	CartActionLogTypeAddCartItem    CartActionLogType = 1 // 商品追加
	CartActionLogTypeRemoveCartItem CartActionLogType = 2 // 商品削除
)

// CartActionLog - カート操作ログ
type CartActionLog struct {
	SessionID string            `gorm:"primaryKey;<-:create"` // セッションID
	CreatedAt time.Time         `gorm:"primaryKey;<-:create"` // 登録日時
	Type      CartActionLogType `gorm:"primaryKey;<-:create"` // カート操作種別
	UserID    string            `gorm:"default:null"`         // ユーザーID
	UserAgent string            `gorm:""`                     // ユーザーエージェント
	ClientIP  string            `gorm:""`                     // 接続元IPアドレス
	ProductID string            `gorm:"default:null"`         // 商品ID
	UpdatedAt time.Time         `gorm:""`                     // 更新日時
}

type CartActionLogs []*CartActionLog

type AddCartItemActionLogParams struct {
	SessionID string
	UserID    string
	UserAgent string
	ClientIP  string
	ProductID string
}

type RemoveCartItemActionLogParams struct {
	SessionID string
	UserID    string
	UserAgent string
	ClientIP  string
	ProductID string
}

func NewAddCartItemActionLog(params *AddCartItemActionLogParams) *CartActionLog {
	return &CartActionLog{
		SessionID: params.SessionID,
		Type:      CartActionLogTypeAddCartItem,
		UserID:    params.UserID,
		UserAgent: params.UserAgent,
		ClientIP:  params.ClientIP,
		ProductID: params.ProductID,
	}
}

func NewRemoveCartItemActionLog(params *RemoveCartItemActionLogParams) *CartActionLog {
	return &CartActionLog{
		SessionID: params.SessionID,
		Type:      CartActionLogTypeRemoveCartItem,
		UserID:    params.UserID,
		UserAgent: params.UserAgent,
		ClientIP:  params.ClientIP,
		ProductID: params.ProductID,
	}
}
