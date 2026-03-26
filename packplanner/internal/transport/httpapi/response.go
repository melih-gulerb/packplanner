package httpapi

import "github.com/labstack/echo/v4"

// BaseResponse keeps all JSON API responses in a single envelope.
type BaseResponse[T any] struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    *T     `json:"data,omitempty"`
}

type healthResponse struct {
	Status string `json:"status"`
}

func respondWithSuccess[T any](c echo.Context, statusCode int, message string, data T) error {
	return c.JSON(statusCode, BaseResponse[T]{
		Message: message,
		Success: true,
		Data:    &data,
	})
}

func respondWithError(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, BaseResponse[any]{
		Message: message,
		Success: false,
	})
}
