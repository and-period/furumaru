//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/media/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"

	"github.com/and-period/furumaru/api/internal/media/entity"
)

var (
	ErrInvalidArgument    = errors.New("database: invalid argument")
	ErrNotFound           = errors.New("database: not found")
	ErrAlreadyExists      = errors.New("database: already exists")
	ErrFailedPrecondition = errors.New("database: failed precondition")
	ErrCanceled           = errors.New("database: canceled")
	ErrDeadlineExceeded   = errors.New("database: deadline exceeded")
	ErrInternal           = errors.New("database: internal error")
	ErrUnknown            = errors.New("database: unknown")
)

type Database struct {
	Broadcast Broadcast
}

type Broadcast interface {
	GetByScheduleID(ctx context.Context, scheduleID string, fields ...string) (*entity.Broadcast, error)
	Create(ctx context.Context, broadcast *entity.Broadcast) error
	Update(ctx context.Context, broadcastID string, params *UpdateBroadcastParams) error
}

type UpdateBroadcastParams struct {
	Status entity.BroadcastStatus
	*InitializeBroadcastParams
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
