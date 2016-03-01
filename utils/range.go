package utils

import (
	"errors"
	"strconv"
	"strings"
)

var errInvalidRange = errors.New("Invalid range")
var errInvalidLimit = errors.New("Invalid range limit")

func GetCompleteRange(shortRange string) ([]string, error) {
	var completeRange []string
	var err error

	from_to := strings.Split(shortRange, "-")
	if len(from_to) < 1 {
		return []string{}, errInvalidRange
	}
	if len(from_to) == 1 {
		_, err := strconv.Atoi(from_to[0])
		if err != nil {
			return []string{}, errInvalidRange
		}

		return from_to, nil
	}

	from, err := strconv.Atoi(from_to[0])
	if err != nil {
		return []string{}, errInvalidLimit
	}

	to, err := strconv.Atoi(from_to[1])
	if err != nil {
		return []string{}, errInvalidLimit
	}

	delta := to - from
	completeRange = make([]string, delta+1, delta+1)
	for i := from; i <= to; i++ {
		value := strconv.Itoa(i)
		index := i - from
		completeRange[index] = value
	}

	return completeRange, nil
}
