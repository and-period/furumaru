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

func TestReportTemplate(t *testing.T) {
	assert.NotNil(t, NewReportTemplate(nil))
}

func TestReportTemplate_Get(t *testing.T) {
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

	_ = m.dbDelete(ctx, reportTemplateTable)
	tmpl := testReportTemplate("report-id", now())
	err = m.db.DB.Create(&tmpl).Error
	require.NoError(t, err)

	type args struct {
		reportID string
	}
	type want struct {
		template *entity.ReportTemplate
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
				reportID: "report-id",
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
				reportID: "other-id",
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

			db := &reportTemplate{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.reportID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreReportTemplateField(actual, now())
			assert.Equal(t, tt.want.template, actual)
		})
	}
}

func testReportTemplate(id string, now time.Time) *entity.ReportTemplate {
	return &entity.ReportTemplate{
		TemplateID: id,
		Template:   "テンプレート: {{.Body}}",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func fillIgnoreReportTemplateField(t *entity.ReportTemplate, now time.Time) {
	if t == nil {
		return
	}
	t.CreatedAt = now
	t.UpdatedAt = now
}
