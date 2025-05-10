package codes

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefectureNames(t *testing.T) {
	t.Parallel()
	tests := map[int32]string{
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
		t.Run(expect, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, expect, PrefectureNames[key])
		})
	}
}

func TestPrefectureJapanese(t *testing.T) {
	t.Parallel()
	tests := map[int32]string{
		0:  "",
		1:  "北海道",
		2:  "青森県",
		3:  "岩手県",
		4:  "宮城県",
		5:  "秋田県",
		6:  "山形県",
		7:  "福島県",
		8:  "茨城県",
		9:  "栃木県",
		10: "群馬県",
		11: "埼玉県",
		12: "千葉県",
		13: "東京都",
		14: "神奈川県",
		15: "新潟県",
		16: "富山県",
		17: "石川県",
		18: "福井県",
		19: "山梨県",
		20: "長野県",
		21: "岐阜県",
		22: "静岡県",
		23: "愛知県",
		24: "三重県",
		25: "滋賀県",
		26: "京都府",
		27: "大阪府",
		28: "兵庫県",
		29: "奈良県",
		30: "和歌山県",
		31: "鳥取県",
		32: "島根県",
		33: "岡山県",
		34: "広島県",
		35: "山口県",
		36: "徳島県",
		37: "香川県",
		38: "愛媛県",
		39: "高知県",
		40: "福岡県",
		41: "佐賀県",
		42: "長崎県",
		43: "熊本県",
		44: "大分県",
		45: "宮崎県",
		46: "鹿児島県",
		47: "沖縄県",
	}
	require.Len(t, PrefectureNames, 47)
	for key, expect := range tests {
		t.Run(expect, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, expect, PrefectureJapanese[key])
		})
	}
}

func TestPrefectureValues(t *testing.T) {
	t.Parallel()
	tests := map[string]int32{
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
		"北海道":       1,
		"青森県":       2,
		"岩手県":       3,
		"宮城県":       4,
		"秋田県":       5,
		"山形県":       6,
		"福島県":       7,
		"茨城県":       8,
		"栃木県":       9,
		"群馬県":       10,
		"埼玉県":       11,
		"千葉県":       12,
		"東京都":       13,
		"神奈川県":      14,
		"新潟県":       15,
		"富山県":       16,
		"石川県":       17,
		"福井県":       18,
		"山梨県":       19,
		"長野県":       20,
		"岐阜県":       21,
		"静岡県":       22,
		"愛知県":       23,
		"三重県":       24,
		"滋賀県":       25,
		"京都府":       26,
		"大阪府":       27,
		"兵庫県":       28,
		"奈良県":       29,
		"和歌山県":      30,
		"鳥取県":       31,
		"島根県":       32,
		"岡山県":       33,
		"広島県":       34,
		"山口県":       35,
		"徳島県":       36,
		"香川県":       37,
		"愛媛県":       38,
		"高知県":       39,
		"福岡県":       40,
		"佐賀県":       41,
		"長崎県":       42,
		"熊本県":       43,
		"大分県":       44,
		"宮崎県":       45,
		"鹿児島県":      46,
		"沖縄県":       47,
	}
	require.Len(t, PrefectureNames, 47)
	for key, expect := range tests {
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
		values    []int32
		expect    []string
		expectErr error
	}{
		{
			name:      "success",
			values:    []int32{1, 2, 3},
			expect:    []string{"hokkaido", "aomori", "iwate"},
			expectErr: nil,
		},
		{
			name:      "failed to convert",
			values:    []int32{0},
			expect:    nil,
			expectErr: ErrUnknownPrefecture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ToPrefectureNames(tt.values...)
			assert.ErrorIs(t, tt.expectErr, err)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestToPrefectureJapanese(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		values    int32
		expect    string
		expectErr error
	}{
		{
			name:      "success",
			values:    1,
			expect:    "北海道",
			expectErr: nil,
		},
		{
			name:      "failed to convert",
			values:    0,
			expect:    "",
			expectErr: ErrUnknownPrefecture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ToPrefectureJapanese(tt.values)
			assert.ErrorIs(t, tt.expectErr, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestToPrefectureValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		values    []string
		expect    []int32
		expectErr error
	}{
		{
			name:      "success",
			values:    []string{"hokkaido", "aomori", "iwate"},
			expect:    []int32{1, 2, 3},
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
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, ValidatePrefectureNames(key), expect)
		})
	}
}

func TestValidatePrefectureValues(t *testing.T) {
	t.Parallel()
	tests := map[int32]error{
		0: ErrUnknownPrefecture,
		1: nil,
	}
	for key, expect := range tests {
		t.Run(strconv.Itoa(int(key)), func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, ValidatePrefectureValues(key), expect)
		})
	}
}
