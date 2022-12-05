package entity

import (
	"strconv"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/postalcode"
)

// PostalCode - 郵便番号詳細情報
type PostalCode struct {
	PostalCode     string // 郵便番号
	PrefectureCode string // 都道府県コード
	Prefecture     string // 都道府県名
	City           string // 市区町村名
	Town           string // 町域名
}

func NewPostalCode(p *postalcode.PostalCode) (*PostalCode, error) {
	code, err := strconv.ParseInt(p.PrefectureCode, 10, 64)
	if err != nil {
		return nil, err
	}
	prefectureCode, err := codes.ToPrefectureName(code)
	if err != nil {
		return nil, err
	}
	return &PostalCode{
		PostalCode:     p.PostalCode,
		PrefectureCode: prefectureCode,
		Prefecture:     p.Prefecture,
		City:           p.City,
		Town:           p.Town,
	}, nil
}
