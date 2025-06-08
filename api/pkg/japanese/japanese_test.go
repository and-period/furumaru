package japanese

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHiraganaToKatakana(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		hiragana string
		expect   string
	}{
		{
			name:     "ひらがな -> カタカナ",
			hiragana: "いーしゃんてん",
			expect:   "イーシャンテン",
		},
		{
			name:     "カタカナ -> カタカナ",
			hiragana: "サンショクドウジュン",
			expect:   "サンショクドウジュン",
		},
		{
			name:     "漢字",
			hiragana: "嶺上開花",
			expect:   "嶺上開花",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := HiraganaToKatakana(tt.hiragana)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
