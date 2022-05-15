package response

// Admin - 管理者情報
type Admin struct {
	ID            string `json:"id"`            // 管理者ID
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	Email         string `json:"email"`         // メールアドレス
	Role          int32  `json:"role"`          // 権限
	ThumbnailURL  string `json:"thumbnailUrl"`  // サムネイルURL
	CreatedAt     int64  `json:"createdAt"`     // 登録日時
	UpdatedAt     int64  `json:"updatedAt"`     // 更新日時
}

type AdminResponse struct {
	*Admin
}

type AdminsResponse struct {
	Admins []*Admin `json:"admins"` // 管理者一覧
}

type AdminMeResponse struct {
	ID            string `json:"id"`            // ユーザーID
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	Email         string `json:"email"`         // メールアドレス
	Role          int32  `json:"role"`          // 権限
	ThumbnailURL  string `json:"thumbnailUrl"`  // サムネイルURL
}
