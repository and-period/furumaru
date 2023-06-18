//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package media

import "context"

type Service interface {
	// コーディネータサムネイル画像を生成
	GenerateCoordinatorThumbnail(ctx context.Context, in *GenerateFileInput) (string, error)
	// コーディネータサムネイル画像アップロード
	UploadCoordinatorThumbnail(ctx context.Context, in *UploadFileInput) (string, error)
	// コーディネータサムネイル画像リサイズ
	ResizeCoordinatorThumbnail(ctx context.Context, in *ResizeFileInput) error
	// コーディネータヘッダー画像を生成
	GenerateCoordinatorHeader(ctx context.Context, in *GenerateFileInput) (string, error)
	// コーディネータヘッダー画像アップロード
	UploadCoordinatorHeader(ctx context.Context, in *UploadFileInput) (string, error)
	// コーディネータヘッダー画像リサイズ
	ResizeCoordinatorHeader(ctx context.Context, in *ResizeFileInput) error
	// コーディネータ紹介映像を生成
	GenerateCoordinatorPromotionVideo(ctx context.Context, in *GenerateFileInput) (string, error)
	// コーディネータ紹介映像アップロード
	UploadCoordinatorPromotionVideo(ctx context.Context, in *UploadFileInput) (string, error)
	// コーディネータ購入特典映像を生成
	GenerateCoordinatorBonusVideo(ctx context.Context, in *GenerateFileInput) (string, error)
	// コーディネータ購入特典映像アップロード
	UploadCoordinatorBonusVideo(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者サムネイル画像を生成
	GenerateProducerThumbnail(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者サムネイル画像アップロード
	UploadProducerThumbnail(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者サムネイル画像リサイズ
	ResizeProducerThumbnail(ctx context.Context, in *ResizeFileInput) error
	// 生産者ヘッダー画像を生成
	GenerateProducerHeader(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者ヘッダー画像アップロード
	UploadProducerHeader(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者ヘッダー画像リサイズ
	ResizeProducerHeader(ctx context.Context, in *ResizeFileInput) error
	// 生産者紹介映像を生成
	GenerateProducerPromotionVideo(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者紹介映像アップロード
	UploadProducerPromotionVideo(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者購入特典映像を生成
	GenerateProducerBonusVideo(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者購入特典映像アップロード
	UploadProducerBonusVideo(ctx context.Context, in *UploadFileInput) (string, error)
	// 商品メディア(画像)を生成
	GenerateProductMediaImage(ctx context.Context, in *GenerateFileInput) (string, error)
	// 商品メディア(映像)を生成
	GenerateProductMediaVideo(ctx context.Context, in *GenerateFileInput) (string, error)
	// 商品メディアアップロード
	UploadProductMedia(ctx context.Context, in *UploadFileInput) (string, error)
	// 商品メディアリサイズ
	ResizeProductMedia(ctx context.Context, in *ResizeFileInput) error
	// 品目アイコン画像を生成
	GenerateProductTypeIcon(ctx context.Context, in *GenerateFileInput) (string, error)
	// 品目アイコン画像アップロード
	UploadProductTypeIcon(ctx context.Context, in *UploadFileInput) (string, error)
	// 品目アイコン画像リサイズ
	ResizeProductTypeIcon(ctx context.Context, in *ResizeFileInput) error
}
