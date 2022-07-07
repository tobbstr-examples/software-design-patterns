package transactionscript

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type sqlDbClient interface {
	Exec(ctx context.Context, query string)
}

type RequestModel struct {
	ID     string `json:"id"`
	Weight int    `json:"weight"`
}

type Controller struct {
	dbClient sqlDbClient
}

func NewController(dbClient sqlDbClient) *Controller {
	return &Controller{dbClient: dbClient}
}

// CreateResource is an HTTP endpoint that creates a fictious resource using the transaction script pattern.
// Note! Error handling and input validation have been omitted in this example, but if they were included
// they would all be part of the same transaction script i.e. the CreateResource method.
func (c *Controller) CreateResource(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var reqModel RequestModel
	json.Unmarshal(body, &reqModel)

	// This paragraph inserts a record in a database. For more advanced scenarios in which a database transaction
	// would be needed, it'd be all be done in this HTTP handler as a single transaction script.
	ctx := r.Context()
	query := fmt.Sprintf("INSERT INTO example_tbl (id, weight) VALUES('%s', '%d');", reqModel.ID, reqModel.Weight)
	c.dbClient.Exec(ctx, query)

	w.WriteHeader(http.StatusCreated)
	w.Write(nil)
}
