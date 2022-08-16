package order

import (
	"fmt"

	"github.com/google/uuid"
)

// Order is an aggregate root and entity
type Order struct {
	id         ID          // the entity's unique id whose implementation is a value object
	customerID CustomerID  // another aggregate root's id, also implemented as a value object
	orderItems []OrderItem // the items this order consists of, which are a slice of value objects
	state      OrderState  // the state of the order which is a value type

	events []Event // A slice of domain events
}

// NewOrder is a factory function for creating a new Order entity which begins its life cycle
func NewOrder(customerID CustomerID, orderItems []OrderItem, state OrderState) (*Order, error) {
	orderID, err := NewID(uuid.NewString())
	if err != nil {
		return nil, err
	}

	return ReconstituteOrder(orderID, customerID, orderItems, state)
}

// ReconstituteOrder is a factory function for instantiating an Order entity in the middle of
// its life cycle.
func ReconstituteOrder(id ID, customerID CustomerID, orderItems []OrderItem, state OrderState) (*Order, error) {
	return &Order{
		id:         id,
		customerID: customerID,
		orderItems: orderItems,
		state:      state,
	}, nil
}

// ID returns a copy of the value object which means not leaking references to Order's internal id
func (o *Order) ID() ID {
	return o.id
}

// CustomerID is not leaking references either
func (o *Order) CustomerID() CustomerID {
	return o.customerID
}

// GetOrderItems returns a copy of the slice to avoid leaking a reference to Order's
// internal slice.
func (o *Order) OrderItems() []OrderItem {
	orderItems := make([]OrderItem, 0, len(o.orderItems))
	copy(orderItems, o.orderItems)

	return orderItems
}

// State is not leaking references either
func (o *Order) State() OrderState {
	return o.state
}

// Events returns a slice of domain events that should be published by the application service
func (o *Order) Events() []Event {
	domainMsgs := make([]Event, 0, len(o.events))
	copy(domainMsgs, o.events)

	return domainMsgs
}

// Submit is an aggregate command
func (o *Order) Submit() error {
	// check business rules and invariants
	if o.state != OrderStatePending {
		return fmt.Errorf("expected state pending found = %d", o.state)
	}
	o.state = OrderStateSubmitted

	// add domain event to domain event queue
	submitEvent := NewEvent("submit")
	o.events = append(o.events, submitEvent)

	return nil
}

// Cancel is an aggregate command
func (o *Order) Cancel() error {
	// check business rules and invariants
	if o.state != OrderStatePending {
		return fmt.Errorf("expected state pending found = %d", o.state)
	}

	o.state = OrderStateCancelled

	// add domain event to domain event queue
	cancelMsg := NewEvent("cancel")
	o.events = append(o.events, cancelMsg)

	return nil
}
