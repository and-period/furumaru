package response

// Coordinator - コーディネータ情報
type Coordinator struct {
	ID                string   `json:"id"`                // コーディネータID
	MarcheName        string   `json:"marcheName"`        // マルシェ名
	Username          string   `json:"username"`          // 表示名
	Profile           string   `json:"profile"`           // 紹介文
	ProductTypeIDs    []string `json:"productTypeIds"`    // 取り扱い品目一覧
	ThumbnailURL      string   `json:"thumbnailUrl"`      // サムネイルURL
	Thumbnails        []*Image `json:"thumbnails"`        // サムネイルURL(リサイズ済み)一覧
	HeaderURL         string   `json:"headerUrl"`         // ヘッダー画像URL
	Headers           []*Image `json:"headers"`           // ヘッダー画像URL(リサイズ済み)一覧
	PromotionVideoURL string   `json:"promotionVideoUrl"` // 紹介映像URL
	InstagramID       string   `json:"instagramId"`       // Instagramアカウント
	FacebookID        string   `json:"facebookId"`        // Facebookアカウント
}
