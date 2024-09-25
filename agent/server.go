package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
)

type (
	Paragraph struct {
		Code string `json:"code" validate:"required"`
	}
	ParagraphResult struct {
		Data string `json:"data" validate:"required"`
		Type string `json:"type" validate:"required"`
	}
)

//go:embed template.mjs
var js_program_template string

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})
	e.PUT("/:notebook_id/:paragraph_id/exe", func(c echo.Context) (err error) {
		notebook_id := c.Param("notebook_id")
		paragraph_id := c.Param("paragraph_id")

		u := new(Paragraph)
		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		fmt.Println(u)

		js_program := fmt.Sprintf(js_program_template, notebook_id, paragraph_id, u.Code)
		tmpFile, err := os.CreateTemp("", "*.mjs")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		defer os.Remove(tmpFile.Name())
		if _, err := tmpFile.Write([]byte(js_program)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err := tmpFile.Close(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		cmd := exec.Command("node", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error running Node.js: %v, Output: %s", err, output))
		}

		// Return the result of the executed JavaScript code
		return c.String(http.StatusOK, string(output))
	})
	e.Logger.Fatal(e.Start(":1323"))
}
