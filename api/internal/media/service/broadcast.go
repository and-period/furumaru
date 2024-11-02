package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/youtube"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListBroadcasts(ctx context.Context, in *media.ListBroadcastsInput) (entity.Broadcasts, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders, err := s.newListBroadcastsOrders(in.Orders)
	if err != nil {
		return nil, 0, fmt.Errorf("service: failed to create broadcast orders: err=%s: %w", err.Error(), exception.ErrInvalidArgument)
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

func (s *service) newListBroadcastsOrders(in []*media.ListBroadcastsOrder) ([]*database.ListBroadcastsOrder, error) {
	res := make([]*database.ListBroadcastsOrder, len(in))
	for i := range in {
		var key database.ListBroadcastsOrderKey
		switch in[i].Key {
		case media.ListBroadcastsOrderByUpdatedAt:
			key = database.ListBroadcastsOrderByUpdatedAt
		default:
			return nil, errors.New("service: invalid broadcast order key")
		}
		res[i] = &database.ListBroadcastsOrder{
			Key:        key,
			OrderByASC: in[i].OrderByASC,
		}
	}
	return res, nil
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

func (s *service) GetBroadcastAuth(ctx context.Context, in *media.GetBroadcastAuthInput) (*entity.BroadcastAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	auth := &entity.BroadcastAuth{SessionID: in.SessionID}
	err := s.cache.Get(ctx, auth)
	if errors.Is(err, dynamodb.ErrNotFound) {
		return nil, fmt.Errorf("service: invalid session: %w", exception.ErrUnauthenticated)
	}
	if err != nil {
		return nil, internalError(err)
	}
	return auth, nil
}

func (s *service) AuthYoutubeBroadcast(ctx context.Context, in *media.AuthYoutubeBroadcastInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	if !strings.HasPrefix(in.YoutubeHandle, "@") {
		return "", fmt.Errorf("service: invalid youtube handle: %w", exception.ErrInvalidArgument)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return "", internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusDisabled {
		return "", fmt.Errorf("service: this broadcast is not disabled: %w", exception.ErrFailedPrecondition)
	}
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: broadcast.ScheduleID,
	}
	schedule, err := s.store.GetSchedule(ctx, scheduleIn)
	if err != nil {
		return "", internalError(err)
	}
	if schedule.Status != sentity.ScheduleStatusWaiting {
		return "", fmt.Errorf("service: this schedule is not waiting: %w", exception.ErrFailedPrecondition)
	}
	sessionID := s.generateID()
	params := &entity.BroadcastAuthParams{
		SessionID:  sessionID,
		Account:    in.YoutubeHandle,
		ScheduleID: in.ScheduleID,
		Now:        s.now(),
		TTL:        s.authYoutubeTTL,
	}
	auth := entity.NewYoutubeBroadcastAuth(params)
	if err := s.cache.Insert(ctx, auth); err != nil {
		return "", internalError(err)
	}
	return s.youtube.NewAuth().GetAuthCodeURL(sessionID), nil
}

func (s *service) AuthYoutubeBroadcastEvent(ctx context.Context, in *media.AuthYoutubeBroadcastEventInput) (*entity.BroadcastAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	auth := &entity.BroadcastAuth{SessionID: in.State}
	err := s.cache.Get(ctx, auth)
	if errors.Is(err, dynamodb.ErrNotFound) {
		return nil, fmt.Errorf("service: invalid session: %w", exception.ErrUnauthenticated)
	}
	if err != nil {
		return nil, internalError(err)
	}
	token, err := s.youtube.NewAuth().GetToken(ctx, in.AuthCode)
	if err != nil {
		return nil, internalError(err)
	}
	if !token.Valid() {
		return nil, fmt.Errorf("service: invalid youtube token: %w", exception.ErrUnauthenticated)
	}
	service, err := s.youtube.NewService(ctx, token)
	if err != nil {
		return nil, internalError(err)
	}
	channels, err := service.ListChannels(ctx)
	if err != nil {
		return nil, internalError(err)
	}
	if !auth.ValidYoutubeAuth(channels) {
		return nil, fmt.Errorf("service: invalid youtube auth: %w", exception.ErrForbidden)
	}
	if err := auth.SetToken(token); err != nil {
		return nil, internalError(err)
	}
	if err := s.cache.Insert(ctx, auth); err != nil {
		return nil, internalError(err)
	}
	return auth, nil
}

func (s *service) CreateYoutubeBroadcast(ctx context.Context, in *media.CreateYoutubeBroadcastInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	auth := &entity.BroadcastAuth{SessionID: in.SessionID}
	err := s.cache.Get(ctx, auth)
	if errors.Is(err, dynamodb.ErrNotFound) {
		return fmt.Errorf("service: invalid session: %w", exception.ErrUnauthenticated)
	}
	if err != nil {
		return internalError(err)
	}
	token, err := auth.GetToken()
	if err != nil {
		return fmt.Errorf("service: invalid youtube token. err=%s: %w", err.Error(), exception.ErrUnauthenticated)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, auth.ScheduleID)
	if err != nil {
		return internalError(err)
	}
	if broadcast.Status != entity.BroadcastStatusDisabled {
		return fmt.Errorf("service: this broadcast is not disabled: %w", exception.ErrFailedPrecondition)
	}
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: broadcast.ScheduleID,
	}
	schedule, err := s.store.GetSchedule(ctx, scheduleIn)
	if err != nil {
		return internalError(err)
	}
	if schedule.Status != sentity.ScheduleStatusWaiting {
		return fmt.Errorf("service: this schedule is not waiting: %w", exception.ErrFailedPrecondition)
	}
	service, err := s.youtube.NewService(ctx, token)
	if err != nil {
		return internalError(err)
	}
	youtubeIn := &youtube.CreateLiveBroadcastParams{
		Title:       in.Title,
		Description: in.Description,
		Public:      in.Public,
		StartAt:     schedule.StartAt,
		EndAt:       schedule.EndAt,
	}
	liveBroadcast, err := service.CreateLiveBroadcast(ctx, youtubeIn)
	if err != nil {
		return internalError(err)
	}
	streamIn := &youtube.CreateLiveStreamParams{
		Title: schedule.Title,
	}
	stream, err := service.CreateLiveStream(ctx, streamIn)
	if err != nil {
		return internalError(err)
	}
	if err := service.BindLiveBroadcast(ctx, liveBroadcast.Id, stream.Id); err != nil {
		return internalError(err)
	}
	params := &database.UpdateBroadcastParams{UpsertYoutubeBroadcastParams: &database.UpsertYoutubeBroadcastParams{
		YoutubeAccount:     auth.Account,
		YoutubeBroadcastID: liveBroadcast.Id,
		YoutubeStreamID:    stream.Id,
		YoutubeStreamURL:   stream.Cdn.IngestionInfo.IngestionAddress,
		YoutubeStreamKey:   stream.Cdn.IngestionInfo.StreamName,
		YoutubeBackupURL:   stream.Cdn.IngestionInfo.BackupIngestionAddress,
	}}
	err = s.db.Broadcast.Update(ctx, broadcast.ID, params)
	return internalError(err)
}
