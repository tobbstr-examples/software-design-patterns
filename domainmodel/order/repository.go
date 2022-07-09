package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type (
	Tx interface {
		Commit() error
		Exec(query string, args ...any) (sql.Result, error)
		Query(query string, args ...any) (*sql.Rows, error)
		Rollback() error
	}

	Repository interface {
		FindByID(ctx context.Context, id string) (*Order, error)
		Upsert(ctx context.Context, order *Order) error
	}

	connPool interface {
		BeginTransaction() *sql.Tx
	}
)

// repository implements the Repository interface.
// NOTE! Error handling is omitted to simplify the examples.
type repository struct {
	tx Tx
}

// NewRepository returns an interface to make the AppService testable.
func NewRepository(tx Tx) Repository {
	return &repository{tx: tx}
}

func (r *repository) FindByID(ctx context.Context, id string) (*Order, error) {
	// run fictional query against database
	rows, _ := r.tx.Query("... query ...")
	defer rows.Close()

	var (
		rowID             string
		rowCustID         string
		rowJsonOrderItems []byte
		rowState          int
	)

	type (
		jsonOrderItem struct {
			ArticleNo string `json:"articleNo"`
			Quantity  int    `json:"quantity"`
		}

		jsonOrderItems struct {
			Data []jsonOrderItem `json:"data"`
		}
	)

	// scan rows to map values returned from the database
	if !rows.Next() {
		return nil, fmt.Errorf("could not find order by id = %s", id)
	}
	rows.Scan(&rowID, &rowCustID, &rowJsonOrderItems, &rowState)

	var rowOrderItems jsonOrderItems
	json.Unmarshal(rowJsonOrderItems, &rowOrderItems)

	// delegate creation/validation/business rules to factories/constructors
	ID, _ := NewID(rowID)
	customerID, _ := NewCustomerID(rowCustID)
	var orderItems []OrderItem
	for _, rowOrderItem := range rowOrderItems.Data {
		orderItem, _ := NewOrderItem(rowOrderItem.ArticleNo, rowOrderItem.Quantity)
		orderItems = append(orderItems, orderItem)
	}
	state, _ := NewOrderState(rowState)
	order, _ := ReconstituteOrder(ID, customerID, orderItems, state)

	return order, nil
}

func (r *repository) Upsert(ctx context.Context, order *Order) error {
	// run fictional query against database
	_, err := r.tx.Exec("... query ...")
	return err
}

type TxMaker struct {
	connPool connPool
}

func NewTxMaker(connPool connPool) *TxMaker {
	return &TxMaker{connPool: connPool}
}

func (f *TxMaker) BeginTransaction() *sql.Tx {
	return f.connPool.BeginTransaction()
}
