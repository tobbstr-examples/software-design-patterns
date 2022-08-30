// Package mock contains useful mock objects for when writing tests for code that uses the
// uow package.
package mock

import (
	"context"

	"github.com/tobbstr-examples/business-logic-patterns/structural/unit-of-work/uow"
)

type Doer struct {
	DoFn func(ctx context.Context, do uow.Do) error
}

func NewDoer(stores uow.Stores) *Doer {
	return &Doer{
		DoFn: func(ctx context.Context, do uow.Do) error {
			return do(ctx, stores)
		},
	}
}

func (m *Doer) Do(ctx context.Context, do uow.Do) error {
	return m.DoFn(ctx, do)
}

type Stores struct {
	AFn func() uow.AStore
	BFn func() uow.BStore
}

func (m *Stores) A() uow.AStore {
	return m.AFn()
}

func (m *Stores) B() uow.BStore {
	return m.BFn()
}

type AStore struct {
	SaveFn func(ctx context.Context, id string) error
}

func (m *AStore) Save(ctx context.Context, id string) error {
	return m.SaveFn(ctx, id)
}

type BStore struct {
	SaveFn func(ctx context.Context, id int) error
}

func (m *BStore) Save(ctx context.Context, id int) error {
	return m.SaveFn(ctx, id)
}
