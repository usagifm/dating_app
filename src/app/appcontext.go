package app

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/usagifm/dating-app/lib/i18n"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/lib/postgres"
	vplusRedis "github.com/usagifm/dating-app/lib/redis"
	"github.com/usagifm/dating-app/lib/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type appContext struct {
	db               *sqlx.DB
	requestValidator *validator.Validate
	cfg              *Configuration
	redisClient      *redis.Client
	tracer           trace.Tracer
}

var appCtx appContext

func Init(ctx context.Context) error {
	logger.Init(ctx)

	logger.GetLogger(ctx).Info("Loading the config...")
	cfg, err := InitConfig(ctx)
	if err != nil {
		return err
	}
	logger.GetLogger(ctx).Info("Config loaded sucessfully...")

	logger.GetLogger(ctx).Info("Connecting to the redis...")
	rds, err := vplusRedis.Init(ctx, cfg.Redis.Host, cfg.Redis.Password)
	if err != nil {
		return err
	}
	logger.GetLogger(ctx).Info("Connected to the redis !")

	if err := i18n.Init(ctx, cfg.Translation.FilePath, cfg.Translation.DefaultLanguage); err != nil {
		panic(err)
	}

	logger.GetLogger(ctx).Info("Connecting to the database...")
	db, err := postgres.InitSQLX(ctx, postgres.PostgresConfig{
		ConnectionUrl:      cfg.Postgres.ConnURI,
		MaxPoolSize:        cfg.Postgres.MaxPoolSize,
		MaxIdleConnections: cfg.Postgres.MaxIdleConnections,
		ConnMaxIdleTime:    cfg.Postgres.MaxIdleTime,
		ConnMaxLifeTime:    cfg.Postgres.MaxLifeTime,
	})
	if err != nil {
		return err
	}
	logger.GetLogger(ctx).Info("Connected to the database !")

	logger.GetLogger(ctx).Info("Initializing tracer...")
	_, err = tracer.Init(ctx, cfg.ServiceName, "0.0.1", cfg.TraceEndpoint, cfg.TraceRate)
	if err != nil {
		return err
	}
	logger.GetLogger(ctx).Info("Tracer initialized !")

	appCtx = appContext{
		db:               db,
		redisClient:      rds,
		requestValidator: validator.New(),
		cfg:              cfg,
		tracer:           otel.Tracer(cfg.ServiceName),
	}

	return nil
}

func RequestValidator() *validator.Validate {
	return appCtx.requestValidator
}

func DB() *sqlx.DB {
	return appCtx.db
}

func Cache() *redis.Client {
	return appCtx.redisClient
}

func Config() Configuration {
	return *appCtx.cfg
}

func Tracer() trace.Tracer {
	return appCtx.tracer
}
