package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestReport_Fields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		report *Report
		expect map[string]string
	}{
		{
			name: "success",
			report: &Report{
				ReportID:   ReportIDReceivedContact,
				Overview:   "レポートの概要です。",
				Link:       "https://and-period.jp",
				ReceivedAt: jst.Date(2022, 7, 14, 18, 30, 0, 0),
			},
			expect: map[string]string{
				"Overview":   "レポートの概要です。",
				"Link":       "https://and-period.jp",
				"ReceivedAt": "2022-07-14 18:30:00",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.report.Fields())
		})
	}
}
