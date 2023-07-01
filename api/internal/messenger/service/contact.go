package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListContacts(ctx context.Context, in *messenger.ListContactsInput) (entity.Contacts, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}

	params := &database.ListContactsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		contacts entity.Contacts
		total    int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		contacts, err = s.db.Contact.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Contact.Count(ectx)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return contacts, total, nil
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
	params := &entity.NewContactParams{
		Title:       in.Title,
		Content:     in.Content,
		Username:    in.Username,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Note:        in.Note,
	}
	contact := entity.NewContact(params)
	contact.Fill(in.CategoryID, in.UserID, in.ResponderID)
	if err := s.db.Contact.Create(ctx, contact); err != nil {
		return nil, exception.InternalError(err)
	}

	return contact, nil
}

func (s *service) UpdateContact(ctx context.Context, in *messenger.UpdateContactInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	if _, err := s.db.Contact.Get(ctx, in.ContactID); err != nil {
		return exception.InternalError(err)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if in.UserID == "" {
			return nil
		}
		userID := &user.GetUserInput{
			UserID: in.UserID,
		}
		_, err := s.user.GetUser(ectx, userID)
		return err
	})
	eg.Go(func() error {
		if in.ResponderID == "" {
			return nil
		}
		adminID := &user.GetAdminInput{
			AdminID: in.ResponderID,
		}
		_, err := s.user.GetAdmin(ectx, adminID)
		return err
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid user id format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateContactParams{
		Title:       in.Title,
		CategoryID:  in.CategoryID,
		Content:     in.Content,
		Username:    in.Username,
		UserID:      in.UserID,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Status:      in.Status,
		ResponderID: in.ResponderID,
		Note:        in.Note,
	}
	if err := s.db.Contact.Update(ctx, in.ContactID, params); err != nil {
		return exception.InternalError(err)
	}
	return nil
}

func (s *service) DeleteContact(ctx context.Context, in *messenger.DeleteContactInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Contact.Delete(ctx, in.ContactID)
	return exception.InternalError(err)
}
