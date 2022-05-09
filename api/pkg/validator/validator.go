//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package validator

import (
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(s interface{}) error                                                               // 構造体のバリデーション検証
	RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error // カスタムバリデーションの登録
}

const (
	hiraganaString    = "^[ぁ-ゔー]*$"
	passwordString    = "^[a-zA-Z0-9_!@#$_%^&*.?()-=+]*$"
	phoneNumberString = "^\\+[0-9]{11,17}$"
)

var (
	hiraganaRegex    = regexp.MustCompile(hiraganaString)
	passwordRegex    = regexp.MustCompile(passwordString)
	phoneNumberRegex = regexp.MustCompile(phoneNumberString)
)

//nolint:errcheck
func NewValidator() Validator {
	v := validator.New()

	// hiragana - 正規表現を使用して平仮名のみであるかの検証
	v.RegisterValidation("hiragana", validateHiragana)
	// password - 正規表現を利用してパスワードに使用不可な文字を含んでいないかの検証
	v.RegisterValidation("password", validatePassword)
	// phone_number - 電話番号のフォーマットが正しいかの検証
	v.RegisterValidation("phone_number", validatePhoneNumber)

	return v
}

func validateHiragana(fl validator.FieldLevel) bool {
	return hiraganaRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	return phoneNumberRegex.MatchString(fl.Field().String())
}
