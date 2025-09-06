package request

type SignInRequest struct {
	Username string `json:"username" binding:"required"` // ユーザー名 (メールアドレス,電話番号)
	Password string `json:"password" binding:"required"` // パスワード
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 更新トークン
}

type UpdateAuthPasswordRequest struct {
	OldPassword          string `json:"oldPassword" binding:"required"`                              // 現在のパスワード
	NewPassword          string `json:"newPassword" binding:"min=8,max=32"`                          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=NewPassword"` // パスワード (確認用)
}

type ForgotAuthPasswordRequest struct {
	Email string `json:"email" binding:"required,email"` // メールアドレス
}

type ResetAuthPasswordRequest struct {
	Email                string `json:"email" binding:"required,email"`                           // メールアドレス
	VerifyCode           string `json:"verifyCode" binding:"required"`                            // 検証コード
	Password             string `json:"password" binding:"min=8,max=32"`                          // パスワード
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"` // パスワード (確認用)
}
