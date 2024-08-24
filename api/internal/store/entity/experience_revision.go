package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/jinzhu/copier"
)

// ExperienceRevision - 体験変更履歴情報
type ExperienceRevision struct {
	ID                    int64     `gorm:"primaryKey;<-:create"` // 体験変更履歴ID
	ExperienceID          string    `gorm:""`                     // 体験ID
	PriceAdult            int64     `gorm:""`                     // 大人料金（高校生以上）
	PriceJuniorHighSchool int64     `gorm:""`                     // 中学生料金
	PriceElementarySchool int64     `gorm:""`                     // 小学生料金
	PricePreschool        int64     `gorm:""`                     // 幼児料金
	PriceSenior           int64     `gorm:""`                     // シニア料金（65歳以上）
	CreatedAt             time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt             time.Time `gorm:""`                     // 更新日時
}

type ExperienceRevisions []*ExperienceRevision

type NewExperienceRevisionParams struct {
	ExperienceID          string
	PriceAdult            int64
	PriceJuniorHighSchool int64
	PriceElementarySchool int64
	PricePreschool        int64
	PriceSenior           int64
}

func NewExperienceRevision(params *NewExperienceRevisionParams) *ExperienceRevision {
	return &ExperienceRevision{
		ExperienceID:          params.ExperienceID,
		PriceAdult:            params.PriceAdult,
		PriceJuniorHighSchool: params.PriceJuniorHighSchool,
		PriceElementarySchool: params.PriceElementarySchool,
		PricePreschool:        params.PricePreschool,
		PriceSenior:           params.PriceSenior,
	}
}

func (rs ExperienceRevisions) ExperienceIDs() []string {
	return set.UniqBy(rs, func(r *ExperienceRevision) string {
		return r.ExperienceID
	})
}

func (rs ExperienceRevisions) MapByExperienceID() map[string]*ExperienceRevision {
	res := make(map[string]*ExperienceRevision, len(rs))
	for _, r := range rs {
		res[r.ExperienceID] = r
	}
	return res
}

func (rs ExperienceRevisions) Merge(experiences map[string]*Experience) (Experiences, error) {
	res := make(Experiences, 0, len(rs))
	for _, r := range rs {
		experience := &Experience{}
		base, ok := experiences[r.ExperienceID]
		if !ok {
			base = &Experience{ID: r.ExperienceID}
		}
		opt := copier.Option{IgnoreEmpty: true, DeepCopy: true}
		if err := copier.CopyWithOption(&experience, &base, opt); err != nil {
			return nil, err
		}
		experience.ExperienceRevision = *r
		res = append(res, experience)
	}
	return res, nil
}
