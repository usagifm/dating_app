package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/usagifm/dating-app/lib/atomic"
	atomicSqlx "github.com/usagifm/dating-app/lib/atomic/sqlx"
	"github.com/usagifm/dating-app/lib/logger"
	vplusSqlx "github.com/usagifm/dating-app/lib/sqlx"
	"github.com/usagifm/dating-app/src/app"
)

const (
	createNewUser = iota + 200
	updateUser
	getUserProfile
	getUserByEmail
	AllFields = `id,is_verified,name,gender,email,password,age,bio,photo_url,created_at,updated_at`
)

var (
	masterQueries = []string{
		getUserProfile: `SELECT id,is_verified,name,gender,email,age,bio,photo_url,created_at,updated_at FROM tbl_users WHERE id = $1 LIMIT 1`,
		getUserByEmail: `SELECT id,is_verified,name,gender,email,password,age,bio,photo_url FROM tbl_users WHERE email = $1 LIMIT 1`,
	}
	masterNamedQueries = []string{
		createNewUser: `INSERT INTO tbl_users (is_verified,name,gender,email,password,age,bio,photo_url)
		VALUES (:is_verified ,:name ,:gender ,:email ,:password ,:age ,:bio ,:photo_url) RETURNING id`,
		updateUser: `UPDATE tbl_users SET name = :name, age = :age, bio = :bio, updated_at = NOW() 
		WHERE id = :user_id`,
	}
)

type UserRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redisClient       *redis.Client
	redisConfig       app.Redis
}

func InitUserRepository(ctx context.Context, db *sqlx.DB, redisClient *redis.Client, redisConfig app.Redis) (*UserRepository, error) {
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

	return &UserRepository{
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		db:                db,
		redisClient:       redisClient,
		redisConfig:       redisConfig,
	}, nil
}

func (r *UserRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *UserRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
