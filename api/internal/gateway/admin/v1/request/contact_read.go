package request

type CreateContactReadRequest struct {
	ContactID string `json:"contactId,omitempty"` // お問い合わせID
	UserID    string `json:"userId"`              // 送信者ID
	UserType  int32  `json:"userType,omitempty"`  // 送信者種別(不明:0, admin:1, uer:2, guest:3)
}
