package service

import (
	"context"
	"encoding/json"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"golang.org/x/sync/errgroup"
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
	report := &entity.ReportConfig{
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
	notification, err := s.db.Notification.Get(ctx, in.NotificationID)
	if err != nil {
		return exception.InternalError(err)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if !notification.HasUserTarget() {
			return
		}
		return s.notifyUserNotification(ectx, notification)
	})
	eg.Go(func() (err error) {
		if !notification.HasAdminTarget() {
			return
		}
		return s.notifyAdminNotification(ectx, notification)
	})
	eg.Go(func() (err error) {
		// TODO: システムレポート周りの実装
		return
	})
	if err := eg.Wait(); err != nil {
		return exception.InternalError(err)
	}
	return nil
}

func (s *service) notifyUserNotification(ctx context.Context, notification *entity.Notification) error {
	// TODO: 後ほどユーザー側への通知も実装する
	return nil
}

func (s *service) notifyAdminNotification(ctx context.Context, notification *entity.Notification) error {
	maker := entity.NewAdminURLMaker(s.adminWebURL())
	message := &entity.MessageConfig{
		MessageID:   entity.MessageIDNotification,
		MessageType: entity.MessageTypeNotification,
		Title:       notification.Title,
		Link:        maker.Notification(notification.ID),
		ReceivedAt:  s.now(),
	}
	payload := &entity.WorkerPayload{
		EventType: entity.EventTypeAdminNotification,
		Message:   message,
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if !notification.HasCoordinatorTarget() {
			return
		}
		return s.sendAllCoordinators(ectx, payload)
	})
	eg.Go(func() (err error) {
		if !notification.HasProducerTarget() {
			return
		}
		return s.sendAllProducers(ectx, payload)
	})
	return eg.Wait()
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

func (s *service) sendAllCoordinators(ctx context.Context, payload *entity.WorkerPayload) error {
	const unit = 200

	var next int64
	for {
		in := &user.ListCoordinatorsInput{
			Limit:  unit,
			Offset: next,
		}
		coordinators, total, err := s.user.ListCoordinators(ctx, in)
		if err != nil {
			return err
		}
		if len(coordinators) == 0 {
			return nil
		}

		payload := *payload // copy
		payload.QueueID = uuid.Base58Encode(uuid.New())
		payload.UserType = entity.UserTypeCoordinator
		payload.UserIDs = coordinators.IDs()
		if err := s.sendMessage(ctx, &payload); err != nil {
			return err
		}

		next += int64(len(coordinators))
		if next >= total {
			return nil
		}
	}
}

func (s *service) sendAllProducers(ctx context.Context, payload *entity.WorkerPayload) error {
	const unit = 200

	var next int64
	for {
		in := &user.ListProducersInput{
			Limit:  unit,
			Offset: next,
		}
		producers, total, err := s.user.ListProducers(ctx, in)
		if err != nil {
			return err
		}
		if len(producers) == 0 {
			return nil
		}

		payload := *payload // copy
		payload.QueueID = uuid.Base58Encode(uuid.New())
		payload.UserType = entity.UserTypeProducer
		payload.UserIDs = producers.IDs()
		if err := s.sendMessage(ctx, &payload); err != nil {
			return err
		}

		next += int64(len(producers))
		if next >= total {
			return nil
		}
	}
}
