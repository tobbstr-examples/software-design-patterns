package domainmodel

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type (
	RequestModel struct {
		ID string `json:"id"`
	}
)

type orderAppService interface {
	SubmitOrder(ctx context.Context, id string) error
}

// OrderController is a simplified version of a controller since it doesn't handle errors.
type OrderController struct {
	orderAppSvc orderAppService
}

func NewController(orderAppService orderAppService) *OrderController {
	return &OrderController{orderAppSvc: orderAppService}
}

// SubmitOrder submits the Order with the given id
func (c *OrderController) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	// Bind request model
	body, _ := io.ReadAll(r.Body)
	var reqModel RequestModel
	json.Unmarshal(body, &reqModel)

	// Use orderAppSvc to coordinate the submission of the given Order
	ctx := r.Context()
	c.orderAppSvc.SubmitOrder(ctx, reqModel.ID)

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
