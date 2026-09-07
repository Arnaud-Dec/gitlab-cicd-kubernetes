package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

type statusResponse struct {
	Status string `json:"status" example:"ok"`
}

// Live godoc
//
//	@Summary		Liveness check
//	@Description	Reports whether the process is up. Never checks the database. Served on HEALTH_PORT, not the main API port.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	statusResponse
//	@Router			/healthz [get]
func (h *HealthHandler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// Ready godoc
//
//	@Summary		Readiness check
//	@Description	Reports whether the app can serve traffic, i.e. the database is reachable. Served on HEALTH_PORT, not the main API port.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	statusResponse
//	@Failure		503	{object}	statusResponse
//	@Router			/readyz [get]
func (h *HealthHandler) Ready(c echo.Context) error {
	sqlDB, err := h.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		return c.JSON(http.StatusServiceUnavailable, statusResponse{Status: "unavailable"})
	}
	return c.JSON(http.StatusOK, statusResponse{Status: "ready"})
}
