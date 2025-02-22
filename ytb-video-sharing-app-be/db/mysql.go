package db

import (
	"context"
	"database/sql"
	"os"
	"ytb-video-sharing-app-be/pkg"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysql struct {
	client *sqlx.DB
}

func NewMySQL() (pkg.Database, error) {
	dbClient, err := sqlx.Open("mysql", os.Getenv("MYSQL_DSN"))

	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Open")
	}

	err = dbClient.Ping()

	if err != nil {
		return nil, errors.Wrap(err, "dbClient.Ping")
	}

	return &mysql{
		client: dbClient,
	}, nil
}

// Begin implements pkg.Database.
func (m *mysql) Begin(ctx context.Context) (pkg.Tx, error) {
	transaction, err := m.client.Begin()

	if err != nil {
		return nil, errors.Wrap(err, "m.client.Begin")
	}

	return &tx{
		transaction: transaction,
	}, nil
}

// Close implements pkg.Database.
func (m *mysql) Close(ctx context.Context) error {
	return m.client.Close()
}

// Exec implements pkg.Database.
func (m *mysql) Exec(ctx context.Context, sql string, args ...any) error {
	res, err := m.client.Exec(sql, args...)

	if err != nil {
		return parseError(err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return pkg.ErrNoRowsAffected
	}
	return nil
}

func parseError(err error) error {
	var mysqlErr *mysqlDriver.MySQLError

	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return pkg.ErrDuplicate
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return pkg.ErrNoRows
	}

	return err
}

// Query implements pkg.Database.
func (m *mysql) Query(ctx context.Context, sql string, args ...any) (pkg.Rows, error) {
	res, err := m.client.Query(sql, args...)

	if err != nil {
		return nil, errors.Wrap(err, "m.client.Query")
	}

	return &rows{
		rows: res,
	}, nil

}

// QueryRow implements pkg.Database.
func (m *mysql) QueryRow(ctx context.Context, sql string, args ...any) pkg.Row {
	return m.client.QueryRow(sql, args...)
}

type tx struct {
	transaction *sql.Tx
}

// Commit implements pkg.Tx.
func (t *tx) Commit(ctx context.Context) error {
	return t.transaction.Commit()
}

// Exec implements pkg.Tx.
func (t *tx) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := t.transaction.Exec(sql, args...)
	return err
}

// QueryRow implements pkg.Tx.
func (t *tx) QueryRow(ctx context.Context, sql string, args ...any) pkg.Row {
	return t.transaction.QueryRow(sql, args...)
}

// Rollback implements pkg.Tx.
func (t *tx) Rollback(ctx context.Context) error {
	return t.transaction.Rollback()
}

type rows struct {
	rows *sql.Rows
}

// Close implements pkg.Rows.
func (r *rows) Close() {
	r.rows.Close()
}

// Next implements pkg.Rows.
func (r *rows) Next() bool {
	return r.rows.Next()
}

// Scan implements pkg.Rows.
func (r *rows) Scan(dst ...any) error {
	return r.rows.Scan(dst...)
}
