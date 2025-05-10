package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestPostalCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		postalCode *entity.PostalCode
		expect     *PostalCode
	}{
		{
			name: "success",
			postalCode: &entity.PostalCode{
				PostalCode:     "1000014",
				PrefectureCode: 13,
				Prefecture:     "東京都",
				City:           "千代田区",
				Town:           "永田町",
			},
			expect: &PostalCode{
				PostalCode: response.PostalCode{
					PostalCode:     "1000014",
					PrefectureCode: 13,
					Prefecture:     "東京都",
					City:           "千代田区",
					Town:           "永田町",
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPostalCode(tt.postalCode)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPostalCode_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		postalCode *PostalCode
		expect     *response.PostalCode
	}{
		{
			name: "success",
			postalCode: &PostalCode{
				PostalCode: response.PostalCode{
					PostalCode:     "1000014",
					PrefectureCode: 13,
					Prefecture:     "東京都",
					City:           "千代田区",
					Town:           "永田町",
				},
			},
			expect: &response.PostalCode{
				PostalCode:     "1000014",
				PrefectureCode: 13,
				Prefecture:     "東京都",
				City:           "千代田区",
				Town:           "永田町",
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.postalCode.Response())
		})
	}
}
