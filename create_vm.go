package main

import (
	"encoding/json"
)

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
