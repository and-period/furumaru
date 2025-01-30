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

func TestReportTemplate(t *testing.T) {
	assert.NotNil(t, NewReportTemplate(nil))
}

func TestReportTemplate_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	tmpl := testReportTemplate("report-id", now())
	err = db.DB.Create(&tmpl).Error
	require.NoError(t, err)

	type args struct {
		reportID entity.ReportTemplateID
	}
	type want struct {
		template *entity.ReportTemplate
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
				reportID: "report-id",
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
				reportID: "other-id",
			},
			want: want{
				template: nil,
				err:      database.ErrNotFound,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &reportTemplate{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.reportID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.template, actual)
		})
	}
}

func testReportTemplate(id entity.ReportTemplateID, now time.Time) *entity.ReportTemplate {
	return &entity.ReportTemplate{
		TemplateID: id,
		Template:   "テンプレート: {{.Body}}",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
