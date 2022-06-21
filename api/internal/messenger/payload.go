package messenger

import "github.com/and-period/furumaru/api/internal/messenger/entity"

type UserType int32

const (
	UserTypeNone          UserType = 0
	UserTypeUser          UserType = 1 // 購入者
	UserTypeAdministrator UserType = 2 // システム管理者
	UserTypeCoordinator   UserType = 3 // 仲介者
	UserTypeProducer      UserType = 4 // 生産者
)

type WorkerPayload struct {
	UserType UserType           `json:"userType,omitempty"` // 送信先ユーザー種別
	UserIDs  []string           `json:"userIds,omitempty"`  // 送信先ユーザー一覧
	Email    *entity.MailConfig `json:"email,omitempty"`    // メール送信設定
}
