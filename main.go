package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		log.Info().Msg("Recived healthcheck")
		return c.String(http.StatusOK, "Service is up")
	})

	e.PUT("/task", func(c echo.Context) error {
		task := new(TaskRequest)
		if err := c.Bind(task); err != nil {
			log.Error().Msgf("%v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Request couldnt be parsed")
		}
		log.Printf("Recived task: %#v", task)
		switch task.Action {
		case "CREATE_VM":
			log.Info().Msg("CREATE_VM case")
			vm_payload, err := parseCreateVmPayload(task.Payload)
			if err != nil {
				log.Error().Msgf("CREATE_VM parse error: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			log.Info().Msgf("%v", vm_payload)

		case "STOP_VM":
			log.Info().Msg("STOP_VM case")
			vm_payload, err := parseStopVmPayload(task.Payload)
			if err != nil {
				log.Error().Msgf("STOP_VM parse error: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			log.Info().Msgf("%v", vm_payload)

		case "RUN_PARAGRAPH":
			log.Info().Msg("RUN_PARAGRAPH case")
			vm_payload, err := parseRunParagraph(task.Payload)
			if err != nil {
				log.Error().Msgf("RUN_PARAGRAPH parse error: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			log.Info().Msgf("%v", vm_payload)

		default:
			log.Error().Msgf("unknown action: %s", task.Action)
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"Unknown Action; Use CREATE_VM|STOP_VM|RUN_PARAGRAPH")
		}
		task_queue_success := &TaskResponse{
			Result: "TASK_QUEUED",
		}
		return c.JSON(http.StatusOK, task_queue_success)
	})

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			log.Trace().Msgf("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))
	e.Logger.Fatal(e.Start(":1323"))
}
