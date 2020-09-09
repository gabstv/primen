package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	txml := `<window xmlns="https://github.com/gabstv/primen" scale="2.0">
    <!-- test comment -->
    <button id="bt1" x="10" y="10" w="50" h="20" on-click="exit">Hello</button>
</window>`
	root, err := ParseXMLString(txml)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(root))
	assert.Equal(t, NodeElement, root[0].Type())
	roote := root[0].(ElementNode)
	assert.Equal(t, "window", roote.TagName())
	assert.Equal(t, "2.0", roote.Attributes()["scale"])
}
