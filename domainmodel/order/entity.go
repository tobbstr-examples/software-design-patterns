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

	messages []Message // A slice of domain messages
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

// Messages returns a slice of domain messages that should be published by the application service
func (o *Order) Messages() []Message {
	domainMsgs := make([]Message, 0, len(o.messages))
	copy(domainMsgs, o.messages)

	return domainMsgs
}

// Submit is an aggregate command
func (o *Order) Submit() error {
	// check business rules and invariants
	if o.state != OrderStatePending {
		return fmt.Errorf("expected state pending found = %d", o.state)
	}
	o.state = OrderStateSubmitted

	// add domain event to domain message queue
	submitMsg := NewMessage("submit")
	o.messages = append(o.messages, submitMsg)

	return nil
}

// Cancel is an aggregate command
func (o *Order) Cancel() error {
	// check business rules and invariants
	if o.state != OrderStatePending {
		return fmt.Errorf("expected state pending found = %d", o.state)
	}

	o.state = OrderStateCancelled

	// add domain event to domain message queue
	cancelMsg := NewMessage("cancel")
	o.messages = append(o.messages, cancelMsg)

	return nil
}
