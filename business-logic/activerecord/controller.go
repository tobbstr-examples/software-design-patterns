package activerecord

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type (
	RequestModel struct {
		ID string `json:"id"`
	}

	ResponseModel struct {
		Bmi string `json:"bmi"`
	}
)

type humanFinder interface {
	FindByID(ctx context.Context, id string) *Human
}

type Controller struct {
	humanFinder humanFinder
}

func NewController(humanFinder humanFinder) *Controller {
	return &Controller{humanFinder: humanFinder}
}

// CalculateBMI calculates BMI for a human given the id in the payload
func (c *Controller) CalculateBMI(w http.ResponseWriter, r *http.Request) {
	// Bind request model
	body, _ := io.ReadAll(r.Body)
	var reqModel RequestModel
	json.Unmarshal(body, &reqModel)

	// Use humanFinder to reconstitute the human given by the id in the request model
	ctx := r.Context()
	human := c.humanFinder.FindByID(ctx, reqModel.ID)

	// Calculate the human's BMI = business logic
	bmi := human.Bmi()

	// Create response model
	respModel := ResponseModel{
		Bmi: strconv.FormatFloat(float64(bmi), 'f', 1, 32),
	}

	// Encode response model as JSON
	body, _ = json.Marshal(&respModel)

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
