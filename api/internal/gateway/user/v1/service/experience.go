package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/shopspring/decimal"
)

// ExperienceStatus - 体験受付状況
type ExperienceStatus int32

const (
	ExperienceStatusUnknown   ExperienceStatus = 0
	ExperienceStatusWaiting   ExperienceStatus = 1 // 販売開始前
	ExperienceStatusAccepting ExperienceStatus = 2 // 体験受付中
	ExperienceStatusSoldOut   ExperienceStatus = 3 // 体験受付終了
	ExperienceStatusFinished  ExperienceStatus = 4 // 販売終了
)

func NewExperienceStatus(status entity.ExperienceStatus) ExperienceStatus {
	switch status {
	case entity.ExperienceStatusPrivate:
		return ExperienceStatusUnknown
	case entity.ExperienceStatusWaiting:
		return ExperienceStatusWaiting
	case entity.ExperienceStatusAccepting:
		return ExperienceStatusAccepting
	case entity.ExperienceStatusSoldOut:
		return ExperienceStatusSoldOut
	case entity.ExperienceStatusFinished:
		return ExperienceStatusFinished
	default:
		return ExperienceStatusUnknown
	}
}

func (s ExperienceStatus) Response() int32 {
	return int32(s)
}

type Experience struct {
	response.Experience
	revisionID int64
}

type Experiences []*Experience

type CalcExperienceParams struct {
	AdultCount            int64
	JuniorHighSchoolCount int64
	ElementarySchoolCount int64
	PreschoolCount        int64
	SeniorCount           int64
	Promotion             *Promotion
}

func NewExperience(experience *entity.Experience, rate *ExperienceRate) *Experience {
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
			Status:                NewExperienceStatus(experience.Status).Response(),
			ThumbnailURL:          experience.ThumbnailURL,
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
			HostPrefecture:        experience.HostPrefecture,
			HostCity:              experience.HostCity,
			HostAddressLine1:      experience.HostAddressLine1,
			HostAddressLine2:      experience.HostAddressLine2,
			HostLongitude:         experience.HostLongitude,
			HostLatitude:          experience.HostLatitude,
			Rate:                  rate.Response(),
			StartAt:               experience.StartAt.Unix(),
			EndAt:                 experience.EndAt.Unix(),
		},
		revisionID: experience.ExperienceRevision.ID,
	}
}

//nolint:nakedret
func (e *Experience) Calc(params *CalcExperienceParams) (subtotal int64, discount int64) {
	if e == nil || params == nil {
		return
	}

	dsub := decimal.Zero
	dsub = dsub.Add(decimal.NewFromInt(e.PriceAdult).Mul(decimal.NewFromInt(params.AdultCount)))
	dsub = dsub.Add(decimal.NewFromInt(e.PriceJuniorHighSchool).Mul(decimal.NewFromInt(params.JuniorHighSchoolCount)))
	dsub = dsub.Add(decimal.NewFromInt(e.PriceElementarySchool).Mul(decimal.NewFromInt(params.ElementarySchoolCount)))
	dsub = dsub.Add(decimal.NewFromInt(e.PricePreschool).Mul(decimal.NewFromInt(params.PreschoolCount)))
	dsub = dsub.Add(decimal.NewFromInt(e.PriceSenior).Mul(decimal.NewFromInt(params.SeniorCount)))
	subtotal = dsub.IntPart()

	if params.Promotion == nil {
		return
	}

	switch DiscountType(params.Promotion.DiscountType) {
	case DiscountTypeAmount:
		if subtotal < params.Promotion.DiscountRate {
			discount = subtotal
		} else {
			discount = params.Promotion.DiscountRate
		}
	case DiscountTypeRate:
		if params.Promotion.DiscountRate <= 0 {
			return
		}
		rate := decimal.NewFromInt(params.Promotion.DiscountRate).Div(decimal.NewFromInt(100))
		discount = dsub.Mul(rate).IntPart()
	}

	return
}

func (e *Experience) Response() *response.Experience {
	if e == nil {
		return nil
	}
	return &e.Experience
}

func NewExperiences(experiences entity.Experiences, rates map[string]*ExperienceRate) Experiences {
	res := make(Experiences, len(experiences))
	for i, experience := range experiences {
		res[i] = NewExperience(experience, rates[experience.ID])
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

func (es Experiences) Response() []*response.Experience {
	res := make([]*response.Experience, len(es))
	for i := range es {
		res[i] = es[i].Response()
	}
	return res
}

type ExperienceMedia struct {
	response.ExperienceMedia
}

type MultiExperienceMedia []*ExperienceMedia

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

func NewMultiExperienceMedia(media []*entity.ExperienceMedia) MultiExperienceMedia {
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

type ExperienceRate struct {
	response.ExperienceRate
	experienceID string
}

type ExperienceRates []*ExperienceRate

func newExperienceRate(review *entity.AggregatedExperienceReview) *ExperienceRate {
	return &ExperienceRate{
		ExperienceRate: response.ExperienceRate{
			Count:   review.Count,
			Average: review.Average,
			Detail: map[int64]int64{
				1: review.Rate1,
				2: review.Rate2,
				3: review.Rate3,
				4: review.Rate4,
				5: review.Rate5,
			},
		},
		experienceID: review.ExperienceID,
	}
}

func newEmptyExperienceRate() *ExperienceRate {
	return &ExperienceRate{
		ExperienceRate: response.ExperienceRate{
			Count:   0,
			Average: 0.0,
			Detail: map[int64]int64{
				1: 0,
				2: 0,
				3: 0,
				4: 0,
				5: 0,
			},
		},
		experienceID: "",
	}
}

func (r *ExperienceRate) Response() *response.ExperienceRate {
	if r == nil {
		return newEmptyExperienceRate().Response()
	}
	return &r.ExperienceRate
}

func NewExperienceRates(reviews entity.AggregatedExperienceReviews) ExperienceRates {
	res := make(ExperienceRates, len(reviews))
	for i, review := range reviews {
		res[i] = newExperienceRate(review)
	}
	return res
}

func (rs ExperienceRates) MapByExperienceID() map[string]*ExperienceRate {
	res := make(map[string]*ExperienceRate, len(rs))
	for _, r := range rs {
		res[r.experienceID] = r
	}
	return res
}

func (rs ExperienceRates) Response() []*response.ExperienceRate {
	res := make([]*response.ExperienceRate, len(rs))
	for i := range rs {
		res[i] = rs[i].Response()
	}
	return res
}
