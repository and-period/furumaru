package types

// AuthUser - 認証済みユーザー情報
type AuthUser struct {
	ID            string `json:"id"`            // ユーザーID
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓（かな）
	FirstnameKana string `json:"firstnameKana"` // 名（かな）
	Email         string `json:"email"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
	LastCheckInAt int64  `json:"lastCheckInAt"` // 最新のチェックイン日時
}

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

type AuthUserResponse struct {
	*AuthUser
}
