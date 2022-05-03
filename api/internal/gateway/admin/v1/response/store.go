package response

// Store - 店舗情報
type Store struct {
	ID           int64    `json:"id"`               // 店舗ID
	Name         string   `json:"name"`             // 店舗名
	ThumbnailURL string   `json:"thumbnailUrl"`     // サムネイルURL
	Staffs       []*Staff `json:"staffs,omitempty"` // 店舗スタッフ一覧
	CreatedAt    int64    `json:"createdAt"`        // 登録日時
	UpdatedAt    int64    `json:"updatedAt"`        // 更新日時
}

// Staff - 店舗スタッフ情報
type Staff struct {
	ID    string `json:"id"`    // 販売者ID
	Name  string `json:"name"`  // 販売者名
	Email string `json:"email"` // メールアドレス
	Role  int32  `json:"role"`  // 権限 (1:管理者, 2:編集者, 3:閲覧者)
}

type StoreResponse struct {
	*Store
}

type StoresResponse struct {
	Stores []*Store `json:"stores"` // 店舗一覧
}
