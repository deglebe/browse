package layout

import (
	"fmt"
	"strings"

	"github.com/deglebe/browse/pkg/dom"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// RenderOp is one draw instruction: draw Text at (X,Y).
type RenderOp struct {
	Text	 string
	X, Y	 int
}

type Context struct {
	Faces			map[string]font.Face
	DefaultFace		font.Face
	LineHeights		map[string]int
	DefaultLineHeight	int
	MaxWidth		int
	ListIndent		int
}

func NewContext(maxWidth int) (*Context, error) {
	// 1) Parse the raw TTF bytes
	ttFont, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("parsing gofont: %w", err)
	}

	// 2) Define the size (pt) you want per tag
	sizes := map[string]float64{
		"h1": 32,
		"h2": 24,
		"h3": 18,
		"p":  14,
		"li": 14,
	}

	faces := make(map[string]font.Face, len(sizes))
	lineHeights := make(map[string]int, len(sizes))

	// 3) Build a font.Face for each tag via opentype.NewFace
	for tag, size := range sizes {
		face, err := opentype.NewFace(ttFont, &opentype.FaceOptions{
			Size:	size,
			DPI:	 72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			return nil, fmt.Errorf("creating face for %s: %w", tag, err)
		}
		faces[tag] = face
		lineHeights[tag] = face.Metrics().Height.Round()
	}

	defaultFH := lineHeights["p"]
	return &Context{
		Faces:			 faces,
		DefaultFace:	   faces["p"],
		LineHeights:	   lineHeights,
		DefaultLineHeight: defaultFH,
		MaxWidth:		  maxWidth,
		ListIndent:		20,
	}, nil
}

func Render(root *dom.Node, ctx *Context) ([]RenderOp, int) {
	var ops []RenderOp
	y := 0
	walk(root, 0, &y, &ops, ctx)
	return ops, y
}

func walk(n *dom.Node, indent int, y *int, ops *[]RenderOp, ctx *Context) {
	if n.Type == dom.ElementNode {
		tag := strings.ToLower(n.Data)
		switch tag {
		case "h1", "h2", "h3", "p":
			*y += ctx.LineHeights[tag]
			emitText(tag, flatten(n), indent, y, ops, ctx)
			*y += ctx.LineHeights[tag]

		case "ul", "ol":
			*y += ctx.DefaultLineHeight
			idx := 1
			for _, c := range n.Children {
				if strings.ToLower(c.Data) == "li" {
					renderListItem(c, indent, y, ops, ctx, tag == "ol", idx)
					if tag == "ol" {
						idx++
					}
				}
			}
			*y += ctx.DefaultLineHeight

		default:
			for _, c := range n.Children {
				walk(c, indent, y, ops, ctx)
			}
		}
	} else if n.Type == dom.TextNode {
		emitText("p", n.Data, indent, y, ops, ctx)
	}
}

func renderListItem(n *dom.Node, indentPtr int, y *int, ops *[]RenderOp, ctx *Context, numbered bool, idx int) {
	prefix := "- "
	// prefix := "• " // in case we want fancy :shrug:
	if numbered {
		prefix = fmt.Sprintf("%d. ", idx)
	}
	level := countLists(n)
	newIndent := indentPtr + ctx.ListIndent*(level-1)
	emitText("li", prefix+flatten(n), newIndent, y, ops, ctx)
}

func flatten(n *dom.Node) string {
	var sb strings.Builder
	var f func(*dom.Node)
	f = func(x *dom.Node) {
		if x.Type == dom.TextNode {
			sb.WriteString(x.Data)
		}
		for _, c := range x.Children {
			f(c)
		}
	}
	f(n)
	return sb.String()
}

func emitText(tag, text string, indent int, y *int, ops *[]RenderOp, ctx *Context) {
	face, ok := ctx.Faces[tag]
	if !ok {
		face = ctx.DefaultFace
	}
	lineH, ok := ctx.LineHeights[tag]
	if !ok {
		lineH = ctx.DefaultLineHeight
	}

	if ctx.MaxWidth <= 0 {
		*ops = append(*ops, RenderOp{Text: text, X: indent, Y: *y})
		*y += lineH
		return
	}

	drawer := font.Drawer{Face: face}
	words := strings.Fields(text)
	if len(words) == 0 {
		*y += lineH
		return
	}
	line := ""
	for _, w := range words {
		candidate := w
		if line != "" {
			candidate = line + " " + w
		}
		adv := drawer.MeasureString(candidate).Ceil()
		if adv+indent > ctx.MaxWidth && line != "" {
			*ops = append(*ops, RenderOp{Text: line, X: indent, Y: *y})
			*y += lineH
			line = w
		} else {
			line = candidate
		}
	}
	if line != "" {
		*ops = append(*ops, RenderOp{Text: line, X: indent, Y: *y})
		*y += lineH
	}
}

func countLists(n *dom.Node) int {
	lvl := 0
	for p := n.Parent; p != nil; p = p.Parent {
		t := strings.ToLower(p.Data)
		if t == "ul" || t == "ol" {
			lvl++
		}
	}
	return lvl
}
