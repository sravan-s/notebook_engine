package main

import (
	"encoding/json"
	"errors"
)

type RunParagraphPayload struct {
	NotebookId  string `json:"notebook_id"`
	ParagraphId string `json:"paragraph_id"`
	Code        string `json:"code"`
}

func parseRunParagraph(payload interface{}) (RunParagraphPayload, error) {
	var runPayload RunParagraphPayload

	bytes, err := json.Marshal(payload)
	if err != nil {
		return RunParagraphPayload{}, err
	}

	if err := json.Unmarshal(bytes, &runPayload); err != nil {
		return RunParagraphPayload{}, err
	}

	if err := validateId(runPayload.NotebookId); err != nil {
		return RunParagraphPayload{}, err
	}
	if err := validateId(runPayload.ParagraphId); err != nil {
		return RunParagraphPayload{}, err
	}
	if len(runPayload.Code) < 1 {
		return RunParagraphPayload{}, errors.New("code cannot be empty")
	}

	return runPayload, nil
}
