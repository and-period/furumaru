package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/mailer"
)

// FIXME: workerへ処理を寄せる
func (s *service) NotifyRegisterAdmin(ctx context.Context, in *messenger.NotifyRegisterAdminInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}

	maker := entity.NewAdminURLMaker(s.adminWebURL())
	builder := entity.NewTemplateDataBuilder().
		Name(in.Name).
		Password(in.Password).
		WebURL(maker.SignIn())

	msg := &entity.MailConfig{
		EmailID:       entity.EmailIDRegisterAdmin,
		Substitutions: builder.Build(),
	}

	personalization := &mailer.Personalization{
		Name:          in.Name,
		Address:       in.Email,
		Type:          mailer.AddressTypeTo,
		Substitutions: mailer.NewSubstitutions(msg.Substitutions),
	}
	err := s.sendInfoMail(ctx, msg, personalization)
	return exception.InternalError(err)
}
