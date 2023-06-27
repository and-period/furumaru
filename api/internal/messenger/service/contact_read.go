package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) CreateContactRead(ctx context.Context, in *messenger.CreateContactReadInput) (*entity.ContactRead, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	contactIn := &messenger.GetContactInput{
		ContactID: in.ContactID,
	}
	_, err := s.GetContact(ctx, contactIn)
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid contact id: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, exception.InternalError(err)
	}
	params := &entity.NewContactReadParams{
		ContactID: in.ContactID,
		UserID:    in.UserID,
		UserType:  in.UserType,
		Read:      false,
	}
	contactRead, err := entity.NewContactRead(params)
	if err != nil {
		return nil, exception.ErrInvalidArgument
	}
	if err := s.db.ContactRead.Create(ctx, contactRead); err != nil {
		return nil, exception.InternalError(err)
	}
	return contactRead, exception.InternalError(err)
}
