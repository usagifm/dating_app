package v1

import (
	"context"

	atomicSQLX "github.com/usagifm/dating-app/lib/atomic/sqlx"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	packageRepo "github.com/usagifm/dating-app/src/repository/package"
	userRepo "github.com/usagifm/dating-app/src/repository/user"
	userMatchRepo "github.com/usagifm/dating-app/src/repository/user_match"
	userPackageRepo "github.com/usagifm/dating-app/src/repository/user_package"
	userPreferenceRepo "github.com/usagifm/dating-app/src/repository/user_preference"
	userSwipeRepo "github.com/usagifm/dating-app/src/repository/user_swipe"
	"github.com/usagifm/dating-app/src/v1/handler"
	authService "github.com/usagifm/dating-app/src/v1/service/auth"
	datingService "github.com/usagifm/dating-app/src/v1/service/dating"
)

type repositories struct {
	AtomicSessionProvider *atomicSQLX.SqlxAtomicSessionProvider
	User                  *userRepo.UserRepository
	UserPreference        *userPreferenceRepo.UserPreferenceRepository
	UserSwipe             *userSwipeRepo.UserSwipeRepository
	UserMatch             *userMatchRepo.UserMatchRepository
	Package               *packageRepo.PackageRepository
	UserPackage           *userPackageRepo.UserPackageRepository
}

type services struct {
	Auth   *authService.AuthService
	Dating *datingService.DatingService
}

type Handlers struct {
	DatingApp handler.DatingAppHandler
}

type Dependency struct {
	Repositories *repositories
	Services     *services
	Handlers     *Handlers
}

func initRepositories(ctx context.Context) *repositories {
	var r repositories
	var err error

	r.AtomicSessionProvider = atomicSQLX.NewSqlxAtomicSessionProvider(app.DB(), app.Tracer())

	logger.GetLogger(ctx).Printf("initializing user repo ")
	r.User, err = userRepo.InitUserRepository(ctx, app.DB(), app.Cache(), app.Config().Redis)
	if err != nil {
		logger.GetLogger(ctx).Fatal("init user repo err: ", err)
	}

	logger.GetLogger(ctx).Printf("initializing user preference repo ")
	r.UserPreference, err = userPreferenceRepo.InitUserPreferenceRepository(ctx, app.DB(), app.Cache(), app.Config().Redis)
	if err != nil {
		logger.GetLogger(ctx).Fatal("init user preference repo err: ", err)
	}

	logger.GetLogger(ctx).Printf("initializing user swipe repo ")
	r.UserSwipe, err = userSwipeRepo.InitUserSwipeRepository(ctx, app.DB(), app.Cache(), app.Config().Redis)
	if err != nil {
		logger.GetLogger(ctx).Fatal("init user preference repo err: ", err)
	}

	logger.GetLogger(ctx).Printf("initializing user match repo ")
	r.UserMatch, err = userMatchRepo.InitUserMatchRepository(ctx, app.DB(), app.Cache(), app.Config().Redis)
	if err != nil {
		logger.GetLogger(ctx).Fatal("init user match repo err: ", err)
	}

	logger.GetLogger(ctx).Printf("initializing package repo ")
	r.Package, err = packageRepo.InitPackageRepository(ctx, app.DB(), app.Cache(), app.Config().Redis)
	if err != nil {
		logger.GetLogger(ctx).Fatal("init package repo err: ", err)
	}

	logger.GetLogger(ctx).Printf("initializing user package repo ")
	r.UserPackage, err = userPackageRepo.InitUserPackageRepository(ctx, app.DB(), app.Cache(), app.Config().Redis)
	if err != nil {
		logger.GetLogger(ctx).Fatal("init user package repo err: ", err)
	}

	return &r
}

func initServices(ctx context.Context, r *repositories) *services {

	return &services{
		Auth:   authService.NewAuthService(r.AtomicSessionProvider, r.User, r.UserPreference),
		Dating: datingService.NewDatingService(r.AtomicSessionProvider, r.User, r.UserPreference, r.UserSwipe, r.UserMatch, r.Package, r.UserPackage),
	}
}

func initHandlers(ctx context.Context, s *services) *Handlers {
	var dep Handlers

	dep.DatingApp = *handler.NewDatingAppHandler(s.Auth, s.Dating)
	return &dep

}

func Dependencies(ctx context.Context) *Dependency {
	repositories := initRepositories(ctx)
	services := initServices(ctx, repositories)
	handlers := initHandlers(ctx, services)

	return &Dependency{
		Repositories: repositories,
		Services:     services,
		Handlers:     handlers,
	}
}
