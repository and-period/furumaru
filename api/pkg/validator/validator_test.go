package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
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
			name: "valid phone_number when mobile phone",
			input: &input{
				PhoneNumber: "090-1234-1234",
			},
			hasErr: false,
		},
		{
			name: "valid phone_number when land phone",
			input: &input{
				PhoneNumber: "03-1234-1234",
			},
			hasErr: false,
		},
		{
			name: "valid phone_number when e164",
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
		{
			name: "invalid phone_number",
			input: &input{
				PhoneNumber: "12345678",
			},
			opts:   []Option{},
			hasErr: true,
		},
		{
			name: "invalid phone_number without hyphen",
			input: &input{
				PhoneNumber: "09012341234",
			},
			opts:   []Option{},
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

func TestValidator_CustomValidation(t *testing.T) {
	t.Parallel()
	type input struct {
		Custom string `validate:"omitempty,custom"`
	}
	tests := []struct {
		name   string
		input  *input
		custom func(v *validator.Validate) error
		hasErr bool
	}{
		{
			name: "valid custom validation",
			input: &input{
				Custom: "12345678",
			},
			custom: func(v *validator.Validate) error {
				fn := func(fl validator.FieldLevel) bool {
					return fl.Field().String() == "12345678"
				}
				return v.RegisterValidation("custom", fn)
			},
			hasErr: false,
		},
		{
			name: "invalid custom validation",
			input: &input{
				Custom: "87654321",
			},
			custom: func(v *validator.Validate) error {
				fn := func(fl validator.FieldLevel) bool {
					return fl.Field().String() == "12345678"
				}
				return v.RegisterValidation("custom", fn)
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(WithCustomValidation(tt.custom))
			err := validator.Struct(tt.input)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}
