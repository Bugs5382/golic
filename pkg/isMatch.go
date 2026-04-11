package pkg

import (
	"path/filepath"
	"regexp"
	"strings"
)

func IsMatch(s string, p string) bool {
	if strings.HasPrefix(p, ".") {
		return filepath.Ext(s) == p
	}

	var res strings.Builder
	res.WriteString("^")

	i := 0
	for i < len(p) {
		switch {
		// Handle recursive directory prefix: **/
		case i+2 < len(p) && p[i:i+3] == "**/":
			res.WriteString("(?:.*/)?")
			i += 3

		// Handle recursive directory middle: /**/
		case i+3 < len(p) && p[i:i+4] == "/**/":
			res.WriteString("/(?:.*/)?")
			i += 4

		// Handle recursive directory suffix: /**
		case i+2 < len(p) && p[i:i+3] == "/**":
			res.WriteString("(?:/.*)?")
			i += 3

		// Handle single-level wildcard: *
		case p[i] == '*':
			res.WriteString("[^/]*")
			i++

		// Handle single-character wildcard: ?
		case p[i] == '?':
			res.WriteString("[^/]")
			i++

		// Handle literal characters (escape regex meta-chars like . or +)
		default:
			res.WriteString(regexp.QuoteMeta(string(p[i])))
			i++
		}
	}
	res.WriteString("$")

	matched, _ := regexp.MatchString(res.String(), s)
	return matched
}
