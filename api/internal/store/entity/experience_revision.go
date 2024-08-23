package entity

// ExperienceRevision - 体験変更履歴情報
type ExperienceRevision struct {
	ID                    int64  `gorm:"primaryKey;<-:create"` // 体験変更履歴ID
	ExperienceID          string `gorm:""`                     // 体験ID
	PriceAdult            int64  `gorm:""`                     // 大人料金（高校生以上）
	PriceJuniorHighSchool int64  `gorm:""`                     // 中学生料金
	PriceElementarySchool int64  `gorm:""`                     // 小学生料金
	PricePreschooler      int64  `gorm:""`                     // 幼児料金
	PriceSenior           int64  `gorm:""`                     // シニア料金（65歳以上）
	CreatedAt             string `gorm:"<-:create"`            // 登録日時
	UpdatedAt             string `gorm:""`                     // 更新日時
}
