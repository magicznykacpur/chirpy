package cleaner

import (
	"strings"
)

const asterisks = "****"

func CleanBodyBy(body, key string) string {
	keyTitle := strings.Join([]string{strings.ToUpper(key[:1]), key[1:]}, "")
	keyLower := strings.ToLower(key)
	keyUpper := strings.ToUpper(key)

	cleaned := strings.ReplaceAll(body, keyTitle, asterisks)
	cleaned = strings.ReplaceAll(cleaned, keyLower, asterisks)
	cleaned = strings.ReplaceAll(cleaned, keyUpper, asterisks)

	return cleaned
}