package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (h *handler) multiGetUsers(ctx context.Context, userIDs []string) (entity.Users, error) {
	if len(userIDs) == 0 {
		return entity.Users{}, nil
	}
	in := &user.MultiGetUsersInput{
		UserIDs: userIDs,
	}
	return h.user.MultiGetUsers(ctx, in)
}
