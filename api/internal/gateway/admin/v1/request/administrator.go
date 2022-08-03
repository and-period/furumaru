package request

type CreateAdministratorRequest struct {
	Lastname      string `json:"lastname,omitempty"`      // 姓
	Firstname     string `json:"firstname,omitempty"`     // 名
	LastnameKana  string `json:"lastnameKana,omitempty"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana,omitempty"` // 名(かな)
	Email         string `json:"email,omitempty"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber,omitempty"`   // 電話番号
}

type UpdateAdministratorRequest struct {
	Lastname      string `json:"lastname,omitempty"`      // 姓
	Firstname     string `json:"firstname,omitempty"`     // 名
	LastnameKana  string `json:"lastnameKana,omitempty"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana,omitempty"` // 名(かな)
	PhoneNumber   string `json:"phoneNumber,omitempty"`   // 電話番号
}

type UpdateAdministratorEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}
