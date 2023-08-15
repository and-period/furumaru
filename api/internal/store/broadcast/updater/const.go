package updater

// CreatePayload - 配信リソース作成時の情報
type CreatePayload struct {
	ScheduleID                string `json:"ScheduleId"`
	CloudFrontDistributionARN string `json:"CloudFrontDistributionArn"`
	CloudFrontURL             string `json:"CloudFrontUrl"`
	MediaLiveChannelARN       string `json:"MediaLiveChannelArn"`
	MediaLiveMp4InputARN      string `json:"MediaLiveMp4InputArn"`
	MediaLiveRtmpInputARN     string `json:"MediaLiveRtmpInputArn"`
	MediaLiveRtmpInputURL     string `json:"MediaLiveRtmpInputUrl"`
	MediaStoreContainerARN    string `json:"MediaStoreContainerArn"`
}

// RemovePayload - 配信リソース削除時の情報
type RemovePayload struct {
	ScheduleID string `json:"ScheduleId"`
}
