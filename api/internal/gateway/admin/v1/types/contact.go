package types

// お問い合わせステータス
type ContactStatus int32

const (
	ContactStatusUnknown    ContactStatus = iota // 不明
	ContactStatusWaiting                         // 未着手
	ContactStatusInprogress                      // 進行中
	ContactStatusDone                            // 完了
	ContactStatusDiscard                         // 対応不要
)

// お問い合わせ情報
type Contact struct {
	ID          string        `json:"id"`          // お問い合わせID
	Title       string        `json:"title"`       // 件名
	CategoryID  string        `json:"categoryId"`  // お問い合わせ種別ID
	Content     string        `json:"content"`     // 内容
	Username    string        `json:"username"`    // 氏名
	UserID      string        `json:"userId"`      // ユーザーID
	Email       string        `json:"email"`       // メールアドレス
	PhoneNumber string        `json:"phoneNumber"` // 電話番号
	Status      ContactStatus `json:"status"`      // お問い合わせステータス
	ResponderID string        `json:"responderId"` // 対応者ID
	Note        string        `json:"note"`        // 対応者メモ
	CreatedAt   int64         `json:"createdAt"`   // 登録日時
	UpdatedAt   int64         `json:"updatedAt"`   // 更新日時
}

type CreateContactRequest struct {
	Title       string `json:"title" validate:"required,max=64"`     // お問い合わせ件名
	CategoryID  string `json:"categoryId" validate:"required"`       // お問い合わせ種別ID
	Content     string `json:"content" validate:"required,max=2000"` // お問い合わせ内容
	Username    string `json:"username" validate:"required,max=32"`  // 氏名
	UserID      string `json:"userId" validate:"omitempty"`          // 問い合わせ作成者ID(null許容)
	Email       string `json:"email" validate:"required,email"`      // メールアドレス
	PhoneNumber string `json:"phoneNumber" validate:"required,e164"` // 電話番号
	ResponderID string `json:"responderId" validate:"omitempty"`     // 対応者ID(null許容)
	Note        string `json:"note" validate:"omitempty,max=500"`    // 対応者メモ
}

type UpdateContactRequest struct {
	Title       string        `json:"title" validate:"required,max=64"`     // お問い合わせ件名
	CategoryID  string        `json:"categoryId" validate:"required"`       // お問い合わせ種別ID
	Content     string        `json:"content" validate:"required,max=2000"` // お問い合わせ内容
	Username    string        `json:"username" validate:"required,max=32"`  // 氏名
	UserID      string        `json:"userId" validate:"omitempty"`          // 問い合わせ作成者ID(null許容)
	Email       string        `json:"email" validate:"required,email"`      // メールアドレス
	PhoneNumber string        `json:"phoneNumber" validate:"required,e164"` // 電話番号
	ResponderID string        `json:"responderId" validate:"omitempty"`     // 対応者ID(null許容)
	Note        string        `json:"note" validate:"omitempty,max=500"`    // 対応者メモ
	Status      ContactStatus `json:"status" validate:"min=0,max=4"`        // お問い合わせステータス
}

type ContactResponse struct {
	Contact   *Contact         `json:"contact"`   // お問い合わせ情報
	Category  *ContactCategory `json:"category"`  // お問い合わせ種別情報
	Threads   []*Thread        `json:"threads"`   // お問い合わせ会話履歴一覧
	User      *User            `json:"user"`      // ユーザー情報
	Responder *Admin           `json:"responder"` // 対応者情報
}

type ContactsResponse struct {
	Contacts   []*Contact         `json:"contacts"`   // お問い合わせ一覧
	Categories []*ContactCategory `json:"categories"` // お問い合わせ種別一覧
	Threads    []*Thread          `json:"threads"`    // お問い合わせ会話履歴一覧
	Users      []*User            `json:"users"`      // ユーザー一覧
	Responders []*Admin           `json:"admins"`     // 管理者一覧
	Total      int64              `json:"total"`      // お問い合わせ合計
}
