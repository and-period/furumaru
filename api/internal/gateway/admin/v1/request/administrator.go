package request

type CreateAdministratorRequest struct {
	Lastname      string `json:"lastname" binding:"required,max=16"`               // 姓
	Firstname     string `json:"firstname" binding:"required,max=16"`              // 名
	LastnameKana  string `json:"lastnameKana" binding:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana" binding:"required,max=32,hiragana"` // 名(かな)
	Email         string `json:"email" binding:"required,email"`                   // メールアドレス
	PhoneNumber   string `json:"phoneNumber" binding:"required,e164"`              // 電話番号
}

type UpdateAdministratorRequest struct {
	Lastname      string `json:"lastname" binding:"required,max=16"`               // 姓
	Firstname     string `json:"firstname" binding:"required,max=16"`              // 名
	LastnameKana  string `json:"lastnameKana" binding:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana" binding:"required,max=32,hiragana"` // 名(かな)
	PhoneNumber   string `json:"phoneNumber" binding:"required,e164"`              // 電話番号
}

type UpdateAdministratorEmailRequest struct {
	Email string `json:"email" binding:"required,email"` // メールアドレス
}
