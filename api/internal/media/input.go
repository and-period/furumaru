package media

import (
	"github.com/and-period/furumaru/api/internal/media/entity"
)

type GenerateUploadURLInput struct {
	FileType string `validate:"required"`
}

type GenerateBroadcastArchiveMP4UploadInput struct {
	GenerateUploadURLInput
	ScheduleID string `validate:"required"`
}

type GetUploadEventInput struct {
	UploadURL string `validate:"required,url"`
}

type ResizeFileInput struct {
	TargetID string   `validate:"required"`
	URLs     []string `validate:"min=1,dive,required,url"`
}

type ListBroadcastsInput struct {
	ScheduleIDs   []string               `validate:"dive,required"`
	CoordinatorID string                 `validate:""`
	OnlyArchived  bool                   `validate:""`
	Limit         int64                  `validate:"required_without=NoLimit,min=0,max=200"`
	Offset        int64                  `validate:"min=0"`
	NoLimit       bool                   `validate:""`
	Orders        []*ListBroadcastsOrder `validate:"dive,required"`
}

type ListBroadcastsOrder struct {
	Key        entity.BroadcastOrderBy `validate:"required"`
	OrderByASC bool                    `validate:""`
}

type GetBroadcastByScheduleIDInput struct {
	ScheduleID string `validate:"required"`
}

type CreateBroadcastInput struct {
	ScheduleID    string `validate:"required"`
	CoordinatorID string `validate:"required"`
}

type UpdateBroadcastArchiveInput struct {
	ScheduleID string `validate:"required"`
	ArchiveURL string `validate:"required,url"`
}

type PauseBroadcastInput struct {
	ScheduleID string `validate:"required"`
}

type UnpauseBroadcastInput struct {
	ScheduleID string `validate:"required"`
}

type ActivateBroadcastRTMPInput struct {
	ScheduleID string `validate:"required"`
}

type ActivateBroadcastMP4Input struct {
	ScheduleID string `validate:"required"`
	InputURL   string `validate:"required,url"`
}

type ActivateBroadcastStaticImageInput struct {
	ScheduleID string `validate:"required"`
}

type DeactivateBroadcastStaticImageInput struct {
	ScheduleID string `validate:"required"`
}

type CreateBroadcastViewerLogInput struct {
	ScheduleID string `validate:"required"`
	SessionID  string `validate:"required"`
	UserID     string `validate:""`
	UserAgent  string `validate:""`
	ClientIP   string `validate:"omitempty,ip_addr"`
}
