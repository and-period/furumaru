package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
)

// ScheduleStatus - 開催状況
type ScheduleStatus types.ScheduleStatus

func NewScheduleStatus(status sentity.ScheduleStatus, archived bool) ScheduleStatus {
	switch status {
	case sentity.ScheduleStatusWaiting:
		return ScheduleStatus(types.ScheduleStatusWaiting)
	case sentity.ScheduleStatusLive:
		return ScheduleStatus(types.ScheduleStatusLive)
	case sentity.ScheduleStatusClosed:
		if archived {
			return ScheduleStatus(types.ScheduleStatusArchived)
		}
		return ScheduleStatus(types.ScheduleStatusClosed)
	default:
		return ScheduleStatus(types.ScheduleStatusUnknown)
	}
}

func (s ScheduleStatus) Response() types.ScheduleStatus {
	return types.ScheduleStatus(s)
}

type Schedule struct {
	types.Schedule
}

type Schedules []*Schedule

func NewSchedule(schedule *sentity.Schedule, broadcast *mentity.Broadcast) *Schedule {
	distributionURL, archived, metadata := newBroadcastDetail(broadcast)
	return &Schedule{
		Schedule: types.Schedule{
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

func newBroadcastDetail(broadcast *mentity.Broadcast) (string, bool, *types.ScheduleDistributionMetadata) {
	metadata := &types.ScheduleDistributionMetadata{
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

func (s *Schedule) Response() *types.Schedule {
	return &s.Schedule
}

func NewSchedules(schedules sentity.Schedules, broadcasts map[string]*mentity.Broadcast) Schedules {
	res := make(Schedules, len(schedules))
	for i := range schedules {
		res[i] = NewSchedule(schedules[i], broadcasts[schedules[i].ID])
	}
	return res
}

func (ss Schedules) Response() []*types.Schedule {
	res := make([]*types.Schedule, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
