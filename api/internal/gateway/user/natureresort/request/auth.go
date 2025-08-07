package request

type SignInRequest struct {
	AuthToken string `json:"authToken" binding:"required"` // LINEの認証トークン
}

type GetAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 更新トークン
}
