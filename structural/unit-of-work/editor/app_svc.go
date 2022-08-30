package editor

import (
	"context"

	"github.com/tobbstr-examples/business-logic-patterns/structural/unit-of-work/uow"
)

type ApplicationService struct {
	uowDoer uow.Doer
}

func NewApplicationService(uowDoer uow.Doer) *ApplicationService {
	return &ApplicationService{uowDoer: uowDoer}
}

func (svc *ApplicationService) OrchestrateWritingToMultipleTables(ctx context.Context) {
	// This is the implementation of the uow.Do() function. In other words,
	// it's the actual transaction the ApplicationService orchestrates in this method.
	doUnitOfWork := func(ctx context.Context, stores uow.Stores) error {
		if err := stores.A().Save(ctx, "hello world"); err != nil {
			// handle error
			return err
		}

		if err := stores.B().Save(ctx, 5); err != nil {
			// handle error
			return err
		}

		return nil
	}

	if err := svc.uowDoer.Do(ctx, doUnitOfWork); err != nil {
		// handle error
	}
}
