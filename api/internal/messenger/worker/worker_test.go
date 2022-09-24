package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	mock_database "github.com/and-period/furumaru/api/mock/messenger/database"
	mock_messaging "github.com/and-period/furumaru/api/mock/pkg/firebase/messaging"
	mock_line "github.com/and-period/furumaru/api/mock/pkg/line"
	mock_mailer "github.com/and-period/furumaru/api/mock/pkg/mailer"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var errmock = errors.New("some error")

type mocks struct {
	db        *dbMocks
	mailer    *mock_mailer.MockClient
	line      *mock_line.MockClient
	messaging *mock_messaging.MockClient
	user      *mock_user.MockService
}

type dbMocks struct {
	Contact         *mock_database.MockContact
	Message         *mock_database.MockMessage
	MessageTemplate *mock_database.MockMessageTemplate
	Notification    *mock_database.MockNotification
	PushTemplate    *mock_database.MockPushTemplate
	ReceivedQueue   *mock_database.MockReceivedQueue
	ReportTemplate  *mock_database.MockReportTemplate
	Schedule        *mock_database.MockSchedule
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
		db:        newDBMocks(ctrl),
		mailer:    mock_mailer.NewMockClient(ctrl),
		line:      mock_line.NewMockClient(ctrl),
		messaging: mock_messaging.NewMockClient(ctrl),
		user:      mock_user.NewMockService(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Contact:         mock_database.NewMockContact(ctrl),
		Message:         mock_database.NewMockMessage(ctrl),
		MessageTemplate: mock_database.NewMockMessageTemplate(ctrl),
		Notification:    mock_database.NewMockNotification(ctrl),
		PushTemplate:    mock_database.NewMockPushTemplate(ctrl),
		ReceivedQueue:   mock_database.NewMockReceivedQueue(ctrl),
		ReportTemplate:  mock_database.NewMockReportTemplate(ctrl),
		Schedule:        mock_database.NewMockSchedule(ctrl),
	}
}

