package worker

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/mailer"
)

func (w *worker) sendInfoMail(ctx context.Context, payload *messenger.WorkerPayload) error {
	var (
		ps  []*mailer.Personalization
		err error
	)
	switch payload.UserType {
	case messenger.UserTypeUser:
		ps, err = w.getPersonalizationsByUser(ctx, payload.UserIDs, payload.Email.Substitutions)
	case messenger.UserTypeAdministrator:
		ps, err = w.getPersonalizationsByAdministrator(ctx, payload.UserIDs, payload.Email.Substitutions)
	case messenger.UserTypeCoordinator:
		ps, err = w.getPersonalizationsByCoordinator(ctx, payload.UserIDs, payload.Email.Substitutions)
	case messenger.UserTypeProducer:
		ps, err = w.getPersonalizationsByProducer(ctx, payload.UserIDs, payload.Email.Substitutions)
	default:
		return fmt.Errorf("worker: unknown user type: %w", exception.ErrInvalidArgument)
	}
	if err != nil {
		return err
	}
	if len(ps) == 0 {
		return nil
	}
	sendFn := func() error {
		return w.mailer.MultiSendFromInfo(ctx, payload.Email.EmailID, ps)
	}
	retry := backoff.NewExponentialBackoff(w.maxRetries)
	return backoff.Retry(ctx, retry, sendFn, backoff.WithRetryablel(w.retryable))
}

func (w *worker) getPersonalizationsByUser(
	ctx context.Context, userIDs []string, msg map[string]string,
) ([]*mailer.Personalization, error) {
	in := &user.MultiGetUsersInput{
		UserIDs: userIDs,
	}
	users, err := w.user.MultiGetUsers(ctx, in)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}

	ps := make([]*mailer.Personalization, 0, len(users))
	for _, u := range users {
		if u.Email == "" {
			continue
		}
		builder := entity.NewTemplateDataBuilder().Data(msg).Name(u.Username)
		p := &mailer.Personalization{
			Name:          u.Username,
			Address:       u.Email,
			Type:          mailer.AddressTypeTo,
			Substitutions: mailer.NewSubstitutions(builder.Build()),
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (w *worker) getPersonalizationsByAdministrator(
	ctx context.Context, administratorIDs []string, msg map[string]string,
) ([]*mailer.Personalization, error) {
	return nil, exception.ErrNotImplemented
}

func (w *worker) getPersonalizationsByCoordinator(
	ctx context.Context, coordinatorIDs []string, msg map[string]string,
) ([]*mailer.Personalization, error) {
	return nil, exception.ErrNotImplemented
}

func (w *worker) getPersonalizationsByProducer(
	ctx context.Context, producerIDs []string, msg map[string]string,
) ([]*mailer.Personalization, error) {
	return nil, exception.ErrNotImplemented
}
