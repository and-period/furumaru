package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) ListContacts(ctx context.Context, in *messenger.ListContactsInput) (entity.Contacts, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	return nil, exception.ErrNotImplemented
}

func (s *service) GetContact(ctx context.Context, in *messenger.GetContactInput) (*entity.Contact, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	return nil, exception.ErrNotImplemented
}

func (s *service) CreateContact(ctx context.Context, in *messenger.CreateContactInput) (*entity.Contact, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	return nil, exception.ErrNotImplemented
}

func (s *service) UpdateContact(ctx context.Context, in *messenger.UpdateContactInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	return exception.ErrNotImplemented
}
