package usecase

import (
	"balance_api/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

// BalanceUseCase keeps all it needs to perform business logic
type BalanceUseCase struct {
	repo   BalanceRepo
	report ReportFile
}

// New is a constructor for BalanceUseCase
func New(r BalanceRepo, f ReportFile) *BalanceUseCase {
	return &BalanceUseCase{
		repo:   r,
		report: f,
	}
}

// GetByID returns entity.Balance of given id from repo, entity.ErrNoID in case if there is no such one
func (uc *BalanceUseCase) GetByID(ctx context.Context, id int) (entity.Balance, error) {
	balance, err := uc.repo.GetByID(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return entity.Balance{}, err
	case err != nil:
		return entity.Balance{}, fmt.Errorf("BalanceUseCase - GetByID: %w", err)
	}
	return balance, nil
}

// CreateOrder puts new order in repo, returns entity.ErrNoID if there is no such user,
// entity.ErrOrderExists if order exists, entity.ErrNotEnoughMoney if user doesn't have
// enough money for this order, entity.ErrNoService if service id is wrong
func (uc *BalanceUseCase) CreateOrder(ctx context.Context, order entity.Order) error {
	balance, err := uc.repo.GetByID(ctx, order.UserID)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("BalanceUseCase - CreateOrder: %w", err)
	}
	err = isLess(balance.Amount, order.Sum)
	if err != nil {
		return err
	}
	err = uc.repo.CheckServiceID(ctx, order.ServiceID)
	switch {
	case errors.Is(err, entity.ErrNoService):
		return err
	case err != nil:
		return fmt.Errorf("BalanceUseCase - CreateOrder: %w", err)
	}
	_, err = uc.repo.GetOrderByID(ctx, order.ID)
	switch {
	case err == nil:
		return entity.ErrOrderExists
	case !errors.Is(err, entity.ErrOrderNoExists):
		return fmt.Errorf("BalanceUseCase - CreateOrder: %w", err)
	}
	err = uc.repo.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("BalanceUseCase - CreateOrder: %w", err)
	}
	return nil
}

// ChangeOrderStatus commits or rollback order, returns entity.ErrOrderNoExists if there is no order with that id,
// entity.ErrOrderMismatch if order in request is not the same as database one, entity.ErrCantChangeStatus if order
// already committed/canceled
func (uc *BalanceUseCase) ChangeOrderStatus(ctx context.Context, order entity.Order) error {
	dbOrder, err := uc.repo.GetOrderByID(ctx, order.ID)
	switch {
	case errors.Is(err, entity.ErrOrderNoExists):
		return err
	case err != nil:
		return fmt.Errorf("BalanceUseCase - ChangeOrderStatus: %w", err)
	}
	if order.ServiceID != dbOrder.ServiceID || order.UserID != dbOrder.UserID || !isEqual(order.Sum, dbOrder.Sum) {
		return entity.ErrOrderMismatch
	}
	if dbOrder.StatusID != 1 {
		return entity.ErrCantChangeStatus
	}
	if order.StatusID == 2 {
		err = uc.repo.CommitOrder(ctx, order)
		if err != nil {
			return fmt.Errorf("BalanceUseCase - ChangeOrderStatus: %w", err)
		}
		return nil
	}
	err = uc.repo.RollbackOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("BalanceUseCase - ChangeOrderStatus: %w", err)
	}
	return nil
}

// Increase adds money to user or creates it if there is no one
func (uc *BalanceUseCase) Increase(ctx context.Context, balance entity.Balance) error {
	_, err := uc.repo.GetByID(ctx, balance.ID)
	switch {
	case errors.Is(err, entity.ErrNoID):
		err = uc.repo.CreateUser(ctx, balance)
		if err != nil {
			return fmt.Errorf("BalanceUseCase - Increase: %w", err)
		}
		return nil
	case err != nil:
		return fmt.Errorf("BalanceUseCase - Increase: %w", err)
	}
	err = uc.repo.Increase(ctx, balance)
	if err != nil {
		return fmt.Errorf("BalanceUseCase - Increase: %w", err)
	}
	return nil
}

// GetHistory gets list of orders from db of given user, returns entity.ErrNoID if there is no such user or
func (uc *BalanceUseCase) GetHistory(ctx context.Context, history entity.History) (entity.History, error) {
	_, err := uc.repo.GetByID(ctx, history.UserID)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return entity.History{}, err
	case err != nil:
		return entity.History{}, fmt.Errorf("BalanceUseCase - GetHistory: %w", err)
	}
	history, err = uc.repo.GetHistory(ctx, history)
	switch {
	case errors.Is(err, entity.ErrEmptyPage):
		return entity.History{}, err
	case err != nil:
		return entity.History{}, fmt.Errorf("BalanceUseCase - GetHistory: %w", err)
	}
	return history, nil
}

// UpdateReport creates a report, returns entity.ErrEmptyReport if report is empty
func (uc *BalanceUseCase) UpdateReport(ctx context.Context, year, month int) (string, error) {
	r, err := uc.repo.GetReport(ctx, year, month)
	switch {
	case errors.Is(err, entity.ErrEmptyReport):
		return "", err
	case err != nil:
		return "", fmt.Errorf("BalanceUseCase - UpdateReport: %w", err)
	}
	zero := ""
	if month < 10 {
		zero = "0"
	}
	name := strconv.Itoa(year) + "-" + zero + strconv.Itoa(month)
	name, err = uc.report.Create(ctx, name, r)
	if err != nil {
		return "", fmt.Errorf("BalanceUseCase - UpdateReport: %w", err)
	}
	return name, nil
}

// GetReportDir is getter of report dir
func (uc *BalanceUseCase) GetReportDir() string {
	return uc.report.GetDir()
}

func isLess(gotStr, decStr string) error {
	got, _ := decimal.NewFromString(gotStr)
	dec, _ := decimal.NewFromString(decStr)
	if got.LessThan(dec) {
		return entity.ErrNotEnoughMoney
	}
	return nil
}

func isEqual(orderStr, dbStr string) bool {
	order, _ := decimal.NewFromString(orderStr)
	db, _ := decimal.NewFromString(dbStr)
	return order.Equal(db)
}
