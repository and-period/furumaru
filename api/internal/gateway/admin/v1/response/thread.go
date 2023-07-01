package response

// お問い合わせ会話履歴
type Thread struct {
	ID        string `json:"id"`        // お問い合わせ会話履歴ID
	ContactID string `json:"contactId"` // お問い合わせID
	Content   string `json:"content"`   // 会話内容
	UserID    string `json:"userId"`    // 送信者ID
	UserType  int32  `json:"userType"`  // 送信者タイプ
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type ThreadResponse struct {
	*Thread
}

type ThreadsResponse struct {
	Threads []*Thread `json:"threads"` // お問い合わせ会話履歴一覧
	Total   int64     `json:"total"`   // 会話履歴合計
}