func newWorker(mocks *mocks, opts ...testOption) *worker {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Mailer:    mocks.mailer,
		Line:      mocks.line,
		Messaging: mocks.messaging,
		DB: &database.Database{
			Contact:         mocks.db.Contact,
			Message:         mocks.db.Message,
			MessageTemplate: mocks.db.MessageTemplate,
			Notification:    mocks.db.Notification,
			PushTemplate:    mocks.db.PushTemplate,
			ReceivedQueue:   mocks.db.ReceivedQueue,
			ReportTemplate:  mocks.db.ReportTemplate,
			Schedule:        mocks.db.Schedule,
		},
		User: mocks.user,
	}
	worker := NewWorker(params).(*worker)
	worker.concurrency = 1
	worker.maxRetries = 1
	worker.now = func() time.Time {
		return dopts.now()
	}
	return worker
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

		w := newWorker(mocks, opts...)
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

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		record    events.SQSMessage
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				queue := &entity.ReceivedQueue{Done: true}
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, gomock.Any()).Return(queue, nil)
			},
			record: events.SQSMessage{
				Body: `{"queueId":"", "eventType":0, "userType":0, "userIds":[]}`,
			},
			expectErr: nil,
		},
		{
			name:      "failed to unmarshall sqs event",
			setup:     func(ctx context.Context, mocks *mocks) {},
			record:    events.SQSMessage{},
			expectErr: nil,
		},
		{
			name: "failed to run with retry",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, gomock.Any()).Return(nil, exception.ErrUnavailable)
			},
			record: events.SQSMessage{
				Body: `{"queueId":"", "eventType":0, "userType":0, "userIds":[]}`,
			},
			expectErr: exception.ErrUnavailable,
		},
		{
			name: "failed to run without retry",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, gomock.Any()).Return(nil, errmock)
			},
			record: events.SQSMessage{
				Body: `{"queueId":"", "eventType":0, "userType":0, "userIds":[]}`,
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.dispatch(ctx, tt.record)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestWorker_Run(t *testing.T) {
	t.Parallel()

	usersIn := &user.MultiGetUsersInput{
		UserIDs: []string{"user-id"},
	}
	users := uentity.Users{{
		Member: uentity.Member{
			Email: "test-user@and-period.jp",
		},
		Customer: uentity.Customer{
			Lastname:  "&.",
			Firstname: "スタッフ",
		},
	}}
	devicesIn := &user.MultiGetAdminDevicesInput{
		AdminIDs: []string{"admin-id"},
	}
	devices := []string{"instance-id"}
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
	message := &messaging.Message{
		Title:    "件名: テストお問い合わせ",
		Body:     "テンプレートです。",
		ImageURL: "https://and-period.jp/image.png",
		Data:     map[string]string{"Title": "テストお問い合わせ"},
	}
	ptemplate := &entity.PushTemplate{
		TemplateID:    entity.PushIDContact,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  "テンプレートです。",
		ImageURL:      "https://and-period.jp/image.png",
		CreatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
	}
	mtemplate := &entity.MessageTemplate{
		TemplateID:    entity.MessageIDNotification,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  `テンプレートです。`,
		CreatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
	}
	rtemplate := &entity.ReportTemplate{
		TemplateID: entity.ReportIDReceivedContact,
		Template:   `{"type":"bubble","body":{"type":"box","contents":[{"type":"text","text":"{{.Overview}}"}]}}`,
		CreatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
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
			name: "success to send mail",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(users, nil)
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
			name: "success to send push",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(nil)
				mocks.db.PushTemplate.EXPECT().Get(gomock.Any(), entity.PushIDContact).Return(ptemplate, nil)
				mocks.user.EXPECT().MultiGetAdminDevices(gomock.Any(), devicesIn).Return(devices, nil)
				mocks.messaging.EXPECT().MultiSend(gomock.Any(), message, devices).Return(int64(1), int64(0), nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					PushID: entity.PushIDContact,
					Data:   map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: nil,
		},
		{
			name: "success to message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(nil)
				mocks.db.MessageTemplate.EXPECT().Get(gomock.Any(), entity.MessageIDNotification).Return(mtemplate, nil)
				mocks.db.Message.EXPECT().MultiCreate(gomock.Any(), gomock.Any()).Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeNotification,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Message: &entity.MessageConfig{
					MessageID:   entity.MessageIDNotification,
					MessageType: entity.MessageTypeNotification,
					Title:       "メッセージタイトル",
					Link:        "https://and-period.jp",
					ReceivedAt:  time.Now(),
				},
			},
			expectErr: nil,
		},
		{
			name: "success to report",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(nil)
				mocks.db.ReportTemplate.EXPECT().Get(gomock.Any(), entity.ReportIDReceivedContact).Return(rtemplate, nil)
				mocks.line.EXPECT().PushMessage(gomock.Any(), gomock.Any()).Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				Report: &entity.ReportConfig{
					ReportID: entity.ReportIDReceivedContact,
					Overview: "お問い合わせ件名",
					Link:     "htts://admin.and-period.jp/contacts/contact-id",
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
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to send mail",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(nil, errmock)
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
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to send push",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.user.EXPECT().MultiGetAdminDevices(gomock.Any(), devicesIn).Return(nil, errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					PushID: entity.PushIDContact,
					Data:   map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to create message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.MessageTemplate.EXPECT().Get(gomock.Any(), entity.MessageIDNotification).Return(nil, errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeNotification,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Message: &entity.MessageConfig{
					MessageID:   entity.MessageIDNotification,
					MessageType: entity.MessageTypeNotification,
					Title:       "メッセージタイトル",
					Link:        "https://and-period.jp",
					ReceivedAt:  time.Now(),
				},
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to report",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReportTemplate.EXPECT().Get(gomock.Any(), entity.ReportIDReceivedContact).Return(nil, errmock)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				Report: &entity.ReportConfig{
					ReportID: entity.ReportIDReceivedContact,
					Overview: "お問い合わせ件名",
					Link:     "htts://admin.and-period.jp/contacts/contact-id",
				},
			},
			expectErr: exception.ErrUnknown,
		},
		{
			name: "failed to update received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id").Return(queue, nil)
				mocks.db.ReceivedQueue.EXPECT().UpdateDone(ctx, "queue-id", true).Return(errmock)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(users, nil)
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
			expectErr: exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
			err := worker.run(ctx, tt.payload)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
