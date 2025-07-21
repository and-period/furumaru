package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/stretchr/testify/assert"
)

func TestSearchPostalCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.SearchPostalCodeInput
		expect    *entity.PostalCode
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				code := &postalcode.PostalCode{
					PostalCode:     "1000014",
					PrefectureCode: "13",
					Prefecture:     "東京都",
					City:           "千代田区",
					Town:           "永田町",
					PrefectureKana: "ﾄｳｷｮｳﾄ",
					CityKana:       "ﾁﾖﾀﾞｸ",
					TownKana:       "ﾅｶﾞﾀﾁｮｳ",
				}
				mocks.postalCode.EXPECT().Search(ctx, "1000014").Return(code, nil)
			},
			input: &store.SearchPostalCodeInput{
				PostlCode: "1000014",
			},
			expect: &entity.PostalCode{
				PostalCode:     "1000014",
				PrefectureCode: 13,
				Prefecture:     "東京都",
				City:           "千代田区",
				Town:           "永田町",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.SearchPostalCodeInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to search postal code",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.postalCode.EXPECT().Search(ctx, "1000014").Return(nil, assert.AnError)
			},
			input: &store.SearchPostalCodeInput{
				PostlCode: "1000014",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to invalid postal code",
			setup: func(ctx context.Context, mocks *mocks) {
				code := &postalcode.PostalCode{
					PostalCode:     "1000014",
					PrefectureCode: "tokyo",
					Prefecture:     "東京都",
					City:           "千代田区",
					Town:           "永田町",
					PrefectureKana: "ﾄｳｷｮｳﾄ",
					CityKana:       "ﾁﾖﾀﾞｸ",
					TownKana:       "ﾅｶﾞﾀﾁｮｳ",
				}
				mocks.postalCode.EXPECT().Search(ctx, "1000014").Return(code, nil)
			},
			input: &store.SearchPostalCodeInput{
				PostlCode: "1000014",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
				actual, err := service.SearchPostalCode(ctx, tt.input)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expect, actual)
			}),
		)
	}
}
