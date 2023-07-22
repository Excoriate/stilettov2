package utils

import (
	"encoding/json"
)

// ConvertStringMapInValidMaps ConvertStringJSONMapIntoMap converts a string in JSON format into a map[string]string
func ConvertStringMapInValidMap(mapInJSONFormat string) ([]map[string]string, error) {
	if mapInJSONFormat == "" {
		return nil, nil
	}

	var mapInValidFormat []map[string]string

	err := json.Unmarshal([]byte(mapInJSONFormat), &mapInValidFormat)
	if err != nil {
		return nil, err
	}

	return mapInValidFormat, nil
}
