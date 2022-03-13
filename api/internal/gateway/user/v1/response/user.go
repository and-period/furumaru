package response

type CreateUserResponse struct {
	ID string `json:"id"` // ユーザーID
}

type UserMeResponse struct {
	ID          string `json:"id"`          // ユーザーID
	Email       string `json:"email"`       // メールアドレス
	PhoneNumber string `json:"phoneNumber"` // 電話番号
}
