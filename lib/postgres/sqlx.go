package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/usagifm/dating-app/lib/logger"
)

func InitSQLX(ctx context.Context, cfg PostgresConfig) (*sqlx.DB, error) {
	xDB, err := sqlx.Connect("postgres", cfg.ConnectionUrl)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to load the database err:%v", err)
		return nil, err
	}

	if err = xDB.Ping(); err != nil {
		logger.GetLogger(ctx).Errorf("failed to ping the database err:%v", err)
		return nil, err
	}

	xDB.SetMaxOpenConns(cfg.MaxPoolSize)
	xDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	xDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	xDB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	return xDB, nil
}
