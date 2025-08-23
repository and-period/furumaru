package service

import (
	"context"
	"errors"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) UpsertGuest(ctx context.Context, in *user.UpsertGuestInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	guest, err := s.db.Guest.GetByEmail(ctx, in.Email)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return "", internalError(err)
	}
	if guest == nil { // 登録処理
		params := &entity.NewUserParams{
			UserType:      entity.UserTypeGuest,
			Registered:    false, // ゲストとして登録
			Lastname:      in.Lastname,
			Firstname:     in.Firstname,
			LastnameKana:  in.LastnameKana,
			FirstnameKana: in.FirstnameKana,
			Email:         in.Email,
		}
		user := entity.NewUser(params)
		err = s.db.Guest.Create(ctx, user)
		guest = &user.Guest
	} else { // 更新処理
		params := &database.UpdateGuestParams{
			Lastname:      in.Lastname,
			Firstname:     in.Firstname,
			LastnameKana:  in.LastnameKana,
			FirstnameKana: in.FirstnameKana,
		}
		err = s.db.Guest.Update(ctx, guest.UserID, params)
	}
	if err != nil {
		return "", internalError(err)
	}
	return guest.UserID, nil
}
