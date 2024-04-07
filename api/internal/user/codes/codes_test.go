package codes

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestRegisterValidations(t *testing.T) {
	t.Parallel()

	type input struct {
		AccountID string `validate:"omitempty,account_id"`
	}
	tests := []struct {
		name   string
		input  *input
		hasErr bool
	}{
		{
			name: "valid all",
			input: &input{
				AccountID: "account-id_1234",
			},
			hasErr: false,
		},
		{
			name: "invalid account_id",
			input: &input{
				AccountID: "アカウントID",
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := validator.NewValidator(validator.WithCustomValidation(RegisterValidations))
			err := v.Struct(tt.input)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}
