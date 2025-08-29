package request

type SignInRequest struct {
	Username string `json:"username"` // ユーザー名 (メールアドレス)
	Password string `json:"password"` // パスワード
}

type ConnectGoogleAccountRequest struct {
	Code        string `json:"code"`        // 認証コード
	Nonce       string `json:"nonce"`       // セキュア文字列（リプレイアタック対策）
	RedirectURI string `json:"redirectUri"` // リダイレクトURI
}

type ConnectLINEAccountRequest struct {
	Code        string `json:"code"`        // 認証コード
	Nonce       string `json:"nonce"`       // セキュア文字列（リプレイアタック対策）
	RedirectURI string `json:"redirectUri"` // リダイレクトURI
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken"` // 更新トークン
}

type RegisterAuthDeviceRequest struct {
	Device string `json:"device"` // デバイスID
}

type UpdateAuthEmailRequest struct {
	Email string `json:"email"` // メールアドレス
}

type VerifyAuthEmailRequest struct {
	VerifyCode string `json:"verifyCode"` // 検証コード
}

type UpdateAuthPasswordRequest struct {
	OldPassword          string `json:"oldPassword"`          // 現在のパスワード
	NewPassword          string `json:"newPassword"`          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation"` // パスワード (確認用)
}
type ForgotAuthPasswordRequest struct {
	Email string `json:"email"` // メールアドレス
}

type ResetAuthPasswordRequest struct {
	Email                string `json:"email"`                // メールアドレス
	VerifyCode           string `json:"verifyCode"`           // 検証コード
	Password             string `json:"password"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation"` // パスワード (確認用)
}
