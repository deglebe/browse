package dom

import (
	"fmt"
	"sort"
	"strings"
)

type NodeType int

const (
	TextNode NodeType = iota
	ElementNode
)

// represents an element or text in the dom
type Node struct {
	Type		NodeType
	Data		string
	Attrs		map[string]string
	Children	[]*Node
	Parent		*Node
	SelfClosing	bool
}

func (n *Node) PrettyPrint(indent int) {
	if n.Parent == nil && n.Data == "root" {
		for _, c := range n.Children {
			c.PrettyPrint(indent)
		}
		return
	}

	pad := strings.Repeat("  ", indent)

	if n.Type == ElementNode {
		parts := make([]string, 0, len(n.Attrs))
		for k, v := range n.Attrs {
			if v == "" {
				parts = append(parts, k)
			} else {
				parts = append(parts, fmt.Sprintf("%s=%q", k, v))
			}
		}
		sort.Strings(parts)
		attrStr := ""
		if len(parts) > 0 {
			attrStr = " " + strings.Join(parts, " ")
		}

		if n.SelfClosing {
			fmt.Printf("%s<Element %s%s/>\n", pad, n.Data, attrStr)
		} else {
			fmt.Printf("%s<Element %s%s>\n", pad, n.Data, attrStr)
			for _, c := range n.Children {
				c.PrettyPrint(indent + 1)
			}
			fmt.Printf("%s</%s>\n", pad, n.Data)
		}
	} else {
		fmt.Printf("%s%q\n", pad, n.Data)
	}
}
