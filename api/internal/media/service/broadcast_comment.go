package service

import (
	"context"

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
	params := &database.ListBroadcastCommentsParams{
		BroadcastID:  broadcast.ID,
		WithDisabled: in.WithDisabled,
		CreatedAtGte: in.CreatedAtGte,
		CreatedAtLt:  in.CreatedAtLt,
		Limit:        in.Limit,
		NextToken:    in.NextToken,
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

func (s *service) CreateBroadcastGuestComment(
	ctx context.Context, in *media.CreateBroadcastGuestCommentInput,
) (*entity.BroadcastComment, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	if err != nil {
		return nil, internalError(err)
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

func (s *service) UpdateBroadcastComment(ctx context.Context, in *media.UpdateBroadcastCommentInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateBroadcastCommentParams{
		Disabled: in.Disabled,
	}
	err := s.db.BroadcastComment.Update(ctx, in.CommentID, params)
	return internalError(err)
}
