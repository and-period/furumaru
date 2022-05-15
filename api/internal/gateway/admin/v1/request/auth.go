package request

type SignInRequest struct {
	Username string `json:"username,omitempty"` // ユーザー名 (メールアドレス)
	Password string `json:"password,omitempty"` // パスワード
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"` // 更新トークン
}
