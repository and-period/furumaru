package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) GetContactRead(ctx context.Context, in *messenger.GetContactReadInput) (*entity.ContactRead, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	contactRead, err := s.db.ContactRead.Get(ctx, in.ContactID, in.UserID)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	return contactRead, nil
}

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

func (s *service) UpdateContactReadFlag(ctx context.Context, in *messenger.UpdateContactReadFlagInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateContactReadFlagParams{
		ContactID: in.ContactID,
		UserID:    in.UserID,
		Read:      in.Read,
	}
	err := s.db.ContactRead.UpdateRead(ctx, params)
	return exception.InternalError(err)
}
