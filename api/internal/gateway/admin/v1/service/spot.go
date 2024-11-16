package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
)

// SpotUserType - 投稿者の種別
type SpotUserType int32

const (
	SpotUserTypeUnknown SpotUserType = 0
	SpotUserTypeUser    SpotUserType = 1 // ユーザー
	SpotUserTypeAdmin   SpotUserType = 2 // 管理者
)

func NewSpotUserTypeFromInt32(userType int32) SpotUserType {
	return SpotUserType(userType)
}

type Spot struct {
	response.Spot
}
