package html

import (
	"strconv"
	"strings"
)

var namedEntities = map[string]string{
	"amp":	"&",
	"lt":	"<",
	"gt":	">",
	"quot":	`"`,
	"nbsp":	"\u00A0",
}

var voidElements = map[string]bool{
	"area":		true,
	"base":		true,
	"br":		true,
	"col":		true,
	"embed":	true,
	"hr":		true,
	"img":		true,
	"input":	true,
	"link":		true,
	"meta":		true,
	"param":	true,
	"source":	true,
	"track":	true,
	"wbr":		true,
}

func Unescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '&' {
			if sem := strings.IndexByte(s[i:], ';'); sem > 0 {
				entity := s[i+1 : i+sem]
				if len(entity) > 0 && entity[0] == '#' {
					if len(entity) > 1 && (entity[1] == 'x' || entity[1] == 'X') {
						if num, err := strconv.ParseInt(entity[2:], 16, 32); err == nil {
							b.WriteRune(rune(num))
							i += sem + 1
							continue
						}
					} else if num, err := strconv.Atoi(entity[1:]); err == nil {
						b.WriteRune(rune(num))
						i += sem + 1
						continue
					}
				}
				if val, ok := namedEntities[entity]; ok {
					b.WriteString(val)
					i += sem + 1
					continue
				}
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

func IsVoidElement(tag string) bool {
	return voidElements[strings.ToLower(tag)]
}
