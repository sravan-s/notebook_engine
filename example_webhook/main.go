package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EventRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

func main() {
	e := echo.New()
	e.PUT("/event", func(c echo.Context) error {
		event := new(EventRequest)
		if err := c.Bind(event); err != nil {
			errrorMsg := fmt.Sprintf("Error %v", err)
			slog.Info(errrorMsg)
			return echo.NewHTTPError(http.StatusInternalServerError, "Request couldnt be parsed")
		}
		eventStr := fmt.Sprintf("Event: %v", event)
		slog.Info(eventStr)
		return c.String(http.StatusOK, "OK")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
