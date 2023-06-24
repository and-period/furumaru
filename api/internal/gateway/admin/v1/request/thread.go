package request

type CreateThreadRequest struct {
	ContactID string `json:"contactId,omitempty"` // お問い合わせID
	UserID    string `json:"userId"`              // 送信者ID
	Content   string `json:"content,omitempty"`   // 内容
	UserType  int32  `json:"userType,omitempty"`  // 送信者種別(不明:0, admin:1, uer:2, guest:3)
}

type UpdateThreadRequest struct {
	ThreadID string `json:"threadId,omitempty"` // お問い合わせID
	UserID   string `json:"userId"`             // 送信者ID
	Content  string `json:"content,omitempty"`  // 内容
	UserType int32  `json:"userType,omitempty"` // 送信者種別(不明:0, admin:1, uer:2, guest:3)
}
