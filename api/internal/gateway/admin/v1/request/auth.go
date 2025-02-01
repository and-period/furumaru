package request

type SignInRequest struct {
	Username string `json:"username,omitempty"` // ユーザー名 (メールアドレス)
	Password string `json:"password,omitempty"` // パスワード
}

type ConnectGoogleAccountRequest struct {
	Code        string `json:"code,omitempty"`        // 認証コード
	Nonce       string `json:"nonce,omitempty"`       // セキュア文字列（リプレイアタック対策）
	RedirectURI string `json:"redirectUri,omitempty"` // リダイレクトURI
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
type ForgotAuthPasswordRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type ResetAuthPasswordRequest struct {
	Email                string `json:"email,omitempty"`                // メールアドレス
	VerifyCode           string `json:"verifyCode,omitempty"`           // 検証コード
	Password             string `json:"password,omitempty"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}
