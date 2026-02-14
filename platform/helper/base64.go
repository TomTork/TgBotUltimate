package helper

import "encoding/base64"

func EncodeStringToBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func DecodeStringToBase64(value string) string {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
