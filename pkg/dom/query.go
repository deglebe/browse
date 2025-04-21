package dom

import (
	"strings"
)

func (n *Node) GetElementsByTagName(name string) []*Node {
	var out []*Node
	name = strings.ToLower(name)
	var walk func(*Node)
	walk = func(cur *Node) {
		if cur.Type == ElementNode && strings.ToLower(cur.Data) == name {
			out = append(out, cur)
		}
		for _, c := range cur.Children {
			walk(c)
		}
	}
	walk(n)
	return out
}

func (n *Node) GetElementByID(id string) *Node {
	var found *Node
	var walk func(*Node)
	walk = func(cur *Node) {
		if found != nil { return }

		if cur.Type == ElementNode {
			if v, ok := cur.Attrs["id"]; ok && v == id {
				found = cur
				return
			}
		}

		for _, c := range cur.Children {
			walk(c)
			if found != nil { return }
		}
	}

	walk(n)
	return found
}

func (n *Node) QuerySelectorAll(sel string) []*Node {
	if sel == "" { return nil }

	switch {
	case sel[0] == '#':
		if node := n.GetElementByID(sel[1:]); node != nil {
			return []*Node{node}
		}
		return nil
	case sel[0] == '.':
		cls := sel[1:]
		var out []*Node
		var walk func(*Node)
		walk = func(cur *Node) {
			if cur.Type == ElementNode {
				if v, ok := cur.Attrs["class"]; ok {
					for _, part := range strings.Fields(v) {
						if part == cls {
							out = append(out, cur)
							break
						}
					}
				}
			}
			for _, c := range cur.Children {
				walk(c)
			}
		}
		walk(n)
		return out
	default:
		return n.GetElementsByTagName(sel)
	}
}

func (n *Node) QuerySelector(sel string) *Node {
	matches := n.QuerySelectorAll(sel)
	if len(matches) > 0 { return matches[0] }
	return nil
}
