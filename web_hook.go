package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Error().Msg("Error loading .env file")
	}

	return os.Getenv(key)
}

func sendToWebHook(webhookurl string, task Task, failed bool) {
	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling task")
		return
	}

	action := string(task.Action)
	if failed {
		action = action + "_FAILED"
	}
	payload, err := json.Marshal(map[string]string{
		"action":  action,
		"payload": string(taskJson), // Convert byte array to string
	})
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling payload")
		return
	}

	client := &http.Client{}

	// 3.
	req, err := http.NewRequest(http.MethodPut, webhookurl, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Error().Msgf("sendToWebHook: %v", err)
	}

	// 4.
	response, err := client.Do(req)
	if err != nil {
		log.Error().Msgf("sendToWebHook: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Error().Msgf("Received non-OK response: %v", response.Status)
	}
}
