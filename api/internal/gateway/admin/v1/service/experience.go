package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// ExperienceStatus - 体験受付状況
type ExperienceStatus int32

const (
	ExperienceStatusUnknown   ExperienceStatus = 0
	ExperienceStatusPrivate   ExperienceStatus = 1 // 非公開
	ExperienceStatusWaiting   ExperienceStatus = 2 // 販売開始前
	ExperienceStatusAccepting ExperienceStatus = 3 // 体験受付中
	ExperienceStatusSoldOut   ExperienceStatus = 4 // 体験受付終了
	ExperienceStatusFinished  ExperienceStatus = 5 // 販売終了
	ExperienceStatusArchived  ExperienceStatus = 6 // アーカイブ済み
)

type Experience struct {
	response.Experience
	revisionID int64
}

type Experiences []*Experience

type ExperienceMedia struct {
	response.ExperienceMedia
}

type MultiExperienceMedia []*ExperienceMedia

func NewExperienceStatus(status entity.ExperienceStatus) ExperienceStatus {
	switch status {
	case entity.ExperienceStatusPrivate:
		return ExperienceStatusPrivate
	case entity.ExperienceStatusWaiting:
		return ExperienceStatusWaiting
	case entity.ExperienceStatusAccepting:
		return ExperienceStatusAccepting
	case entity.ExperienceStatusSoldOut:
		return ExperienceStatusSoldOut
	case entity.ExperienceStatusFinished:
		return ExperienceStatusFinished
	case entity.ExperienceStatusArchived:
		return ExperienceStatusArchived
	default:
		return ExperienceStatusUnknown
	}
}

func (s ExperienceStatus) Response() int32 {
	return int32(s)
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
		Experience: response.Experience{
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

func (e *Experience) Response() *response.Experience {
	return &e.Experience
}

func NewExperiences(experiences entity.Experiences) Experiences {
	res := make(Experiences, len(experiences))
	for i := range experiences {
		res[i] = NewExperience(experiences[i])
	}
	return res
}

func (e Experiences) Response() []*response.Experience {
	res := make([]*response.Experience, len(e))
	for i := range e {
		res[i] = e[i].Response()
	}
	return res
}

func NewExperienceMedia(media *entity.ExperienceMedia) *ExperienceMedia {
	return &ExperienceMedia{
		ExperienceMedia: response.ExperienceMedia{
			URL:         media.URL,
			IsThumbnail: media.IsThumbnail,
		},
	}
}

func (m *ExperienceMedia) Response() *response.ExperienceMedia {
	return &m.ExperienceMedia
}

func NewMultiExperienceMedia(media entity.MultiExperienceMedia) MultiExperienceMedia {
	res := make(MultiExperienceMedia, len(media))
	for i := range media {
		res[i] = NewExperienceMedia(media[i])
	}
	return res
}

func (m MultiExperienceMedia) Response() []*response.ExperienceMedia {
	res := make([]*response.ExperienceMedia, len(m))
	for i := range m {
		res[i] = m[i].Response()
	}
	return res
}
