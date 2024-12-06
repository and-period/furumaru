package codes

import (
	"errors"

	"github.com/and-period/furumaru/api/pkg/set"
)

var ErrUnknownPrefecture = errors.New("entity: unknown prefecture")

var PrefectureNames = map[int32]string{
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

var PrefectureJapanese = map[int32]string{
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

var PrefectureValues = map[string]int32{
	// 英語 -> 都道府県コード
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
	// 日本語 -> 都道府県コード
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

func ToPrefectureName(value int32) (string, error) {
	name, ok := PrefectureNames[value]
	if !ok {
		return "", ErrUnknownPrefecture
	}
	return name, nil
}

func ToPrefectureNames(values ...int32) ([]string, error) {
	return set.UniqWithErr(values, ToPrefectureName)
}

func ToPrefectureJapanese(value int32) (string, error) {
	name, ok := PrefectureJapanese[value]
	if !ok {
		return "", ErrUnknownPrefecture
	}
	return name, nil
}

func ToPrefectureValue(name string) (int32, error) {
	value, ok := PrefectureValues[name]
	if !ok {
		return 0, ErrUnknownPrefecture
	}
	return value, nil
}

func ToPrefectureValues(names ...string) ([]int32, error) {
	return set.UniqWithErr(names, ToPrefectureValue)
}

func ValidatePrefectureNames(names ...string) error {
	for _, name := range names {
		if _, ok := PrefectureValues[name]; !ok {
			return ErrUnknownPrefecture
		}
	}
	return nil
}

func ValidatePrefectureValues(values ...int32) error {
	for _, value := range values {
		if _, ok := PrefectureNames[value]; !ok {
			return ErrUnknownPrefecture
		}
	}
	return nil
}
