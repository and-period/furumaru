package response

// AuthUser - 認証済みユーザー情報
type AuthUser struct {
	ID                  string `json:"id"`                  // ユーザーID
	Username            string `json:"username"`            // ユーザー名 (表示名)
	AccountID           string `json:"accountId"`           // ユーザー名 (検索用)
	ThumbnailURL        string `json:"thumbnailUrl"`        // サムネイルURL
	Lastname            string `json:"lastname"`            // 姓
	Firstname           string `json:"firstname"`           // 名
	LastnameKana        string `json:"lastnameKana"`        // 姓（かな）
	FirstnameKana       string `json:"firstnameKana"`       // 名（かな）
	Email               string `json:"email"`               // メールアドレス
	NotificationEnabled bool   `json:"notificationEnabled"` // 通知の有効化設定
}

type AuthUserResponse struct {
	*AuthUser
}

type CreateAuthUserResponse struct {
	ID string `json:"id"` // ユーザーID
}
