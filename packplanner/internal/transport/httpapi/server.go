package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"packplanner/internal/application/packapp"
)

// NewServer wires the HTTP transport to the application use cases.
func NewServer(service packapp.Service, allowedOrigins []string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	packHandler := NewPackHandler(service)
	registerSwaggerRoutes(e)
	healthHandler := func(c echo.Context) error {
		return respondWithSuccess(c, http.StatusOK, "service is healthy", healthResponse{
			Status: "ok",
		})
	}

	e.GET("/", healthHandler)
	e.GET("/health", healthHandler)

	// Versioned routes make future API changes easier to introduce safely.
	api := e.Group("/api/v1")
	api.GET("/pack-sizes", packHandler.ListPackSizes)
	api.PUT("/pack-sizes", packHandler.UpdatePackSizes)
	api.POST("/pack-plans", packHandler.CalculatePackPlan)

	return e
}
