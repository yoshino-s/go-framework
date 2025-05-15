package ent

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/google/uuid"
	"github.com/yoshino-s/go-framework/log"
	"go.uber.org/zap"
)

// DebugDriver is a driver that logs all driver operations.
type DebugDriver struct {
	dialect.Driver             // underlying driver.
	logger         *zap.Logger // logger.
}

// Debug gets a driver and an optional logging function, and returns
// a new debugged-driver that prints all outgoing operations.
func DebugEntDriver(d dialect.Driver, logger *zap.Logger) dialect.Driver {
	drv := &DebugDriver{d, logger}
	return drv
}

func (d *DebugDriver) log(ctx context.Context, prefix string, query string, args any) {
	if d.logger.Level().Enabled(zap.DebugLevel) {
		d.logger.Debug(fmt.Sprintf("%s: %s", prefix, query), zap.String("query", query), zap.Any("args", args), log.Context(ctx))
	}
}

// Exec logs its params and calls the underlying driver Exec method.
func (d *DebugDriver) Exec(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "driver.Exec", query, args.([]any))
	return d.Driver.Exec(ctx, query, args, v)
}

// ExecContext logs its params and calls the underlying driver ExecContext method if it is supported.
func (d *DebugDriver) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Driver.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	d.log(ctx, "driver.ExecContext", query, args)
	return drv.ExecContext(ctx, query, args...)
}

// Query logs its params and calls the underlying driver Query method.
func (d *DebugDriver) Query(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "driver.Query", query, args.([]any))
	return d.Driver.Query(ctx, query, args, v)
}

// QueryContext logs its params and calls the underlying driver QueryContext method if it is supported.
func (d *DebugDriver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Driver.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	d.log(ctx, "driver.QueryContext", query, args)
	return drv.QueryContext(ctx, query, args...)
}

// Tx adds an log-id for the transaction and calls the underlying driver Tx command.
func (d *DebugDriver) Tx(ctx context.Context) (dialect.Tx, error) {
	tx, err := d.Driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	d.logger.Debug("driver.Tx", zap.String("id", id), log.Context(ctx))
	return &DebugTx{tx, id, d.logger, ctx}, nil
}

// BeginTx adds an log-id for the transaction and calls the underlying driver BeginTx command if it is supported.
func (d *DebugDriver) BeginTx(ctx context.Context, opts *sql.TxOptions) (dialect.Tx, error) {
	drv, ok := d.Driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.BeginTx is not supported")
	}
	tx, err := drv.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	d.logger.Debug("driver.BeginTx", zap.String("id", id), log.Context(ctx))
	return &DebugTx{tx, id, d.logger, ctx}, nil
}

// DebugTx is a transaction implementation that logs all transaction operations.
type DebugTx struct {
	dialect.Tx                 // underlying transaction.
	id         string          // transaction logging id.
	logger     *zap.Logger     // logger.
	ctx        context.Context // underlying transaction context.
}

func (d *DebugTx) log(ctx context.Context, prefix string, query string, args any) {
	if d.logger.Level().Enabled(zap.DebugLevel) {
		d.logger.Debug(fmt.Sprintf("%s: %s", prefix, query), zap.String("query", query), zap.Any("args", args), log.Context(ctx))
	}
}

// Exec logs its params and calls the underlying transaction Exec method.
func (d *DebugTx) Exec(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "Tx.Exec", query, args)
	return d.Tx.Exec(ctx, query, args, v)
}

// ExecContext logs its params and calls the underlying transaction ExecContext method if it is supported.
func (d *DebugTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Tx.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.ExecContext is not supported")
	}
	d.log(ctx, "Tx.ExecContext", query, args)
	return drv.ExecContext(ctx, query, args...)
}

// Query logs its params and calls the underlying transaction Query method.
func (d *DebugTx) Query(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "Tx.Query", query, args)
	return d.Tx.Query(ctx, query, args, v)
}

// QueryContext logs its params and calls the underlying transaction QueryContext method if it is supported.
func (d *DebugTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Tx.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.QueryContext is not supported")
	}
	d.log(ctx, "Tx.QueryContext", query, args)
	return drv.QueryContext(ctx, query, args...)
}

// Commit logs this step and calls the underlying transaction Commit method.
func (d *DebugTx) Commit() error {
	d.logger.Debug("Tx.Commit", zap.String("id", d.id), log.Context(d.ctx))
	return d.Tx.Commit()
}

// Rollback logs this step and calls the underlying transaction Rollback method.
func (d *DebugTx) Rollback() error {
	d.logger.Debug("Tx.Rollback", zap.String("id", d.id), log.Context(d.ctx))
	return d.Tx.Rollback()
}
