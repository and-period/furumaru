package worker

import (
	"context"
	"sync"
	"testing"
	"time"

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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type mocks struct {
	db        *dbMocks
	mailer    *mock_mailer.MockClient
	line      *mock_line.MockClient
	messaging *mock_messaging.MockClient
	user      *mock_user.MockService
}

type dbMocks struct {
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
		WaitGroup:      &sync.WaitGroup{},
		Mailer:         mocks.mailer,
		Line:           mocks.line,
		AdminMessaging: mocks.messaging,
		UserMessaging:  mocks.messaging,
		DB: &database.Database{
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
		ctx, cancel := context.WithCancel(t.Context())
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
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
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
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, gomock.Any(), gomock.Any()).
					Return(nil, context.Canceled)
			},
			record: events.SQSMessage{
				Body: `{"queueId":"", "eventType":0, "userType":0, "userIds":[], "email":{}}`,
			},
			expectErr: context.Canceled,
		},
		{
			name: "failed to run without retry",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			record: events.SQSMessage{
				Body: `{"queueId":"", "eventType":0, "userType":0, "userIds":[], "email":{}}`,
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
				err := worker.dispatch(ctx, tt.record)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}

func TestWorker_Run(t *testing.T) {
	t.Parallel()

	usersIn := &user.MultiGetUsersInput{
		UserIDs: []string{"user-id"},
	}
	users := uentity.Users{{
		ID:         "user-id",
		Registered: true,
		Member: uentity.Member{
			Username:      "username",
			Lastname:      "username",
			Firstname:     "",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "りようしゃ",
			Email:         "test-user@and-period.jp",
		},
	}}
	devicesIn := &user.MultiGetAdminDevicesInput{
		AdminIDs: []string{"admin-id"},
	}
	devices := []string{"instance-id"}
	personalizations := []*mailer.Personalization{
		{
			Name:    "username",
			Address: "test-user@and-period.jp",
			Type:    mailer.AddressTypeTo,
			Substitutions: map[string]interface{}{
				"key": "value",
				"氏名":  "username",
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
		TemplateID:    entity.PushTemplateIDContact,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  "テンプレートです。",
		ImageURL:      "https://and-period.jp/image.png",
		CreatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
	}
	mtemplate := &entity.MessageTemplate{
		TemplateID:    entity.MessageTemplateIDNotificationLive,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  `テンプレートです。`,
		CreatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:     jst.Date(2022, 7, 14, 18, 30, 0, 0),
	}
	rtemplate := &entity.ReportTemplate{
		TemplateID: entity.ReportTemplateIDReceivedContact,
		Template:   `{"type":"bubble","body":{"type":"box","contents":[{"type":"text","text":"{{.Overview}}"}]}}`,
		CreatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
		UpdatedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
	}
	queue := func(notifyType entity.NotifyType) *entity.ReceivedQueue {
		return &entity.ReceivedQueue{
			ID:         "queue-id",
			NotifyType: notifyType,
			EventType:  entity.EventTypeUnknown,
			UserType:   entity.UserTypeAdmin,
			UserIDs:    []string{"user-id"},
			Done:       false,
			CreatedAt:  jst.Date(2022, 7, 10, 18, 30, 0, 0),
			UpdatedAt:  jst.Date(2022, 7, 10, 18, 30, 0, 0),
		}
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
				const notifyType = entity.NotifyTypeEmail
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.db.ReceivedQueue.EXPECT().
					UpdateDone(ctx, "queue-id", notifyType, true).
					Return(nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(users, nil)
				mocks.mailer.EXPECT().
					MultiSendFromInfo(gomock.Any(), "admin-register", personalizations).
					Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					TemplateID: entity.EmailTemplateIDAdminRegister,
					Substitutions: map[string]interface{}{
						"key": "value",
						"氏名":  "username",
					},
				},
			},
			expectErr: nil,
		},
		{
			name: "success to send push",
			setup: func(ctx context.Context, mocks *mocks) {
				const notifyType = entity.NotifyTypePush
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.db.ReceivedQueue.EXPECT().
					UpdateDone(ctx, "queue-id", notifyType, true).
					Return(nil)
				mocks.db.PushTemplate.EXPECT().
					Get(gomock.Any(), entity.PushTemplateIDContact).
					Return(ptemplate, nil)
				mocks.user.EXPECT().
					MultiGetAdminDevices(gomock.Any(), devicesIn).
					Return(devices, nil)
				mocks.messaging.EXPECT().
					MultiSend(gomock.Any(), message, devices).
					Return(int64(1), int64(0), nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: nil,
		},
		{
			name: "success to message",
			setup: func(ctx context.Context, mocks *mocks) {
				const notifyType = entity.NotifyTypeMessage
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.db.ReceivedQueue.EXPECT().
					UpdateDone(ctx, "queue-id", notifyType, true).
					Return(nil)
				mocks.db.MessageTemplate.EXPECT().
					Get(gomock.Any(), entity.MessageTemplateIDNotificationLive).
					Return(mtemplate, nil)
				mocks.db.Message.EXPECT().MultiCreate(gomock.Any(), gomock.Any()).Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeNotification,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Message: &entity.MessageConfig{
					TemplateID:  entity.MessageTemplateIDNotificationLive,
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
				const notifyType = entity.NotifyTypeReport
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.db.ReceivedQueue.EXPECT().
					UpdateDone(ctx, "queue-id", notifyType, true).
					Return(nil)
				mocks.db.ReportTemplate.EXPECT().
					Get(gomock.Any(), entity.ReportTemplateIDReceivedContact).
					Return(rtemplate, nil)
				mocks.line.EXPECT().PushMessage(gomock.Any(), gomock.Any()).Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				Report: &entity.ReportConfig{
					TemplateID: entity.ReportTemplateIDReceivedContact,
					Overview:   "お問い合わせ件名",
					Link:       "htts://admin.and-period.jp/contacts/contact-id",
				},
			},
			expectErr: nil,
		},
		{
			name: "success already done",
			setup: func(ctx context.Context, mocks *mocks) {
				const notifyType = entity.NotifyTypeEmail
				queue := &entity.ReceivedQueue{Done: true, NotifyType: notifyType}
				mocks.db.ReceivedQueue.EXPECT().Get(ctx, "queue-id", notifyType).Return(queue, nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					TemplateID:    entity.EmailTemplateIDAdminRegister,
					Substitutions: map[string]interface{}{"key": "value"},
				},
			},
			expectErr: nil,
		},
		{
			name:  "success empty payload",
			setup: func(ctx context.Context, mocks *mocks) {},
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
				const notifyType = entity.NotifyTypeEmail
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					TemplateID:    entity.EmailTemplateIDAdminRegister,
					Substitutions: map[string]interface{}{"key": "value"},
				},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to send mail",
			setup: func(ctx context.Context, mocks *mocks) {
				const notifyType = entity.NotifyTypeEmail
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					TemplateID:    entity.EmailTemplateIDAdminRegister,
					Substitutions: map[string]interface{}{"key": "value"},
				},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to send push",
			setup: func(ctx context.Context, mocks *mocks) {
				const notifyType = entity.NotifyTypePush
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.user.EXPECT().
					MultiGetAdminDevices(gomock.Any(), devicesIn).
					Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeReceivedContact,
				UserType:  entity.UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Push: &entity.PushConfig{
					TemplateID: entity.PushTemplateIDContact,
					Data:       map[string]string{"Title": "テストお問い合わせ"},
				},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to create message",
			setup: func(ctx context.Context, mocks *mocks) {
				const notifyType = entity.NotifyTypeMessage
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.db.MessageTemplate.EXPECT().
					Get(gomock.Any(), entity.MessageTemplateIDNotificationLive).
					Return(nil, assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeNotification,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Message: &entity.MessageConfig{
					TemplateID:  entity.MessageTemplateIDNotificationLive,
					MessageType: entity.MessageTypeNotification,
					Title:       "メッセージタイトル",
					Link:        "https://and-period.jp",
					ReceivedAt:  time.Now(),
				},
			},
			expectErr: assert.AnError,
		},
		{
			name: "failed to update received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				const notifyType = entity.NotifyTypeEmail
				mocks.db.ReceivedQueue.EXPECT().
					Get(ctx, "queue-id", notifyType).
					Return(queue(notifyType), nil)
				mocks.db.ReceivedQueue.EXPECT().
					UpdateDone(ctx, "queue-id", notifyType, true).
					Return(assert.AnError)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(users, nil)
				mocks.mailer.EXPECT().
					MultiSendFromInfo(gomock.Any(), "admin-register", personalizations).
					Return(nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email: &entity.MailConfig{
					TemplateID:    entity.EmailTemplateIDAdminRegister,
					Substitutions: map[string]interface{}{"key": "value"},
				},
			},
			expectErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			testWorker(tt.setup, func(ctx context.Context, t *testing.T, worker *worker) {
				err := worker.run(ctx, tt.payload)
				assert.ErrorIs(t, err, tt.expectErr)
			}),
		)
	}
}
