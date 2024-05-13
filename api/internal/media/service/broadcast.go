package service

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (s *service) ListBroadcasts(ctx context.Context, in *media.ListBroadcastsInput) (entity.Broadcasts, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders := make([]*database.ListBroadcastsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListBroadcastsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListBroadcastsParams{
		ScheduleIDs:   in.ScheduleIDs,
		CoordinatorID: in.CoordinatorID,
		OnlyArchived:  in.OnlyArchived,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
		Orders:        orders,
	}
	var (
		broadcasts entity.Broadcasts
		total      int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		broadcasts, err = s.db.Broadcast.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Broadcast.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return broadcasts, total, nil
}

func (s *service) GetBroadcastByScheduleID(ctx context.Context, in *media.GetBroadcastByScheduleIDInput) (*entity.Broadcast, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	return broadcast, internalError(err)
}

func (s *service) CreateBroadcast(ctx context.Context, in *media.CreateBroadcastInput) (*entity.Broadcast, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewBroadcastParams{
		ScheduleID:    in.ScheduleID,
		CoordinatorID: in.CoordinatorID,
	}
	broadcast := entity.NewBroadcast(params)
	if err := s.db.Broadcast.Create(ctx, broadcast); err != nil {
		return nil, internalError(err)
	}
	return broadcast, nil
}

func (s *service) UpdateBroadcastArchive(ctx context.Context, in *media.UpdateBroadcastArchiveInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusDisabled {
		return fmt.Errorf("service: this broadcast is not disabled: %w", exception.ErrFailedPrecondition)
	}
	params := &database.UpdateBroadcastParams{UploadBroadcastArchiveParams: &database.UploadBroadcastArchiveParams{
		ArchiveURL:   in.ArchiveURL,
		ArchiveFixed: true, // ライブ配信時のコメントとの対応が取れなくなるため、編集済みにする
	}}
	err = s.db.Broadcast.Update(ctx, broadcast.ID, params)
	return internalError(err)
}

func (s *service) PauseBroadcast(ctx context.Context, in *media.PauseBroadcastInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusActive {
		return fmt.Errorf("service: this broadcase is not activated: %w", exception.ErrFailedPrecondition)
	}
	settings := []*medialive.ScheduleSetting{{
		Name:       fmt.Sprintf("%s immediate-pause", jst.Format(s.now(), time.DateTime)),
		ActionType: medialive.ScheduleActionTypePauseState,
		StartType:  medialive.ScheduleStartTypeImmediate,
		Reference:  string(medialive.PipelineIDPipeline0),
	}}
	params := &medialive.CreateScheduleParams{
		ChannelID: broadcast.MediaLiveChannelID,
		Settings:  settings,
	}
	err = s.media.CreateSchedule(ctx, params)
	return internalError(err)
}

func (s *service) UnpauseBroadcast(ctx context.Context, in *media.UnpauseBroadcastInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusActive {
		return fmt.Errorf("service: this broadcase is not activated: %w", exception.ErrFailedPrecondition)
	}
	settings := []*medialive.ScheduleSetting{{
		Name:       fmt.Sprintf("%s immediate-unpause", jst.Format(s.now(), time.DateTime)),
		ActionType: medialive.ScheduleActionTypeUnpauseState,
		StartType:  medialive.ScheduleStartTypeImmediate,
	}}
	params := &medialive.CreateScheduleParams{
		ChannelID: broadcast.MediaLiveChannelID,
		Settings:  settings,
	}
	err = s.media.CreateSchedule(ctx, params)
	return internalError(err)
}

func (s *service) ActivateBroadcastRTMP(ctx context.Context, in *media.ActivateBroadcastRTMPInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusActive {
		return fmt.Errorf("service: this broadcase is not activated: %w", exception.ErrFailedPrecondition)
	}
	settings := []*medialive.ScheduleSetting{{
		Name:       fmt.Sprintf("%s immediate-input-rtmp", jst.Format(s.now(), time.DateTime)),
		ActionType: medialive.ScheduleActionTypeInputSwitch,
		StartType:  medialive.ScheduleStartTypeImmediate,
		Reference:  broadcast.MediaLiveRTMPInputName,
	}}
	params := &medialive.CreateScheduleParams{
		ChannelID: broadcast.MediaLiveChannelID,
		Settings:  settings,
	}
	err = s.media.CreateSchedule(ctx, params)
	return internalError(err)
}

func (s *service) ActivateBroadcastMP4(ctx context.Context, in *media.ActivateBroadcastMP4Input) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusActive {
		return fmt.Errorf("service: this broadcase is not activated: %w", exception.ErrFailedPrecondition)
	}
	videoURI, err := s.tmp.ReplaceURLToS3URI(in.InputURL)
	if err != nil {
		return internalError(err)
	}
	settings := []*medialive.ScheduleSetting{{
		Name:       fmt.Sprintf("%s immediate-input-mp4", jst.Format(s.now(), time.DateTime)),
		ActionType: medialive.ScheduleActionTypeInputSwitch,
		StartType:  medialive.ScheduleStartTypeImmediate,
		Reference:  broadcast.MediaLiveMP4InputName,
		Source:     videoURI,
	}}
	params := &medialive.CreateScheduleParams{
		ChannelID: broadcast.MediaLiveChannelID,
		Settings:  settings,
	}
	err = s.media.CreateSchedule(ctx, params)
	return internalError(err)
}

func (s *service) ActivateBroadcastStaticImage(ctx context.Context, in *media.ActivateBroadcastStaticImageInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusActive {
		return fmt.Errorf("service: this broadcase is not activated: %w", exception.ErrFailedPrecondition)
	}
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: in.ScheduleID,
	}
	schedule, err := s.store.GetSchedule(ctx, scheduleIn)
	if err != nil {
		return internalError(err)
	}
	imageURI, err := s.storage.ReplaceURLToS3URI(schedule.ImageURL)
	if err != nil {
		return internalError(err)
	}
	err = s.media.ActivateStaticImage(ctx, broadcast.MediaLiveChannelID, imageURI)
	return internalError(err)
}

