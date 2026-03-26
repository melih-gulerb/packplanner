package httpapi

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed swagger/index.html swagger/openapi.json
var swaggerAssets embed.FS

func registerSwaggerRoutes(e *echo.Echo) {
	e.GET("/swagger", serveSwaggerUI)
	e.GET("/swagger/", serveSwaggerUI)
	e.GET("/swagger/openapi.json", serveOpenAPISpec)
}

func serveSwaggerUI(c echo.Context) error {
	page, err := swaggerAssets.ReadFile("swagger/index.html")
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "failed to load swagger ui")
	}

	return c.HTMLBlob(http.StatusOK, page)
}

func serveOpenAPISpec(c echo.Context) error {
	spec, err := swaggerAssets.ReadFile("swagger/openapi.json")
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "failed to load openapi spec")
	}

	return c.Blob(http.StatusOK, "application/json; charset=utf-8", spec)
}
