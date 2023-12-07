package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	t.Parallel()

	type input struct {
		Name     string `validate:"omitempty,max=4"`
		Hiragana string `validate:"omitempty,hiragana"`
		Password string `validate:"omitempty,password"`
	}
	tests := []struct {
		name   string
		input  *input
		opts   []Option
		hasErr bool
	}{
		{
			name:   "valid all",
			input:  &input{},
			opts:   []Option{},
			hasErr: false,
		},
		{
			name: "valid name",
			input: &input{
				Name: "1234",
			},
			opts:   []Option{},
			hasErr: false,
		},
		{
			name: "valid hiragana",
			input: &input{
				Hiragana: "あいうえお",
			},
			opts:   []Option{},
			hasErr: false,
		},
		{
			name: "valid password",
			input: &input{
				Password: "abcxyzABCXYZ012789_!@#$_%^&*.?()-=+",
			},
			opts:   []Option{},
			hasErr: false,
		},
		{
			name: "valid password with optional",
			input: &input{
				Password: "abcxyzABCXYZ012789_!@#$_%^&*.?()-=+",
			},
			opts: []Option{WithPasswordValidation(&PasswordParams{
				RequireNumbers:   true,
				RequireSymbols:   true,
				RequireUppercase: true,
				RequireLowercase: true,
			})},
			hasErr: false,
		},
		{
			name: "invalid name",
			input: &input{
				Name: "12345",
			},
			opts:   []Option{},
			hasErr: true,
		},
		{
			name: "invalid hiragana",
			input: &input{
				Hiragana: "アイウエオ",
			},
			opts:   []Option{},
			hasErr: true,
		},
		{
			name: "invalid password",
			input: &input{
				Password: "[];:{}/",
			},
			opts:   []Option{},
			hasErr: true,
		},
		{
			name: "invalid password only numbers",
			input: &input{
				Password: "12345678",
			},
			opts: []Option{WithPasswordValidation(&PasswordParams{
				RequireNumbers:   true,
				RequireSymbols:   false,
				RequireUppercase: false,
				RequireLowercase: true,
			})},
			hasErr: true,
		},
		{
			name: "invalid password with numbers",
			input: &input{
				Password: "abcxyzABCXYZ_!@#$_%^&*.?()-=+",
			},
			opts: []Option{WithPasswordValidation(&PasswordParams{
				RequireNumbers:   true,
				RequireSymbols:   false,
				RequireUppercase: false,
				RequireLowercase: false,
			})},
			hasErr: true,
		},
		{
			name: "invalid password with synbols",
			input: &input{
				Password: "abcxyzABCXYZ012789",
			},
			opts: []Option{WithPasswordValidation(&PasswordParams{
				RequireNumbers:   false,
				RequireSymbols:   true,
				RequireUppercase: false,
				RequireLowercase: false,
			})},
			hasErr: true,
		},
		{
			name: "invalid password with uppercase",
			input: &input{
				Password: "abcxyz012789_!@#$_%^&*.?()-=+",
			},
			opts: []Option{WithPasswordValidation(&PasswordParams{
				RequireNumbers:   false,
				RequireSymbols:   false,
				RequireUppercase: true,
				RequireLowercase: false,
			})},
			hasErr: true,
		},
		{
			name: "invalid password with lowercase",
			input: &input{
				Password: "ABCXYZ012789_!@#$_%^&*.?()-=+",
			},
			opts: []Option{WithPasswordValidation(&PasswordParams{
				RequireNumbers:   false,
				RequireSymbols:   false,
				RequireUppercase: false,
				RequireLowercase: true,
			})},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.opts...)
			err := validator.Struct(tt.input)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}
