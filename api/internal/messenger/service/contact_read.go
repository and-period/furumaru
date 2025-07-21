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

func (s *service) CreateContactRead(
	ctx context.Context,
	in *messenger.CreateContactReadInput,
) (*entity.ContactRead, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	_, err := s.db.Contact.Get(ctx, in.ContactID)
	if errors.Is(err, database.ErrNotFound) {
		return nil, fmt.Errorf(
			"service: invalid contact id: %s: %w",
			err.Error(),
			exception.ErrInvalidArgument,
		)
	}
	if err != nil {
		return nil, internalError(err)
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
		return nil, internalError(err)
	}
	return contactRead, nil
}
