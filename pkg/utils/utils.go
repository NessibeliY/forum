package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ParsePositiveIntID(id string) (int, error) {
	matched, err := regexp.MatchString(`^[1-9]\d*$`, id)
	if err != nil || !matched {
		return 0, fmt.Errorf("match string: %s", id)
	}

	num, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("convert to int: %s", id)
	}

	return num, nil
}
