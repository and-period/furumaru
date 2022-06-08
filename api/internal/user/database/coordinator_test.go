package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinator(t *testing.T) {
	assert.NotNil(t, NewCoordinator(nil))
}
