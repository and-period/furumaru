package scheduler

const streamName = "/live/a"

// CreatePayload - 配信リソース作成リクエスト
type CreatePayload struct {
	ScheduleID   string                `json:"ScheduleId"`
	ChannelInput *CreateChannelPayload `json:"ChannelInput"`
	MP4Input     *CreateMp4Payload     `json:"MP4Input"`
	RtmpInput    *CreateRtmpPayload    `json:"RtmpInput"`
}

// CreateChannelPayload - 配信リソース(MediaLive チャンネル)
type CreateChannelPayload struct {
	Name                   string `json:"Name"`
	StartTime              string `json:"StartTime"`
	InputLossImageSlateURI string `json:"InputLossImageSlateUri"`
}

// CreateMp4Payload - 配信リソース(MediaLive MP4インプット)
type CreateMp4Payload struct {
	OpeningVideoURL string `json:"OpeningVideoUrl"`
}

// CreateRtmpPayload - 配信リソース(MediaLive プッシュRTMPインプット)
type CreateRtmpPayload struct {
	StreamName string `json:"StreamName"`
}
