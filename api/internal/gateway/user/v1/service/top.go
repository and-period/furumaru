package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type TopCommonLive struct {
	response.TopCommonLive
}

type TopCommonLives []*TopCommonLive

type TopCommonArchive struct {
	response.TopCommonArchive
}

type TopCommonArchives []*TopCommonArchive

func NewTopCommonLive(schedule *entity.Schedule, products entity.Products) *TopCommonLive {
	return &TopCommonLive{
		TopCommonLive: response.TopCommonLive{
			ScheduleID:    schedule.ID,
			CoordinatorID: schedule.CoordinatorID,
			Status:        NewScheduleStatus(schedule.Status, false).Response(),
			Title:         schedule.Title,
			ThumbnailURL:  schedule.ThumbnailURL,
			Thumbnails:    NewImages(schedule.Thumbnails).Response(),
			StartAt:       schedule.StartAt.Unix(),
			EndAt:         schedule.EndAt.Unix(),
			Products:      newTopCommonLiveProducts(products),
		},
	}
}

func newTopCommonLiveProducts(products entity.Products) []*response.TopCommonLiveProduct {
	res := make([]*response.TopCommonLiveProduct, len(products))
	for i := range products {
		res[i] = &response.TopCommonLiveProduct{
			ProductID:    products[i].ID,
			Name:         products[i].Name,
			Price:        products[i].Price,
			Inventory:    products[i].Inventory,
			ThumbnailURL: "",
			Thumbnails:   []*response.Image{},
		}
		for _, media := range products[i].Media {
			if !media.IsThumbnail {
				continue
			}
			res[i].ThumbnailURL = media.URL
			res[i].Thumbnails = NewImages(media.Images).Response()
			break
		}
	}
	return res
}

func (l *TopCommonLive) Response() *response.TopCommonLive {
	return &l.TopCommonLive
}

func NewTopCommonLives(schedules entity.Schedules, lives entity.Lives, products entity.Products) TopCommonLives {
	livesMap := lives.GroupByScheduleID()
	res := make(TopCommonLives, len(schedules))
	for i := range schedules {
		productIDs := livesMap[schedules[i].ID].ProductIDs()
		res[i] = NewTopCommonLive(schedules[i], products.Filter(productIDs...))
	}
	return res
}

func (ls TopCommonLives) Response() []*response.TopCommonLive {
	res := make([]*response.TopCommonLive, len(ls))
	for i := range ls {
		res[i] = ls[i].Response()
	}
	return res
}

func NewTopCommonArchive(schedule *entity.Schedule) *TopCommonArchive {
	return &TopCommonArchive{
		TopCommonArchive: response.TopCommonArchive{
			ScheduleID:    schedule.ID,
			CoordinatorID: schedule.CoordinatorID,
			Title:         schedule.Title,
			ThumbnailURL:  schedule.ThumbnailURL,
			Thumbnails:    NewImages(schedule.Thumbnails).Response(),
		},
	}
}

func (a *TopCommonArchive) Response() *response.TopCommonArchive {
	return &a.TopCommonArchive
}

func NewTopCommonArchives(schedules entity.Schedules) TopCommonArchives {
	res := make(TopCommonArchives, len(schedules))
	for i := range schedules {
		res[i] = NewTopCommonArchive(schedules[i])
	}
	return res
}

func (as TopCommonArchives) Response() []*response.TopCommonArchive {
	res := make([]*response.TopCommonArchive, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}
