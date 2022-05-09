package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	t.Parallel()

	type input struct {
		Name        string `validate:"omitempty,max=4"`
		Hiragana    string `validate:"omitempty,hiragana"`
		Password    string `validate:"omitempty,password"`
		PhoneNumber string `validate:"omitempty,phone_number"`
	}
	tests := []struct {
		name   string
		input  *input
		hasErr bool
	}{
		{
			name:   "valid all",
			input:  &input{},
			hasErr: false,
		},
		{
			name: "valid name",
			input: &input{
				Name: "1234",
			},
			hasErr: false,
		},
		{
			name: "valid hiragana",
			input: &input{
				Hiragana: "あいうえお",
			},
			hasErr: false,
		},
		{
			name: "valid password",
			input: &input{
				Password: "abcxyzABCXYZ012789_!@#$_%^&*.?()-=+",
			},
			hasErr: false,
		},
		{
			name: "valid phone number",
			input: &input{
				PhoneNumber: "+819012341234",
			},
			hasErr: false,
		},
		{
			name: "invalid name",
			input: &input{
				Name: "12345",
			},
			hasErr: true,
		},
		{
			name: "invalid hiragana",
			input: &input{
				Hiragana: "アイウエオ",
			},
			hasErr: true,
		},
		{
			name: "invalid password",
			input: &input{
				Password: "[];:{}/",
			},
			hasErr: true,
		},
		{
			name: "invalid phone number",
			input: &input{
				PhoneNumber: "090-1234-1234",
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			validator := NewValidator()
			err := validator.Struct(tt.input)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}
