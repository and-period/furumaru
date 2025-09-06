package request

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl" binding:"required,url"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl" binding:"required,url"` // 配信動画URL
}

type AuthYoutubeBroadcastRequest struct {
	YoutubeHandle string `json:"youtubeHandle" binding:"required"` // 連携先Youtubeアカウント
}

type CallbackAuthYoutubeBroadcastRequest struct {
	State    string `json:"state" binding:"required"`    // Google認証時に取得したstate
	AuthCode string `json:"authCode" binding:"required"` // Google認証時に取得したcode
}

type CreateYoutubeBroadcastRequest struct {
	Title       string `json:"title" binding:"required,max=100"`         // ライブ配信タイトル
	Description string `json:"description" binding:"omitempty,max=1000"` // ライブ配信説明
	Public      bool   `json:"public" binding:""`                        // 公開設定
}
