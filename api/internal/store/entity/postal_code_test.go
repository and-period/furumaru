package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/stretchr/testify/assert"
)

func TestPostalCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		postalCode *postalcode.PostalCode
		expect     *PostalCode
		hasErr     bool
	}{
		{
			name: "success",
			postalCode: &postalcode.PostalCode{
				PostalCode:     "1000014",
				PrefectureCode: "13",
				Prefecture:     "東京都",
				City:           "千代田区",
				Town:           "永田町",
				PrefectureKana: "ﾄｳｷｮｳﾄ",
				CityKana:       "ﾁﾖﾀﾞｸ",
				TownKana:       "ﾅｶﾞﾀﾁｮｳ",
			},
			expect: &PostalCode{
				PostalCode:     "1000014",
				PrefectureCode: 13,
				Prefecture:     "東京都",
				City:           "千代田区",
				Town:           "永田町",
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture code",
			postalCode: &postalcode.PostalCode{
				PostalCode:     "1000014",
				PrefectureCode: "tokyo",
				Prefecture:     "東京都",
				City:           "千代田区",
				Town:           "永田町",
				PrefectureKana: "ﾄｳｷｮｳﾄ",
				CityKana:       "ﾁﾖﾀﾞｸ",
				TownKana:       "ﾅｶﾞﾀﾁｮｳ",
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewPostalCode(tt.postalCode)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
