package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DB struct {
	l  *zap.SugaredLogger
	db *pgxpool.Pool
}

func NewDB(db *pgxpool.Pool, l *zap.SugaredLogger) DBClient {
	return &DB{db: db, l: l}
}

func (d *DB) Ping(ctx context.Context) error {
	d.l.Infow("Pinging database")

	err := d.db.Ping(ctx)
	if err != nil {
		d.l.Errorw("Database ping failed", "error", err)
		return err
	}

	d.l.Infow("Database ping successful")
	return nil
}

func (d *DB) Close(ctx context.Context) {
	d.l.Infow("Closing database connection")
	d.db.Close()
}

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	d.l.Debugw("Exec", "query", query)
	tag, err := d.db.Exec(ctx, query, args...)
	if err != nil {
		d.l.Errorw("Exec failed", "error", err, "query", query)
		return pgconn.CommandTag{}, err
	}
	return tag, nil
}

func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	d.l.Debugw("Querying database", "query", query, "args", args)

	rows, err := d.db.Query(ctx, query, args...)
	if err != nil {
		d.l.Errorw("Error querying database", "error", err, "query", query, "args", args)
		return nil, err
	}

	return rows, nil
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	d.l.Debugw("QueryRow", "query", query)
	return d.db.QueryRow(ctx, query, args...)
}
