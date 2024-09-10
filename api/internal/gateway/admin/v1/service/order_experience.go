package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type OrderExperience struct {
	response.OrderExperience
	orderID string
}

type OrderExperiences []*OrderExperience

type OrderExperienceRemarks struct {
	response.OrderExperienceRemarks
}

func NewOrderExperience(item *entity.OrderExperience, experience *Experience) *OrderExperience {
	var (
		experienceID                                                                          string
		adultPrice, juniorHighSchoolPrice, elementarySchoolPrice, preschoolPrice, seniorPrice int64
	)
	if item == nil {
		return nil
	}
	if experience != nil {
		experienceID = experience.ID
		adultPrice = experience.PriceAdult
		juniorHighSchoolPrice = experience.PriceJuniorHighSchool
		elementarySchoolPrice = experience.PriceElementarySchool
		preschoolPrice = experience.PricePreschool
		seniorPrice = experience.PriceSenior
	}
	return &OrderExperience{
		OrderExperience: response.OrderExperience{
			ExperienceID:          experienceID,
			AdultCount:            item.AdultCount,
			AdultPrice:            adultPrice,
			JuniorHighSchoolCount: item.JuniorHighSchoolCount,
			JuniorHighSchoolPrice: juniorHighSchoolPrice,
			ElementarySchoolCount: item.ElementarySchoolCount,
			ElementarySchoolPrice: elementarySchoolPrice,
			PreschoolCount:        item.PreschoolCount,
			PreschoolPrice:        preschoolPrice,
			SeniorCount:           item.SeniorCount,
			SeniorPrice:           seniorPrice,
			Remarks:               NewOrderExperienceRemarks(&item.Remarks).Response(),
		},
		orderID: item.OrderID,
	}
}

func (e *OrderExperience) Response() *response.OrderExperience {
	if e == nil {
		return nil
	}
	return &e.OrderExperience
}

func NewOrderExperiences(items entity.OrderExperiences, experiences map[int64]*Experience) OrderExperiences {
	res := make(OrderExperiences, len(items))
	for i, v := range items {
		res[i] = NewOrderExperience(v, experiences[v.ExperienceRevisionID])
	}
	return res
}

func (es OrderExperiences) Response() []*response.OrderExperience {
	res := make([]*response.OrderExperience, len(es))
	for i := range es {
		experience := es[i].Response()
		if experience == nil {
			continue
		}
		res[i] = experience
	}
	return res
}

func NewOrderExperienceRemarks(remarks *entity.OrderExperienceRemarks) *OrderExperienceRemarks {
	var requestedDate, requestedTime string
	if !remarks.RequestedDate.IsZero() {
		requestedDate = jst.FormatYYYYMMDD(remarks.RequestedDate)
	}
	if !remarks.RequestedTime.IsZero() {
		requestedTime = jst.FormatHHMM(remarks.RequestedTime)
	}
	return &OrderExperienceRemarks{
		response.OrderExperienceRemarks{
			Transportation: remarks.Transportation,
			RequestedDate:  requestedDate,
			RequestedTime:  requestedTime,
		},
	}
}

func (r *OrderExperienceRemarks) Response() *response.OrderExperienceRemarks {
	return &r.OrderExperienceRemarks
}
