package worker

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"go.uber.org/zap"
)

func (w *worker) multiSendMail(ctx context.Context, payload *entity.WorkerPayload) error {
	ps, err := w.newPersonalizations(ctx, payload)
	if err != nil {
		return err
	}
	return w.sendMail(ctx, payload.Email.TemplateID, ps...)
}

func (w *worker) sendMail(ctx context.Context, templateID entity.EmailTemplateID, ps ...*mailer.Personalization) error {
	if len(ps) == 0 {
		w.logger.Debug("Personalizations is empty", zap.String("templateId", string(templateID)))
		return nil
	}
	w.logger.Debug("Send email", zap.String("templateId", string(templateID)), zap.Any("personalizations", ps))
	sendFn := func() error {
		return w.mailer.MultiSendFromInfo(ctx, string(templateID), ps)
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(w.isRetryable))
}

func (w *worker) newPersonalizations(
	ctx context.Context, payload *entity.WorkerPayload,
) ([]*mailer.Personalization, error) {
	ps := make([]*mailer.Personalization, 0, len(payload.UserIDs))
	execute := func(name, email string) {
		if email == "" {
			return
		}
		builder := entity.NewTemplateDataBuilder().
			Data(payload.Email.Substitutions).
			Name(name)
		p := &mailer.Personalization{
			Name:          name,
			Address:       email,
			Type:          mailer.AddressTypeTo,
			Substitutions: mailer.NewSubstitutions(builder.Build()),
		}
		ps = append(ps, p)
	}
	var err error
	switch payload.UserType {
	case entity.UserTypeAdmin:
		err = w.fetchAdmins(ctx, payload.UserIDs, execute)
	case entity.UserTypeAdministrator:
		err = w.fetchAdministrators(ctx, payload.UserIDs, execute)
	case entity.UserTypeCoordinator:
		err = w.fetchCoordinators(ctx, payload.UserIDs, execute)
	case entity.UserTypeProducer:
		err = w.fetchProducers(ctx, payload.UserIDs, execute)
	case entity.UserTypeUser:
		err = w.fetchUsers(ctx, payload.UserIDs, execute)
	case entity.UserTypeGuest:
		err = w.fetchUsers(ctx, payload.UserIDs, execute)
	default:
		err = fmt.Errorf("worker: failed to multi send mail: %w", errUnknownUserType)
	}
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (w *worker) fetchAdmins(ctx context.Context, adminIDs []string, execute func(name, email string)) error {
	in := &user.MultiGetAdminsInput{
		AdminIDs: adminIDs,
	}
	admins, err := w.user.MultiGetAdmins(ctx, in)
	if err != nil {
		return err
	}
	for i := range admins {
		execute(admins[i].Name(), admins[i].Email)
	}
	return nil
}

func (w *worker) fetchAdministrators(
	ctx context.Context, administratorIDs []string, execute func(name, email string),
) error {
	in := &user.MultiGetAdministratorsInput{
		AdministratorIDs: administratorIDs,
	}
	administrators, err := w.user.MultiGetAdministrators(ctx, in)
	if err != nil {
		return err
	}
	for i := range administrators {
		execute(administrators[i].Name(), administrators[i].Email)
	}
	return nil
}

func (w *worker) fetchCoordinators(
	ctx context.Context, coordinatorIDs []string, execute func(name, email string),
) error {
	in := &user.MultiGetCoordinatorsInput{
		CoordinatorIDs: coordinatorIDs,
	}
	coordinators, err := w.user.MultiGetCoordinators(ctx, in)
	if err != nil {
		return err
	}
	for i := range coordinators {
		execute(coordinators[i].Username, coordinators[i].Email)
	}
	return nil
}

func (w *worker) fetchProducers(ctx context.Context, producerIDs []string, execute func(name, email string)) error {
	in := &user.MultiGetProducersInput{
		ProducerIDs: producerIDs,
	}
	producers, err := w.user.MultiGetProducers(ctx, in)
	if err != nil {
		return err
	}
	for i := range producers {
		execute(producers[i].Username, producers[i].Email)
	}
	return nil
}

func (w *worker) fetchUsers(ctx context.Context, userIDs []string, execute func(name, email string)) error {
	in := &user.MultiGetUsersInput{
		UserIDs: userIDs,
	}
	users, err := w.user.MultiGetUsers(ctx, in)
	if err != nil {
		return err
	}
	for i := range users {
		execute(users[i].Name(), users[i].Email())
	}
	return nil
}
