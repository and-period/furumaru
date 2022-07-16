package response

// Contact - お問い合わせ情報
type Contact struct {
	ID          string `json:"id"`          // お問い合わせID
	Title       string `json:"title"`       // 件名
	Content     string `json:"content"`     // 内容
	Username    string `json:"username"`    // 氏名
	Email       string `json:"email"`       // メールアドレス
	PhoneNumber string `json:"phoneNumber"` // 電話番号
	Status      int32  `json:"status"`      // 対応状況
	Priority    int32  `json:"priority"`    // 優先度
	Note        string `json:"note"`        // 対応時メモ
	CreatedAt   int64  `json:"createdAt"`   // 登録日時
	UpdatedAt   int64  `json:"updatedAt"`   // 更新日時
}

type ContactResponse struct {
	*Contact
}

type ContactsResponse struct {
	Contacts []*Contact `json:"contacts"` // お問い合わせ一覧
	Total    int64      `json:"total"`    // 合計数
}
