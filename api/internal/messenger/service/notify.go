package service

import (
	"context"
	"encoding/json"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"golang.org/x/sync/errgroup"
)

// NotifyOrderAuthorized - 支払い完了
func (s *service) NotifyOrderAuthorized(ctx context.Context, in *messenger.NotifyOrderAuthorizedInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	orderIn := &store.GetOrderInput{
		OrderID: in.OrderID,
	}
	order, err := s.store.GetOrder(ctx, orderIn)
	if err != nil {
		return internalError(err)
	}
	builder := entity.NewTemplateDataBuilder().Order(order)
	mail := &entity.MailConfig{
		EmailID:       entity.EmailIDUserOrderAuthorized,
		Substitutions: builder.Build(),
	}
	payload := &entity.WorkerPayload{
		QueueID:   uuid.Base58Encode(uuid.New()),
		EventType: entity.EventTypeOrderAuthorized,
		UserType:  entity.UserTypeUser,
		UserIDs:   []string{order.UserID},
		Email:     mail,
	}
	err = s.sendMessage(ctx, payload)
	return internalError(err)
}

// NotifyRegisterAdmin - 管理者登録
func (s *service) NotifyRegisterAdmin(ctx context.Context, in *messenger.NotifyRegisterAdminInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
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
		EventType: entity.EventTypeRegisterAdmin,
		UserType:  entity.UserTypeAdmin,
		UserIDs:   []string{in.AdminID},
		Email:     mail,
	}
	err := s.sendMessage(ctx, payload)
	return internalError(err)
}

// NotifyResetAdminPassword - 管理者パスワードリセット
func (s *service) NotifyResetAdminPassword(ctx context.Context, in *messenger.NotifyResetAdminPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	maker := entity.NewAdminURLMaker(s.adminWebURL())
	builder := entity.NewTemplateDataBuilder().
		WebURL(maker.SignIn()).
		Password(in.Password)
	mail := &entity.MailConfig{
		EmailID:       entity.EmailIDAdminResetPassword,
		Substitutions: builder.Build(),
	}
	payload := &entity.WorkerPayload{
		QueueID:   uuid.Base58Encode(uuid.New()),
		EventType: entity.EventTypeResetAdminPassword,
		UserType:  entity.UserTypeAdmin,
		UserIDs:   []string{in.AdminID},
		Email:     mail,
	}
	err := s.sendMessage(ctx, payload)
	return internalError(err)
}

// NotifyNotification - お知らせ発行
func (s *service) NotifyNotification(ctx context.Context, in *messenger.NotifyNotificationInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	notification, err := s.db.Notification.Get(ctx, in.NotificationID)
	if err != nil {
		return internalError(err)
	}
	if notification.Type == entity.NotificationTypePromotion {
		in := &store.GetPromotionInput{
			PromotionID: notification.PromotionID,
		}
		promotion, err := s.store.GetPromotion(ctx, in)
		if err != nil {
			return internalError(err)
		}
		notification.Title = promotion.Title
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
	if err := eg.Wait(); err != nil {
		return internalError(err)
	}
	maker := entity.NewAdminURLMaker(s.adminWebURL())
	report := &entity.ReportConfig{
		ReportID:    entity.ReportIDNotification,
		Overview:    notification.Title,
		Detail:      notification.Body,
		Link:        maker.Notification(notification.ID),
		PublishedAt: notification.PublishedAt,
	}
	payload := &entity.WorkerPayload{
		QueueID:   uuid.Base58Encode(uuid.New()),
		EventType: entity.EventTypeNotification,
		Report:    report,
	}
	return s.sendMessage(ctx, payload)
}

func (s *service) notifyUserNotification(ctx context.Context, notification *entity.Notification) error {
	if !notification.HasUserTarget() {
		return nil
	}
	message := &entity.MessageConfig{
		MessageID:   entity.MessageIDNotification,
		MessageType: entity.MessageTypeNotification,
		Title:       notification.Title,
		ReceivedAt:  s.now(),
	}
	payload := &entity.WorkerPayload{
		EventType: entity.EventTypeNotification,
		Message:   message,
	}
	return s.sendAllUsers(ctx, payload)
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
		EventType: entity.EventTypeNotification,
		Message:   message,
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if !notification.HasAdministratorTarget() {
			return
		}
		return s.sendAllAdministrators(ectx, payload)
	})
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

/*
 * private
 */
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

func (s *service) sendAll(
	ctx context.Context,
	payload *entity.WorkerPayload,
	userType entity.UserType,
	listFn func(limit, offset int64) ([]string, int64, error),
) error {
	const unit = 200

	var next int64
	for {
		userIDs, total, err := listFn(unit, next)
		if err != nil {
			return err
		}
		if len(userIDs) == 0 {
			return nil
		}

		payload := *payload // copy
		payload.QueueID = uuid.Base58Encode(uuid.New())
		payload.UserType = userType
		payload.UserIDs = userIDs
		if err := s.sendMessage(ctx, &payload); err != nil {
			return err
		}

		next += int64(len(userIDs))
		if next >= total {
			return nil
		}
	}
}

func (s *service) sendAllUsers(ctx context.Context, payload *entity.WorkerPayload) error {
	listFn := func(limit, offset int64) ([]string, int64, error) {
		in := &user.ListUsersInput{
			Limit:          limit,
			Offset:         offset,
			OnlyRegistered: true,
		}
		users, total, err := s.user.ListUsers(ctx, in)
		if err != nil || len(users) == 0 {
			return nil, 0, err
		}
		return users.IDs(), total, nil
	}
	return s.sendAll(ctx, payload, entity.UserTypeUser, listFn)
}

func (s *service) sendAllAdministrators(ctx context.Context, payload *entity.WorkerPayload) error {
	listFn := func(limit, offset int64) ([]string, int64, error) {
		in := &user.ListAdministratorsInput{
			Limit:  limit,
			Offset: offset,
		}
		administrators, total, err := s.user.ListAdministrators(ctx, in)
		if err != nil || len(administrators) == 0 {
			return nil, 0, err
		}
		return administrators.IDs(), total, nil
	}
	return s.sendAll(ctx, payload, entity.UserTypeAdministrator, listFn)
}

func (s *service) sendAllCoordinators(ctx context.Context, payload *entity.WorkerPayload) error {
	listFn := func(limit, offset int64) ([]string, int64, error) {
		in := &user.ListCoordinatorsInput{
			Limit:  limit,
			Offset: offset,
		}
		coordinators, total, err := s.user.ListCoordinators(ctx, in)
		if err != nil || len(coordinators) == 0 {
			return nil, 0, err
		}
		return coordinators.IDs(), total, nil
	}
	return s.sendAll(ctx, payload, entity.UserTypeCoordinator, listFn)
}

func (s *service) sendAllProducers(ctx context.Context, payload *entity.WorkerPayload) error {
	listFn := func(limit, offset int64) ([]string, int64, error) {
		in := &user.ListProducersInput{
			Limit:  limit,
			Offset: offset,
		}
		producers, total, err := s.user.ListProducers(ctx, in)
		if err != nil || len(producers) == 0 {
			return nil, 0, err
		}
		return producers.IDs(), total, nil
	}
	return s.sendAll(ctx, payload, entity.UserTypeProducer, listFn)
}
