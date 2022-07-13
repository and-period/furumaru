package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	mock_database "github.com/and-period/furumaru/api/mock/messenger/database"
	mock_mailer "github.com/and-period/furumaru/api/mock/pkg/mailer"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var errmock = errors.New("some error")

type mocks struct {
	db     *dbMocks
	mailer *mock_mailer.MockClient
	user   *mock_user.MockService
}

type dbMocks struct {
	Contact        *mock_database.MockContact
	Notification   *mock_database.MockNotification
	ReceivedQueue  *mock_database.MockReceivedQueue
	ReportTemplate *mock_database.MockReportTemplate
}

type testOptions struct {
	now func() time.Time
}

type testOption func(opts *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

type testCaller func(ctx context.Context, t *testing.T, worker *worker)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:     newDBMocks(ctrl),
		mailer: mock_mailer.NewMockClient(ctrl),
		user:   mock_user.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Contact:        mock_database.NewMockContact(ctrl),
		Notification:   mock_database.NewMockNotification(ctrl),
		ReceivedQueue:  mock_database.NewMockReceivedQueue(ctrl),
		ReportTemplate: mock_database.NewMockReportTemplate(ctrl),
	}
}

func newWorker(mocks *mocks, opts ...testOption) *worker {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &worker{
		now:       dopts.now,
		logger:    zap.NewNop(),
		waitGroup: &sync.WaitGroup{},
		mailer:    mocks.mailer,
		db: &database.Database{
			Contact:        mocks.db.Contact,
			Notification:   mocks.db.Notification,
			ReceivedQueue:  mocks.db.ReceivedQueue,
			ReportTemplate: mocks.db.ReportTemplate,
		},
		user:        mocks.user,
		concurrency: 1,
		maxRetries:  1,
	}
}

func testWorker(
	setup func(ctx context.Context, mocks *mocks),
	testFunc testCaller,
	opts ...testOption,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := newMocks(ctrl)

		w := newWorker(mocks)
		setup(ctx, mocks)

		testFunc(ctx, t, w)
		w.waitGroup.Wait()
	}
}

func TestWorker(t *testing.T) {
	t.Parallel()
	w := NewWorker(&Params{}, WithLogger(zap.NewNop()), WithConcurrency(1), WithMaxRetries(3))
	assert.NotNil(t, w)
}

func TestWorker_Dispatch(t *testing.T) {
	t.Parallel()

	in := &user.MultiGetUsersInput{
		UserIDs: []string{"user-id"},
	}
	users := uentity.Users{{
		Username: "&. スタッフ",
		Email:    "test-user@and-period.jp",
	}}
	personalizations := []*mailer.Personalization{
		{
			Name:    "&. スタッフ",
			Address: "test-user@and-period.jp",
			Type:    mailer.AddressTypeTo,
			Substitutions: map[string]interface{}{
				"key": "value",
				"氏名":  "&. スタッフ",
			},
		},
	}
	queue := &entity.ReceivedQueue{
		ID:        "queue-id",
		EventType: entity.EventTypeUnknown,
		UserType:  entity.UserTypeAdmin,
		UserIDs:   []string{"user-id"},
		Done:      false,
		CreatedAt: jst.Date(2022, 7, 10, 18, 30, 0, 0),
		UpdatedAt: jst.Date(2022, 7, 10, 18, 30, 0, 0),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		payload   *entity.WorkerPayload
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), in).Return(users, nil)
				mocks.mailer.EXPECT().MultiSendFromInfo(gomock.Any(), entity.EmailIDAdminRegister, personalizations).Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: nil,
		},
		{
			name: "success already done",
			setup: func(ctx context.Context, mocks *mocks) {
				queue := &entity.ReceivedQueue{Done: true}
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: nil,
		},
		{
			name: "success empty payload",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email:     nil,
			},
			expectErr: nil,
		},
		{
			name: "failed to get received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(nil, errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: errmock,
		},
		{
			name: "failed to multi send mail",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), in).Return(nil, errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: errmock,
		},
		{
			name: "failed to update received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(errmock)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), in).Return(users, nil)
				mocks.mailer.EXPECT().MultiSendFromInfo(gomock.Any(), entity.EmailIDAdminRegister, personalizations).Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					EmailID:       entity.EmailIDAdminRegister,
					Substitutions: map[string]string{"key": "value"},
				},
			},
			expectErr: errmock,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.dispatch(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
