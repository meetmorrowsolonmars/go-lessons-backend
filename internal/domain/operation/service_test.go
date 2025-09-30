package operation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name          string
		service       func(mc *minimock.Controller) *Service
		operation     model.Operation
		wantOperation model.Operation
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				userStore := NewUserStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)
				categoryStore := NewCategoryStoreMock(mc)

				categoryStore.GetByIDMock.
					When(minimock.AnyContext, 5).
					Then(model.Category{}, nil)

				userStore.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.User{}, nil)

				accountStore.GetByIDMock.
					When(minimock.AnyContext, 100).
					Then(model.Account{}, nil)

				operation := model.Operation{
					ID:          uuid.MustParse("da270184-3a9e-4082-9509-c48a6f8c505f"),
					UserID:      10,
					AccountID:   100,
					Type:        model.OperationTypeDebit,
					CategoryID:  5,
					Amount:      decimal.NewFromFloat(160.5),
					Description: "Salary",
					CreateTime:  now,
				}

				store.CreateMock.
					When(minimock.AnyContext, model.Operation{
						UserID:      operation.UserID,
						AccountID:   operation.AccountID,
						Type:        operation.Type,
						CategoryID:  operation.CategoryID,
						Amount:      operation.Amount,
						Description: operation.Description,
					}).
					Then(operation, nil)

				return &Service{
					store:         store,
					userStore:     userStore,
					accountStore:  accountStore,
					categoryStore: categoryStore,
				}
			},
			operation: model.Operation{
				UserID:      10,
				AccountID:   100,
				Type:        model.OperationTypeDebit,
				CategoryID:  5,
				Amount:      decimal.NewFromFloat(160.5),
				Description: "Salary",
			},
			wantOperation: model.Operation{
				ID:          uuid.MustParse("da270184-3a9e-4082-9509-c48a6f8c505f"),
				UserID:      10,
				AccountID:   100,
				Type:        model.OperationTypeDebit,
				CategoryID:  5,
				Amount:      decimal.NewFromFloat(160.5),
				Description: "Salary",
				CreateTime:  now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "get category error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				userStore := NewUserStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)
				categoryStore := NewCategoryStoreMock(mc)

				categoryStore.GetByIDMock.
					When(minimock.AnyContext, 5).
					Then(model.Category{}, errors.New("some error"))

				return &Service{
					store:         store,
					userStore:     userStore,
					accountStore:  accountStore,
					categoryStore: categoryStore,
				}
			},
			operation: model.Operation{
				CategoryID: 5,
			},
			wantOperation: model.Operation{},
			wantErr:       assert.Error,
		},
		{
			name: "get user error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				userStore := NewUserStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)
				categoryStore := NewCategoryStoreMock(mc)

				userStore.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.User{}, errors.New("some error"))

				return &Service{
					store:         store,
					userStore:     userStore,
					accountStore:  accountStore,
					categoryStore: categoryStore,
				}
			},
			operation: model.Operation{
				UserID: 10,
			},
			wantOperation: model.Operation{},
			wantErr:       assert.Error,
		},
		{
			name: "get account error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				userStore := NewUserStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)
				categoryStore := NewCategoryStoreMock(mc)

				userStore.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.User{}, nil)

				accountStore.GetByIDMock.
					When(minimock.AnyContext, 100).
					Then(model.Account{}, errors.New("some error"))

				return &Service{
					store:         store,
					userStore:     userStore,
					accountStore:  accountStore,
					categoryStore: categoryStore,
				}
			},
			operation: model.Operation{
				UserID:    10,
				AccountID: 100,
			},
			wantOperation: model.Operation{},
			wantErr:       assert.Error,
		},
		{
			name: "create operation error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				userStore := NewUserStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)
				categoryStore := NewCategoryStoreMock(mc)

				userStore.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.User{}, nil)

				accountStore.GetByIDMock.
					When(minimock.AnyContext, 100).
					Then(model.Account{}, nil)

				operation := model.Operation{
					UserID:    10,
					AccountID: 100,
				}

				store.CreateMock.
					When(minimock.AnyContext, model.Operation{
						UserID:    operation.UserID,
						AccountID: operation.AccountID,
					}).
					Then(model.Operation{}, errors.New("some error"))

				return &Service{
					store:         store,
					userStore:     userStore,
					accountStore:  accountStore,
					categoryStore: categoryStore,
				}
			},
			operation: model.Operation{
				UserID:    10,
				AccountID: 100,
			},
			wantOperation: model.Operation{},
			wantErr:       assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			service := tc.service(mc)

			gotOperation, gotErr := service.Create(context.Background(), tc.operation)

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantOperation, gotOperation)
		})
	}
}

func TestService_GetByAccountID(t *testing.T) {
	t.Parallel()

	var (
		now    = time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)
		limit  = int64(10)
		offset = int64(0)
	)

	cases := []struct {
		name           string
		service        func(mc *minimock.Controller) *Service
		accountID      int64
		wantOperations []model.Operation
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)

				accountStore.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.Account{}, nil)

				operations := []model.Operation{
					{
						ID:          uuid.MustParse("da270184-3a9e-4082-9509-c48a6f8c505f"),
						UserID:      10,
						AccountID:   10,
						Type:        model.OperationTypeDebit,
						Amount:      decimal.NewFromFloat(160.5),
						Description: "Salary",
						CreateTime:  now,
					},
					{
						ID:          uuid.MustParse("516f29da-d014-4594-a28e-771ca71b80c2"),
						UserID:      10,
						AccountID:   10,
						Type:        model.OperationTypeCredit,
						Amount:      decimal.NewFromFloat(160.5),
						Description: "Travel",
						CreateTime:  now,
					},
				}

				store.GetByAccountIDMock.
					When(minimock.AnyContext, 10, limit, offset).
					Then(operations, nil)

				return &Service{
					store:        store,
					accountStore: accountStore,
				}
			},
			accountID: 10,
			wantOperations: []model.Operation{
				{
					ID:          uuid.MustParse("da270184-3a9e-4082-9509-c48a6f8c505f"),
					UserID:      10,
					AccountID:   10,
					Type:        model.OperationTypeDebit,
					Amount:      decimal.NewFromFloat(160.5),
					Description: "Salary",
					CreateTime:  now,
				},
				{
					ID:          uuid.MustParse("516f29da-d014-4594-a28e-771ca71b80c2"),
					UserID:      10,
					AccountID:   10,
					Type:        model.OperationTypeCredit,
					Amount:      decimal.NewFromFloat(160.5),
					Description: "Travel",
					CreateTime:  now,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "get account error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)

				accountStore.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.Account{}, errors.New("some error"))

				return &Service{
					store:        store,
					accountStore: accountStore,
				}
			},
			accountID:      10,
			wantOperations: nil,
			wantErr:        assert.Error,
		},
		{
			name: "get operations error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)
				accountStore := NewAccountStoreMock(mc)

				accountStore.GetByIDMock.
					When(minimock.AnyContext, 10).
					Then(model.Account{}, nil)

				store.GetByAccountIDMock.
					When(minimock.AnyContext, 10, limit, offset).
					Then(nil, errors.New("some error"))

				return &Service{
					store:        store,
					accountStore: accountStore,
				}
			},
			accountID:      10,
			wantOperations: nil,
			wantErr:        assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			service := tc.service(mc)

			gotOperations, gotErr := service.GetByAccountID(context.Background(), tc.accountID, limit, offset)

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantOperations, gotOperations)
		})
	}
}
