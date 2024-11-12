package user_swipe

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
	createUserSwipe = iota + 200
	getUserSwipe
	getUserSwipePerToday
	getMatchedUserSwipe
	AllFields = `id,swiper_id,swiped_id,swipe_type,created_at,updated_at`
)

var (
	masterQueries = []string{
		getUserSwipe:         fmt.Sprintf(`SELECT %s FROM tbl_user_swipes WHERE swiper_id = $1`, AllFields),
		getUserSwipePerToday: fmt.Sprintf(`SELECT %s FROM tbl_user_swipes WHERE swiper_id = $1 AND created_at > $2 AND created_at < $3`, AllFields),
		getMatchedUserSwipe:  fmt.Sprintf(`SELECT %s FROM tbl_user_swipes WHERE swiped_id = $1 AND swiper_id = $2 LIMIT 1`, AllFields),
	}
	masterNamedQueries = []string{
		createUserSwipe: `INSERT INTO tbl_user_swipes (
			swiper_id,swiped_id,swipe_type
		) VALUES (
			:swiper_id ,:swiped_id ,:swipe_type
		)RETURNING id;`,
	}
)

type UserSwipeRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redisClient       *redis.Client
	redisConfig       app.Redis
}

func InitUserSwipeRepository(ctx context.Context, db *sqlx.DB, redisClient *redis.Client, redisConfig app.Redis) (*UserSwipeRepository, error) {
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

	return &UserSwipeRepository{
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		db:                db,
		redisClient:       redisClient,
		redisConfig:       redisConfig,
	}, nil
}

func (r *UserSwipeRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *UserSwipeRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
