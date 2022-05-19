package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestListStaffsByStoreID(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	staffs := entity.Staffs{
		{
			StoreID:   1,
			UserID:    "user-id01",
			Role:      entity.StoreRoleAdministrator,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			StoreID:   1,
			UserID:    "user-id02",
			Role:      entity.StoreRoleEditor,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListStaffsByStoreIDInput
		expect    entity.Staffs
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Staff.EXPECT().ListByStoreID(ctx, int64(1)).Return(staffs, nil)
			},
			input: &store.ListStaffsByStoreIDInput{
				StoreID: 1,
			},
			expect:    staffs,
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ListStaffsByStoreIDInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get stores",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Staff.EXPECT().ListByStoreID(ctx, int64(1)).Return(nil, errmock)
			},
			input: &store.ListStaffsByStoreIDInput{
				StoreID: 1,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			actual, err := service.ListStaffsByStoreID(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}
