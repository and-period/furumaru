package response

// Auth - 認証情報
type Auth struct {
	UserID       string `json:"userId"`       // ユーザーID
	AccessToken  string `json:"accessToken"`  // アクセストークン
	RefreshToken string `json:"refreshToken"` // 更新トークン
	ExpiresIn    int32  `json:"expiresIn"`    // 有効期限
	TokenType    string `json:"tokenType"`    // トークン種別
}

// AuthUser - 認証中ユーザー情報
type AuthUser struct {
	ID            string `json:"id"`                      // ユーザーID
	Username      string `json:"username"`                // ユーザー名 (表示名)
	ThumbnailURL  string `json:"thumbnailUrl"`            // サムネイルURL
	Lastname      string `json:"lastname,omitempty"`      // 姓
	Firstname     string `json:"firstname,omitempty"`     // 名
	LastnameKana  string `json:"lastnameKana,omitempty"`  // 姓（かな）
	FirstnameKana string `json:"firstnameKana,omitempty"` // 名（かな）
}

type AuthResponse struct {
	*Auth
}

type AuthUserResponse struct {
	*AuthUser
}

type CreateAuthResponse struct {
	ID string `json:"id"` // ユーザーID
}
