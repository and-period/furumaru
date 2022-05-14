package request

type CreateAdminRequest struct {
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	Email         string `json:"email"`         // メールアドレス
	Role          int32  `json:"role"`          // 権限
}

type UpdateAdminEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type VerifyAdminEmailRequest struct {
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type UpdateAdminPasswordRequest struct {
	OldPassword          string `json:"oldPassword,omitempty"`          // 現在のパスワード
	NewPassword          string `json:"newPassword,omitempty"`          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}
