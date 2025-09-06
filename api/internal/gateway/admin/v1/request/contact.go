package request

type CreateContactRequest struct {
	Title       string `json:"title" binding:"required,max=64"`     // お問い合わせ件名
	CategoryID  string `json:"categoryId" binding:"required"`       // お問い合わせ種別ID
	Content     string `json:"content" binding:"required,max=2000"` // お問い合わせ内容
	Username    string `json:"username" binding:"required,max=32"`  // 氏名
	UserID      string `json:"userId" binding:"omitempty"`          // 問い合わせ作成者ID(null許容)
	Email       string `json:"email" binding:"required,email"`      // メールアドレス
	PhoneNumber string `json:"phoneNumber" binding:"required,e164"` // 電話番号
	ResponderID string `json:"responderId" binding:"omitempty"`     // 対応者ID(null許容)
	Note        string `json:"note" binding:"omitempty,max=500"`    // 対応者メモ
}

// お問い合わせステータス(作成時は不明)
type ContactStatus int32

const (
	ContactStatusUnknown    ContactStatus = iota // 不明
	ContactStatusWaiting                         // 未着手
	ContactStatusInprogress                      // 進行中
	ContactStatusDone                            // 完了
	ContactStatusDiscard                         // 対応不要
)

type UpdateContactRequest struct {
	Title       string        `json:"title" binding:"required,max=64"`     // お問い合わせ件名
	CategoryID  string        `json:"categoryId" binding:"required"`       // お問い合わせ種別ID
	Content     string        `json:"content" binding:"required,max=2000"` // お問い合わせ内容
	Username    string        `json:"username" binding:"required,max=32"`  // 氏名
	UserID      string        `json:"userId" binding:"omitempty"`          // 問い合わせ作成者ID(null許容)
	Email       string        `json:"email" binding:"required,email"`      // メールアドレス
	PhoneNumber string        `json:"phoneNumber" binding:"required,e164"` // 電話番号
	ResponderID string        `json:"responderId" binding:"omitempty"`     // 対応者ID(null許容)
	Note        string        `json:"note" binding:"omitempty,max=500"`    // 対応者メモ
	Status      ContactStatus `json:"status" binding:"min=0,max=4"`        // お問い合わせステータス
}
