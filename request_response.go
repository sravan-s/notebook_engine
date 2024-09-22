package main

type TaskRequest struct {
	Action  string `json:"action"`
	Payload interface{} `json:"payload"`
}

type TaskResponse struct {
	Result string `json:"result"`
}
