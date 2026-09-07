package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"todoapp/internal/version"
)

type versionResponse struct {
	Version string `json:"version" example:"2026-09-06T21:50"`
}

// Version godoc
//
//	@Summary		Build version
//	@Description	The build timestamp this binary was compiled with (see "Build with a version stamp" in the README), or "dev" for a plain go run/go build.
//	@Tags			meta
//	@Produce		json
//	@Success		200	{object}	versionResponse
//	@Router			/version [get]
func Version(c echo.Context) error {
	return c.JSON(http.StatusOK, versionResponse{Version: version.Version})
}
