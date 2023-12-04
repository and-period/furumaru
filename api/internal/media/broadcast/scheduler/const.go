package scheduler

const (
	streamName         = "live/a"
	archiveFilename    = "original.mp4"
	playlistFilename   = "live.m3u8"
	dynamicMP4InputURL = "$urlPath$"
)

// CreatePayload - 配信リソース作成リクエスト
type CreatePayload struct {
	ScheduleID   string                `json:"ScheduleId"`
	ChannelInput *CreateChannelPayload `json:"ChannelInput"`
	MP4Input     *CreateMp4Payload     `json:"MP4Input"`
	RtmpInput    *CreateRtmpPayload    `json:"RtmpInput"`
	ArchiveInput *CreateArchivePayload `json:"ArchiveInput"`
}

// CreateChannelPayload - 配信リソース(MediaLive チャンネル)
type CreateChannelPayload struct {
	Name                   string `json:"Name"`
	StartTime              string `json:"StartTime"`
	InputLossImageSlateURI string `json:"InputLossImageSlateUri"`
}

// CreateMp4Payload - 配信リソース(MediaLive MP4インプット)
type CreateMp4Payload struct {
	InputURL string `json:"InputUrl"`
}

// CreateRtmpPayload - 配信リソース(MediaLive プッシュRTMPインプット)
type CreateRtmpPayload struct {
	StreamName string `json:"StreamName"`
}

// CreateArchivePayload - 配信リソース(MediaLive アーカイブグループ)
type CreateArchivePayload struct {
	BucketName string `json:"BucketName"` // 保管先S3バケット名
	Path       string `json:"Path"`       // 保管先S3バケットPath
}

// RemovePayload - 配信リソース削除リクエスト
type RemovePayload struct {
	CloudFrontDistributionARN string `json:"CloudFrontDistributionArn"`
	MediaLiveChannelID        string `json:"MediaLiveChannelId"`
	MediaStoreContainerARN    string `json:"MediaStoreContainerArn"`
}
