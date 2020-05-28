package atlaspacker

import (
	"context"
	"errors"
	"image"
	"sort"
	"sync"
	"sync/atomic"
)

var (
	ErrNoFit   = errors.New("does not fit")
	ErrNoNodes = errors.New("no nodes")
)

type PackerInput struct {
	MarginLeft   int
	MarginRight  int
	MarginTop    int
	MarginBottom int
	Padding      int
	FixedWidth   int
	FixedHeight  int
	MaxWidth     int
	MaxHeight    int
	Count        int
}

type PackerAtlas struct {
	Width  int
	Height int
	Nodes  []*RectPackerNode
}

type rectPackerTrie struct {
	node   *RectPackerNode
	right  *rectPackerTrie
	bottom *rectPackerTrie
	used   bool
	x      int
	y      int
	width  int
	height int
}

type RectPackerNode struct {
	X      int
	Y      int
	Width  int
	Height int
	id     int
}

func (n *RectPackerNode) ID() int {
	return n.id
}

func (n *RectPackerNode) R() image.Rectangle {
	return image.Rect(n.X, n.Y, n.X+n.Width, n.Y+n.Height)
}

type ImgRectPacker interface {
	AddRect(r image.Rectangle) *RectPackerNode
	AddRects(r ...image.Rectangle) []*RectPackerNode
}

type RectPacker interface {
	Add(width, height int) *RectPackerNode
	Adds(rpwh ...int) []*RectPackerNode
}

type BinTreeRectPacker struct {
	nodes  []*RectPackerNode
	previd int32
	lock   sync.Mutex
}

