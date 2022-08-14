package entity

const (
	PushIDContact = "contact" // お問い合わせ受信
)

// PushConfig - プッシュ通知作成設定
type PushConfig struct {
	PushID string            `json:"pushId"` // プッシュテンプレートID
	Data   map[string]string `json:"data"`   // 動的な設定
}
