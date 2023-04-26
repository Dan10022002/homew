//go:generate mockgen -source ./database.go -destination=./mocks/database.go -package=mock_database
package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBops interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	GetPool(ctx context.Context) *pgxpool.Pool
}

type Database struct {
	cluster *pgxpool.Pool
}

func newDatabase(cluster *pgxpool.Pool) *Database {
	return &Database{cluster: cluster}
}

func (db Database) GetPool(_ context.Context) *pgxpool.Pool {
	return db.cluster
}

func (db Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.cluster, dest, query, args...)
}

func (db Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.cluster, dest, query, args...)
}

func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

func (db Database) Truncate(ctx context.Context) {
	var tables []string

	err := db.Select(ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public' "+
		"AND table_type = 'BASE TYPE' AND table_name != 'goose_db_version'")
	if err != nil {
		panic(err)
	}

	if len(tables) == 0 {
		panic("Run migrations plz")
	}

	q := fmt.Sprintf("Truncate tables %s", strings.Join(tables, ", "))
	if _, err := db.Exec(ctx, q); err != nil {
		panic(err)
	}
}
