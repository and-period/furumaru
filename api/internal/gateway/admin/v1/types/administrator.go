package types

// Administrator - システム管理者情報
type Administrator struct {
	ID            string      `json:"id"`            // 管理者ID
	Status        AdminStatus `json:"status"`        // 管理者の状態
	Lastname      string      `json:"lastname"`      // 姓
	Firstname     string      `json:"firstname"`     // 名
	LastnameKana  string      `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string      `json:"firstnameKana"` // 名(かな)
	Email         string      `json:"email"`         // メールアドレス
	PhoneNumber   string      `json:"phoneNumber"`   // 電話番号
	CreatedAt     int64       `json:"createdAt"`     // 登録日時
	UpdatedAt     int64       `json:"updatedAt"`     // 更新日時
}

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

type AdministratorResponse struct {
	Administrator *Administrator `json:"administrator"` // システム管理者情報
}

type AdministratorsResponse struct {
	Administrators []*Administrator `json:"administrators"` // システム管理者一覧
	Total          int64            `json:"total"`          // 合計数
}
