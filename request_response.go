package main

import (
	"encoding/json"
	"errors"
)

type TaskRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

type TaskResponse struct {
	Result string `json:"result"`
}

type CreateVMPayload struct {
	NotebookId string `json:"notebook_id"`
}

func parseCreateVmPayload(payload interface{}) (CreateVMPayload, error) {
	var vmPayload CreateVMPayload

	bytes, err := json.Marshal(payload)
	if err != nil {
		return CreateVMPayload{}, err
	}

	if err := json.Unmarshal(bytes, &vmPayload); err != nil {
		return CreateVMPayload{}, err
	}

	if err := validateId(vmPayload.NotebookId); err != nil {
		return CreateVMPayload{}, err
	}

	return vmPayload, nil
}

type StopVMPayload struct {
	NotebookId string `json:"notebook_id"`
}

func parseStopVmPayload(payload interface{}) (StopVMPayload, error) {
	var vmPayload StopVMPayload

	bytes, err := json.Marshal(payload)
	if err != nil {
		return StopVMPayload{}, err
	}

	if err := json.Unmarshal(bytes, &vmPayload); err != nil {
		return StopVMPayload{}, err
	}

	if err := validateId(vmPayload.NotebookId); err != nil {
		return StopVMPayload{}, err
	}

	return vmPayload, nil
}

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
