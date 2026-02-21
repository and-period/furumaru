package uuid

import (
	"github.com/google/uuid"
	base58 "github.com/jbenet/go-base58"
)

const base58Alphabet = base58.FlickrAlphabet

// New - uuidの生成
func New() string {
	return uuid.New().String()
}

// Base58Encode - uuidをbase58エンコードして返す
func Base58Encode(v4 string) string {
	id, _ := uuid.Parse(v4)
	return base58.EncodeAlphabet(id[:], base58Alphabet)
}
