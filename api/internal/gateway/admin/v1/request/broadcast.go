package request

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl" validate:"required,url"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl" validate:"required,url"` // 配信動画URL
}

type AuthYoutubeBroadcastRequest struct {
	YoutubeHandle string `json:"youtubeHandle" validate:"required"` // 連携先Youtubeアカウント
}

type CallbackAuthYoutubeBroadcastRequest struct {
	State    string `json:"state" validate:"required"`    // Google認証時に取得したstate
	AuthCode string `json:"authCode" validate:"required"` // Google認証時に取得したcode
}

type CreateYoutubeBroadcastRequest struct {
	Title       string `json:"title" validate:"required,max=100"`         // ライブ配信タイトル
	Description string `json:"description" validate:"omitempty,max=1000"` // ライブ配信説明
	Public      bool   `json:"public" validate:""`                        // 公開設定
}
