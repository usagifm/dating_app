package handler

import (
	"context"

	"github.com/usagifm/dating-app/src/entity"
	"github.com/usagifm/dating-app/src/v1/contract"
)

type AuthService interface {
	Login(ctx context.Context, param contract.LoginRequest) (string, error)
	SignUp(ctx context.Context, param contract.SignUpRequest) error
	GetProfile(ctx context.Context) (*entity.User, error)
}

type DatingService interface {
	GetUserPreference(ctx context.Context) (*entity.UserPreference, error)
	UpdateUserPreference(ctx context.Context, param contract.UpdateUserPreferenceRequest) error
	GetProfilesByPreference(ctx context.Context) ([]*entity.User, error)
	Swipe(ctx context.Context, param contract.SwipeRequest) (bool, error)
	GetUserMatches(ctx context.Context) ([]*entity.User, error)
	GetPackages(ctx context.Context) ([]*entity.Package, error)
	BuyPackage(ctx context.Context, param contract.BuyPackageRequest) error
}

type DatingAppHandler struct {
	sAuth   AuthService
	sDating DatingService
}

func NewDatingAppHandler(sAuth AuthService, sDating DatingService) *DatingAppHandler {
	return &DatingAppHandler{
		sAuth:   sAuth,
		sDating: sDating,
	}
}
