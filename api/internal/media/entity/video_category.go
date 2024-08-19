package entity

// VideoCategory - オンデマンド配信カテゴリ情報
type VideoCategory struct {
	ID        string `gorm:"primaryKey;<-:create"` // オンデマンドカテゴリID
	Name      string `gorm:""`                     // カテゴリ名
	CreatedAt string `gorm:"<-:create"`            // 作成日時
	UpdatedAt string `gorm:""`                     // 更新日時
}

type VideoCategories []*VideoCategory
