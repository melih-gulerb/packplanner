package packapp

import (
	"context"

	"packplanner/internal/domain/pack"
)

// Service exposes the application use cases without leaking transport or storage details.
type Service interface {
	ListPackSizes(ctx context.Context) ([]int, error)
	UpdatePackSizes(ctx context.Context, packSizes []int) ([]int, error)
	CalculateShipment(ctx context.Context, orderQuantity int) (pack.ShipmentPlan, error)
}

type service struct {
	repository pack.Repository
	planner    pack.Planner
}

func NewService(repository pack.Repository, planner pack.Planner) Service {
	return service{
		repository: repository,
		planner:    planner,
	}
}

func (s service) ListPackSizes(ctx context.Context) ([]int, error) {
	return s.repository.List(ctx)
}

func (s service) UpdatePackSizes(ctx context.Context, packSizes []int) ([]int, error) {
	// Persist a canonical version of the list so every caller works with the same data shape.
	normalized, err := pack.NormalizePackSizes(packSizes)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Replace(ctx, normalized); err != nil {
		return nil, err
	}

	return normalized, nil
}

func (s service) CalculateShipment(ctx context.Context, orderQuantity int) (pack.ShipmentPlan, error) {
	// The planner stays stateless; the active pack sizes come from the configured repository.
	packSizes, err := s.repository.List(ctx)
	if err != nil {
		return pack.ShipmentPlan{}, err
	}

	return s.planner.Calculate(orderQuantity, packSizes)
}
