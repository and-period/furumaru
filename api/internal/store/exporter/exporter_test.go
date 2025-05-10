package exporter

import (
	"bytes"
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExporter(t *testing.T) {
	t.Parallel()
	receipt := &testReceipt{
		field1: "バリュー1",
		field2: "バリュー2",
		field3: "バリュー3",
	}
	tests := []struct {
		name         string
		encodingType codes.CharacterEncodingType
		receipt      *testReceipt
		expect       string
	}{
		{
			name:         "success utf-8",
			encodingType: codes.CharacterEncodingTypeUTF8,
			receipt:      receipt,
			expect: "フィールド1,フィールド2,フィールド3\n" +
				"バリュー1,バリュー2,バリュー3\n",
		},
		{
			name:         "success shift-jis",
			encodingType: codes.CharacterEncodingTypeShiftJIS,
			receipt:      receipt,
			expect: "\x83t\x83B\x81[\x83\x8b\x83h1,\x83t\x83B\x81[\x83\x8b\x83h2,\x83t\x83B\x81[\x83\x8b\x83h3\n" +
				"\x83o\x83\x8a\x83\x85\x81[1,\x83o\x83\x8a\x83\x85\x81[2,\x83o\x83\x8a\x83\x85\x81[3\n",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			writer := NewExporter(buf, tt.encodingType)
			err := writer.WriteHeader(&testReceipt{})
			require.NoError(t, err)
			err = writer.WriteBody(tt.receipt)
			require.NoError(t, err)
			err = writer.Flush()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, buf.String())
		})
	}
}

type testReceipt struct {
	field1 string
	field2 string
	field3 string
}

func (r *testReceipt) Header() []string {
	return []string{"フィールド1", "フィールド2", "フィールド3"}
}

func (r *testReceipt) Record() []string {
	return []string{r.field1, r.field2, r.field3}
}
