package request

type UpdateLivePublicRequest struct {
	Published   bool   `json:"published,omitempty"`   // 公開フラグ
	Canceled    bool   `json:"canceled,omitempty"`    // 配信中止フラグ
	ChannelName string `json:"channelName,omitempty"` // チャンネル名
}
