package csv

// 文字エンコード種別
type CharacterEncodingType int32

const (
	CharacterEncodingTypeUTF8 CharacterEncodingType = iota
	CharacterEncodingTypeShiftJIS
)
