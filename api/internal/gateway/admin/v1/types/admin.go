package types

// AdminType - 管理者ロール
type AdminType int32

const (
	AdminTypeUnknown       AdminType = 0
	AdminTypeAdministrator AdminType = 1 // 管理者
	AdminTypeCoordinator   AdminType = 2 // コーディネータ
	AdminTypeProducer      AdminType = 3 // 生産者
)

// AdminStatus - 管理者ステータス
type AdminStatus int32

const (
	AdminStatusUnknown     AdminStatus = 0
	AdminStatusInvited     AdminStatus = 1 // 招待中
	AdminStatusActivated   AdminStatus = 2 // 有効
	AdminStatusDeactivated AdminStatus = 3 // 無効
)

// Admin - 管理者情報
type Admin struct {
	ID            string    `json:"id"`            // 管理者ID
	Type          AdminType `json:"role"`          // 管理者種別
	Lastname      string    `json:"lastname"`      // 姓
	Firstname     string    `json:"firstname"`     // 名
	LastnameKana  string    `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string    `json:"firstnameKana"` // 名(かな)
	Email         string    `json:"email"`         // メールアドレス
	CreatedAt     int64     `json:"createdAt"`     // 登録日時
	UpdatedAt     int64     `json:"updateAt"`      // 更新日時
}
