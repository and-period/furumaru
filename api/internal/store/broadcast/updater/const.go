package updater

// CreatePayload - 配信リソース作成時の情報
type CreatePayload struct {
	ScheduleID                string `json:"ScheduleId"`
	CloudFrontDistributionARN string `json:"CloudFrontDistributionArn"`
	CloudFrontURL             string `json:"CloudFrontUrl"`
	MediaLiveChannelARN       string `json:"MediaLiveChannelArn"`
	MediaLiveChannelID        string `json:"MediaLiveChannelId"`
	MediaLiveMp4InputARN      string `json:"MediaLiveMp4InputArn"`
	MediaLiveMp4InputName     string `json:"MediaLiveMp4InputName"`
	MediaLiveRtmpInputARN     string `json:"MediaLiveRtmpInputArn"`
	MediaLiveRtmpInputName    string `json:"MediaLiveRtmpInputName"`
	MediaLiveRtmpInputURL     string `json:"MediaLiveRtmpInputUrl"`
	MediaLiveRtmpStreamName   string `json:"MediaLiveRtmpStreamName"`
	MediaStoreContainerARN    string `json:"MediaStoreContainerArn"`
}

// RemovePayload - 配信リソース削除時の情報
type RemovePayload struct {
	ScheduleID string `json:"ScheduleId"`
}
