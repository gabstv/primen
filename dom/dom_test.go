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

func TestFindByID(t *testing.T) {
	root := Element("window", map[string]string{"id": "main"},
		Element("button", map[string]string{"id": "bt1"}),
		Element("button", map[string]string{"id": "bt2"}),
		Element("button", map[string]string{"id": "bt3"}),
		Element("div", map[string]string{"id": "div1"}, Element("button", map[string]string{"id": "bx"})))
	assert.NotNil(t, root.FindChildByID("bt1"))
	assert.NotNil(t, root.FindChildByID("bt3"))
	assert.NotNil(t, root.FindChildByID("div1"))
	assert.NotNil(t, root.FindChildByID("bx"))
	assert.Nil(t, root.FindChildByID("bt9"))
}
