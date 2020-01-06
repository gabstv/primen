package spriteutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	nodes := []*RectPackerNode{
		&RectPackerNode{
			Width:  10,
			Height: 10,
			id:     1,
		},
		&RectPackerNode{
			Width:  40,
			Height: 10,
			id:     2,
		},
		&RectPackerNode{
			Width:  10,
			Height: 100,
			id:     3,
		},
		&RectPackerNode{
			Width:  60,
			Height: 30,
			id:     4,
		},
		&RectPackerNode{
			Width:  80,
			Height: 40,
			id:     5,
		},
	}
	SortNodes(nodes)
	assert.Equal(t, 3, nodes[0].id)
	assert.Equal(t, 5, nodes[1].id)
	assert.Equal(t, 4, nodes[2].id)
	assert.Equal(t, 2, nodes[3].id)
	assert.Equal(t, 1, nodes[4].id)
}

func TestBinTreePacker(t *testing.T) {
	pkr := &BinTreeRectPacker{}
	n0 := pkr.Add(300, 300)
	n1 := pkr.Add(100, 100)
	_ = pkr.Add(100, 100)
	_ = pkr.Add(100, 100)
	atls, err := pkr.Pack(context.TODO(), PackerInput{
		MarginLeft:   2,
		MarginTop:    1,
		MarginBottom: 3,
		MarginRight:  3,
		Padding:      1,
		FixedWidth:   512,
		FixedHeight:  512,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, n0.X, "n0 x")
	assert.Equal(t, 1, n0.Y, "n0 y")
	assert.Equal(t, 303, n1.X, "n1 x")
	assert.Equal(t, 1, len(atls))
}
