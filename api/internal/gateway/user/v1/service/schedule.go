package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
)

// ScheduleStatus - 開催状況
type ScheduleStatus int32

const (
	ScheduleStatusUnknown  ScheduleStatus = 0
	ScheduleStatusWaiting  ScheduleStatus = 1 // ライブ配信前
	ScheduleStatusLive     ScheduleStatus = 2 // ライブ配信中
	ScheduleStatusClosed   ScheduleStatus = 3 // ライブ配信終了
	ScheduleStatusArchived ScheduleStatus = 4 // アーカイブ配信
)

func NewScheduleStatus(status sentity.ScheduleStatus, archived bool) ScheduleStatus {
	switch status {
	case sentity.ScheduleStatusWaiting:
		return ScheduleStatusWaiting
	case sentity.ScheduleStatusLive:
		return ScheduleStatusLive
	case sentity.ScheduleStatusClosed:
		if archived {
			return ScheduleStatusArchived
		}
		return ScheduleStatusClosed
	default:
		return ScheduleStatusUnknown
	}
}

func (s ScheduleStatus) Response() int32 {
	return int32(s)
}

type Schedule struct {
	response.Schedule
}

type Schedules []*Schedule

func NewSchedule(schedule *sentity.Schedule, broadcast *mentity.Broadcast) *Schedule {
	distributionURL, archived, metadata := newBroadcastDetail(broadcast)
	return &Schedule{
		Schedule: response.Schedule{
			ID:                   schedule.ID,
			CoordinatorID:        schedule.CoordinatorID,
			Status:               NewScheduleStatus(schedule.Status, archived).Response(),
			Title:                schedule.Title,
			Description:          schedule.Description,
			ThumbnailURL:         schedule.ThumbnailURL,
			DistributionURL:      distributionURL,
			DistributionMetadata: metadata,
			StartAt:              schedule.StartAt.Unix(),
			EndAt:                schedule.EndAt.Unix(),
		},
	}
}

func newBroadcastDetail(broadcast *mentity.Broadcast) (string, bool, *response.ScheduleDistributionMetadata) {
	metadata := &response.ScheduleDistributionMetadata{
		Subtitles: map[string]string{},
	}
	if broadcast == nil {
		return "", false, metadata
	}
	switch {
	case broadcast.ArchiveURL != "":
		if broadcast.ArchiveMetadata != nil {
			metadata.Subtitles = broadcast.ArchiveMetadata.Subtitles
		}
		return broadcast.ArchiveURL, true, metadata
	case broadcast.Status == mentity.BroadcastStatusActive:
		return broadcast.OutputURL, false, metadata
	default:
		return "", false, metadata
	}
}

func (s *Schedule) Response() *response.Schedule {
	return &s.Schedule
}

func NewSchedules(schedules sentity.Schedules, broadcasts map[string]*mentity.Broadcast) Schedules {
	res := make(Schedules, len(schedules))
	for i := range schedules {
		res[i] = NewSchedule(schedules[i], broadcasts[schedules[i].ID])
	}
	return res
}

func (ss Schedules) Response() []*response.Schedule {
	res := make([]*response.Schedule, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
