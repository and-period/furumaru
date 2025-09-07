package request

type SignInRequest struct {
	Username string `json:"username" validate:"required"` // ユーザー名 (メールアドレス,電話番号)
	Password string `json:"password" validate:"required"` // パスワード
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"` // 更新トークン
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
