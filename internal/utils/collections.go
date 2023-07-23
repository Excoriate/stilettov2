package utils

import "fmt"

func FindInSlice(toFind string, data []string) error {
	for _, item := range data {
		if item == toFind {
			return nil
		}
	}

	return fmt.Errorf("could not find %s in slice", toFind)
}
