//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/media/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"

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
	Broadcast Broadcast
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

type Error struct {
	err error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}
