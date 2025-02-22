package pkg

import (
	"context"
)

type Database interface {
	// Closes the connection to the database, returning an error if the operation fails.
	Close(ctx context.Context) error

	// Starts a new db transaction, returning a Tx object or an error.
	Begin(ctx context.Context) (Tx, error)

	// Executes a SQL query (e.g, INSERT, DELETE, UPDATE) with optional args, returning an errors if the operation fails.
	Exec(ctx context.Context, sql string, args ...any) error

	// Executes a SQL query expected to return a single row, returning a Row object for result scanning.
	QueryRow(ctx context.Context, sql string, args ...any) Row

	// Executes a SQL query expected to return multiple rows, returning a Rows object or an error
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) error
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Row interface {
	Scan(dst ...any) error
}

type Rows interface {
	Scan(dst ...any) error
	Close()
	Next() bool
}