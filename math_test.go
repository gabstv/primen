package tau

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyOrigin(t *testing.T) {
	assert.Equal(t, 0.0, applyOrigin(10, 0))
	assert.Equal(t, -5., applyOrigin(10, 0.5))
	assert.Equal(t, -10.0, applyOrigin(10, 1))
}
