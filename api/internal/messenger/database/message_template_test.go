package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageTemplate(t *testing.T) {
	assert.NotNil(t, NewMessageTemplate(nil))
}

func TestMessageTemplate_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, messageTemplateTable)
	tmpl := testMessageTemplate("message-id", now())
	err = m.db.DB.Create(&tmpl).Error
	require.NoError(t, err)

	type args struct {
		messageID string
	}
	type want struct {
		template *entity.MessageTemplate
		hasErr   bool
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
				messageID: "message-id",
			},
			want: want{
				template: tmpl,
				hasErr:   false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				messageID: "other-id",
			},
			want: want{
				template: nil,
				hasErr:   true,
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

			db := &messageTemplate{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.messageID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreMessageTemplateField(actual, now())
			assert.Equal(t, tt.want.template, actual)
		})
	}
}

func testMessageTemplate(id string, now time.Time) *entity.MessageTemplate {
	return &entity.MessageTemplate{
		TemplateID:    id,
		TitleTemplate: "件名: {{.Title}}",
		BodyTemplate:  "内容: {{.Body}}",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func fillIgnoreMessageTemplateField(t *entity.MessageTemplate, now time.Time) {
	if t == nil {
		return
	}
	t.CreatedAt = now
	t.UpdatedAt = now
}
