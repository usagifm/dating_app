package user_match

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/usagifm/dating-app/lib/atomic"
	atomicSqlx "github.com/usagifm/dating-app/lib/atomic/sqlx"
	"github.com/usagifm/dating-app/lib/logger"
	vplusSqlx "github.com/usagifm/dating-app/lib/sqlx"
	"github.com/usagifm/dating-app/src/app"
)

const (
	createUserMatch = iota + 200
	getUserMatches
	AllFields = `id,user_id_1,user_id_2,created_at,updated_at`
)

var (
	masterQueries = []string{
		getUserMatches: fmt.Sprintf(`SELECT %s FROM tbl_user_matches WHERE (user_id_1 = $1 OR user_id_2 = $1)`, AllFields),
	}
	masterNamedQueries = []string{
		createUserMatch: `INSERT INTO tbl_user_matches (
			user_id_1,user_id_2
		) VALUES (
			:user_id_1 ,:user_id_2
		)RETURNING id;`,
	}
)

type UserMatchRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redisClient       *redis.Client
	redisConfig       app.Redis
}

func InitUserMatchRepository(ctx context.Context, db *sqlx.DB, redisClient *redis.Client, redisConfig app.Redis) (*UserMatchRepository, error) {
	stmpts, err := vplusSqlx.PrepareQueries(db, masterQueries)
	if err != nil {
		logger.GetLogger(ctx).Error("PrepareQueries err:", err)
		return nil, err
	}

	namedStmpts, err := vplusSqlx.PrepareNamedQueries(db, masterNamedQueries)
	if err != nil {
		logger.GetLogger(ctx).Error("PrepareNamedQueries err:", err)
		return nil, err
	}

	return &UserMatchRepository{
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		db:                db,
		redisClient:       redisClient,
		redisConfig:       redisConfig,
	}, nil
}

func (r *UserMatchRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = r.masterStmts[queryId]
	}
	return statement, err
}

func (r *UserMatchRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = r.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
