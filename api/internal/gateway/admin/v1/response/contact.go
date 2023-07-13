package response

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

type ContactResponse struct {
	Contact  *Contact         `json:"contact"`  // お問い合わせ情報
	Category *ContactCategory `json:"category"` // お問い合わせ種別情報
	Threads  []*Thread        `json:"threads"`  // お問い合わせ会話履歴一覧
}

type ContactsResponse struct {
	Contacts []*Contact `json:"contacts"` // お問い合わせ一覧
	Total    int64      `json:"total"`    // お問い合わせ合計
}
