package user

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
	getUserPackages = iota + 200
	getUserPackage
	createOrUpdateUserPackage
	AllFields = `id,user_id,package_id,valid_date,created_at,updated_at`
)

var (
	masterQueries = []string{
		getUserPackages: fmt.Sprintf(`SELECT %s FROM tbl_user_packages`, AllFields),
		getUserPackage:  fmt.Sprintf(`SELECT %s FROM tbl_user_packages LIMIT 1`, AllFields),
	}
	masterNamedQueries = []string{
		createOrUpdateUserPackage: `INSERT INTO tbl_user_packages(user_id,package_id,valid_date)
		VALUES (:user_id, :package_id, :valid_date)
		ON CONFLICT (user_id,package_id) DO UPDATE
			SET valid_date = EXCLUDED.valid_date
		RETURNING id`,
	}
)

type UserPackageRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redisClient       *redis.Client
	redisConfig       app.Redis
}

func InitUserPackageRepository(ctx context.Context, db *sqlx.DB, redisClient *redis.Client, redisConfig app.Redis) (*UserPackageRepository, error) {
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

	return &UserPackageRepository{
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		db:                db,
		redisClient:       redisClient,
		redisConfig:       redisConfig,
	}, nil
}

func (r *UserPackageRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = r.masterStmts[queryId]
	}
	return statement, err
}

func (r *UserPackageRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = r.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
