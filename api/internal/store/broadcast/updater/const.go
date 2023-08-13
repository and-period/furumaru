package updater

// CreateOutput - 配信リソース作成レスポンス
type CreateOutput struct {
	CloudFrontURL             string `json:"CloudFrontUrl"`
	CloudFrontDistributionARN string `json:"CloudFrontDistributionArn"`
	MediaLiveRtmpInputARN     string `json:"MediaLiveRtmpInputArn"`
	MediaLiveRtmpInputURL     string `json:"MediaLiveRtmpInputUrl"`
	MediaLiveMp4InputARN      string `json:"MediaLiveMp4InputArn"`
	MediaLiveChannelARN       string `json:"MediaLiveChannelArn"`
	MediaStoreContainerARN    string `json:"MediaStoreContainerArn"`
}
