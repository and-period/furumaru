//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package media

import "context"

type Service interface {
	// 仲介者サムネイル画像を生成
	GenerateCoordinatorThumbnail(ctx context.Context, in *GenerateFileInput) (string, error)
	// 仲介者サムネイル画像アップロード
	UploadCoordinatorThumbnail(ctx context.Context, in *UploadFileInput) (string, error)
	// 仲介者ヘッダー画像を生成
	GenerateCoordinatorHeader(ctx context.Context, in *GenerateFileInput) (string, error)
	// 仲介者ヘッダー画像アップロード
	UploadCoordinatorHeader(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者サムネイル画像を生成
	GenerateProducerThumbnail(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者サムネイル画像アップロード
	UploadProducerThumbnail(ctx context.Context, in *UploadFileInput) (string, error)
	// 生産者ヘッダー画像を生成
	GenerateProducerHeader(ctx context.Context, in *GenerateFileInput) (string, error)
	// 生産者ヘッダー画像アップロード
	UploadProducerHeader(ctx context.Context, in *UploadFileInput) (string, error)
	// 商品画像を生成
	GenerateProductImage(ctx context.Context, in *GenerateFileInput) (string, error)
	// 商品画像アップロード
	UploadProductImage(ctx context.Context, in *UploadFileInput) (string, error)
	// 商品映像を生成
	GenerateProductVideo(ctx context.Context, in *GenerateFileInput) (string, error)
	// 商品映像アップロード
	UploadProductVideo(ctx context.Context, in *UploadFileInput) (string, error)
	// 品目アイコン画像を生成
	GenerateProductTypeIcon(ctx context.Context, in *GenerateFileInput) (string, error)
	// 品目アイコン画像アップロード
	UploadProductTypeIcon(ctx context.Context, in *UploadFileInput) (string, error)
}
