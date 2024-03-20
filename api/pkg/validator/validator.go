//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package validator

import (
	"strings"

	regexp "github.com/dlclark/regexp2"
	validator "github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(s interface{}) error // 構造体のバリデーション検証
}

type options struct {
	password *PasswordParams
}

type Option func(opts *options)

// PasswordParams - 追加の検証項目
type PasswordParams struct {
	RequireNumbers   bool // 少なくとも１つの数字を含む
	RequireSymbols   bool // 少なくとも１つの特殊文字を含む
	RequireUppercase bool // 少なくとも１つの大文字を含む
	RequireLowercase bool // 少なくとも１つの小文字を含む
}

func WithPasswordValidation(params *PasswordParams) Option {
	return func(opts *options) {
		opts.password = params
	}
}

//nolint:gosec
const (
	hiraganaString    = "^[ぁ-ゔー]*$"
	passwordString    = "^[a-zA-Z0-9_!@#$_%^&*.?()-=+]*$"
	phoneNumberString = "^0\\d{1,4}-\\d{1,4}-\\d{4}$"
	e164String        = "^\\+[1-9]?[0-9]{7,14}$"
)

//nolint:errcheck
func NewValidator(opts ...Option) Validator {
	v := validator.New()

	// オプション値の追加
	dopts := &options{}
	for i := range opts {
		opts[i](dopts)
	}

	hiraganaRegex := regexp.MustCompile(hiraganaString, 0)
	passwordRegex := compilePasswordRegex(dopts.password)
	phoneNumberRegex := regexp.MustCompile(phoneNumberString, 0)

	// hiragana - 正規表現を使用して平仮名のみであるかの検証
	v.RegisterValidation("hiragana", validateHiragana(hiraganaRegex))
	// password - 正規表現を利用してパスワードに使用不可な文字を含んでいないかの検証
	v.RegisterValidation("password", validatePassword(passwordRegex))
	// phone_number - 正規表現を利用して電話番号（ハイフンあり）の形式であるかの検証
	v.RegisterValidation("phone_number", validatePhoneNumber(phoneNumberRegex))

	return v
}

func validateHiragana(regex *regexp.Regexp) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		match, _ := regex.MatchString(fl.Field().String())
		return match
	}
}

func validatePassword(regex *regexp.Regexp) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		match, _ := regex.MatchString(fl.Field().String())
		return match
	}
}

func validatePhoneNumber(regex *regexp.Regexp) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		match, _ := regex.MatchString(fl.Field().String())
		if match {
			return true
		}
		// TODO: フロントエンドの移行が完了し次第、後続処理は削除する
		e164Regex := regexp.MustCompile(e164String, 0)
		match, _ = e164Regex.MatchString(fl.Field().String())
		return match
	}
}

func compilePasswordRegex(params *PasswordParams) *regexp.Regexp {
	if params == nil {
		return regexp.MustCompile(passwordString, 0)
	}
	b := &strings.Builder{}
	b.WriteString("^")
	if params.RequireNumbers {
		b.WriteString("(?=.*[0-9])")
	}
	if params.RequireSymbols {
		b.WriteString("(?=.*[_!@#$_%^&*.?()\\-=+])")
	}
	if params.RequireUppercase {
		b.WriteString("(?=.*[A-Z])")
	}
	if params.RequireLowercase {
		b.WriteString("(?=.*[a-z])")
	}
	b.WriteString(passwordString[1:]) // はじめの「^」を除いた文字列を代入
	return regexp.MustCompile(b.String(), 0)
}
