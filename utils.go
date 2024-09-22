package main

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	SUPPORTED_PATTERN = `^[a-zA-Z0-9]+$`
	MAX_LEN           = 24
	MIN_LEN           = 4
)

func validateId(id string) error {
	if len(id) > MAX_LEN {
		maxLenError := fmt.Sprintf("id cannot be longer than %d characters", MAX_LEN)
		return errors.New(maxLenError)
	}

	if len(id) < MIN_LEN {
		minLenError := fmt.Sprintf("id cannot be less than %d characters", MIN_LEN)
		return errors.New(minLenError)
	}

	re := regexp.MustCompile(SUPPORTED_PATTERN)
	if !re.MatchString(id) {
		return errors.New("id must contain only a..z | A..Z | 0..9")
	}
	return nil
}
