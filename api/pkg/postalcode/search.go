package postalcode

import (
	"context"
	"fmt"
	"net/http"
)

// PostalCode - 郵便番号詳細情報
type PostalCode struct {
	PostalCode     string `json:"zipcode"`  // 郵便番号
	PrefectureCode string `json:"prefcode"` // 都道府県コード
	Prefecture     string `json:"address1"` // 都道府県名
	City           string `json:"address2"` // 市区町村名
	Town           string `json:"address3"` // 町域名
	PrefectureKana string `json:"kana1"`    // 都道府県名(カナ)
	CityKana       string `json:"kana2"`    // 市区町村名(カナ)
	TownKana       string `json:"kana3"`    // 町域名(カナ)
}

type searchResponse struct {
	Status      int           `json:"status"`            // HTTPステータスコード
	Message     string        `json:"message,omitempty"` // エラーメッセージ
	PostalCodes []*PostalCode `json:"results"`           // 住所一覧
}

// reference: http://zipcloud.ibsnet.co.jp/doc/api
func (c *client) Search(ctx context.Context, code string) (*PostalCode, error) {
	const format = "https://zipcloud.ibsnet.co.jp/api/search?zipcode=%s&limit=1"
	url := fmt.Sprintf(format, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, c.newError(err)
	}

	res := &searchResponse{}
	if err := c.do(req, res); err != nil {
		return nil, c.newError(err)
	}
	if len(res.PostalCodes) == 0 {
		return nil, ErrNotFound
	}
	return res.PostalCodes[0], nil
}
