package media

import (
	"time"

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
	Key string `validate:"required"`
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

type GetBroadcastAuthInput struct {
	SessionID string `validate:"required"`
}

type AuthYoutubeBroadcastInput struct {
	ScheduleID    string `validate:"required"`
	YoutubeHandle string `validate:"required"`
}

type AuthYoutubeBroadcastEventInput struct {
	State    string `validate:"required"`
	AuthCode string `validate:"required"`
}

type CreateYoutubeBroadcastInput struct {
	SessionID   string `validate:"required"`
	Title       string `validate:"required,max=100"`
	Description string `validate:"required,max=1000"`
	Public      bool   `validate:""`
}

type CreateBroadcastViewerLogInput struct {
	ScheduleID string `validate:"required"`
	SessionID  string `validate:"required"`
	UserID     string `validate:""`
	UserAgent  string `validate:""`
	ClientIP   string `validate:"omitempty,ip_addr"`
}

type AggregateBroadcastViewerLogsInput struct {
	ScheduleID   string                                     `validate:"required"`
	Interval     entity.AggregateBroadcastViewerLogInterval `validate:"required"`
	CreatedAtGte time.Time                                  `validate:""`
	CreatedAtLt  time.Time                                  `validate:""`
}

type ListBroadcastCommentsInput struct {
	ScheduleID   string                        `validate:"required"`
	CreatedAtGte time.Time                     `validate:""`
	CreatedAtLt  time.Time                     `validate:""`
	Limit        int64                         `validate:"max=200"`
	NextToken    string                        `validate:""`
	Orders       []*ListBroadcastCommentsOrder `validate:"dive,required"`
}

type ListBroadcastCommentsOrder struct {
	Key        entity.BroadcastCommentOrderBy `validate:"required"`
	OrderByASC bool                           `validate:""`
}

type CreateBroadcastCommentInput struct {
	ScheduleID string `validate:"required"`
	UserID     string `validate:"required"`
	Content    string `validate:"required,max=200"`
}

type CreateBroadcastGuestCommentInput struct {
	ScheduleID string `validate:"required"`
	Content    string `validate:"required,max=200"`
}

type UpdateBroadcastCommentInput struct {
	CommentID string `validate:"required"`
	Disabled  bool   `validate:""`
}

type ListVideosInput struct {
	Name                  string `validate:""`
	CoordinatorID         string `validate:""`
	OnlyPublished         bool   `validate:""`
	OnlyDisplayProduct    bool   `validate:""`
	OnlyDisplayExperience bool   `validate:""`
	ExcludeLimited        bool   `validate:""`
	Limit                 int64  `validate:"required_without=NoLimit,min=0,max=200"`
	Offset                int64  `validate:"min=0"`
	NoLimit               bool   `validate:""`
}

type GetVideoInput struct {
	VideoID string `validate:"required"`
}

type CreateVideoInput struct {
	Title             string    `validate:"required,max=128"`
	Description       string    `validate:"required,max=2000"`
	CoordinatorID     string    `validate:"required"`
	ProductIDs        []string  `validate:"dive,required"`
	ExperienceIDs     []string  `validate:"dive,required"`
	ThumbnailURL      string    `validate:"required,url"`
	VideoURL          string    `validate:"required,url"`
	Public            bool      `validate:""`
	Limited           bool      `validate:""`
	DisplayProduct    bool      `validate:""`
	DisplayExperience bool      `validate:""`
	PublishedAt       time.Time `validate:"required"`
}

type UpdateVideoInput struct {
	VideoID           string    `validate:"required"`
	Title             string    `validate:"required,max=128"`
	Description       string    `validate:"required,max=2000"`
	ProductIDs        []string  `validate:"dive,required"`
	ExperienceIDs     []string  `validate:"dive,required"`
	ThumbnailURL      string    `validate:"required,url"`
	VideoURL          string    `validate:"required,url"`
	Public            bool      `validate:""`
	Limited           bool      `validate:""`
	DisplayProduct    bool      `validate:""`
	DisplayExperience bool      `validate:""`
	PublishedAt       time.Time `validate:"required"`
}

type DeleteVideoInput struct {
	VideoID string `validate:"required"`
}
