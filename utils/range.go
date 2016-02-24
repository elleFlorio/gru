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
	if len(from_to) != 2 {
		return []string{}, errInvalidRange
	}

	from, err := strconv.Atoi(from_to[0])
	to, err := strconv.Atoi(from_to[1])

	if err != nil {
		return []string{}, errInvalidLimit
	}

	completeRange = make([]string, to-from, to-from)
	for i := from; i <= to; i++ {
		value := strconv.Itoa(i)
		index := i - from
		completeRange[index] = value
	}

	return completeRange, nil
}
