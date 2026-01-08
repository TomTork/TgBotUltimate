package helper

import (
	"fmt"
	"strings"
)

func ConvertValuesToSQLUpdate(fields []string, values []interface{}) string {
	__values := make([]string, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case int, int64, uint64:
			__values[i] = fmt.Sprintf("%s = ", fields[i]) + fmt.Sprintf("%d", v)
		default:
			__values[i] = fmt.Sprintf("%s = ", fields[i]) + fmt.Sprintf("'%s'", v)
		}
	}
	return strings.Join(fields, "\n")
}
