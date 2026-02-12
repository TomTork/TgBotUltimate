package helper

import (
	"fmt"
	"strings"
)

func ConvertValuesToSQLUpdate(fields []string, values []interface{}) string {
	__values := make([]string, len(values))
	for i, value := range values {
		if value == nil {
			__values[i] = fields[i] + "=NULL"
			continue
		}
		switch v := value.(type) {
		case int, int64, uint64, uint8:
			__values[i] = fmt.Sprintf("%s = ", fields[i]) + fmt.Sprintf("%d", v)
		case float32:
			__values[i] = fmt.Sprintf("%s = %f", fields[i], v)
		default:
			__values[i] = fmt.Sprintf("%s = ", fields[i]) + fmt.Sprintf("'%s'", v)
		}
	}
	return strings.Join(__values, ",")
}
