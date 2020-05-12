package aseprite

import (
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
	assert.Equal(t, FileTypeMap, x.Type())
	assert.Equal(t, "1", x.GetMetadata().Scale)
	assert.Equal(t, "RGBA8888", x.GetMetadata().Format)

	x, err = Parse(getFile(t, "background.json"))
	assert.NoError(t, err)
	assert.Equal(t, FileTypeSlice, x.Type())
}
