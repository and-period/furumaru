package request

type CreateAdministratorRequest struct {
	Lastname      string `json:"lastname" validate:"required,max=16"`               // 姓
	Firstname     string `json:"firstname" validate:"required,max=16"`              // 名
	LastnameKana  string `json:"lastnameKana" validate:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana" validate:"required,max=32,hiragana"` // 名(かな)
	Email         string `json:"email" validate:"required,email"`                   // メールアドレス
	PhoneNumber   string `json:"phoneNumber" validate:"required,e164"`              // 電話番号
}

type UpdateAdministratorRequest struct {
	Lastname      string `json:"lastname" validate:"required,max=16"`               // 姓
	Firstname     string `json:"firstname" validate:"required,max=16"`              // 名
	LastnameKana  string `json:"lastnameKana" validate:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana" validate:"required,max=32,hiragana"` // 名(かな)
	PhoneNumber   string `json:"phoneNumber" validate:"required,e164"`              // 電話番号
}

type UpdateAdministratorEmailRequest struct {
	Email string `json:"email" validate:"required,email"` // メールアドレス
}
