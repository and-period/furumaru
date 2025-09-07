package request

type CreateThreadRequest struct {
	ContactID string `json:"contactId" validate:"required"`        // お問い合わせID
	UserID    string `json:"userId" validate:"required"`           // 送信者ID
	Content   string `json:"content" validate:"required,max=1000"` // 内容
	UserType  int32  `json:"userType" validate:"min=0,max=3"`      // 送信者種別(不明:0, admin:1, uer:2, guest:3)
}

type UpdateThreadRequest struct {
	ThreadID string `json:"threadId" validate:"required"`         // お問い合わせID
	UserID   string `json:"userId" validate:"required"`           // 送信者ID
	Content  string `json:"content" validate:"required,max=1000"` // 内容
	UserType int32  `json:"userType" validate:"min=0,max=3"`      // 送信者種別(不明:0, admin:1, uer:2, guest:3)
}
