package utils

import (
	"github.com/satori/go.uuid"
	"strings"
)

func NormaliseStringUpper(target string) string {
	return strings.TrimSpace(strings.ToUpper(target))
}

func NormaliseStringLower(target string) string {
	return strings.TrimSpace(strings.ToLower(target))
}

func NormaliseNoSpaces(target string) string {
	return strings.TrimSpace(target)
}

func RemoveDoubleQuotes(target string) string {
	return strings.Trim(target, "\"")
}

func IsKeyInMapOptional(key string, optionalKeys []string) bool {
	for _, optionalKey := range optionalKeys {
		if key == optionalKey {
			return true
		}
	}
	return false
}

func MapIsNulOrEmpty(target map[string]string) bool {
	return target == nil || len(target) == 0
}

func GetUUID() string {
	id := uuid.NewV4()
	return id.String()
}

func MergeSlices(slices ...[]string) []string {
	var merged []string
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}
