package request

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl"` // 配信動画URL
}

type AuthYoutubeBroadcastRequest struct {
	YoutubeHandle string `json:"youtubeHandle"` // 連携先Youtubeアカウント
}

type CallbackAuthYoutubeBroadcastRequest struct {
	State    string `json:"state"`    // Google認証時に取得したstate
	AuthCode string `json:"authCode"` // Google認証時に取得したcode
}

type CreateYoutubeBroadcastRequest struct {
	Title       string `json:"title"`       // ライブ配信タイトル
	Description string `json:"description"` // ライブ配信説明
	Public      bool   `json:"public"`      // 公開設定
}
