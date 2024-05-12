package request

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl,omitempty"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl,omitempty"` // 配信動画URL
}

type CreateYoutubeBroadcastRequest struct {
	AuthCode string `json:"authCode,omitempty"` // 認証コード
}
