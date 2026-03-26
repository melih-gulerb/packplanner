package httpapi

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"packplanner/internal/application/packapp"
	"packplanner/internal/domain/pack"
)

// PackHandler translates HTTP requests into application use case calls.
type PackHandler struct {
	service packapp.Service
}

type updatePackSizesRequest struct {
	PackSizes []int `json:"pack_sizes"`
}

type calculatePackPlanRequest struct {
	OrderQuantity int `json:"order_quantity"`
}

type packSizesResponse struct {
	PackSizes []int `json:"pack_sizes"`
}

type calculatePackPlanResponse struct {
	OrderQuantity int                    `json:"order_quantity"`
	TotalItems    int                    `json:"total_items"`
	TotalPacks    int                    `json:"total_packs"`
	Packs         []calculatePackDetails `json:"packs"`
}

type calculatePackDetails struct {
	PackSize int `json:"pack_size"`
	Quantity int `json:"quantity"`
}

// NewPackHandler creates the HTTP adapter for pack-related endpoints.
func NewPackHandler(service packapp.Service) PackHandler {
	return PackHandler{service: service}
}

func (h PackHandler) ListPackSizes(c echo.Context) error {
	packSizes, err := h.service.ListPackSizes(c.Request().Context())
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	return respondWithSuccess(c, http.StatusOK, "pack sizes retrieved successfully", packSizesResponse{
		PackSizes: packSizes,
	})
}

func (h PackHandler) UpdatePackSizes(c echo.Context) error {
	var request updatePackSizesRequest
	if err := c.Bind(&request); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	packSizes, err := h.service.UpdatePackSizes(c.Request().Context(), request.PackSizes)
	if err != nil {
		return respondWithDomainError(c, err)
	}

	return respondWithSuccess(c, http.StatusOK, "pack sizes updated successfully", packSizesResponse{
		PackSizes: packSizes,
	})
}

func (h PackHandler) CalculatePackPlan(c echo.Context) error {
	var request calculatePackPlanRequest
	if err := c.Bind(&request); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	plan, err := h.service.CalculateShipment(c.Request().Context(), request.OrderQuantity)
	if err != nil {
		return respondWithDomainError(c, err)
	}

	return respondWithSuccess(c, http.StatusOK, "pack plan calculated successfully", toCalculatePackPlanResponse(plan))
}

func toCalculatePackPlanResponse(plan pack.ShipmentPlan) calculatePackPlanResponse {
	// Keep the HTTP response shape separate from the domain model.
	packs := make([]calculatePackDetails, 0, len(plan.Packs))
	for _, shipmentPack := range plan.Packs {
		packs = append(packs, calculatePackDetails{
			PackSize: shipmentPack.PackSize,
			Quantity: shipmentPack.Quantity,
		})
	}

	return calculatePackPlanResponse{
		OrderQuantity: plan.OrderQuantity,
		TotalItems:    plan.TotalItems,
		TotalPacks:    plan.TotalPacks,
		Packs:         packs,
	}
}

func respondWithDomainError(c echo.Context, err error) error {
	// Validation failures map to 400 responses; everything else is treated as an internal error.
	switch {
	case errors.Is(err, pack.ErrInvalidOrderQuantity), errors.Is(err, pack.ErrInvalidPackSizes):
		return respondWithError(c, http.StatusBadRequest, err.Error())
	default:
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}
}
