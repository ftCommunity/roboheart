package marshallers

import "strings"

func extractString(raw []byte) string {
	return strings.TrimSuffix(strings.TrimPrefix(string(raw), "\""), "\"")
}

func MakeByteString(raw string) []byte {
	return []byte("\"" + raw + "\"")
}
