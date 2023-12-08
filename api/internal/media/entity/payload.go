package entity

// FileType - ファイル種別
type FileType int32

const (
	FileTypeUnknown              FileType = 0
	FileTypeCoordinatorThumbnail FileType = 1 // コーディネータサムネイル画像
	FileTypeCoordinatorHeader    FileType = 2 // コーディネータヘッダー画像
	FileTypeProducerThumbnail    FileType = 3 // 生産者サムネイル画像
	FileTypeProducerHeader       FileType = 4 // 生産者ヘッダー画像
	FileTypeProductMedia         FileType = 5 // 商品メディア
	FileTypeProductTypeIcon      FileType = 6 // 品目アイコン画像
	FileTypeScheduleThumbnail    FileType = 7 // 開催スケジュールサムネイル画像
	FileTypeUserThumbnail        FileType = 8 // 購入者サムネイル画像
)

// ResizerPayload - メディア
type ResizerPayload struct {
	TargetID string   `json:"id"`       // 対象データID
	FileType FileType `json:"fileType"` // ファイル種別
	URLs     []string `json:"urls"`     // ファイル参照先URL一覧
}
