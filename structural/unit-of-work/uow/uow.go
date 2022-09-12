/*
Package uow is a library that holds the interfaces used by all parts of the application that uses
transactions. Its clients normally are application services and controllers.
*/
package uow

import "context"

type (
	// AStore describes the behaviour of a store of type A
	AStore interface {
		Save(ctx context.Context, id string) error
	}

	// BStore describes the behaviour of a store of type B
	BStore interface {
		Save(ctx context.Context, id int) error
	}

	// Stores is a wrapper object that enables access to all of the stores used in the application.
	// This means that whenever a new store is created its interface has to be defined above and then
	// added to this interface.
	Stores interface {
		A() AStore
		B() BStore
	}
)

// Do is the function that contains the transactional logic. It's instantiated and defined
// inside the function it's used such as inside a controller method or application service
// method.
// See app_svc.go file for more information.
type Do func(ctx context.Context, stores Stores) error

// Doer is the dependency any object needs to have in order to perform transactional logic.
// See app_svc.go file for more information.
type Doer interface {
	// Atomically executes do atomically.
	Atomically(ctx context.Context, do Do) error
}
