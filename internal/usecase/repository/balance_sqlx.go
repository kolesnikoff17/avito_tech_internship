package repository

import (
	"balance_api/internal/entity"
	"balance_api/pkg/postgres"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// BalanceRepo keeps db connection pool
type BalanceRepo struct {
	*postgres.Db
}

// New is a constructor for BalanceRepo
func New(db *postgres.Db) *BalanceRepo {
	return &BalanceRepo{db}
}

// GetByID returns entity.Balance of a given id, entity.ErrNoID in case if there is no such one
func (r *BalanceRepo) GetByID(ctx context.Context, id int) (entity.Balance, error) {
	var res entity.Balance
	err := r.Pool.GetContext(ctx, &res, `SELECT user_id, amount FROM users WHERE user_id = $1`, id)
	if err != nil {
		return entity.Balance{}, entity.ErrNoID
	}
	return res, nil
}

// CreateOrder creates new order and transfers money from user's balance to special account
func (r *BalanceRepo) CreateOrder(ctx context.Context, order entity.Order) error {
	tx, err := r.Pool.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("BalanceRepository - CreateOrder: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO orders (order_id, service_id, user_id, order_sum, status_id)
						VALUES (:order_id, :service_id, :user_id, :order_sum, 1)`, order)
	if err != nil {
		return fmt.Errorf("BalanceRepository - CreateOrder: %w", err)
	}
	_, err = tx.NamedExecContext(ctx,
		`UPDATE users SET amount = amount - :order_sum, reserved = reserved + :order_sum 
             WHERE user_id = :user_id`, order)
	if err != nil {
		return fmt.Errorf("BalanceRepository - CreateOrder: %w", err)
	}
	return tx.Commit()
}

// GetOrderByID returns order with given id, entity.ErrOrderNoExists if there is no one
func (r *BalanceRepo) GetOrderByID(ctx context.Context, id int) (entity.Order, error) {
	var res entity.Order
	err := r.Pool.GetContext(ctx, &res,
		`SELECT order_id, service_id, user_id, status_id, order_sum FROM orders WHERE order_id = $1`, id)
	if err != nil {
		return entity.Order{}, entity.ErrOrderNoExists
	}
	return res, nil
}

// CheckServiceID returns nil if service exists in db, entity.ErrNoService otherwise
func (r *BalanceRepo) CheckServiceID(ctx context.Context, id int) error {
	var service struct {
		ID int `db:"service_id"`
	}
	err := r.Pool.GetContext(ctx, &service,
		`SELECT service_id FROM services WHERE service_id = $1`, id)
	if err != nil {
		return entity.ErrNoService
	}
	return nil
}

// CommitOrder updates order and reduce amount of user's reserved money
func (r *BalanceRepo) CommitOrder(ctx context.Context, order entity.Order) error {
	tx, err := r.Pool.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("BalanceRepository - CommitOrder: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.NamedExecContext(ctx,
		`UPDATE orders SET status_id = :status_id, modified = now() WHERE order_id = :order_id`, order)
	if err != nil {
		return fmt.Errorf("BalanceRepository - CommitOrder: %w", err)
	}
	_, err = tx.NamedExecContext(ctx,
		`UPDATE users SET reserved = reserved - :order_sum WHERE user_id = :user_id`, order)
	if err != nil {
		return fmt.Errorf("BalanceRepository - CommitOrder: %w", err)
	}
	return tx.Commit()
}

// RollbackOrder updates order and transfers money back from reserved to main account
func (r *BalanceRepo) RollbackOrder(ctx context.Context, order entity.Order) error {
	tx, err := r.Pool.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("BalanceRepository - RollbackOrder: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.NamedExecContext(ctx,
		`UPDATE orders SET status_id = :status_id, modified = now() WHERE order_id = :order_id`, order)
	if err != nil {
		return fmt.Errorf("BalanceRepository - RollbackOrder: %w", err)
	}
	_, err = tx.NamedExecContext(ctx,
		`UPDATE users SET reserved = reserved - :order_sum, amount = amount + :order_sum 
             WHERE user_id = :user_id`, order)
	if err != nil {
		return fmt.Errorf("BalanceRepository - RollbackOrder: %w", err)
	}
	return tx.Commit()
}

// CreateUser creates new user and put order with initial replenishment
func (r *BalanceRepo) CreateUser(ctx context.Context, balance entity.Balance) error {
	tx, err := r.Pool.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("BalanceRepository - CreateUser: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO users (user_id, amount) VALUES (:user_id, :amount)`, balance)
	if err != nil {
		return fmt.Errorf("BalanceRepository - CreateUser: %w", err)
	}
	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO replenishments (user_id, amount) 
						VALUES (:user_id, :amount)`, balance)
	if err != nil {
		return fmt.Errorf("BalanceRepository - CreateUser: %w", err)
	}
	return tx.Commit()
}

// Increase updates user's account and  put order with this replenishment
func (r *BalanceRepo) Increase(ctx context.Context, balance entity.Balance) error {
	tx, err := r.Pool.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("BalanceRepository - Increase: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.NamedExecContext(ctx,
		`UPDATE users SET amount = amount + :amount WHERE user_id = :user_id`, balance)
	if err != nil {
		return fmt.Errorf("BalanceRepository - Increase: %w", err)
	}
	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO replenishments (user_id, amount) 
						VALUES (:user_id, :amount)`, balance)
	if err != nil {
		return fmt.Errorf("BalanceRepository - Increase: %w", err)
	}
	return tx.Commit()
}

