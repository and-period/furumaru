package service

import (
	"context"
	"encoding/json"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) NotifyRegisterAdmin(ctx context.Context, in *messenger.NotifyRegisterAdminInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	builder := entity.NewTemplateDataBuilder().Password(in.Password)
	msg := &entity.MailConfig{
		EmailID:       entity.EmailIDRegisterAdmin,
		Substitutions: builder.Build(),
	}
	payload := &messenger.WorkerPayload{
		EventType: messenger.EventTypeRegisterAdmin,
		UserType:  messenger.UserTypeAdmin,
		UserIDs:   []string{in.AdminID},
		Email:     msg,
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return exception.InternalError(err)
	}
	// TODO: messageIdの重複排除の実装
	_, err = s.producer.SendMessage(ctx, buf)
	return exception.InternalError(err)
}
