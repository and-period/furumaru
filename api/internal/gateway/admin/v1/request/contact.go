package request

type CreateContactRequest struct {
	Title       string `json:"title,omitempty"`       // お問い合わせ件名
	Content     string `json:"content,omitempty"`     // お問い合わせ内容
	Username    string `json:"username,omitempty"`    // 氏名
	Email       string `json:"email,omitempty"`       // メールアドレス
	PhoneNumber string `json:"phoneNumber,omitempty"` // 電話番号
	Note        string `json:"note"`                  // 対応者メモ
}
