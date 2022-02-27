package response

import "github.com/and-period/marche/api/internal/gateway/user/v1/entity"

type AuthResponse struct {
	*entity.Auth
}
