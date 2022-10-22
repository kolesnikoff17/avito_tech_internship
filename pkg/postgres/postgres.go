package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib" // needs for connection via sqlx
	"github.com/jmoiron/sqlx"
)

const (
	defaultMaxConns = 5
)

// Db keeps pool of connections to db
type Db struct {
	maxConns int
	Pool     *sqlx.DB
}

// New is constructor for Db
func New(uri string, opts ...Option) (*Db, error) {
	c, err := sqlx.Connect("pgx", uri)
	if err != nil {
		return nil, err
	}
	pg := &Db{Pool: c,
		maxConns: defaultMaxConns}
	for _, opt := range opts {
		opt(pg)
	}
	pg.Pool.SetMaxIdleConns(pg.maxConns)
	pg.Pool.SetMaxOpenConns(pg.maxConns)
	return pg, nil
}

// Close closes db's connection pool
func (db *Db) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
