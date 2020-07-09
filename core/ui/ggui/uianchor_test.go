package ggui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUIPosition(t *testing.T) {
	pct, locked, offs := parseUIPosition("50%-20")
	assert.Equal(t, 50., pct)
	assert.True(t, locked)
	assert.Equal(t, -20., offs)

	pct, locked, offs = parseUIPosition("50% - 20")
	assert.Equal(t, 50., pct)
	assert.True(t, locked)
	assert.Equal(t, -20., offs)
}
