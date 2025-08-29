package request

type CreateContactRequest struct {
	Title       string `json:"title"`       // お問い合わせ件名
	CategoryID  string `json:"categoryId"`  // お問い合わせ種別ID
	Content     string `json:"content"`     // お問い合わせ内容
	Username    string `json:"username"`    // 氏名
	UserID      string `json:"userId"`      // 問い合わせ作成者ID(null許容)
	Email       string `json:"email"`       // メールアドレス
	PhoneNumber string `json:"phoneNumber"` // 電話番号
	ResponderID string `json:"responderId"` // 対応者ID(null許容)
	Note        string `json:"note"`        // 対応者メモ
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
	Title       string        `json:"title"`       // お問い合わせ件名
	CategoryID  string        `json:"categoryId"`  // お問い合わせ種別ID
	Content     string        `json:"content"`     // お問い合わせ内容
	Username    string        `json:"username"`    // 氏名
	UserID      string        `json:"userId"`      // 問い合わせ作成者ID(null許容)
	Email       string        `json:"email"`       // メールアドレス
	PhoneNumber string        `json:"phoneNumber"` // 電話番号
	ResponderID string        `json:"responderId"` // 対応者ID(null許容)
	Note        string        `json:"note"`        // 対応者メモ
	Status      ContactStatus `json:"status"`      // お問い合わせステータス
}
