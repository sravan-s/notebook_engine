package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"
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

func copyFile(from string, to string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		log.Error().Msgf("%v", err)
		return err
	}
	err = os.WriteFile(to, data, 0o644)
	if err != nil {
		log.Error().Msgf("%v", err)
		return err
	}
	return nil
}

func deleteFileIfExists(filePath string) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, so delete it
		err = os.Remove(filePath)
		if err != nil {
			log.Error().Msgf("failed to delete file: %s", err)
			return err
		}
		fmt.Println("File deleted successfully")
	} else if os.IsNotExist(err) {
		// File does not exist, do nothing
		fmt.Println("File does not exist")
	} else {
		// Some other error occurred while checking
		log.Error().Msgf("failed to check if file exists: %s", err)
		return err
	}
	return nil
}

func httpPut(url string, payload []byte) (*http.Response, error) {
	log.Info().Msgf("httpPut to %v %v", url, payload)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	response.Body.Close()

	return response, err
}
