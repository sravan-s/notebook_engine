package main

import (
	"encoding/json"
)

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
