package codes

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefectureNames(t *testing.T) {
	t.Parallel()
	tests := map[int64]string{
		0:  "",
		1:  "hokkaido",
		2:  "aomori",
		3:  "iwate",
		4:  "miyagi",
		5:  "akita",
		6:  "yamagata",
		7:  "fukushima",
		8:  "ibaraki",
		9:  "tochigi",
		10: "gunma",
		11: "saitama",
		12: "chiba",
		13: "tokyo",
		14: "kanagawa",
		15: "niigata",
		16: "toyama",
		17: "ishikawa",
		18: "fukui",
		19: "yamanashi",
		20: "nagano",
		21: "gifu",
		22: "shizuoka",
		23: "aichi",
		24: "mie",
		25: "shiga",
		26: "kyoto",
		27: "osaka",
		28: "hyogo",
		29: "nara",
		30: "wakayama",
		31: "tottori",
		32: "shimane",
		33: "okayama",
		34: "hiroshima",
		35: "yamaguchi",
		36: "tokushima",
		37: "kagawa",
		38: "ehime",
		39: "kochi",
		40: "fukuoka",
		41: "saga",
		42: "nagasaki",
		43: "kumamoto",
		44: "oita",
		45: "miyazaki",
		46: "kagoshima",
		47: "okinawa",
	}
	require.Len(t, PrefectureNames, 47)
	for key, expect := range tests {
		key, expect := key, expect
		t.Run(expect, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, expect, PrefectureNames[key])
		})
	}
}

func TestPrefectureValues(t *testing.T) {
	t.Parallel()
	tests := map[string]int64{
		"":          0,
		"hokkaido":  1,
		"aomori":    2,
		"iwate":     3,
		"miyagi":    4,
		"akita":     5,
		"yamagata":  6,
		"fukushima": 7,
		"ibaraki":   8,
		"tochigi":   9,
		"gunma":     10,
		"saitama":   11,
		"chiba":     12,
		"tokyo":     13,
		"kanagawa":  14,
		"niigata":   15,
		"toyama":    16,
		"ishikawa":  17,
		"fukui":     18,
		"yamanashi": 19,
		"nagano":    20,
		"gifu":      21,
		"shizuoka":  22,
		"aichi":     23,
		"mie":       24,
		"shiga":     25,
		"kyoto":     26,
		"osaka":     27,
		"hyogo":     28,
		"nara":      29,
		"wakayama":  30,
		"tottori":   31,
		"shimane":   32,
		"okayama":   33,
		"hiroshima": 34,
		"yamaguchi": 35,
		"tokushima": 36,
		"kagawa":    37,
		"ehime":     38,
		"kochi":     39,
		"fukuoka":   40,
		"saga":      41,
		"nagasaki":  42,
		"kumamoto":  43,
		"oita":      44,
		"miyazaki":  45,
		"kagoshima": 46,
		"okinawa":   47,
	}
	require.Len(t, PrefectureNames, 47)
	for key, expect := range tests {
		key, expect := key, expect
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, expect, PrefectureValues[key])
		})
	}
}

func TestToPrefectureNames(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		values    []int64
		expect    []string
		expectErr error
	}{
		{
			name:      "success",
			values:    []int64{1, 2, 3},
			expect:    []string{"hokkaido", "aomori", "iwate"},
			expectErr: nil,
		},
		{
			name:      "failed to convert",
			values:    []int64{0},
			expect:    nil,
			expectErr: ErrUnknownPrefecture,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ToPrefectureNames(tt.values...)
			assert.ErrorIs(t, tt.expectErr, err)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestToPrefectureValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		values    []string
		expect    []int64
		expectErr error
	}{
		{
			name:      "success",
			values:    []string{"hokkaido", "aomori", "iwate"},
			expect:    []int64{1, 2, 3},
			expectErr: nil,
		},
		{
			name:      "failed to convert",
			values:    []string{""},
			expect:    nil,
			expectErr: ErrUnknownPrefecture,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ToPrefectureValues(tt.values...)
			assert.ErrorIs(t, tt.expectErr, err)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestValidatePrefectureNames(t *testing.T) {
	t.Parallel()
	tests := map[string]error{
		"tokyo":      nil,
		"california": ErrUnknownPrefecture,
	}
	for key, expect := range tests {
		key, expect := key, expect
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, ValidatePrefectureNames(key), expect)
		})
	}
}

func TestValidatePrefectureValues(t *testing.T) {
	t.Parallel()
	tests := map[int64]error{
		0: ErrUnknownPrefecture,
		1: nil,
	}
	for key, expect := range tests {
		key, expect := key, expect
		t.Run(strconv.FormatInt(key, 10), func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, ValidatePrefectureValues(key), expect)
		})
	}
}
