package request

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl,omitempty"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl,omitempty"` // 配信動画URL
}

type AuthYoutubeBroadcastRequest struct {
	GoogleAccount string `json:"googleAccount,omitempty"` // 連携先Googleアカウント
}

type CreateYoutubeBroadcastRequest struct {
	State    string `json:"state,omitempty"`    // 
	AuthCode string `json:"authCode,omitempty"` // 認証コード
	Public   bool   `json:"public,omitempty"`   // 公開設定
}
