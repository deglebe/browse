package html

import (
	"io"

	"github.com/deglebe/browse/pkg/dom"
)

type Parser struct {
	tz *Tokenizer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{tz: NewTokenizer(r)}
}

func (p *Parser) Parse() (*dom.Node, error) {
	// dummy root, can likely be removed eventually
	root := &dom.Node{
		Type:	 dom.ElementNode,
		Data:	 "root",
		Attrs:	map[string]string{},
		Children: []*dom.Node{},
	}
	stack := []*dom.Node{root}

	for {
		tok, err := p.tz.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		current := stack[len(stack)-1]
		switch tok.Type {

		case StartTagToken:
			isVoid := IsVoidElement(tok.Data)

			node := &dom.Node{
				Type:		dom.ElementNode,
				Data:		tok.Data,
				Attrs:	   tok.Attr,
				Children:	[]*dom.Node{},
				Parent:	  current,
				SelfClosing: tok.SelfClosing || isVoid,
			}
			current.Children = append(current.Children, node)
			if !node.SelfClosing {
				stack = append(stack, node)
			}

		case EndTagToken:
			if len(stack) > 1 && stack[len(stack)-1].Data == tok.Data {
				stack = stack[:len(stack)-1]
			}

		case TextToken:
			textNode := &dom.Node{
				Type:   dom.TextNode,
				Data:   tok.Data,
				Parent: current,
			}
			current.Children = append(current.Children, textNode)
		}
	}

	return root, nil
}
