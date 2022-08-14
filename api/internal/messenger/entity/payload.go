package entity

// EventType - Worker実行種別
type EventType int32

const (
	EventTypeUnknown            EventType = 0
	EventTypeRegisterAdmin      EventType = 1 // 管理者登録通知
	EventTypeResetAdminPassword EventType = 2 // 管理者パスワードリセット通知
	EventTypeReceivedContact    EventType = 3 // お問い合わせ受領通知
	EventTypeNotification       EventType = 4 // お知らせ発行通知
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
	UserTypeGuest         UserType = 6 // 未登録ユーザー
)

// WorkerPayload - Worker実行内容
type WorkerPayload struct {
	QueueID   string         `json:"queueId"`           // メッセージキューID(重複実行抑止用)
	EventType EventType      `json:"eventType"`         // Worker実行種別
	UserType  UserType       `json:"userType"`          // 送信先ユーザー種別
	UserIDs   []string       `json:"userIds"`           // 送信先ユーザー一覧
	Guest     *Guest         `json:"guest,omitempty"`   // 未登録ユーザー情報
	Email     *MailConfig    `json:"email,omitempty"`   // メール送信設定
	Push      *PushConfig    `json:"push,omitempty"`    // プッシュ通知設定
	Message   *MessageConfig `json:"message,omitempty"` // メッセージ作成設定
	Report    *ReportConfig  `json:"report,omitempty"`  // システムレポート送信設定
}

type Guest struct {
	Name  string `json:"name"`  // 氏名
	Email string `json:"email"` // メールアドレス
}
