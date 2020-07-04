package z

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRs(t *testing.T) {
	r1 := Rs()
	assert.Equal(t, 16, len(r1))
}

func BenchmarkRs1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Rs()
	}
}
