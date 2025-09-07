package request

type SignInRequest struct {
	AuthToken string `json:"authToken" validate:"required"` // LINEの認証トークン
}

type GetAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"` // 更新トークン
}
