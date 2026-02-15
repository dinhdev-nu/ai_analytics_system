package handlers

import (
	"net/http"
	"strconv"

	"ai-analytics/backend/internal/models"
	"ai-analytics/backend/internal/services"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

// GetRevenueForecast godoc
// @Summary Get revenue forecast
// @Description Get revenue forecast with predictions for a restaurant
// @Tags analytics
// @Param restaurant_id query string true "Restaurant ID"
// @Param months query int false "Number of months" default(12)
// @Success 200 {object} models.ForecastResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/analytics/forecast [get]
func (h *AnalyticsHandler) GetRevenueForecast(c *gin.Context) {
	restaurantID := c.Query("restaurant_id")
	if restaurantID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "restaurant_id is required",
			Code:    400,
		})
		return
	}

	monthsStr := c.DefaultQuery("months", "12")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months < 1 || months > 24 {
		months = 12
	}

	log.Infof("GET /api/v1/analytics/forecast - restaurant_id: %s, months: %d", restaurantID, months)

	forecast, err := h.service.GetRevenueForecast(c.Request.Context(), restaurantID, months)
	if err != nil {
		log.Errorf("Failed to get forecast: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetDashboard godoc
// @Summary Get dashboard data
// @Description Get complete dashboard data including summary, charts, and insights
// @Tags analytics
// @Param restaurant_id query string true "Restaurant ID"
// @Success 200 {object} models.DashboardResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/analytics/dashboard [get]
func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	restaurantID := c.Query("restaurant_id")
	if restaurantID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "bad_request",
			Message: "restaurant_id is required",
			Code:    400,
		})
		return
	}

	log.Infof("GET /api/v1/analytics/dashboard - restaurant_id: %s", restaurantID)

	dashboard, err := h.service.GetDashboard(c.Request.Context(), restaurantID)
	if err != nil {
		log.Errorf("Failed to get dashboard: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if API is running
// @Tags health
// @Success 200 {object} models.SuccessResponse
// @Router /health [get]
func (h *AnalyticsHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "API is running",
	})
}
