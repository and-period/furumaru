package types

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
	ShopIDs      []string `json:"shopIds"`      // 店舗ID一覧
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

type SignInRequest struct {
	Username string `json:"username" validate:"required"` // ユーザー名 (メールアドレス)
	Password string `json:"password" validate:"required"` // パスワード
}

type ConnectGoogleAccountRequest struct {
	Code        string `json:"code" validate:"required"`            // 認証コード
	Nonce       string `json:"nonce" validate:"required"`           // セキュア文字列（リプレイアタック対策）
	RedirectURI string `json:"redirectUri" validate:"required,url"` // リダイレクトURI
}

type ConnectLINEAccountRequest struct {
	Code        string `json:"code" validate:"required"`            // 認証コード
	Nonce       string `json:"nonce" validate:"required"`           // セキュア文字列（リプレイアタック対策）
	RedirectURI string `json:"redirectUri" validate:"required,url"` // リダイレクトURI
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"` // 更新トークン
}

type RegisterAuthDeviceRequest struct {
	Device string `json:"device" validate:"required"` // デバイスID
}

type UpdateAuthEmailRequest struct {
	Email string `json:"email" validate:"required,email"` // メールアドレス
}

type VerifyAuthEmailRequest struct {
	VerifyCode string `json:"verifyCode" validate:"required"` // 検証コード
}

type UpdateAuthPasswordRequest struct {
	OldPassword          string `json:"oldPassword" validate:"required"`                              // 現在のパスワード
	NewPassword          string `json:"newPassword" validate:"min=8,max=32"`                          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=NewPassword"` // パスワード (確認用)
}
type ForgotAuthPasswordRequest struct {
	Email string `json:"email" validate:"required,email"` // メールアドレス
}

type ResetAuthPasswordRequest struct {
	Email                string `json:"email" validate:"required,email"`                           // メールアドレス
	VerifyCode           string `json:"verifyCode" validate:"required"`                            // 検証コード
	Password             string `json:"password" validate:"min=8,max=32"`                          // パスワード
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"` // パスワード (確認用)
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
