package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

func (s *service) ListVideoComments(ctx context.Context, in *media.ListVideoCommentsInput) (entity.VideoComments, string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, "", internalError(err)
	}
	params := &database.ListVideoCommentsParams{
		VideoID:      in.VideoID,
		WithDisabled: in.WithDisabled,
		CreatedAtGte: in.CreatedAtGte,
		CreatedAtLt:  in.CreatedAtLt,
		Limit:        in.Limit,
		NextToken:    in.NextToken,
	}
	comments, token, err := s.db.VideoComment.List(ctx, params)
	return comments, token, internalError(err)
}

func (s *service) CreateVideoComment(ctx context.Context, in *media.CreateVideoCommentInput) (*entity.VideoComment, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	video, err := s.db.Video.Get(ctx, in.VideoID)
	if err != nil {
		return nil, internalError(err)
	}
	if !video.Published() {
		return nil, fmt.Errorf("service: this video is not published: %w", exception.ErrFailedPrecondition)
	}
	params := &entity.NewVideoCommentParams{
		VideoID: in.VideoID,
		UserID:  in.UserID,
		Content: in.Content,
	}
	comment := entity.NewVideoComment(params)
	if err := s.db.VideoComment.Create(ctx, comment); err != nil {
		return nil, internalError(err)
	}
	return comment, nil
}

func (s *service) CreateVideoGuestComment(ctx context.Context, in *media.CreateVideoGuestCommentInput) (*entity.VideoComment, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	video, err := s.db.Video.Get(ctx, in.VideoID)
	if err != nil {
		return nil, internalError(err)
	}
	if !video.Published() {
		return nil, fmt.Errorf("service: this video is not published: %w", exception.ErrFailedPrecondition)
	}
	params := &entity.NewVideoCommentParams{
		VideoID: in.VideoID,
		Content: in.Content,
	}
	comment := entity.NewVideoComment(params)
	if err := s.db.VideoComment.Create(ctx, comment); err != nil {
		return nil, internalError(err)
	}
	return comment, nil
}

func (s *service) UpdateVideoComment(ctx context.Context, in *media.UpdateVideoCommentInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &database.UpdateVideoCommentParams{
		Disabled: in.Disabled,
	}
	err := s.db.VideoComment.Update(ctx, in.CommentID, params)
	return internalError(err)
}
