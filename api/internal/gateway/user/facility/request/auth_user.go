package request

type CreateAuthUserRequest struct {
	AuthToken     string `json:"authToken" binding:"required"`            // LINEの認証トークン
	Lastname      string `json:"lastname" binding:"required,max=16"`      // 姓
	Firstname     string `json:"firstname" binding:"required,max=16"`     // 名
	LastnameKana  string `json:"lastnameKana" binding:"required,max=32"`  // 姓 かな
	FirstnameKana string `json:"firstnameKana" binding:"required,max=32"` // 名 かな
	PhoneNumber   string `json:"phoneNumber" binding:"required,e164"`     // 電話番号
	LastCheckInAt int64  `json:"lastCheckInAt" binding:"min=0"`           // 最新のチェックイン日時
}

type UpdateAuthUserRequest struct {
	Lastname      string `json:"lastname" binding:"required,max=16"`      // 姓
	Firstname     string `json:"firstname" binding:"required,max=16"`     // 名
	LastnameKana  string `json:"lastnameKana" binding:"required,max=32"`  // 姓 かな
	FirstnameKana string `json:"firstnameKana" binding:"required,max=32"` // 名 かな
	PhoneNumber   string `json:"phoneNumber" binding:"required,e164"`     // 電話番号
	LastCheckInAt int64  `json:"lastCheckInAt" binding:"min=0"`           // 最新のチェックイン日時
}
