package medialive

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/medialive"
	"github.com/aws/aws-sdk-go-v2/service/medialive/types"
	"go.uber.org/zap"
)

// ScheduleStartType - 開始タイプ
// @see - https://docs.aws.amazon.com/ja_jp/medialive/latest/ug/ips-switch-types.html
type ScheduleStartType string

const (
	ScheduleStartTypeImmediate ScheduleStartType = "IMMEDIATE" // 即時
	ScheduleStartTypeFollow    ScheduleStartType = "FOLLOW"    // フォロー
	ScheduleStartTypeFixed     ScheduleStartType = "FIXED"     // 固定
)

// ScheduleActionType - 設定種別
// @see - https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/medialive@v1.34.1/types#ScheduleActionSettings
type ScheduleActionType string

const (
	ScheduleActionTypeStaticImageActivate   ScheduleActionType = "STATIC_IMAGE_ACTIVE"     // 静的イメージのアクティブ化
	ScheduleActionTypeStaticImageDeactivate ScheduleActionType = "STATIC_IMAGE_DEACTIVATE" // 静的イメージの非アクティブ化
	ScheduleActionTypeInputSwitch           ScheduleActionType = "INPUT_SWITCH"            // 入力スイッチ
)

type CreateScheduleParams struct {
	ChannelID string             // MediaLiveチャンネルID
	Settings  []*ScheduleSetting // スケジュール設定一覧
}

type ScheduleSetting struct {
	Name       string             // 設定名
	ActionType ScheduleActionType // 設定種別
	StartType  ScheduleStartType  // 開始タイプ
	ExecutedAt time.Time          // 実行時間(開始タイプが固定の場合のみ指定)
	Reference  string             // インプット名
}

func (c *client) CreateSchedule(ctx context.Context, params *CreateScheduleParams) error {
	in := &medialive.BatchUpdateScheduleInput{
		ChannelId: aws.String(params.ChannelID),
		Creates:   &types.BatchScheduleActionCreateRequest{ScheduleActions: c.newScheduleActions(params.Settings)},
	}
	_, err := c.media.BatchUpdateSchedule(ctx, in)
	return err
}

func (c *client) newScheduleActions(params []*ScheduleSetting) []types.ScheduleAction {
	actions := make([]types.ScheduleAction, len(params))
	for i, p := range params {
		actions[i] = types.ScheduleAction{
			ActionName:                  aws.String(p.Name),
			ScheduleActionStartSettings: &types.ScheduleActionStartSettings{},
			ScheduleActionSettings:      &types.ScheduleActionSettings{},
		}
		// 開始設定
		switch p.StartType {
		case ScheduleStartTypeImmediate:
			actions[i].ScheduleActionStartSettings.ImmediateModeScheduleActionStartSettings = &types.ImmediateModeScheduleActionStartSettings{}
		case ScheduleStartTypeFixed:
			actions[i].ScheduleActionStartSettings.FixedModeScheduleActionStartSettings = &types.FixedModeScheduleActionStartSettings{
				Time: aws.String(p.ExecutedAt.Format(time.RFC3339)),
			}
		case ScheduleStartTypeFollow:
			c.logger.Warn("Not implemented start type", zap.Any("setting", p))
		}
		// アクションの設定
		switch p.ActionType {
		case ScheduleActionTypeInputSwitch:
			actions[i].ScheduleActionSettings.InputSwitchSettings = &types.InputSwitchScheduleActionSettings{
				InputAttachmentNameReference: aws.String(p.Reference),
			}
		case ScheduleActionTypeStaticImageActivate:
			c.logger.Warn("Not implemented action type", zap.Any("setting", p))
		case ScheduleActionTypeStaticImageDeactivate:
			c.logger.Warn("Not implemented action type", zap.Any("setting", p))
		}
	}
	return actions
}
