package request

type UpdateBroadcastArchiveRequest struct {
	ArchiveURL string `json:"archiveUrl,omitempty"` // アーカイブ動画URL
}

type ActivateBroadcastMP4Request struct {
	InputURL string `json:"inputUrl,omitempty"` // 配信動画URL
}
