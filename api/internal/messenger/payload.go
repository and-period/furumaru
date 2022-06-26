package messenger

import "github.com/and-period/furumaru/api/internal/messenger/entity"

// EventType - Worker実行種別
type EventType int32

const (
	EventTypeUnknown       EventType = 0
	EventTypeRegisterAdmin EventType = 1 // 管理者登録通知
)

// UserType - 通知先ユーザー種別
type UserType int32

const (
	UserTypeNone          UserType = 0
	UserTypeUser          UserType = 1 // 購入者
	UserTypeAdmin         UserType = 2 // 管理者(システム管理者,仲介者,生産者)
	UserTypeAdministrator UserType = 3 // システム管理者
	UserTypeCoordinator   UserType = 4 // 仲介者
	UserTypeProducer      UserType = 5 // 生産者
)

type WorkerPayload struct {
	EventType EventType          `json:"eventType,omitempty"` // Worker実行種別
	UserType  UserType           `json:"userType,omitempty"`  // 送信先ユーザー種別
	UserIDs   []string           `json:"userIds,omitempty"`   // 送信先ユーザー一覧
	Email     *entity.MailConfig `json:"email,omitempty"`     // メール送信設定
}
