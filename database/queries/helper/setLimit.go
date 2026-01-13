package helper

import "strings"

func SetLimit(query string) string {
	if !strings.Contains(strings.ToLower(query), "limit") {
		if query[len(query)-1] == ';' {
			query = query[:len(query)-1]
			query += " LIMIT 3;"
		}
	}
	return query
}
