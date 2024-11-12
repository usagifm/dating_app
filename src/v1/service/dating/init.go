package dating

import (
	"context"

	"github.com/usagifm/dating-app/lib/atomic"
	"github.com/usagifm/dating-app/src/entity"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, param entity.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserProfile(ctx context.Context, userId int) (*entity.User, error)
	UpdateUser(ctx context.Context, param entity.User) error
}

type UserPreferenceRepository interface {
	UpdateUserPreference(ctx context.Context, param entity.UserPreference) error
	GetUserPreference(ctx context.Context, userId int) (*entity.UserPreference, error)
	GetAnotherUserPreferenceByPreference(ctx context.Context, notIncludedUserId []int, minAge int, maxAge int, preferredGender string) ([]*entity.User, error)
}

type UserSwipeRepository interface {
	CreateUserSwipe(ctx context.Context, param entity.UserSwipe) (int, error)
	GetUserSwipe(ctx context.Context, userId int) ([]*entity.UserSwipe, error)
	GetMatchedUserSwipe(ctx context.Context, swipedId int, swiperId int) (*entity.UserSwipe, error)
	GetTodaySwipesUserId(ctx context.Context, userId int) ([]int, error)
}

type UserMatchRepository interface {
	InvalidateUserMatchesRedis(ctx context.Context, userId int) error
	CreateUserMatch(ctx context.Context, param entity.UserMatch) (int, error)
	GetUserMatchesUserId(ctx context.Context, userId int) ([]int, error)
}

type PackageRepository interface {
	GetPackageById(ctx context.Context, packageId int) (*entity.Package, error)
	GetPackages(ctx context.Context) ([]*entity.Package, error)
}

type UserPackageRepository interface {
	GetUserPackage(ctx context.Context, userId int) (*entity.UserPackage, error)
	CreateOrUpdateUserPackage(ctx context.Context, param entity.UserPackage) (int, error)
}

type DatingService struct {
	atomicSession   atomic.AtomicSessionProvider
	rUser           UserRepository
	rUserPreference UserPreferenceRepository
	rUserSwipe      UserSwipeRepository
	rUserMatch      UserMatchRepository
	rPackage        PackageRepository
	rUserPackage    UserPackageRepository
}

func NewDatingService(aSession atomic.AtomicSessionProvider, rUser UserRepository, rUserPreference UserPreferenceRepository, rUserSwipe UserSwipeRepository, rUserMatch UserMatchRepository, rPackage PackageRepository, rUserPackage UserPackageRepository) *DatingService {
	return &DatingService{
		atomicSession:   aSession,
		rUser:           rUser,
		rUserPreference: rUserPreference,
		rUserSwipe:      rUserSwipe,
		rUserMatch:      rUserMatch,
		rPackage:        rPackage,
		rUserPackage:    rUserPackage,
	}
}
