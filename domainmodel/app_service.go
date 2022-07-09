package domainmodel

import (
	"context"
	"fmt"
	"time"

	"github.com/tobbstr-examples/business-logic--patterns/domainmodel/order"
)

type (
	messagePublisher interface {
		Publish(ctx context.Context, msg order.Message) error
	}

	txMaker interface {
		BeginTransaction(ctx context.Context) order.Tx
	}
)

type Service struct {
	messagePublisher messagePublisher
	txMaker          txMaker
}

func NewService(txMaker txMaker, messagePublisher messagePublisher) *Service {
	return &Service{
		messagePublisher: messagePublisher,
		txMaker:          txMaker,
	}
}

// SubmitOrder coordinates the submission of an Order. This example is a simplified version since
// it takes a shortcut. It's missing the Outbox pattern for making sure domain messages get
// delivered at least once.
func (s *Service) SubmitOrder(ctx context.Context, id string) error {
	// begin database transaction and instantiate a new order repository
	tx := s.txMaker.BeginTransaction(ctx)
	defer func() {
		tx.Rollback()
	}()
	orderRepo := order.NewRepository(tx)

	// use repository to reconstitute an existing Order
	repoFindByIDCtx, cancelRepoFindByID := context.WithTimeout(ctx, 5*time.Second)
	defer cancelRepoFindByID()
	order, err := orderRepo.FindByID(repoFindByIDCtx, id)
	if err != nil {
		return fmt.Errorf("could not find order by id = %s: %w", id, err)
	}

	// perform business logic
	if err = order.Submit(); err != nil {
		return fmt.Errorf("could not submit order by id = %s: %w", id, err)
	}

	// use repository to store the aggregate
	repoUpsertCtx, cancelRepoUpsert := context.WithTimeout(ctx, 5*time.Second)
	defer cancelRepoUpsert()
	if err = orderRepo.Upsert(repoUpsertCtx, order); err != nil {
		return fmt.Errorf("could not upsert order after submission: %w", err)
	}

	// publish domain messages to communicate the change(s) to other aggregates no matter if they belong
	// to the same monolith or some other application. This change in the system will be eventually consistent.
	messagesCtx, cancelMessages := context.WithTimeout(ctx, 10*time.Second)
	defer cancelMessages()
	for _, msg := range order.Messages() {
		// publish message
		s.messagePublisher.Publish(messagesCtx, msg)
	}

	tx.Commit()

	return nil
}
