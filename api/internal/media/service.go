//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package media

import "context"

type Service interface {
	// コーディネータサムネイル画像を生成
	GenerateCoordinatorThumbnail(ctx context.Context, in *GenerateFileInput) (string, error)
	// コーディネータサムネイル画像アップロード
	UploadCoordinatorThumbnail(ctx context.Context, in *UploadFileInput) (string, error)
	// コーディネータヘッダー画像を生成
	GenerateCoordinatorHeader(ctx context.Context, in *GenerateFileInput) (string, error)
	// コーディネータヘッダー画像アップロード
	UploadCoordinatorHeader(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者サムネイル画像を生成
	GenerateProducerThumbnail(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者サムネイル画像アップロード
	UploadProducerThumbnail(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者ヘッダー画像を生成
	GenerateProducerHeader(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者ヘッダー画像アップロード
	UploadProducerHeader(ctx context.Context, in *UploadFileInput) (string, error)
	// 商品メディア(画像)を生成
	GenerateProductMediaImage(ctx context.Context, in *GenerateFileInput) (string, error)
	// 商品メディア(映像)を生成
	GenerateProductMediaVideo(ctx context.Context, in *GenerateFileInput) (string, error)
	// 商品メディアアップロード
	UploadProductMedia(ctx context.Context, in *UploadFileInput) (string, error)
	// 品目アイコン画像を生成
	GenerateProductTypeIcon(ctx context.Context, in *GenerateFileInput) (string, error)
	// 品目アイコン画像アップロード
	UploadProductTypeIcon(ctx context.Context, in *UploadFileInput) (string, error)
}
