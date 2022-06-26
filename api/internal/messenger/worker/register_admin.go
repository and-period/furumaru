package worker

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/mailer"
)

func (w *worker) registerAdmin(ctx context.Context, payload *messenger.WorkerPayload) error {
	if len(payload.UserIDs) == 0 {
		return fmt.Errorf("%w: least one user ids is required", errInvalidPayloadFormat)
	}
	if payload.Email == nil {
		return fmt.Errorf("%w: email config is required", errInvalidPayloadFormat)
	}
	in := &user.GetAdminInput{
		AdminID: payload.UserIDs[0],
	}
	admin, err := w.user.GetAdmin(ctx, in)
	if err != nil {
		return err
	}
	if admin.Email == "" {
		return nil
	}
	maker := entity.NewAdminURLMaker(w.adminWebURL())
	builder := entity.NewTemplateDataBuilder().
		Data(payload.Email.Substitutions).
		Name(admin.Name()).
		WebURL(maker.SignIn())
	p := &mailer.Personalization{
		Name:          admin.Name(),
		Address:       admin.Email,
		Type:          mailer.AddressTypeTo,
		Substitutions: mailer.NewSubstitutions(builder.Build()),
	}
	return w.sendMail(ctx, payload.Email.EmailID, p)
}
