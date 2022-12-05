package codes

import (
	"errors"

	"github.com/and-period/furumaru/api/pkg/set"
)

var ErrUnknownPrefecture = errors.New("entity: unknown prefecture")

var PrefectureNames = map[int64]string{
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

var PrefectureValues = map[string]int64{
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

func ToPrefectureName(value int64) (string, error) {
	name, ok := PrefectureNames[value]
	if !ok {
		return "", ErrUnknownPrefecture
	}
	return name, nil
}

func ToPrefectureNames(values ...int64) ([]string, error) {
	return set.UniqWithErr(values, ToPrefectureName)
}

func ToPrefectureValue(name string) (int64, error) {
	value, ok := PrefectureValues[name]
	if !ok {
		return 0, ErrUnknownPrefecture
	}
	return value, nil
}

func ToPrefectureValues(names ...string) ([]int64, error) {
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

func ValidatePrefectureValues(values ...int64) error {
	for _, value := range values {
		if _, ok := PrefectureNames[value]; !ok {
			return ErrUnknownPrefecture
		}
	}
	return nil
}
