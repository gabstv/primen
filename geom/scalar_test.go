package geom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScalarEqualsEpsilon(t *testing.T) {
	a := .1
	b := .2
	c := a + b // 0.30000000000000004
	require.True(t, ScalarEqualsEpsilon(.3, c, Epsilon))
}
