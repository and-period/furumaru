package response

// リサイズ画像一覧
type Image struct {
	URL  string `json:"url"`  // 画像URL
	Size int32  `json:"size"` // 画像サイズ
}
