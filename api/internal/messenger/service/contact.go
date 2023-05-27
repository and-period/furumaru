package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) GetContact(ctx context.Context, in *messenger.GetContactInput) (*entity.Contact, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	contact, err := s.db.Contact.Get(ctx, in.ContactID)
	return contact, exception.InternalError(err)
}
