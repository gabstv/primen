package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewElements(t *testing.T) {
	root := Element("window", map[string]string{
		"id":    "main",
		"class": "default scaled2x",
	}, Element("button", nil, Text("PLAY")))
	assert.Equal(t, "window", root.TagName())
	assert.Equal(t, "main", root.ID())
	assert.Equal(t, "scaled2x", root.Classes()[1])
	assert.Equal(t, "button", root.Children()[0].(ElementNode).TagName())
	assert.Equal(t, NodeText, root.Children()[0].(ElementNode).Children()[0].Type())
	assert.Equal(t, "PLAY", root.Children()[0].(ElementNode).Children()[0].(TextNode).Text())
}
