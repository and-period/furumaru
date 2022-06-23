package messenger

import "github.com/and-period/furumaru/api/internal/messenger/entity"

// EventType - Worker実行種別
type EventType string

const (
	EventTypeRegisterAdmin EventType = "register-admin" // 管理者登録通知
)

// UserType - 通知先ユーザー種別
type UserType int32

const (
	UserTypeNone          UserType = 0
	UserTypeUser          UserType = 1 // 購入者
	UserTypeAdministrator UserType = 2 // システム管理者
	UserTypeCoordinator   UserType = 3 // 仲介者
	UserTypeProducer      UserType = 4 // 生産者
)

type WorkerPayload struct {
	EventType EventType          `json:"eventType,omitempty"` // Worker実行種別
	UserType  UserType           `json:"userType,omitempty"`  // 送信先ユーザー種別
	UserIDs   []string           `json:"userIds,omitempty"`   // 送信先ユーザー一覧
	Email     *entity.MailConfig `json:"email,omitempty"`     // メール送信設定
}
