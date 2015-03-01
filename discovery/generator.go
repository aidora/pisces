package discovery

import (
	"fmt"
	"regexp"
)

func incString(s string) string {
	if s == "" {
		return "1"
	}
	i := len(s) - 1
	if s[i] == '9' {
		return incString(s[:len(s)-1]) + "0"

	} else if s[i] == 'Z' {
		return incString(s[:len(s)-1]) + "A"

	} else if s[i] == 'z' {
		return incString(s[:len(s)-1]) + "a"

	} else {
		return s[:len(s)-1] + string(s[i]+1)
	}
}

//
// IP and Hostname generator
//
func Generate(pattern string) []string {
	re, _ := regexp.Compile(`\[(.+):(.+)\]`)
	submatch := re.FindStringSubmatch(pattern)
	if submatch == nil {
		return []string{pattern}
	}

	from := submatch[1]
	to := submatch[2]

	template := re.ReplaceAllString(pattern, "%s")

	result := make([]string, 0)
	for val := from; ; val = incString(val) {
		entry := fmt.Sprintf(template, val)
		result = append(result, entry)
		if val == to {
			break
		}
	}

	return result
}
