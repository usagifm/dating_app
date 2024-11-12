package atomic

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/usagifm/dating-app/lib/logger"
	"go.opentelemetry.io/otel/trace"
)

type SqlxAtomicSession struct {
	tx    *sqlx.Tx
	trace trace.Tracer
}

func NewAtomicSession(tx *sqlx.Tx, tr trace.Tracer) *SqlxAtomicSession {
	return &SqlxAtomicSession{
		tx:    tx,
		trace: tr,
	}
}

func (s SqlxAtomicSession) Commit(ctx context.Context) error {
	ctx, span := s.trace.Start(ctx, "SqlxAtomicSession/Commit")
	defer span.End()

	err := s.tx.Commit()
	if err != nil {
		logger.GetLogger(ctx).Error("commit err: ", err)
	}
	return err
}

func (s SqlxAtomicSession) Rollback(ctx context.Context) error {
	ctx, span := s.trace.Start(ctx, "SqlxAtomicSession/Rollback")
	defer span.End()

	err := s.tx.Rollback()
	if err != nil {
		logger.GetLogger(ctx).Error("rollback err: ", err)
	}
	return err
}

func (s SqlxAtomicSession) Tx() *sqlx.Tx {
	return s.tx
}
