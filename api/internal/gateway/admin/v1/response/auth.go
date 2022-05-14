package response

// Auth - 認証情報
type Auth struct {
	AdminID      string `json:"adminId"`      // 管理者ID
	Role         int32  `json:"role"`         // 権限
	AccessToken  string `json:"accessToken"`  // アクセストークン
	RefreshToken string `json:"refreshToken"` // 更新トークン
	ExpiresIn    int32  `json:"expiresIn"`    // 有効期限
	TokenType    string `json:"tokenType"`    // トークン種別
}

type AuthResponse struct {
	*Auth
}
