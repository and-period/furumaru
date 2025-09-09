package types

// Message - メッセージ情報
type Message struct {
	ID         string `json:"id"`         // メッセージID
	Type       int32  `json:"type"`       // メッセージ種別
	Title      string `json:"title"`      // メッセージ件名
	Body       string `json:"body"`       // メッセージ内容
	Link       string `json:"link"`       // 遷移先リンク
	Read       bool   `json:"read"`       // 既読フラグ
	ReceivedAt int64  `json:"receivedAt"` // 受信日時
	CreatedAt  int64  `json:"createdAt"`  // 登録日時
	UpdatedAt  int64  `json:"updatedAt"`  // 更新日時
}

type MessageResponse struct {
	Message *Message `json:"message"` // メッセージ情報
}

type MessagesResponse struct {
	Messages []*Message `json:"messages"` // メッセージ一覧
	Total    int64      `json:"total"`    // 合計数
}
