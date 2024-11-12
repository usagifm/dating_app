package atomic

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/usagifm/dating-app/lib/atomic"
	"github.com/usagifm/dating-app/lib/logger"
	"go.opentelemetry.io/otel/trace"
)

type SqlxAtomicSessionProvider struct {
	db    *sqlx.DB
	trace trace.Tracer
}

func NewSqlxAtomicSessionProvider(db *sqlx.DB, tr trace.Tracer) *SqlxAtomicSessionProvider {
	return &SqlxAtomicSessionProvider{
		db:    db,
		trace: tr,
	}
}

func (r *SqlxAtomicSessionProvider) BeginSession(ctx context.Context) (*atomic.AtomicSessionContext, error) {
	ctx, span := r.trace.Start(ctx, "SqlxAtomicSessionProvider/BeginSession")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.GetLogger(ctx).Error("begin tx err: ", err)
		return nil, err
	}

	atomicSession := NewAtomicSession(tx, r.trace)
	return atomic.NewAtomicSessionContext(ctx, atomicSession), nil
}
