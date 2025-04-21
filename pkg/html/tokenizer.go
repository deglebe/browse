package html

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type TokenType int

const (
	TextToken TokenType = iota
	StartTagToken
	EndTagToken
)

// hold one raw token
type Token struct {
	Type		TokenType
	Data		string
	Attr		map[string]string
	SelfClosing	bool
}

type Tokenizer struct {
	r *bufio.Reader
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{r: bufio.NewReader(r)}
}

func (t *Tokenizer) Next() (Token, error) {
	for {
		bs, err := t.r.Peek(4)
		if err != nil {
			break
		}
		if string(bs[:4]) == "<!--" {
			if err := t.skipComment(); err != nil {
				return Token{}, err
			}
			continue
		}
		break
	}

	b, err := t.r.Peek(1)
	if err != nil {
		return Token{}, err
	}
	if b[0] == '<' {
		return t.readTag()
	}
	return t.readText()
}

func (t *Tokenizer) skipComment() error {
	for i := 0; i < 4; i++ {
		if _, err := t.r.ReadByte(); err != nil {
			return err
		}
	}
	for {
		bs, err := t.r.Peek(3)
		if err != nil {
			return err
		}
		if string(bs[:3]) == "-->" {
			t.r.ReadByte()
			t.r.ReadByte()
			t.r.ReadByte()
			return nil
		}
		if _, err := t.r.ReadByte(); err != nil {
			return err
		}
	}
}

func (t *Tokenizer) readTag() (Token, error) {
	t.r.ReadByte()

	closing := false
	if b, _ := t.r.Peek(1); b[0] == '/' {
		closing = true
		t.r.ReadByte()
	}

	name := t.readTagName()

	attrs, selfClosing := t.readAttributes()

	tok := Token{
		Data:		name,
		Attr:		attrs,
		SelfClosing: selfClosing,
	}
	if closing {
		tok.Type = EndTagToken
	} else {
		tok.Type = StartTagToken
	}
	return tok, nil
}

func (t *Tokenizer) readTagName() string {
	var sb strings.Builder
	for {
		b, err := t.r.Peek(1)
		if err != nil {
			break
		}
		c := rune(b[0])
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			sb.WriteByte(b[0])
			t.r.ReadByte()
		} else {
			break
		}
	}
	return sb.String()
}

func (t *Tokenizer) readAttributes() (map[string]string, bool) {
	attrs := make(map[string]string)
	selfClosing := false

	for {
		t.skipWhitespace()
		b, err := t.r.Peek(1)
		if err != nil {
			break
		}
		switch b[0] {
		case '/':
			selfClosing = true
			t.r.ReadByte()
		case '>':
			t.r.ReadByte()
			return attrs, selfClosing
		default:
			name := t.readAttrName()
			if name == "" {
				t.r.ReadByte()
				continue
			}
			t.skipWhitespace()
			val := ""
			if b2, _ := t.r.Peek(1); b2[0] == '=' {
				t.r.ReadByte()
				t.skipWhitespace()
				val = t.readAttrValue()
			}
			attrs[name] = val
		}
	}

	return attrs, selfClosing
}

func (t *Tokenizer) readAttrName() string {
	var sb strings.Builder
	for {
		b, err := t.r.Peek(1)
		if err != nil {
			break
		}
		c := b[0]
		if unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '-' {
			sb.WriteByte(c)
			t.r.ReadByte()
		} else {
			break
		}
	}
	return sb.String()
}

func (t *Tokenizer) readAttrValue() string {
	b, err := t.r.Peek(1)
	if err != nil {
		return ""
	}
	if b[0] == '"' || b[0] == '\'' {
		quote := b[0]
		t.r.ReadByte()
		var sb strings.Builder
		for {
			c, err := t.r.ReadByte()
			if err != nil || c == quote {
				break
			}
			sb.WriteByte(c)
		}
		return Unescape(sb.String())
	}
	var sb strings.Builder
	for {
		b2, err := t.r.Peek(1)
		if err != nil {
			break
		}
		c := b2[0]
		if unicode.IsSpace(rune(c)) || c == '/' || c == '>' {
			break
		}
		sb.WriteByte(c)
		t.r.ReadByte()
	}
	return Unescape(sb.String())
}

func (t *Tokenizer) skipWhitespace() {
	for {
		b, err := t.r.Peek(1)
		if err != nil || !unicode.IsSpace(rune(b[0])) {
			return
		}
		t.r.ReadByte()
	}
}

func (t *Tokenizer) readText() (Token, error) {
	var sb strings.Builder
	for {
		b, err := t.r.Peek(1)
		if err != nil || b[0] == '<' {
			break
		}
		sb.WriteByte(b[0])
		t.r.ReadByte()
	}
	text := strings.TrimSpace(sb.String())
	if text == "" {
		return t.Next()
	}
	return Token{Type: TextToken, Data: Unescape(text)}, nil
}