// GetHistory return's user's transaction history, entity.ErrEmptyPage if page and limit are wrong
func (r *BalanceRepo) GetHistory(ctx context.Context, history entity.History) (entity.History, error) {
	var OrdersSet []entity.Order
	query := queryConstructor(history)
	err := r.Pool.SelectContext(ctx, &OrdersSet, query, history.UserID)
	if err != nil {
		return entity.History{}, fmt.Errorf("BalanceRepository - GetHistory: %w", err)
	}
	if len(OrdersSet) == 0 {
		return entity.History{}, entity.ErrEmptyPage
	}
	history.Orders = OrdersSet
	return history, nil
}

func queryConstructor(history entity.History) string {
	var str strings.Builder
	str.WriteString(`SELECT serv.service_name, o.order_sum, st.status_name, o.created 
											FROM orders AS o
											JOIN services AS serv ON o.service_id = serv.service_id
											JOIN status AS st ON o.status_id = st.status_id
											WHERE o.user_id = $1
											UNION
											SELECT 'Replenishment' AS service_name, amount AS order_sum, 'Approved' AS status_name, created
											FROM replenishments
											WHERE user_id = $1
											ORDER BY `)
	switch history.OrderBy {
	case "date":
		str.WriteString("created ")
	case "sum":
		str.WriteString("order_sum ")
	}
	if history.Desc {
		str.WriteString("DESC\n")
	} else {
		str.WriteString("\n")
	}
	offset := history.Limit * (history.Page - 1)
	str.WriteString(`LIMIT `)
	if history.Limit != 0 {
		str.WriteString(strconv.Itoa(history.Limit))
	} else {
		str.WriteString("ALL")
	}
	str.WriteString(` OFFSET `)
	str.WriteString(strconv.Itoa(offset))
	return str.String()
}

// GetReport returns report with given period, entity.ErrEmptyReport if there were no operations in this period
func (r *BalanceRepo) GetReport(ctx context.Context, year, month int) (entity.Report, error) {
	var Sums []entity.SumByService
	err := r.Pool.SelectContext(ctx, &Sums,
		`SELECT sum(o.order_sum) AS sums, s.service_name FROM orders AS o
						JOIN services AS s ON s.service_id = o.service_id
						WHERE o.status_id = 2
						AND EXTRACT(YEAR FROM o.modified) = $1 AND EXTRACT(MONTH FROM o.modified) = $2
						GROUP BY s.service_name`, year, month)
	if err != nil {
		return entity.Report{}, fmt.Errorf("BalanceRepository - GetReport: %w", err)
	}
	if len(Sums) == 0 {
		return entity.Report{}, entity.ErrEmptyReport
	}
	return entity.Report{Sums: Sums}, nil
}
