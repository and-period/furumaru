package response

// Auth - 認証情報
type Auth struct {
	AdminID      string `json:"adminId"`      // 管理者ID
	Role         int32  `json:"role"`         // 権限
	AccessToken  string `json:"accessToken"`  // アクセストークン
	RefreshToken string `json:"refreshToken"` // 更新トークン
	ExpiresIn    int32  `json:"expiresIn"`    // 有効期限
	TokenType    string `json:"tokenType"`    // トークン種別
}

// AuthUser - 認証中ユーザー情報
type AuthUser struct {
	ID            string `json:"id"`            // ユーザーID
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	StoreName     string `json:"storeName"`     // 店舗名
	ThumbnailURL  string `json:"thumbnailUrl"`  // サムネイルURL
	Role          int32  `json:"role"`          // 権限
}

type AuthResponse struct {
	*Auth
}

type AuthUserResponse struct {
	*AuthUser
}
