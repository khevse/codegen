package stringspkg

import (
	"strings"
	"unicode"
)

func ToSnakeCase(val string) string {
	buf := strings.Builder{}
	for _, r := range val {
		if unicode.IsUpper(r) {
			if buf.Len() > 0 {
				buf.WriteByte('_')
			}

			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(unicode.ToLower(r))
		}
	}

	return buf.String()
}
