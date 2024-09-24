package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	e := echo.New()
	appState := initTaskManager()

	go appState.setupEventLoop()
	go appState.setupChannels()

	e.GET("/", func(c echo.Context) error {
		log.Info().Msg("Recived healthcheck")
		return c.String(http.StatusOK, "Service is up")
	})

	e.PUT("/close", func(c echo.Context) error {
		go func() {
			time.Sleep(1 * time.Second) // Simulate work
			os.Exit(0)
		}()
		return c.String(http.StatusOK, "Service shutting down")
	})

	e.PUT("/task", func(c echo.Context) error {
		task := new(TaskRequest)
		if err := c.Bind(task); err != nil {
			log.Error().Msgf("%v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Request couldnt be parsed")
		}
		log.Printf("Recived task: %#v", task)
		switch task.Action {
		case string(CREATE_VM):
			log.Info().Msg("CREATE_VM case")
			vm_payload, err := parseCreateVmPayload(task.Payload)
			if err != nil {
				log.Error().Msgf("CREATE_VM parse error: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			appTask := Task{
				Id:           time.Now().UnixNano(),
				Action:       CREATE_VM,
				notebook_id:  vm_payload.NotebookId,
				paragraph_id: "",
				code:         "",
			}
			err_add_queue := appState.addTask(appTask)
			if err_add_queue != nil {
				log.Error().Msgf("CREATE_VM addTask error: %v", err_add_queue)
				return echo.NewHTTPError(http.StatusInternalServerError, err_add_queue)
			}

		case string(STOP_VM):
			log.Info().Msg("STOP_VM case")
			vm_payload, err := parseStopVmPayload(task.Payload)
			if err != nil {
				log.Error().Msgf("STOP_VM parse error: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			appTask := Task{
				Id:           time.Now().UnixNano(),
				Action:       STOP_VM,
				notebook_id:  vm_payload.NotebookId,
				paragraph_id: "",
				code:         "",
			}
			err_add_queue := appState.addTask(appTask)
			if err_add_queue != nil {
				log.Error().Msgf("STOP_VM addTask error: %v", err_add_queue)
				return echo.NewHTTPError(http.StatusInternalServerError, err_add_queue)
			}

		case string(RUN_PARAGRAPH):
			log.Info().Msg("RUN_PARAGRAPH case")
			vm_payload, err := parseRunParagraph(task.Payload)
			if err != nil {
				log.Error().Msgf("RUN_PARAGRAPH parse error: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			appTask := Task{
				Id:           time.Now().UnixNano(),
				Action:       RUN_PARAGRAPH,
				notebook_id:  vm_payload.NotebookId,
				paragraph_id: vm_payload.ParagraphId,
				code:         vm_payload.Code,
			}
			err_add_queue := appState.addTask(appTask)
			if err_add_queue != nil {
				log.Error().Msgf("RUN_PARAGRAPH addTask error: %v", err_add_queue)
				return echo.NewHTTPError(http.StatusInternalServerError, err_add_queue)
			}

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
