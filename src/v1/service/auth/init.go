package auth

import (
	"context"

	"github.com/usagifm/dating-app/lib/atomic"
	"github.com/usagifm/dating-app/src/entity"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, param entity.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserProfile(ctx context.Context, userId int) (*entity.User, error)
}

type UserReferenceRepository interface {
	CreateNewUserPreference(ctx context.Context, param entity.UserPreference) (int, error)
}

type AuthService struct {
	atomicSession  atomic.AtomicSessionProvider
	rUser          UserRepository
	rUserReference UserReferenceRepository
}

func NewAuthService(aSession atomic.AtomicSessionProvider, rUser UserRepository, rUserReference UserReferenceRepository) *AuthService {
	return &AuthService{
		atomicSession:  aSession,
		rUser:          rUser,
		rUserReference: rUserReference,
	}
}
