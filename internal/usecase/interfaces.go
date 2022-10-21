package usecase

import (
  "balance_api/internal/entity"
  "context"
)

// Balance is an interface for model layer
type Balance interface {
  GetByID(ctx context.Context, id int) (entity.Balance, error)
  CreateOrder(ctx context.Context, order entity.Order) error
  CommitOrder(ctx context.Context, order entity.Order) error
  RollbackOrder(ctx context.Context, order entity.Order) error
  Increase(ctx context.Context, balance entity.Balance) error
  GetHistory(ctx context.Context, history entity.History) (entity.History, error)
  UpdateReport(ctx context.Context, year, month int) (string, error)
}

// BalanceRepo is an interface for repository layer
type BalanceRepo interface {
  GetByID(ctx context.Context, id int) (entity.Balance, error)
  CreateOrder(ctx context.Context, order entity.Order) error
  VerifyOrder(ctx context.Context, order entity.Order) error
  ChangeOrderStatus(ctx context.Context, order entity.Order) error
  Increase(ctx context.Context, balance entity.Balance) error
  GetHistory(ctx context.Context, history entity.History) (entity.History, error)
  GetReport(ctx context.Context, year, month int) (entity.Report, error)
}
