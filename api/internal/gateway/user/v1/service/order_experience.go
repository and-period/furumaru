package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type OrderExperience struct {
	response.OrderExperience
	orderID string
}

type OrderExperiences []*OrderExperience

func NewOrderExperience(item *entity.OrderExperience, experience *Experience) *OrderExperience {
	var (
		experienceID                                                                          string
		adultPrice, juniorHighSchoolPrice, elementarySchoolPrice, preschoolPrice, seniorPrice int64
		requestedDate, requestedTime                                                          string
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
	if !item.Remarks.RequestedDate.IsZero() {
		requestedDate = jst.FormatYYYYMMDD(item.Remarks.RequestedDate)
	}
	if !item.Remarks.RequestedTime.IsZero() {
		requestedTime = jst.FormatHHMM(item.Remarks.RequestedTime)
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
			Transportation:        item.Remarks.Transportation,
			RequestedDate:         requestedDate,
			RequestedTime:         requestedTime,
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
	res := make(OrderExperiences, 0, len(items))
	for _, v := range items {
		experience := NewOrderExperience(v, experiences[v.ExperienceRevisionID])
		if experience == nil {
			continue
		}
		res = append(res, experience)
	}
	return res
}

func (es OrderExperiences) Response() []*response.OrderExperience {
	res := make([]*response.OrderExperience, 0, len(es))
	for _, e := range es {
		if e == nil {
			continue
		}
		res = append(res, e.Response())
	}
	return res
}
