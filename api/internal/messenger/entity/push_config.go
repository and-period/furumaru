package entity

// PushTemplateID - プッシュテンプレートID
type PushTemplateID string

const (
	PushTemplateIDContact PushTemplateID = "contact" // お問い合わせ受信
)

// PushConfig - プッシュ通知作成設定
type PushConfig struct {
	TemplateID PushTemplateID    `json:"pushId"` // プッシュテンプレートID
	Data       map[string]string `json:"data"`   // 動的な設定
}
