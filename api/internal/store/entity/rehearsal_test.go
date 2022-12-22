package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRehearsal_TableName(t *testing.T) {
	t.Parallel()
	r := &Rehearsal{}
	assert.Equal(t, "rehearsals", r.TableName())
}

func TestRehearsal_PrimaryKey(t *testing.T) {
	t.Parallel()
	r := &Rehearsal{LiveID: "live-id"}
	assert.Equal(t, map[string]interface{}{"live_id": "live-id"}, r.PrimaryKey())
}
