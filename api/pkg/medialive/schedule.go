package medialive

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/medialive"
	"github.com/aws/aws-sdk-go-v2/service/medialive/types"
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
	ScheduleActionTypePauseState            ScheduleActionType = "PAUSE_STATE"             // 一時停止
	ScheduleActionTypeUnpauseState          ScheduleActionType = "UNPAUSE_STATE"           // 一時停止を解除
)

// PipelineID パイプラインID
// @see - https://docs.aws.amazon.com/ja_jp/medialive/latest/ug/x-actions-in-schedule-pause.html
type PipelineID string

const (
	PipelineIDPipeline0 PipelineID = "PIPELINE_0" // パイプライン 0
	PipelineIDPipeline1 PipelineID = "PIPELINE_1" // パイプライン 1
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
	Source     string             // 参照先情報
}

func (c *client) CreateSchedule(ctx context.Context, params *CreateScheduleParams) error {
	in := &medialive.BatchUpdateScheduleInput{
		ChannelId: aws.String(params.ChannelID),
		Creates:   &types.BatchScheduleActionCreateRequest{ScheduleActions: c.newScheduleActions(params.Settings)},
	}
	_, err := c.media.BatchUpdateSchedule(ctx, in)
	return err
}

func (c *client) ActivateStaticImage(ctx context.Context, channelID, imageURL string) error {
	name := fmt.Sprintf("%s immediate-static-image-active", c.now().Format(time.DateTime))
	settings := []*ScheduleSetting{{
		Name:       name,
		ActionType: ScheduleActionTypeStaticImageActivate,
		StartType:  ScheduleStartTypeImmediate,
		Source:     imageURL,
	}}
	in := &medialive.BatchUpdateScheduleInput{
		ChannelId: aws.String(channelID),
		Creates:   &types.BatchScheduleActionCreateRequest{ScheduleActions: c.newScheduleActions(settings)},
	}
	_, err := c.media.BatchUpdateSchedule(ctx, in)
	return err
}

func (c *client) DeactivateStaticImage(ctx context.Context, channelID string) error {
	name := fmt.Sprintf("%s immediate-static-image-deactive", c.now().Format(time.DateTime))
	settings := []*ScheduleSetting{{
		Name:       name,
		ActionType: ScheduleActionTypeStaticImageDeactivate,
		StartType:  ScheduleStartTypeImmediate,
	}}
	in := &medialive.BatchUpdateScheduleInput{
		ChannelId: aws.String(channelID),
		Creates:   &types.BatchScheduleActionCreateRequest{ScheduleActions: c.newScheduleActions(settings)},
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
				Time: aws.String(p.ExecutedAt.UTC().Format(time.RFC3339)),
			}
		case ScheduleStartTypeFollow:
			slog.Warn("Not implemented start type", slog.Any("setting", p))
		}
		// アクションの設定
		switch p.ActionType {
		case ScheduleActionTypeInputSwitch:
			urlPath := make([]string, 0, 1)
			if p.Source != "" {
				urlPath = append(urlPath, p.Source)
			}
			actions[i].ScheduleActionSettings.InputSwitchSettings = &types.InputSwitchScheduleActionSettings{
				InputAttachmentNameReference: aws.String(p.Reference),
				UrlPath:                      urlPath,
			}
		case ScheduleActionTypeStaticImageActivate:
			actions[i].ScheduleActionSettings.StaticImageActivateSettings = &types.StaticImageActivateScheduleActionSettings{
				Image:  &types.InputLocation{Uri: aws.String(p.Source)},
				FadeIn: aws.Int32(1000), // 1.0sec
				Width:  aws.Int32(1920), // フルHDサイズ
				Height: aws.Int32(1080), // フルHDサイズ
				ImageX: aws.Int32(0),
				ImageY: aws.Int32(0),
			}
		case ScheduleActionTypeStaticImageDeactivate:
			actions[i].ScheduleActionSettings.StaticImageDeactivateSettings = &types.StaticImageDeactivateScheduleActionSettings{
				FadeOut: aws.Int32(1000), // 1.0sec
			}
		case ScheduleActionTypePauseState:
			actions[i].ScheduleActionSettings.PauseStateSettings = &types.PauseStateScheduleActionSettings{
				Pipelines: []types.PipelinePauseStateSettings{{
					PipelineId: types.PipelineId(p.Reference),
				}},
			}
		case ScheduleActionTypeUnpauseState:
			actions[i].ScheduleActionSettings.PauseStateSettings = &types.PauseStateScheduleActionSettings{
				Pipelines: []types.PipelinePauseStateSettings{},
			}
		}
	}
	return actions
}
