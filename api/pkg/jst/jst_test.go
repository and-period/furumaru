package jst

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, Now())
}

func TestDate(t *testing.T) {
	t.Parallel()
	now := Now()
	assert.Equal(t, now, Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()))
}

func TestToTime(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		now                    time.Time
		expectBeginningOfDay   time.Time
		expectBeginningOfMonth time.Time
		expectEndOfDay         time.Time
		expectEndOfMonth       time.Time
	}{
		{
			name:                   "success",
			now:                    time.Date(2021, 8, 2, 18, 30, 0, 0, jst),
			expectBeginningOfDay:   time.Date(2021, 8, 2, 0, 0, 0, 0, jst),
			expectBeginningOfMonth: time.Date(2021, 8, 1, 0, 0, 0, 0, jst),
			expectEndOfDay:         time.Date(2021, 8, 3, 0, 0, 0, 0, jst).Add(-time.Nanosecond),
			expectEndOfMonth:       time.Date(2021, 9, 1, 0, 0, 0, 0, jst).Add(-time.Nanosecond),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expectBeginningOfDay, BeginningOfDay(tt.now))
			assert.Equal(t, tt.expectBeginningOfMonth, BeginningOfMonth(tt.now.Year(), int(tt.now.Month())))
			assert.Equal(t, tt.expectEndOfDay, EndOfDay(tt.now))
			assert.Equal(t, tt.expectEndOfMonth, EndOfMonth(tt.now.Year(), int(tt.now.Month())))
		})
	}
}

func TestWithInPeriod(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name    string
		target  time.Time
		startAt time.Time
		endAt   time.Time
		expect  bool
	}{
		{
			name:    "before start at",
			target:  now.Add(-time.Hour),
			startAt: now,
			endAt:   now.Add(time.Hour),
			expect:  false,
		},
		{
			name:    "equal start at",
			target:  now,
			startAt: now,
			endAt:   now.Add(time.Hour),
			expect:  true,
		},
		{
			name:    "with in the period",
			target:  now,
			startAt: now.Add(-time.Hour),
			endAt:   now.Add(time.Hour),
			expect:  true,
		},
		{
			name:    "equal end at",
			target:  now,
			startAt: now.Add(-time.Hour),
			endAt:   now,
			expect:  false,
		},
		{
			name:    "after end at",
			target:  now.Add(time.Hour),
			startAt: now.Add(-time.Hour),
			endAt:   now,
			expect:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, WithInPeriod(tt.target, tt.startAt, tt.endAt))
		})
	}
}

func TestToString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		now            time.Time
		format         string
		expectFormat   string
		expectYYYYMMDD string
		expectYYYYMM   string
		expectHHMMSS   string
		expectHHMM     string
	}{
		{
			name:           "success",
			now:            time.Date(2021, 8, 2, 18, 30, 0, 0, jst),
			format:         "2006/01/02 15:04:05",
			expectFormat:   "2021/08/02 18:30:00",
			expectYYYYMMDD: "20210802",
			expectYYYYMM:   "202108",
			expectHHMMSS:   "183000",
			expectHHMM:     "1830",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expectFormat, Format(tt.now, tt.format))
			assert.Equal(t, tt.expectYYYYMMDD, FormatYYYYMMDD(tt.now))
			assert.Equal(t, tt.expectYYYYMM, FormatYYYYMM(tt.now))
			assert.Equal(t, tt.expectHHMMSS, FormatHHMMSS(tt.now))
			assert.Equal(t, tt.expectHHMM, FormatHHMM(tt.now))
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		format    string
		target    string
		expect    time.Time
		expectErr bool
	}{
		{
			name:      "success",
			format:    "20060102",
			target:    "20210802",
			expect:    time.Date(2021, 8, 2, 0, 0, 0, 0, jst),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := Parse(tt.format, tt.target)
			assert.Equal(t, tt.expectErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestParseFromYYYYMMDD(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		target    string
		expect    time.Time
		expectErr bool
	}{
		{
			name:      "success",
			target:    "20210802",
			expect:    time.Date(2021, 8, 2, 0, 0, 0, 0, jst),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ParseFromYYYYMMDD(tt.target)
			assert.Equal(t, tt.expectErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestParseFromYYYYMM(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		target    string
		expect    time.Time
		expectErr bool
	}{
		{
			name:      "success",
			target:    "202108",
			expect:    time.Date(2021, 8, 1, 0, 0, 0, 0, jst),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ParseFromYYYYMM(tt.target)
			assert.Equal(t, tt.expectErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestParseFromHHMMSS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		target    string
		expect    time.Time
		expectErr bool
	}{
		{
			name:      "success",
			target:    "183000",
			expect:    time.Date(0, 1, 1, 18, 30, 0, 0, jst),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ParseFromHHMMSS(tt.target)
			assert.Equal(t, tt.expectErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestParseFromHHMM(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		target    string
		expect    time.Time
		expectErr bool
	}{
		{
			name:      "success",
			target:    "1830",
			expect:    time.Date(0, 1, 1, 18, 30, 0, 0, jst),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ParseFromHHMM(tt.target)
			assert.Equal(t, tt.expectErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestParseFromUnix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		target    int64
		expect    time.Time
		expectErr bool
	}{
		{
			name:      "success",
			target:    1627896600,
			expect:    time.Date(2021, 8, 2, 18, 30, 0, 0, jst),
			expectErr: false,
		},
		{
			name:      "success to zero time",
			target:    0,
			expect:    time.Time{}.In(jst),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, ParseFromUnix(tt.target))
		})
	}
}

func TestFiscalYear(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		target time.Time
		expect int
	}{
		{
			name:   "success to return this year",
			target: time.Date(2022, 4, 1, 0, 0, 0, 0, jst),
			expect: 2022,
		},
		{
			name:   "success to return last year",
			target: time.Date(2022, 1, 1, 0, 0, 0, 0, jst),
			expect: 2021,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, FiscalYear(tt.target))
		})
	}
}
