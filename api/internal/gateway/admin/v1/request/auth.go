package request

type SignInRequest struct {
	Username string `json:"username" validate:"required"` // ユーザー名 (メールアドレス)
	Password string `json:"password" validate:"required"` // パスワード
}

type ConnectGoogleAccountRequest struct {
	Code        string `json:"code" validate:"required"`            // 認証コード
	Nonce       string `json:"nonce" validate:"required"`           // セキュア文字列（リプレイアタック対策）
	RedirectURI string `json:"redirectUri" validate:"required,url"` // リダイレクトURI
}

type ConnectLINEAccountRequest struct {
	Code        string `json:"code" validate:"required"`            // 認証コード
	Nonce       string `json:"nonce" validate:"required"`           // セキュア文字列（リプレイアタック対策）
	RedirectURI string `json:"redirectUri" validate:"required,url"` // リダイレクトURI
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"` // 更新トークン
}

type RegisterAuthDeviceRequest struct {
	Device string `json:"device" validate:"required"` // デバイスID
}

type UpdateAuthEmailRequest struct {
	Email string `json:"email" validate:"required,email"` // メールアドレス
}

type VerifyAuthEmailRequest struct {
	VerifyCode string `json:"verifyCode" validate:"required"` // 検証コード
}

type UpdateAuthPasswordRequest struct {
	OldPassword          string `json:"oldPassword" validate:"required"`                              // 現在のパスワード
	NewPassword          string `json:"newPassword" validate:"min=8,max=32"`                          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=NewPassword"` // パスワード (確認用)
}
type ForgotAuthPasswordRequest struct {
	Email string `json:"email" validate:"required,email"` // メールアドレス
}

type ResetAuthPasswordRequest struct {
	Email                string `json:"email" validate:"required,email"`                           // メールアドレス
	VerifyCode           string `json:"verifyCode" validate:"required"`                            // 検証コード
	Password             string `json:"password" validate:"min=8,max=32"`                          // パスワード
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"` // パスワード (確認用)
}
