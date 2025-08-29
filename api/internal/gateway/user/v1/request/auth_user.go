package request

type CreateAuthUserRequest struct {
	Username             string `json:"username"`             // ユーザー名(表示名)
	AccountID            string `json:"accountId"`            // ユーザーID(検索名)
	Lastname             string `json:"lastname"`             // 姓
	Firstname            string `json:"firstname"`            // 名
	LastnameKana         string `json:"lastnameKana"`         // 姓（かな）
	FirstnameKana        string `json:"firstnameKana"`        // 名（かな）
	Email                string `json:"email"`                // メールアドレス
	PhoneNumber          string `json:"phoneNumber"`          // 電話番号
	Password             string `json:"password"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation"` // パスワード (確認用)
}

type VerifyAuthUserRequest struct {
	ID         string `json:"id"`         // ユーザーID
	VerifyCode string `json:"verifyCode"` // 検証コード
}

type CreateAuthUserWithGoogleRequest struct {
	Code          string `json:"code"`          // 認証コード
	Nonce         string `json:"nonce"`         // セキュア文字列（リプレイアタック対策）
	RedirectURI   string `json:"redirectUri"`   // リダイレクトURI
	Username      string `json:"username"`      // ユーザー名(表示名)
	AccountID     string `json:"accountId"`     // ユーザーID(検索名)
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓（かな）
	FirstnameKana string `json:"firstnameKana"` // 名（かな）
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
}

type CreateAuthUserWithLINERequest struct {
	Code          string `json:"code"`          // 認証コード
	Nonce         string `json:"nonce"`         // セキュア文字列（リプレイアタック対策）
	RedirectURI   string `json:"redirectUri"`   // リダイレクトURI
	Username      string `json:"username"`      // ユーザー名(表示名)
	AccountID     string `json:"accountId"`     // ユーザーID(検索名)
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓（かな）
	FirstnameKana string `json:"firstnameKana"` // 名（かな）
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
}

type UpdateAuthUserEmailRequest struct {
	Email string `json:"email"` // メールアドレス
}

type VerifyAuthUserEmailRequest struct {
	VerifyCode string `json:"verifyCode"` // 検証コード
}

type UpdateAuthUserUsernameRequest struct {
	Username string `json:"username"` // ユーザー名(表示名)
}

type UpdateAuthUserAccountIDRequest struct {
	AccountID string `json:"accountId"` // ユーザーID(検索名)
}

type UpdateAuthUserNotificationRequest struct {
	Enabled bool `json:"enabled"` // 通知の有効化設定
}

type UpdateAuthUserThumbnailRequest struct {
	ThumbnailURL string `json:"thumbnailUrl"` // サムネイルURL
}
