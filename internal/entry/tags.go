package entry

import (
	"regexp"
	"strings"
)

// Identificar tag en base a #
var tagRe = regexp.MustCompile(`(?:^|\s)#([\p{L}\p{N}_-]+)`)

func ParseTags(input string) (string, []string) {
	var tags []string
	seen := map[string]bool{}

	for _, m := range tagRe.FindAllStringSubmatch(input, -1) {
		tag := strings.ToLower(m[1])
		if !seen[tag] {
			seen[tag] = true
			tags = append(tags, tag)
		}
	}
	clean := tagRe.ReplaceAllString(input, "")
	clean = strings.Join(strings.Fields(clean), " ")
	return clean, tags
}
