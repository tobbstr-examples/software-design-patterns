package order

import (
	"fmt"

	"github.com/google/uuid"
)

type ID string

// NewID is a factory function for instantiating an ID value object
func NewID(id string) (ID, error) {
	// validate input
	if err := tryParseUUID(id); err != nil {
		return "", fmt.Errorf("could not parse id: %w", err)
	}

	// return valid value object
	return ID(id), nil
}

type CustomerID string

// NewCustomerID is a factory function for instantiating a CustomerID value object
func NewCustomerID(id string) (CustomerID, error) {
	// validate input
	if err := tryParseUUID(id); err != nil {
		return "", fmt.Errorf("could not parse customer id: %w", err)
	}

	// return valid value object
	return CustomerID(id), nil
}

type OrderItem struct {
	ArticleNo string
	Quantity  int
}

// NewOrderItem is a factory function for instantiating an OrderItem value object
func NewOrderItem(articleNo string, quantity int) (OrderItem, error) {
	// validate input and enforce business rules
	if articleNo == "" || len(articleNo) > 8 {
		return OrderItem{}, fmt.Errorf("invalid articleNo")
	}

	// validate input and enforce business rules
	if quantity <= 0 {
		return OrderItem{}, fmt.Errorf("invalid quantity")
	}

	// validate input and enforce business rules
	if quantity > 100 {
		return OrderItem{}, fmt.Errorf("quantity too large")
	}

	// return valid value object
	return OrderItem{ArticleNo: articleNo, Quantity: quantity}, nil
}

const (
	OrderStateCancelled OrderState = 0
	OrderStatePending   OrderState = 1
	OrderStateSubmitted OrderState = 2
)

type OrderState int

// NewOrderState is a factory function for instantiating an OrderState value object
func NewOrderState(state int) (OrderState, error) {
	if state < 0 || state > 2 {
		return OrderStateCancelled, fmt.Errorf("invalid order state")
	}

	return OrderState(state), nil
}

func tryParseUUID(id string) error {
	_, err := uuid.Parse(id)
	return err
}
