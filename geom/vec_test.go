package geom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVecDot(t *testing.T) {
	require.Equal(t, float64(66), Vec{-6, 8}.Dot(Vec{5, 12}))
}

func TestVecMag(t *testing.T) {
	require.Equal(t, float64(5), Vec{3, 4}.Magnitude())
}
