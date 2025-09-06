package request

type CreateContactReadRequest struct {
	ContactID string `json:"contactId" binding:"required"`   // お問い合わせID
	UserID    string `json:"userId" binding:"required"`      // 送信者ID
	UserType  int32  `json:"userType" binding:"min=0,max=3"` // 送信者種別(不明:0, admin:1, uer:2, guest:3)
}