func (s *service) DeactivateBroadcastStaticImage(ctx context.Context, in *media.DeactivateBroadcastStaticImageInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusActive {
		return fmt.Errorf("service: this broadcase is not activated: %w", exception.ErrFailedPrecondition)
	}
	err = s.media.DeactivateStaticImage(ctx, broadcast.MediaLiveChannelID)
	return internalError(err)
}

func (s *service) AuthYoutubeBroadcast(ctx context.Context, in *media.AuthYoutubeBroadcastInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return "", internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusDisabled {
		return "", fmt.Errorf("service: this broadcast is not disabled: %w", exception.ErrFailedPrecondition)
	}
	config := s.newYoutubeBroadcastConfig()
	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("schedule_id", broadcast.ScheduleID),
	}
	return config.AuthCodeURL(in.State, opts...), nil
}

func (s *service) CreateYoutubeBroadcast(ctx context.Context, in *media.CreateYoutubeBroadcastInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return fmt.Errorf("service: failed to get broadcast: %s: %w", err.Error(), exception.ErrInternal)
		// return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusDisabled {
		return fmt.Errorf("service: this broadcast is not disabled: %w", exception.ErrFailedPrecondition)
	}
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: in.ScheduleID,
	}
	schedule, err := s.store.GetSchedule(ctx, scheduleIn)
	if err != nil {
		return internalError(err)
	}
	// TODO: youtubeクライアントを作成する(動作検証のため、いったん雑に実装)
	config := s.newYoutubeBroadcastConfig()
	token, err := config.Exchange(ctx, in.AuthCode)
	if err != nil {
		return fmt.Errorf("service: failed to exchange token: %s: %w", err.Error(), exception.ErrForbidden)
	}
	source := config.TokenSource(ctx, token)
	s.logger.Debug("youtube broadcast token", zap.Any("token", token), zap.Any("source", source))
	service, err := youtube.NewService(ctx, option.WithTokenSource(source))
	if err != nil {
		return fmt.Errorf("service: failed to create youtube service: %s: %w", err.Error(), exception.ErrInternal)
		// return internalError(err)
	}
	youtubeIn := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              schedule.Title,
			Description:        schedule.Description,
			ScheduledStartTime: schedule.StartAt.Format(time.RFC3339),
			ScheduledEndTime:   schedule.EndAt.Format(time.RFC3339),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: "unlisted",
		},
		ContentDetails: &youtube.LiveBroadcastContentDetails{
			EnableAutoStart: false,
			EnableAutoStop:  true,
		},
	}
	part := []string{"id", "snippet", "contentDetails", "status"}
	youtubeOut, err := service.LiveBroadcasts.Insert(part, youtubeIn).Do()
	if err != nil {
		return fmt.Errorf("service: failed to insert broadcast: %s: %w", err.Error(), exception.ErrInternal)
		// return internalError(err)
	}
	streamOut, err := service.LiveStreams.List([]string{"id", "snippet", "cdn"}).Id(youtubeOut.Id).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("service: failed to list streams: %s: %w", err.Error(), exception.ErrInternal)
		// return internalError(err)
	}
	if len(streamOut.Items) == 0 {
		return fmt.Errorf("service: list streams is empty: %w", exception.ErrInternal)
	}
	params := &database.UpdateBroadcastParams{UpsertYoutubeBroadcastParams: &database.UpsertYoutubeBroadcastParams{
		YoutubeStreamURL: streamOut.Items[0].Cdn.IngestionInfo.IngestionAddress,
		YoutubeStreamKey: streamOut.Items[0].Cdn.IngestionInfo.StreamName,
		YoutubeBackupURL: streamOut.Items[0].Cdn.IngestionInfo.BackupIngestionAddress,
	}}
	err = s.db.Broadcast.Update(ctx, broadcast.ID, params)
	return internalError(err)
}

func (s *service) newYoutubeBroadcastConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.googleClientID,
		ClientSecret: s.googleClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{youtube.YoutubeScope},
		RedirectURL:  entity.NewAdminURLMaker(s.adminWebURL()).AuthYoutubeCallback(),
	}
}
