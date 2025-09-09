package types

// Auth - 認証情報
type Auth struct {
	UserID       string `json:"userId"`       // ユーザーID
	AccessToken  string `json:"accessToken"`  // アクセストークン
	RefreshToken string `json:"refreshToken"` // 更新トークン
	ExpiresIn    int32  `json:"expiresIn"`    // 有効期限
	TokenType    string `json:"tokenType"`    // トークン種別
}

type SignInRequest struct {
	AuthToken string `json:"authToken" validate:"required"` // LINEの認証トークン
}

type GetAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"` // 更新トークン
}

type AuthResponse struct {
	*Auth
}
