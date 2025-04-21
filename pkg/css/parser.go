package css

import (
	"io"
	"strings"

	"github.com/deglebe/browse/pkg/dom"
)

// stub, read css text and build ss
func ParseStyleSheet(r io.Reader) (*Stylesheet, error) {
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	// TODO: lex buf.String()

	return &Stylesheet{}, nil
}

// stub, walk dom and apply stylesheet
func ApplyStyles(el *dom.Node, sheet *Stylesheet) {
	for _, child := range el.Children {
		ApplyStyles(child, sheet)
	}
}
