package activerecord

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// humanQuerier is an implementation of the local humanFinder interface.
type humanQuerier struct {
	dbClient sqlDbClient
}

func NewHumanQuerier(dbClient sqlDbClient) *humanQuerier {
	return &humanQuerier{dbClient: dbClient}
}

// Useful methods to reconstitute humans ...
func (q *humanQuerier) FindByID(ctx context.Context, id string) *Human {
	query := fmt.Sprintf(humanFindByID, id)
	row := q.dbClient.Query(ctx, query)

	var (
		idField uuid.UUID
		name    string
		weight  int
		height  int
	)

	_ = row.Scan(&idField, &name, &weight, &height)
	return &Human{
		ID:     idField,
		Name:   name,
		Weight: weight,
		Height: height,
	}
}
