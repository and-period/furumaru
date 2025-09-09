package types

// お問い合わせ既読管理
type ContactRead struct {
	ID        string `json:"id"`        // お問い合わせ既読管理ID
	ContactID string `json:"contactId"` // お問い合わせID
	UserID    string `json:"userId"`    // 既読ユーザーID
	UserType  int32  `json:"userType"`  // 既読ユーザータイプ
	Read      bool   `json:"read"`      // 既読フラグ
	CreatedAt int64  `json:"createdAt"` // 登録日時
	UpdatedAt int64  `json:"updatedAt"` // 更新日時
}

type CreateContactReadRequest struct {
	ContactID string `json:"contactId" validate:"required"`   // お問い合わせID
	UserID    string `json:"userId" validate:"required"`      // 送信者ID
	UserType  int32  `json:"userType" validate:"min=0,max=3"` // 送信者種別(不明:0, admin:1, uer:2, guest:3)
}

type ContactReadResponse struct {
	*ContactRead
}

type ContactReadsResponse struct {
	ContactReads []*ContactRead `json:"contactReads"` // お問い合わせ既読管理一覧
}
