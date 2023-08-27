//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/media/$GOPACKAGE/$GOFILE
package database

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
)

type Params struct {
	Database *database.Client
	DynamoDB dynamodb.Client
}

type Database struct {
	Broadcast Broadcast
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Broadcast: NewBroadcast(params.Database),
	}
}

/**
 * interface
 */
type Broadcast interface {
	GetByScheduleID(ctx context.Context, scheduleID string, fields ...string) (*entity.Broadcast, error)
	Create(ctx context.Context, broadcast *entity.Broadcast) error
	Update(ctx context.Context, broadcastID string, params *UpdateBroadcastParams) error
}

/**
 * params
 */
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
