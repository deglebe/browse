package layout

import (
	"fmt"
	"strings"

	"github.com/deglebe/browse/pkg/dom"

	"golang.org/x/image/font"
)

type RenderOp struct {
	Text	string
	X, Y	int
}

type Context struct {
	Face		font.Face
	LineHeight	int
	MaxWidth	int
	ListIndent	int
}

// get total content height
func Render(root *dom.Node, ctx *Context) ([]RenderOp, int) {
	var ops []RenderOp
	y := 0
	walk(root, 0, &y, &ops, ctx)
	return ops, y
}

func walk(n *dom.Node, indent int, y *int, ops *[]RenderOp, ctx *Context) {
	switch n.Type {
	case dom.ElementNode:
		tag := strings.ToLower(n.Data)
		switch tag {
		case "p", "div", "h1", "h2", "h3":
			*y += ctx.LineHeight
			emitText(flattenText(n), indent, y, ops, ctx)
			*y += ctx.LineHeight
		case "ul":
			*y += ctx.LineHeight
			for _, li := range n.GetElementsByTagName("li") {
				renderListItem(li, indent, y, ops, ctx, false, 0)
			}
			*y += ctx.LineHeight
		case "ol":
			*y += ctx.LineHeight
			idx := 1
			for _, li := range n.GetElementsByTagName("li") {
				renderListItem(li, indent, y, ops, ctx, true, idx)
				idx++
			}
			*y += ctx.LineHeight
		default:
			for _, c := range n.Children {
				walk(c, indent, y, ops, ctx)
			}
		}
	case dom.TextNode:
		emitText(n.Data, indent, y, ops, ctx)
	}
}

func flattenText(n *dom.Node) string {
	var sb strings.Builder
	var f func(*dom.Node)
	f = func(x *dom.Node) {
		if x.Type == dom.TextNode {
			sb.WriteString(x.Data)
		}
		for _, c := range x.Children { f(c) }
	}
	f(n)
	return sb.String()
}

func renderListItem(n *dom.Node, indent int, y *int, ops *[]RenderOp, ctx *Context, numbered bool, idx int) {
	var prefix string
	if numbered {
		prefix = fmt.Sprintf("%d. ", idx)
	} else {
		prefix = "- "
		// prefix = "• " // in case we want fanciness :shrug: 
	}
	level := countAncestorLists(n)
	newIndent := indent + ctx.ListIndent * (level - 1)
	text := prefix + flattenText(n)
	emitText(text, newIndent, y, ops, ctx)
}

func emitText(s string, indent int, y *int, ops *[]RenderOp, ctx *Context) {
	// guard against 0 width
	if ctx.MaxWidth <= 0 {
		*ops = append(*ops, RenderOp{Text: s, X: indent, Y: *y})
		*y += ctx.LineHeight
		return
	}

	words := strings.Fields(s)
	if len(words) == 0 {
		*y += ctx.LineHeight
		return
	}
	drawer := &font.Drawer{Face: ctx.Face}
	line := ""
	for _, w := range words {
		test := w
		if line != "" { test = line + " " + w }
		adv := drawer.MeasureString(test).Round()
		if adv + indent > ctx.MaxWidth && line != "" {
			*ops = append(*ops, RenderOp{Text: line, X: indent, Y: *y})
			*y += ctx.LineHeight
			line = w
		} else { line = test }
	}

	if line != "" {
		*ops = append(*ops, RenderOp{Text: line, X: indent, Y: *y})
		*y += ctx.LineHeight
	}
}

func countAncestorLists(n *dom.Node) int {
	cnt := 0
	for p := n.Parent; p != nil; p = p.Parent {
		t := strings.ToLower(p.Data)
		if t == "ul" || t == "ol" {
			cnt++
		}
	}
	return cnt
}
