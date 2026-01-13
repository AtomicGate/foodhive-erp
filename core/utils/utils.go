package utils

import (
	"fmt"
)

func CheckCondition(condition bool, errorMessage string) error {
	if !condition {
		return fmt.Errorf("condition failed: %s", errorMessage)
	}
	return nil
}
