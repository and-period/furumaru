package request

type CreateContactRequest struct {
	Title       string `json:"title,omitempty"`       // 件名
	Content     string `json:"content,omitempty"`     // 内容
	Username    string `json:"username,omitempty"`    // 氏名
	Email       string `json:"email,omitempty"`       // メールアドレス
	PhoneNumber string `json:"phoneNumber,omitempty"` // 電話番号
}
