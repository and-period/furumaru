package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageTemplate(t *testing.T) {
	assert.NotNil(t, NewMessageTemplate(nil))
}

func TestMessageTemplate_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	tmpl := testMessageTemplate("message-id", now())
	err = db.DB.Create(&tmpl).Error
	require.NoError(t, err)

	type args struct {
		messageID entity.MessageTemplateID
	}
	type want struct {
		template *entity.MessageTemplate
		err      error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				messageID: "message-id",
			},
			want: want{
				template: tmpl,
				err:      nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				messageID: "other-id",
			},
			want: want{
				template: nil,
				err:      database.ErrNotFound,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &messageTemplate{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.messageID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.template, actual)
		})
	}
}

func testMessageTemplate(id entity.MessageTemplateID, now time.Time) *entity.MessageTemplate {
	return &entity.MessageTemplate{
		TemplateID:    id,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  "内容: {{.Body}}",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
