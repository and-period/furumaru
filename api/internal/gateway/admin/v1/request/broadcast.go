package request

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl,omitempty"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl,omitempty"` // 配信動画URL
}

type AuthYoutubeBroadcastRequest struct {
	YoutubeHandle string `json:"youtubeHandle,omitempty"` // 連携先Youtubeアカウント
}

type CallbackAuthYoutubeBroadcastRequest struct {
	State    string `json:"state,omitempty"`    // Google認証時に取得したstate
	AuthCode string `json:"authCode,omitempty"` // Google認証時に取得したcode
}

type CreateYoutubeBroadcastRequest struct {
	Title       string `json:"title,omitempty"`       // ライブ配信タイトル
	Description string `json:"description,omitempty"` // ライブ配信説明
	Public      bool   `json:"public,omitempty"`      // 公開設定
}
