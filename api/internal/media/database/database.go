//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/media/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/media/entity"
)

var (
	ErrInvalidArgument    = &Error{err: errors.New("database: invalid argument")}
	ErrNotFound           = &Error{err: errors.New("database: not found")}
	ErrAlreadyExists      = &Error{err: errors.New("database: already exists")}
	ErrFailedPrecondition = &Error{err: errors.New("database: failed precondition")}
	ErrCanceled           = &Error{err: errors.New("database: canceled")}
	ErrDeadlineExceeded   = &Error{err: errors.New("database: deadline exceeded")}
	ErrInternal           = &Error{err: errors.New("database: internal error")}
	ErrUnknown            = &Error{err: errors.New("database: unknown")}
)

type Database struct {
	Broadcast          Broadcast
	BroadcastComment   BroadcastComment
	BroadcastViewerLog BroadcastViewerLog
	Video              Video
}

type Broadcast interface {
	List(ctx context.Context, params *ListBroadcastsParams, fields ...string) (entity.Broadcasts, error)
	Count(ctx context.Context, params *ListBroadcastsParams) (int64, error)
	GetByScheduleID(ctx context.Context, scheduleID string, fields ...string) (*entity.Broadcast, error)
	Create(ctx context.Context, broadcast *entity.Broadcast) error
	Update(ctx context.Context, broadcastID string, params *UpdateBroadcastParams) error
}

type ListBroadcastsParams struct {
	ScheduleIDs   []string
	CoordinatorID string
	OnlyArchived  bool
	Limit         int
	Offset        int
	Orders        []*ListBroadcastsOrder
}

type ListBroadcastsOrder struct {
	Key        entity.BroadcastOrderBy
	OrderByASC bool
}

type UpdateBroadcastParams struct {
	Status entity.BroadcastStatus
	*InitializeBroadcastParams
	*UploadBroadcastArchiveParams
	*UpsertYoutubeBroadcastParams
}

type InitializeBroadcastParams struct {
	InputURL                  string
	OutputURL                 string
	CloudFrontDistributionArn string
	MediaLiveChannelArn       string
	MediaLiveChannelID        string
	MediaLiveRTMPInputArn     string
	MediaLiveRTMPInputName    string
	MediaLiveMP4InputArn      string
	MediaLiveMP4InputName     string
	MediaStoreContainerArn    string
}

type UploadBroadcastArchiveParams struct {
	ArchiveURL   string
	ArchiveFixed bool
}

type UpsertYoutubeBroadcastParams struct {
	YoutubeAccount     string
	YoutubeBroadcastID string
	YoutubeStreamID    string
	YoutubeStreamURL   string
	YoutubeStreamKey   string
	YoutubeBackupURL   string
}

type BroadcastComment interface {
	List(ctx context.Context, params *ListBroadcastCommentsParams, fields ...string) (entity.BroadcastComments, string, error)
	Create(ctx context.Context, comment *entity.BroadcastComment) error
	Update(ctx context.Context, commentID string, params *UpdateBroadcastCommentParams) error
}

type ListBroadcastCommentsParams struct {
	BroadcastID  string
	WithDisabled bool
	CreatedAtGte time.Time
	CreatedAtLt  time.Time
	Limit        int64
	NextToken    string
	Orders       []*ListBroadcastCommentsOrder
}

type ListBroadcastCommentsOrder struct {
	Key        entity.BroadcastCommentOrderBy
	OrderByASC bool
}

type UpdateBroadcastCommentParams struct {
	Disabled bool
}

type BroadcastViewerLog interface {
	Create(ctx context.Context, log *entity.BroadcastViewerLog) error
	GetTotal(ctx context.Context, params *GetBroadcastTotalViewersParams) (int64, error)
	Aggregate(ctx context.Context, params *AggregateBroadcastViewerLogsParams) (entity.AggregatedBroadcastViewerLogs, error)
}

type GetBroadcastTotalViewersParams struct {
	BroadcastID  string
	CreatedAtGte time.Time
	CreatedAtLt  time.Time
}

type AggregateBroadcastViewerLogsParams struct {
	BroadcastID  string
	Interval     entity.AggregateBroadcastViewerLogInterval
	CreatedAtGte time.Time
	CreatedAtLt  time.Time
}

type Video interface {
	List(ctx context.Context, params *ListVideosParams, fields ...string) (entity.Videos, error)
	Count(ctx context.Context, params *ListVideosParams) (int64, error)
	Get(ctx context.Context, videoID string, fields ...string) (*entity.Video, error)
	Create(ctx context.Context, video *entity.Video) error
	Update(ctx context.Context, videoID string, params *UpdateVideoParams) error
	Delete(ctx context.Context, videoID string) error
}

type ListVideosParams struct {
	CoordinatorID string
	Limit         int
	Offset        int
}

type UpdateVideoParams struct {
	Title         string
	Description   string
	ProductIDs    []string
	ExperienceIDs []string
	ThumbnailURL  string
	VideoURL      string
	Public        bool
	Limited       bool
	PublishedAt   time.Time
}

type Error struct {
	err error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}
