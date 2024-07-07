package util

import "strings"

func MapToQuerystring(query map[string]string) string {
	var result []string
	b := strings.Builder{}

	for k, v := range query {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)

		result = append(result, b.String())
		b.Reset()
	}

	return strings.Join(result, "&")
}
