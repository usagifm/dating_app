package migration

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/usagifm/dating-app/lib/logger"
	pg "github.com/usagifm/dating-app/lib/postgres"
	"github.com/usagifm/dating-app/src/app"
)

const (
	migrateLogIdentifier = "otp"
)

type MigrationService interface {
	Up(context.Context) error
	Rollback(context.Context) error
	Version(context.Context) (int, bool, error)
}

type migrationService struct {
	driver  database.Driver
	migrate *migrate.Migrate
}

func New(ctx context.Context, cfg app.Postgres) (MigrationService, error) {
	pgCfg := pg.PostgresConfig{
		ConnectionUrl: cfg.ConnURI,
	}

	sqlxDB, err := pg.InitSQLX(ctx, pgCfg)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error connecting to sqlxDB url:%s, err: %v", pgCfg.ConnectionUrl, err)
		return nil, err
	}

	databaseInstance, err := postgres.WithInstance(sqlxDB.DB, &postgres.Config{})
	if err != nil {
		logger.GetLogger(ctx).Errorf("go-migrate postgres drv init failed: %v", err)
		return nil, err
	}

	migrate, err := migrate.NewWithDatabaseInstance("file://migration/sql",
		migrateLogIdentifier, databaseInstance)
	if err != nil {
		logger.GetLogger(ctx).Errorf("migrate init failed %v", err)
		return nil, err
	}

	return migrationService{
		driver:  databaseInstance,
		migrate: migrate,
	}, nil
}

func (s migrationService) Up(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		logger.GetLogger(ctx).Error("Failed get current version err: ", err)
		return err
	}

	logger.GetLogger(ctx).Infof("Running migration from version: %d", currVersion)
	if err := s.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.GetLogger(ctx).Info("No Changes")
			return nil
		}
		logger.GetLogger(ctx).Error("Failed run migrate err: ", err)
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	logger.GetLogger(ctx).Info("Migration success, current version:", currVersion)
	return nil
}

func (s migrationService) Rollback(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		logger.GetLogger(ctx).Error("Failed get current version err: ", err)
		return err
	}

	logger.GetLogger(ctx).Infof("Rollingback 1 step from version: %d", currVersion)

	if err := s.migrate.Steps(-1); err != nil {
		logger.GetLogger(ctx).Errorf("Failed to rollback, err:%v", err)
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	logger.GetLogger(ctx).Infof("Rollback success, current version:%d", currVersion)
	return nil
}

func (s migrationService) Version(ctx context.Context) (int, bool, error) {
	currVersion, dirty, err := s.driver.Version()
	if err != nil {
		logger.GetLogger(ctx).Errorf("Failed to get version:%w", err)
		return 0, false, err
	}
	return currVersion, dirty, nil
}
