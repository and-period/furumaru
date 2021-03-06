package service

import (
	"context"
	"encoding/json"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"go.uber.org/zap"
)

func (s *service) NotifyRegisterAdmin(ctx context.Context, in *messenger.NotifyRegisterAdminInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	maker := entity.NewAdminURLMaker(s.adminWebURL())
	builder := entity.NewTemplateDataBuilder().
		WebURL(maker.SignIn()).
		Password(in.Password)
	mail := &entity.MailConfig{
		EmailID:       entity.EmailIDAdminRegister,
		Substitutions: builder.Build(),
	}
	payload := &entity.WorkerPayload{
		QueueID:   uuid.Base58Encode(uuid.New()),
		EventType: entity.EventTypeAdminRegister,
		UserType:  entity.UserTypeAdmin,
		UserIDs:   []string{in.AdminID},
		Email:     mail,
	}
	err := s.sendMessage(ctx, payload)
	return exception.InternalError(err)
}

func (s *service) NotifyReceivedContact(ctx context.Context, in *messenger.NotifyReceivedContactInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	contact, err := s.db.Contact.Get(ctx, in.ContactID)
	if err != nil {
		return exception.InternalError(err)
	}
	builder := entity.NewTemplateDataBuilder().
		Name(in.Username).
		Email(in.Email).
		Contact(contact.Title, contact.Content)
	guest := &entity.Guest{
		Name:  in.Username,
		Email: in.Email,
	}
	mail := &entity.MailConfig{
		EmailID:       entity.EmailIDUserReceivedContact,
		Substitutions: builder.Build(),
	}
	maker := entity.NewAdminURLMaker(s.adminWebURL())
	report := &entity.Report{
		ReportID:   entity.ReportIDReceivedContact,
		Overview:   contact.Title,
		Detail:     contact.Content,
		Link:       maker.Contact(contact.ID),
		ReceivedAt: contact.CreatedAt,
	}
	payload := &entity.WorkerPayload{
		QueueID:   uuid.Base58Encode(uuid.New()),
		EventType: entity.EventTypeUserReceivedContact,
		UserType:  entity.UserTypeGuest,
		Guest:     guest,
		Email:     mail,
		Report:    report,
	}
	err = s.sendMessage(ctx, payload)
	return exception.InternalError(err)
}

func (s *service) NotifyNotification(ctx context.Context, in *messenger.NotifyNotificationInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	// TODO: ???????????????
	s.logger.Debug("Notify Notification", zap.String("notificationId", in.NotificationID), zap.Time("now", s.now()))
	return nil
}

func (s *service) sendMessage(ctx context.Context, payload *entity.WorkerPayload) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	queue := entity.NewReceivedQueue(payload)
	if err := s.db.ReceivedQueue.Create(ctx, queue); err != nil {
		return err
	}
	if _, err := s.producer.SendMessage(ctx, buf); err != nil {
		return err
	}
	return nil
}
