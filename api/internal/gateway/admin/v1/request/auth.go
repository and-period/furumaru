package request

type SignInRequest struct {
	Username string `json:"username,omitempty"` // ユーザー名 (メールアドレス)
	Password string `json:"password,omitempty"` // パスワード
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"` // 更新トークン
}

type RegisterAuthDeviceRequest struct {
	Device string `json:"device,omitempty"` // デバイスID
}

type UpdateAuthEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type VerifyAuthEmailRequest struct {
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type UpdateAuthPasswordRequest struct {
	OldPassword          string `json:"oldPassword,omitempty"`          // 現在のパスワード
	NewPassword          string `json:"newPassword,omitempty"`          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}
