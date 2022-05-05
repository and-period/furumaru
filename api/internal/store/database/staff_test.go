package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaff(t *testing.T) {
	assert.NotNil(t, NewStaff(nil))
}

func TestStaff_ListByStoreID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, staffTable, storeTable)
	s := testStore(1, "&.農園", now())
	err = m.db.DB.Create(&s).Error
	require.NoError(t, err)
	staffs := make(entity.Staffs, 2)
	staffs[0] = testStaff(1, "user-id01", now())
	staffs[1] = testStaff(1, "user-id02", now())
	err = m.db.DB.Create(&staffs).Error
	require.NoError(t, err)

	type args struct {
		storeID int64
	}
	type want struct {
		staffs entity.Staffs
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				storeID: 1,
			},
			want: want{
				staffs: staffs,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &staff{db: m.db, now: now}
			actual, err := db.ListByStoreID(ctx, tt.args.storeID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreStaffsField(actual, now())
			assert.ElementsMatch(t, tt.want.staffs, actual)
		})
	}
}

func testStaff(storeID int64, userID string, now time.Time) *entity.Staff {
	return &entity.Staff{
		StoreID:   storeID,
		UserID:    userID,
		Role:      entity.StoreRoleAdministrator,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func fillIgnoreStaffField(s *entity.Staff, now time.Time) {
	if s == nil {
		return
	}
	s.CreatedAt = now
	s.UpdatedAt = now
}

func fillIgnoreStaffsField(ss entity.Staffs, now time.Time) {
	for i := range ss {
		fillIgnoreStaffField(ss[i], now)
	}
}
