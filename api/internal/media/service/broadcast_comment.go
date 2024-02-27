package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) ListBroadcastComments(
	ctx context.Context,
	in *media.ListBroadcastCommentsInput,
) (entity.BroadcastComments, string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, "", internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, "", internalError(err)
	}
	if broadcast.ArchiveFixed {
		// 編集済みアーカイブ動画がアップロードされている場合、対応が取れないためコメントは返さなくする
		return entity.BroadcastComments{}, "", nil
	}
	orders := make([]*database.ListBroadcastCommentsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListBroadcastCommentsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListBroadcastCommentsParams{
		BroadcastID:  broadcast.ID,
		WithDisabled: false,
		CreatedAtGte: in.CreatedAtGte,
		CreatedAtLt:  in.CreatedAtLt,
		Limit:        in.Limit,
		NextToken:    in.NextToken,
		Orders:       orders,
	}
	comments, token, err := s.db.BroadcastComment.List(ctx, params)
	return comments, token, internalError(err)
}

func (s *service) CreateBroadcastComment(ctx context.Context, in *media.CreateBroadcastCommentInput) (*entity.BroadcastComment, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, internalError(err)
	}
	if broadcast.Status == entity.BroadcastStatusDisabled {
		return nil, fmt.Errorf("service: broadcast is disabled: %w", exception.ErrFailedPrecondition)
	}
	params := &entity.BroadcastCommentParams{
		BroadcastID: broadcast.ID,
		UserID:      in.UserID,
		Content:     in.Content,
	}
	comment := entity.NewBroadcastComment(params)
	if err := s.db.BroadcastComment.Create(ctx, comment); err != nil {
		return nil, internalError(err)
	}
	return comment, nil
}

func (s *service) CreateBroadcastGuestComment(ctx context.Context, in *media.CreateBroadcastGuestCommentInput) (*entity.BroadcastComment, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, internalError(err)
	}
	if broadcast.Status == entity.BroadcastStatusDisabled {
		return nil, fmt.Errorf("service: broadcast is disabled: %w", exception.ErrFailedPrecondition)
	}
	params := &entity.BroadcastCommentParams{
		BroadcastID: broadcast.ID,
		Content:     in.Content,
	}
	comment := entity.NewBroadcastComment(params)
	if err := s.db.BroadcastComment.Create(ctx, comment); err != nil {
		return nil, internalError(err)
	}
	return comment, nil
}
