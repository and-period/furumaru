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
		"":     0,
		"北海道":  1,
		"青森県":  2,
		"岩手県":  3,
		"宮城県":  4,
		"秋田県":  5,
		"山形県":  6,
		"福島県":  7,
		"茨城県":  8,
		"栃木県":  9,
		"群馬県":  10,
		"埼玉県":  11,
		"千葉県":  12,
		"東京都":  13,
		"神奈川県": 14,
		"新潟県":  15,
		"富山県":  16,
		"石川県":  17,
		"福井県":  18,
		"山梨県":  19,
		"長野県":  20,
		"岐阜県":  21,
		"静岡県":  22,
		"愛知県":  23,
		"三重県":  24,
		"滋賀県":  25,
		"京都府":  26,
		"大阪府":  27,
		"兵庫県":  28,
		"奈良県":  29,
		"和歌山県": 30,
		"鳥取県":  31,
		"島根県":  32,
		"岡山県":  33,
		"広島県":  34,
		"山口県":  35,
		"徳島県":  36,
		"香川県":  37,
		"愛媛県":  38,
		"高知県":  39,
		"福岡県":  40,
		"佐賀県":  41,
		"長崎県":  42,
		"熊本県":  43,
		"大分県":  44,
		"宮崎県":  45,
		"鹿児島県": 46,
		"沖縄県":  47,
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

func TestValidatePrefectureNames(t *testing.T) {
	t.Parallel()
	tests := map[string]error{
		"東京都":     nil,
		"カリフォルニア": ErrUnknownPrefecture,
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
