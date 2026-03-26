package pack

import "errors"

var (
	ErrInvalidOrderQuantity = errors.New("order quantity must be greater than zero")
	ErrInvalidPackSizes     = errors.New("pack sizes must contain at least one positive value")
)
