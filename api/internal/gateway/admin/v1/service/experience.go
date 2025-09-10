package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// ExperienceStatus - 体験受付状況
type ExperienceStatus types.ExperienceStatus

type Experience struct {
	types.Experience
	revisionID int64
}

type Experiences []*Experience

type ExperienceMedia struct {
	types.ExperienceMedia
}

type MultiExperienceMedia []*ExperienceMedia

func NewExperienceStatus(status entity.ExperienceStatus) ExperienceStatus {
	switch status {
	case entity.ExperienceStatusPrivate:
		return ExperienceStatus(types.ExperienceStatusPrivate)
	case entity.ExperienceStatusWaiting:
		return ExperienceStatus(types.ExperienceStatusWaiting)
	case entity.ExperienceStatusAccepting:
		return ExperienceStatus(types.ExperienceStatusAccepting)
	case entity.ExperienceStatusSoldOut:
		return ExperienceStatus(types.ExperienceStatusSoldOut)
	case entity.ExperienceStatusFinished:
		return ExperienceStatus(types.ExperienceStatusFinished)
	case entity.ExperienceStatusArchived:
		return ExperienceStatus(types.ExperienceStatusArchived)
	default:
		return ExperienceStatus(types.ExperienceStatusUnknown)
	}
}

func (s ExperienceStatus) Response() types.ExperienceStatus {
	return types.ExperienceStatus(s)
}

func NewExperience(experience *entity.Experience) *Experience {
	var point1, point2, point3 string
	if len(experience.RecommendedPoints) > 0 {
		point1 = experience.RecommendedPoints[0]
	}
	if len(experience.RecommendedPoints) > 1 {
		point2 = experience.RecommendedPoints[1]
	}
	if len(experience.RecommendedPoints) > 2 {
		point3 = experience.RecommendedPoints[2]
	}
	return &Experience{
		Experience: types.Experience{
			ID:                    experience.ID,
			CoordinatorID:         experience.CoordinatorID,
			ProducerID:            experience.ProducerID,
			ExperienceTypeID:      experience.TypeID,
			Title:                 experience.Title,
			Description:           experience.Description,
			Public:                experience.Public,
			SoldOut:               experience.SoldOut,
			Status:                NewExperienceStatus(experience.Status).Response(),
			Media:                 NewMultiExperienceMedia(experience.Media).Response(),
			PriceAdult:            experience.PriceAdult,
			PriceJuniorHighSchool: experience.PriceJuniorHighSchool,
			PriceElementarySchool: experience.PriceElementarySchool,
			PricePreschool:        experience.PricePreschool,
			PriceSenior:           experience.PriceSenior,
			RecommendedPoint1:     point1,
			RecommendedPoint2:     point2,
			RecommendedPoint3:     point3,
			PromotionVideoURL:     experience.PromotionVideoURL,
			Duration:              experience.Duration,
			Direction:             experience.Direction,
			BusinessOpenTime:      experience.BusinessOpenTime,
			BusinessCloseTime:     experience.BusinessCloseTime,
			HostPostalCode:        experience.HostPostalCode,
			HostPrefectureCode:    experience.HostPrefectureCode,
			HostCity:              experience.HostCity,
			HostAddressLine1:      experience.HostAddressLine1,
			HostAddressLine2:      experience.HostAddressLine2,
			StartAt:               experience.StartAt.Unix(),
			EndAt:                 experience.EndAt.Unix(),
			CreatedAt:             experience.CreatedAt.Unix(),
			UpdatedAt:             experience.UpdatedAt.Unix(),
		},
		revisionID: experience.ExperienceRevision.ID,
	}
}

func (e *Experience) Response() *types.Experience {
	if e == nil {
		return nil
	}
	return &e.Experience
}

func NewExperiences(experiences entity.Experiences) Experiences {
	res := make(Experiences, len(experiences))
	for i := range experiences {
		res[i] = NewExperience(experiences[i])
	}
	return res
}

func (es Experiences) MapByRevision() map[int64]*Experience {
	res := make(map[int64]*Experience, len(es))
	for _, e := range es {
		res[e.revisionID] = e
	}
	return res
}

func (es Experiences) Response() []*types.Experience {
	res := make([]*types.Experience, len(es))
	for i := range es {
		res[i] = es[i].Response()
	}
	return res
}

func NewExperienceMedia(media *entity.ExperienceMedia) *ExperienceMedia {
	return &ExperienceMedia{
		ExperienceMedia: types.ExperienceMedia{
			URL:         media.URL,
			IsThumbnail: media.IsThumbnail,
		},
	}
}

func (m *ExperienceMedia) Response() *types.ExperienceMedia {
	return &m.ExperienceMedia
}

func NewMultiExperienceMedia(media entity.MultiExperienceMedia) MultiExperienceMedia {
	res := make(MultiExperienceMedia, len(media))
	for i := range media {
		res[i] = NewExperienceMedia(media[i])
	}
	return res
}

func (m MultiExperienceMedia) Response() []*types.ExperienceMedia {
	res := make([]*types.ExperienceMedia, len(m))
	for i := range m {
		res[i] = m[i].Response()
	}
	return res
}
