package scheduler

// CreateInput - 配信リソース作成リクエスト
type CreateInput struct {
	RtmpInput    *CreateChannelInput `json:"RtmpInput"`
	MP4Input     *CreateMp4Input     `json:"MP4Input"`
	ChannelInput *CreateChannelInput `json:"ChannelInput"`
}

// CreateRtmpInput - 配信リソース(MediaLive プッシュRTMPインプット)
type CreateRtmpInput struct {
	StreamName string `json:"StreamName"`
}

// CreateMp4Input - 配信リソース(MediaLive MP4インプット)
type CreateMp4Input struct {
	OpeningVideoURL string `json:"OpeningVideoUrl"`
}

// CreateChannelInput - 配信リソース(MediaLive チャンネル)
type CreateChannelInput struct {
	Name                   string `json:"Name"`
	StartTime              string `json:"StartTime"`
	InputLossImageSlateURI string `json:"InputLossImageSlateUri"`
}
