package helper

import (
	"fmt"
	"strings"
)

func ConvertValuesToSQLCreate(values []interface{}) string {
	__values := make([]string, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case int, int64, uint64:
			__values[i] = fmt.Sprintf("%d", v)
		default:
			__values[i] = fmt.Sprintf("'%s'", v)
		}
	}
	return strings.Join(__values, ",")
}
