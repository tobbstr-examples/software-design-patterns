package activerecord

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const (
	humanFindByID = "SELECT id, name, weight, height from human WHERE id = '%s';"
	humanInsert   = "INSERT INTO human (id, name, weight, height) VALUES('%s','%s','%d','%d');"
	humanUpdate   = "UPDATE human SET id = %s, name = %s, weight = %d, height = %d;"
	humanDelete   = "DELETE FROM human WHERE id = %s;"
)

type sqlDbClient interface {
	Exec(ctx context.Context, query string)
	Query(ctx context.Context, query string) *sql.Row
}

type Human struct {
	dbClient sqlDbClient

	ID     uuid.UUID
	Name   string
	Weight int // kg
	Height int // centimeters
}

func NewHuman(dbClient sqlDbClient, name string, weight, height int) (*Human, error) {
	// Validate input to only allow instantiation of humans with sane values
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	if weight < 4 || weight > 300 {
		return nil, fmt.Errorf("weight outside valid range")
	}

	if height < 35 || height > 250 {
		return nil, fmt.Errorf("height outside valid range")
	}

	return &Human{
		dbClient: dbClient,
		ID:       uuid.New(),
		Name:     name,
		Weight:   weight,
		Height:   height,
	}, nil
}

// Encapuslation of persistence mechanism interaction
func (h *Human) Insert(ctx context.Context) {
	query := fmt.Sprintf(humanInsert, h.ID, h.Name, h.Weight, h.Height)
	h.dbClient.Exec(ctx, query)
}

// Encapuslation of persistence mechanism interaction
func (h *Human) Update(ctx context.Context) {
	query := fmt.Sprintf(humanUpdate, h.ID, h.Name, h.Weight, h.Height)
	h.dbClient.Exec(ctx, query)
}

// Encapuslation of persistence mechanism interaction
func (h *Human) Delete(ctx context.Context) {
	query := fmt.Sprintf(humanDelete, h.ID)
	h.dbClient.Exec(ctx, query)
}

// Business logic
func (h *Human) Bmi() float32 {
	return float32(h.Weight) / ((float32(h.Height) / 100) * (float32(h.Height) / 100))
}
