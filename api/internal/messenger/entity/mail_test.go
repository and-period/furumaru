package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestTemplateBuilder(t *testing.T) {
	builder := NewTemplateDataBuilder().
		Data(map[string]string{"key": "value"}).
		YearMonth(jst.Date(2022, 1, 2, 18, 30, 0, 0)).
		Name("中村 広大").
		Email("test-user@and-period.jp").
		Password("!Qaz2wsx").
		WebURL("http://example.com").
		Contact("件名", "本文")
	data := builder.Build()
	assert.Equal(t, "value", data["key"])
	assert.Equal(t, "2022年01月", data["年月"])
	assert.Equal(t, "中村 広大", data["氏名"])
	assert.Equal(t, "test-user@and-period.jp", data["メールアドレス"])
	assert.Equal(t, "!Qaz2wsx", data["パスワード"])
	assert.Equal(t, "http://example.com", data["サイトURL"])
	assert.Equal(t, "件名", data["件名"])
	assert.Equal(t, "本文", data["本文"])
}
