package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyOrigin(t *testing.T) {
	assert.Equal(t, 0.0, ApplyOrigin(10, 0))
	assert.Equal(t, -5., ApplyOrigin(10, 0.5))
	assert.Equal(t, -10.0, ApplyOrigin(10, 1))
}
