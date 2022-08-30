package editor

import (
	"context"
	"fmt"
	"testing"

	"github.com/tobbstr-examples/business-logic-patterns/structural/unit-of-work/uow"
	uowmock "github.com/tobbstr-examples/business-logic-patterns/structural/unit-of-work/uow/mock"
)

func TestApplicationService_OrchestrateWritingToMultipleTables(t *testing.T) {
	errAny := fmt.Errorf("any-error")

	storeA := uowmock.AStore{SaveFn: func(ctx context.Context, id string) error { return nil }}
	storeAFailingOnSave := storeA
	storeAFailingOnSave.SaveFn = func(ctx context.Context, id string) error { return errAny }

	storeB := uowmock.BStore{SaveFn: func(ctx context.Context, id int) error { return nil }}
	storeBFailingOnSave := storeB
	storeBFailingOnSave.SaveFn = func(ctx context.Context, id int) error { return errAny }

	stores := uowmock.Stores{
		AFn: func() uow.AStore {
			return &storeA
		},
		BFn: func() uow.BStore {
			return &storeB
		},
	}

	storesWithStoreAFailingOnSave := stores
	storesWithStoreAFailingOnSave.AFn = func() uow.AStore { return &storeAFailingOnSave }

	storesWithStoreBFailingOnSave := stores
	storesWithStoreBFailingOnSave.BFn = func() uow.BStore { return &storeBFailingOnSave }

	type fields struct {
		uowDoer uow.Doer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "should return error when storeA.Save() fails",
			fields: fields{
				uowDoer: uowmock.NewDoer(&storesWithStoreAFailingOnSave),
			},
		},
		{
			name: "should return error when storeB.Save() fails",
			fields: fields{
				uowDoer: uowmock.NewDoer(&storesWithStoreBFailingOnSave),
			},
		},
		{
			name: "happy path",
			fields: fields{
				uowDoer: uowmock.NewDoer(&stores),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &ApplicationService{
				uowDoer: tt.fields.uowDoer,
			}
			svc.OrchestrateWritingToMultipleTables(tt.args.ctx)
		})
	}
}
