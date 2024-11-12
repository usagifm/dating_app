package user_preference

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
	createNewUserPreference = iota + 200
	updateUserPreference
	getUserPreference
	AllFields = `id,user_id,preferred_gender,min_age,max_age,created_at,updated_at`
)

var (
	masterQueries = []string{
		getUserPreference: fmt.Sprintf(`SELECT %s FROM tbl_user_preferences WHERE user_id = $1 LIMIT 1`, AllFields),
	}
	masterNamedQueries = []string{
		createNewUserPreference: `INSERT INTO tbl_user_preferences (user_id,preferred_gender,min_age,max_age)
		VALUES (:user_id ,:preferred_gender ,:min_age ,:max_age) RETURNING id`,
		updateUserPreference: `UPDATE tbl_user_preferences SET preferred_gender = :preferred_gender, min_age = :min_age, max_age = :max_age, updated_at = NOW() 
		WHERE user_id = :user_id`,
	}
)

type UserPreferenceRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redisClient       *redis.Client
	redisConfig       app.Redis
}

func InitUserPreferenceRepository(ctx context.Context, db *sqlx.DB, redisClient *redis.Client, redisConfig app.Redis) (*UserPreferenceRepository, error) {
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

	return &UserPreferenceRepository{
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		db:                db,
		redisClient:       redisClient,
		redisConfig:       redisConfig,
	}, nil
}

func (r *UserPreferenceRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *UserPreferenceRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
