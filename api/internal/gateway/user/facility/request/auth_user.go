package request

type CreateAuthUserRequest struct {
	AuthToken     string `json:"authToken" validate:"required"`            // LINEの認証トークン
	Lastname      string `json:"lastname" validate:"required,max=16"`      // 姓
	Firstname     string `json:"firstname" validate:"required,max=16"`     // 名
	LastnameKana  string `json:"lastnameKana" validate:"required,max=32"`  // 姓 かな
	FirstnameKana string `json:"firstnameKana" validate:"required,max=32"` // 名 かな
	PhoneNumber   string `json:"phoneNumber" validate:"required,e164"`     // 電話番号
	LastCheckInAt int64  `json:"lastCheckInAt" validate:"min=0"`           // 最新のチェックイン日時
}

type UpdateAuthUserRequest struct {
	Lastname      string `json:"lastname" validate:"required,max=16"`      // 姓
	Firstname     string `json:"firstname" validate:"required,max=16"`     // 名
	LastnameKana  string `json:"lastnameKana" validate:"required,max=32"`  // 姓 かな
	FirstnameKana string `json:"firstnameKana" validate:"required,max=32"` // 名 かな
	PhoneNumber   string `json:"phoneNumber" validate:"required,e164"`     // 電話番号
	LastCheckInAt int64  `json:"lastCheckInAt" validate:"min=0"`           // 最新のチェックイン日時
}
