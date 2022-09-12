package store

import (
	"context"
	"database/sql"

	"github.com/tobbstr-examples/business-logic-patterns/structural/unit-of-work/uow"
)

type aStore struct {
	tx *sql.Tx
}

func (s *aStore) Save(ctx context.Context, id string) error {
	// Implementation omitted
	return nil
}

type bStore struct {
	tx *sql.Tx
}

func (s *bStore) Save(ctx context.Context, id int) error {
	// Implementation omitted
	return nil
}

type unitOfWorkStores struct {
	aStore *aStore
	bStore *bStore
}

func (s *unitOfWorkStores) A() uow.AStore {
	return s.aStore
}

func (s *unitOfWorkStores) B() uow.BStore {
	return s.bStore
}

type UnitOfWorkDoer struct {
	db sql.DB
}

func NewUoWDoer(db sql.DB) *UnitOfWorkDoer {
	return &UnitOfWorkDoer{db: db}
}

func (w *UnitOfWorkDoer) Atomically(ctx context.Context, do uow.Do) error {
	tx, _ := w.db.Begin()

	unitOfWorkStores := &unitOfWorkStores{
		aStore: &aStore{tx: tx},
		bStore: &bStore{tx: tx},
	}

	if err := do(ctx, unitOfWorkStores); err != nil {
		// The error handling could be improved here
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
