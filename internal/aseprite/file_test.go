package aseprite

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getFile(t *testing.T, name string) []byte {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	require.NotNil(t, b)
	return b
}

func TestParse(t *testing.T) {
	x, err := Parse(getFile(t, "player.json"))
	assert.NoError(t, err)
	assert.Equal(t, "1", x.GetMetadata().Scale)
	assert.Equal(t, "RGBA8888", x.GetMetadata().Format)

	playerimg, _, err := image.Decode(bytes.NewReader(getFile(t, x.GetMetadata().Image)))
	assert.NoError(t, err)

	fi, ok := x.GetFrameByName("player (person) 0.aseprite")
	assert.True(t, ok)
	c1 := color.RGBA{
		A: 255,
	}
	c2 := playerimg.At(fi.Frame.X, fi.Frame.Y)
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	assert.Equal(t, r1, r2)
	assert.Equal(t, g1, g2)
	assert.Equal(t, b1, b2)
	assert.Equal(t, a1, a2)

	x, err = Parse(getFile(t, "background.json"))
	assert.NoError(t, err)
}
