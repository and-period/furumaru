package types

// PostalCode - 郵便番号情報
type PostalCode struct {
	PostalCode     string `json:"postalCode"`     // 郵便番号
	PrefectureCode int32  `json:"prefectureCode"` // 都道府県コード
	Prefecture     string `json:"prefecture"`     // 都道府県名
	City           string `json:"city"`           // 市区町村名
	Town           string `json:"town"`           // 町域名
}

type PostalCodeResponse struct {
	*PostalCode
}
