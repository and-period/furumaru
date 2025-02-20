package response

// Auth - 認証情報
type Auth struct {
	AdminID      string `json:"adminId"`      // 管理者ID
	Type         int32  `json:"type"`         // 管理者種別
	AccessToken  string `json:"accessToken"`  // アクセストークン
	RefreshToken string `json:"refreshToken"` // 更新トークン
	ExpiresIn    int32  `json:"expiresIn"`    // 有効期限
	TokenType    string `json:"tokenType"`    // トークン種別
}

// AuthUser - ログイン中管理者情報
type AuthUser struct {
	AdminID      string   `json:"id"`           // 管理者ID
	ShopIDs      []string `json:"shopId"`       // 店舗ID一覧
	Type         int32    `json:"type"`         // 管理者種別
	Username     string   `json:"username"`     // 表示名
	Email        string   `json:"email"`        // メールアドレス
	ThumbnailURL string   `json:"thumbnailUrl"` // サムネイルURL
}

// AuthProvider - 認証プロバイダ
type AuthProvider struct {
	Type        int32 `json:"type"`        // プロバイダ種別
	ConnectedAt int64 `json:"connectedAt"` // 連携日時
}

type AuthResponse struct {
	*Auth
}

type AuthUserResponse struct {
	*AuthUser
}

type AuthProvidersResponse struct {
	Providers []*AuthProvider `json:"providers"` // プロバイダ一覧
}

type AuthGoogleAccountResponse struct {
	URL string `json:"url"` // Googleアカウント連携URL
}

type AuthLINEAccountResponse struct {
	URL string `json:"url"` // LINEアカウント連携URL
}
