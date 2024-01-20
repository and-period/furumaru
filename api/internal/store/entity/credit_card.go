package entity

import (
	"errors"
	"strconv"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"golang.org/x/text/width"
)

var (
	errCreditCardExpired       = errors.New("entity: this credit card has already expired")
	errCreditCardInvalidCVV    = errors.New("entity: invalid credit card security code")
	errCreditCardInvalidNumber = errors.New("entity: invalid credit card number")
)

// CreditCard - クレジットカード情報
type CreditCard struct {
	Name   string
	Number string
	Month  int64
	Year   int64
	CVV    string
}

type NewCreditCardParams struct {
	Name   string
	Number string
	Month  int64
	Year   int64
	CVV    string
}

func NewCreditCard(params *NewCreditCardParams) *CreditCard {
	return &CreditCard{
		Name:   width.Narrow.String(params.Name),
		Number: params.Number,
		Month:  params.Month,
		Year:   params.Year,
		CVV:    params.CVV,
	}
}

func (c *CreditCard) Validate(now time.Time) error {
	if err := c.validateNumber(); err != nil {
		return err
	}
	if err := c.validateCVV(); err != nil {
		return err
	}
	return c.validateExpired(now)
}

func (c *CreditCard) validateNumber() error {
	length := len(c.Number)
	if length < 13 || 16 < length {
		return errCreditCardInvalidNumber
	}
	// Luhnアルゴリズムによる検証
	// https://ja.wikipedia.org/wiki/Luhn%E3%82%A2%E3%83%AB%E3%82%B4%E3%83%AA%E3%82%BA%E3%83%A0
	var sum int64
	var alternate bool
	for i := length - 1; i >= 0; i-- {
		num, err := strconv.ParseInt(string(c.Number[i]), 10, 64)
		if err != nil {
			return err
		}
		if alternate {
			double := num * 2
			num = (double / 10) + (double % 10)
		}
		alternate = !alternate
		sum += num
	}
	if sum%10 > 0 {
		return errCreditCardInvalidNumber
	}
	return nil
}

func (c *CreditCard) validateExpired(now time.Time) error {
	expired := jst.EndOfMonth(int(c.Year), int(c.Month))
	if expired.Before(now) {
		return errCreditCardExpired
	}
	return nil
}

func (c *CreditCard) validateCVV() error {
	length := len(c.CVV)
	if length < 3 || 4 < length {
		return errCreditCardInvalidCVV
	}
	return nil
}
