package request

type CreateAuthUserRequest struct {
	Username             string `json:"username" binding:"required,max=32"`                       // ユーザー名(表示名)
	AccountID            string `json:"accountId" binding:"required,max=32,account_id"`           // ユーザーID(検索名)
	Lastname             string `json:"lastname" binding:"required,max=16"`                       // 姓
	Firstname            string `json:"firstname" binding:"required,max=16"`                      // 名
	LastnameKana         string `json:"lastnameKana" binding:"required,max=32,hiragana"`          // 姓（かな）
	FirstnameKana        string `json:"firstnameKana" binding:"required,max=32,hiragana"`         // 名（かな）
	Email                string `json:"email" binding:"required,email"`                           // メールアドレス
	PhoneNumber          string `json:"phoneNumber" binding:"required,e164"`                      // 電話番号
	Password             string `json:"password" binding:"min=8,max=32"`                          // パスワード
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"` // パスワード (確認用)
}

type VerifyAuthUserRequest struct {
	ID         string `json:"id" binding:"required"`         // ユーザーID
	VerifyCode string `json:"verifyCode" binding:"required"` // 検証コード
}

type CreateAuthUserWithGoogleRequest struct {
	Code          string `json:"code" binding:"required"`                                // 認証コード
	Nonce         string `json:"nonce" binding:"required"`                               // セキュア文字列（リプレイアタック対策）
	RedirectURI   string `json:"redirectUri" binding:"required,url"`                     // リダイレクトURI
	Username      string `json:"username" binding:"required,max=32"`                     // ユーザー名(表示名)
	AccountID     string `json:"accountId" binding:"required,min=4,max=32,alphanumeric"` // ユーザーID(検索名)
	Lastname      string `json:"lastname" binding:"required,max=16"`                     // 姓
	Firstname     string `json:"firstname" binding:"required,max=16"`                    // 名
	LastnameKana  string `json:"lastnameKana" binding:"required,max=32,hiragana"`        // 姓（かな）
	FirstnameKana string `json:"firstnameKana" binding:"required,max=32,hiragana"`       // 名（かな）
	PhoneNumber   string `json:"phoneNumber" binding:"required,e164"`                    // 電話番号
}

type CreateAuthUserWithLINERequest struct {
	Code          string `json:"code" binding:"required"`                                // 認証コード
	Nonce         string `json:"nonce" binding:"required"`                               // セキュア文字列（リプレイアタック対策）
	RedirectURI   string `json:"redirectUri" binding:"required,url"`                     // リダイレクトURI
	Username      string `json:"username" binding:"required,max=32"`                     // ユーザー名(表示名)
	AccountID     string `json:"accountId" binding:"required,min=4,max=32,alphanumeric"` // ユーザーID(検索名)
	Lastname      string `json:"lastname" binding:"required,max=16"`                     // 姓
	Firstname     string `json:"firstname" binding:"required,max=16"`                    // 名
	LastnameKana  string `json:"lastnameKana" binding:"required,max=32,hiragana"`        // 姓（かな）
	FirstnameKana string `json:"firstnameKana" binding:"required,max=32,hiragana"`       // 名（かな）
	PhoneNumber   string `json:"phoneNumber" binding:"required,e164"`                    // 電話番号
}

type UpdateAuthUserEmailRequest struct {
	Email string `json:"email" binding:"required,email"` // メールアドレス
}

type VerifyAuthUserEmailRequest struct {
	VerifyCode string `json:"verifyCode" binding:"required"` // 検証コード
}

type UpdateAuthUserUsernameRequest struct {
	Username string `json:"username" binding:"required,max=32"` // ユーザー名(表示名)
}

type UpdateAuthUserAccountIDRequest struct {
	AccountID string `json:"accountId" binding:"required,max=32,account_id"` // ユーザーID(検索名)
}

type UpdateAuthUserNotificationRequest struct {
	Enabled bool `json:"enabled" binding:""` // 通知の有効化設定
}

type UpdateAuthUserThumbnailRequest struct {
	ThumbnailURL string `json:"thumbnailUrl" binding:"omitempty,url"` // サムネイルURL
}
