package request

type CreateAdministratorRequest struct {
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	Email         string `json:"email"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
}

type UpdateAdministratorRequest struct {
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
}

type UpdateAdministratorEmailRequest struct {
	Email string `json:"email"` // メールアドレス
}
