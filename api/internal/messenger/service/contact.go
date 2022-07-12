package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"go.uber.org/zap"
)

func (s *service) ListContacts(ctx context.Context, in *messenger.ListContactsInput) (entity.Contacts, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListContactsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	contacts, err := s.db.Contact.List(ctx, params)
	return contacts, exception.InternalError(err)
}

func (s *service) GetContact(ctx context.Context, in *messenger.GetContactInput) (*entity.Contact, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	contact, err := s.db.Contact.Get(ctx, in.ContactID)
	return contact, exception.InternalError(err)
}

func (s *service) CreateContact(ctx context.Context, in *messenger.CreateContactInput) (*entity.Contact, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	contact := entity.NewContact(in.Title, in.Content, in.Username, in.Email, in.PhoneNumber)
	if err := s.db.Contact.Create(ctx, contact); err != nil {
		return nil, exception.InternalError(err)
	}
	s.waitGroup.Add(1)
	go func(contactID, name, email string) {
		defer s.waitGroup.Done()
		in := &messenger.NotifyReceivedContactInput{
			ContactID: contactID,
			Username:  name,
			Email:     email,
		}
		if err := s.NotifyReceivedContact(context.Background(), in); err != nil {
			s.logger.Error("Failed to notify received contact", zap.String("contactId", contactID), zap.Error(err))
		}
	}(contact.ID, in.Username, in.Email)
	return contact, nil
}

func (s *service) UpdateContact(ctx context.Context, in *messenger.UpdateContactInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateContactParams{
		Status:   in.Status,
		Priority: in.Priority,
		Note:     in.Note,
	}
	err := s.db.Contact.Update(ctx, in.ContactID, params)
	return exception.InternalError(err)
}