func (p *BinTreeRectPacker) Add(width, height int) *RectPackerNode {
	id := atomic.AddInt32(&p.previd, 1)
	newnode := &RectPackerNode{
		Width:  width,
		Height: height,
		id:     int(id),
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.nodes == nil {
		p.nodes = make([]*RectPackerNode, 0)
	}
	p.nodes = append(p.nodes, newnode)
	return newnode
}

func (p *BinTreeRectPacker) Adds(rpwh ...int) []*RectPackerNode {
	if len(rpwh) < 2 || len(rpwh)%2 != 0 {
		return nil
	}
	nodes := make([]*RectPackerNode, 0, len(rpwh)/2)
	for i := 0; i < len(rpwh); i += 2 {
		width := rpwh[i]
		height := rpwh[i+1]
		id := atomic.AddInt32(&p.previd, 1)
		newnode := &RectPackerNode{
			Width:  width,
			Height: height,
			id:     int(id),
		}
		nodes = append(nodes, newnode)
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.nodes == nil {
		p.nodes = make([]*RectPackerNode, 0, len(rpwh)/2)
	}
	p.nodes = append(p.nodes, nodes...)
	return nodes
}

func (p *BinTreeRectPacker) AddRect(r image.Rectangle) *RectPackerNode {
	return p.Add(r.Dx(), r.Dy())
}

func (p *BinTreeRectPacker) AddRects(rects ...image.Rectangle) []*RectPackerNode {
	whs := make([]int, 0, len(rects)*2)
	for _, r := range rects {
		whs = append(whs, r.Dx(), r.Dy())
	}
	return p.Adds(whs...)
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rtriefind(base *rectPackerTrie, node *RectPackerNode) *rectPackerTrie {
	if base == nil {
		return nil
	}
	if base.used {
		if x := rtriefind(base.right, node); x != nil {
			return x
		}
		return rtriefind(base.bottom, node)
	}
	if base.width >= node.Width && base.height >= node.Height {
		return base
	}
	return nil
}

func rtrieuse(t *rectPackerTrie, node *RectPackerNode, padding int) {
	t.used = true
	t.node = node
	node.X = t.x
	node.Y = t.y
	if node.Width > node.Height {
		t.bottom = &rectPackerTrie{
			x:      t.x,
			y:      t.y + node.Height + padding,
			width:  t.width,
			height: t.height - node.Height - padding,
		}
		t.right = &rectPackerTrie{
			x:      t.x + node.Width + padding,
			y:      t.y,
			width:  t.width - node.Width - padding,
			height: node.Height,
		}
	} else {
		t.bottom = &rectPackerTrie{
			x:      t.x,
			y:      t.y + node.Height + padding,
			width:  node.Width,
			height: t.height - node.Height - padding,
		}
		t.right = &rectPackerTrie{
			x:      t.x + node.Width + padding,
			y:      t.y,
			width:  t.width - node.Width - padding,
			height: t.height,
		}
	}
}

func canGrow(enabled bool, root *rectPackerTrie, node *RectPackerNode, maxw, maxh, pad int) bool {
	if !enabled || root == nil || node == nil {
		return false
	}
	return root.width+node.Width+pad <= maxw || root.height+node.Height+pad <= maxh
}

func doGrow(root *rectPackerTrie, node *RectPackerNode, padding int) (newroot *rectPackerTrie) {
	ww := root.width + padding + node.Width
	hh := root.height + padding + node.Height

	if hh > ww {
		// grow to the right
		newroot = &rectPackerTrie{
			used:   true,
			width:  ww,
			height: root.height,
			bottom: root,
			right: &rectPackerTrie{
				x:      root.width + padding,
				width:  node.Width,
				height: root.height,
			},
		}
	} else {
		// grow to the bottom
		newroot = &rectPackerTrie{
			used:   true,
			width:  root.width,
			height: hh,
			right:  root,
			bottom: &rectPackerTrie{
				y:      root.height + padding,
				width:  root.width,
				height: node.Height,
			},
		}
	}
	return
}

func fetchAtlas(root *rectPackerTrie) PackerAtlas {
	output := PackerAtlas{}
	if root == nil {
		return output
	}
	output.Width = root.width
	output.Height = root.height
	output.Nodes = findSprites(root)
	return output
}

func findSprites(trie *rectPackerTrie) []*RectPackerNode {
	v := make([]*RectPackerNode, 0)
	if trie.node != nil {
		v = append(v, trie.node)
	}
	if trie.right != nil {
		v = append(v, findSprites(trie.right)...)
	}
	if trie.bottom != nil {
		v = append(v, findSprites(trie.bottom)...)
	}
	return v
}

func (p *BinTreeRectPacker) Pack(ctx context.Context, input PackerInput) ([]PackerAtlas, error) {
	p.lock.Lock()
	clone := make([]*RectPackerNode, len(p.nodes))
	copy(clone, p.nodes)
	p.lock.Unlock()
	if len(clone) < 1 {
		return []PackerAtlas{}, ErrNoNodes
	}
	SortNodes(clone)
	if input.Count <= 0 {
		input.Count = 1
	}
	if input.MaxWidth <= 0 {
		input.MaxWidth = 4096
	}
	if input.MaxHeight <= 0 {
		input.MaxHeight = 4096
	}
	//
	grow := true
	var root *rectPackerTrie
	newRoot := func(input PackerInput) *rectPackerTrie {
		v := &rectPackerTrie{
			x: input.MarginLeft,
			y: input.MarginTop,
		}
		if input.FixedWidth > 0 && input.FixedHeight > 0 {
			v.width = input.FixedWidth - input.MarginLeft - input.MarginRight
			v.height = input.FixedHeight - input.MarginTop - input.MarginBottom
		} else {
			v.width = clone[0].Width
			v.height = clone[0].Height
		}
		return v
	}
	root = newRoot(input)
	if input.FixedWidth > 0 && input.FixedHeight > 0 {
		grow = false
	}
	output := make([]PackerAtlas, 0, 1)
	t0 := root
	for _, node := range clone {
		t0 = rtriefind(t0, node)
		if t0 != nil {
			rtrieuse(t0, node, input.Padding)
		} else {
			// none was found!
			if canGrow(grow, root, node, input.MaxWidth, input.MaxHeight, input.Padding) {
				root = doGrow(root, node, input.Padding)
				t0 = rtriefind(root, node)
				if t0 != nil {
					rtrieuse(t0, node, input.Padding)
				}
			} else {
				// TODO: swap to a new atlas and replace root
				output = append(output, fetchAtlas(root))
				root = newRoot(input)
				t0 = root
			}
			// TODO: deter
		}
	}
	output = append(output, fetchAtlas(root))
	//
	return output, nil
}

func SortNodes(nodes []*RectPackerNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return maxint(nodes[i].Width, nodes[i].Height) > maxint(nodes[j].Width, nodes[j].Height)
	})
}
