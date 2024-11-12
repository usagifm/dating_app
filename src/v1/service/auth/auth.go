package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/usagifm/dating-app/lib/helper"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
	"github.com/usagifm/dating-app/src/middleware/auth"
	"github.com/usagifm/dating-app/src/v1/contract"
)

func (s *AuthService) SignUp(ctx context.Context, param contract.SignUpRequest) error {
	ctx, span := app.Tracer().Start(ctx, "AuthService/SignUp")
	defer span.End()

	user, err := s.rUser.GetUserByEmail(ctx, param.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		logger.GetLogger(ctx).Error("GetUserByEmail err: ", err)
		return err
	}

	if user != nil {
		logger.GetLogger(ctx).Error("user already exists")
		return errors.New("user already exists")
	}

	password := param.Password
	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		logger.GetLogger(ctx).Error("SignUp err: ", err)
		return err
	}

	newUser := entity.User{
		IsVerified: false,
		Name:       param.Name,
		Gender:     param.Gender,
		Email:      param.Email,
		Password:   hashedPassword,
		Age:        param.Age,
		Bio:        param.Bio,
	}

	newUserId, err := s.rUser.CreateNewUser(ctx, newUser)
	if err != nil {
		logger.GetLogger(ctx).Error("CreateNewUser err: ", err)
		return err
	}

	newUserPreference := entity.UserPreference{
		UserId:          newUserId,
		PreferredGender: param.PreferredGender,
		MinAge:          param.MinAge,
		MaxAge:          param.MaxAge,
	}

	_, err = s.rUserReference.CreateNewUserPreference(ctx, newUserPreference)
	if err != nil {
		logger.GetLogger(ctx).Error("CreateNewUserPreference err: ", err)
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, param contract.LoginRequest) (string, error) {
	ctx, span := app.Tracer().Start(ctx, "AuthService/LoginRequest")
	defer span.End()

	user, err := s.rUser.GetUserByEmail(ctx, param.Email)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserByEmail err: ", err)
		return "", err
	}

	isPasswordMatch := helper.ComparePassword(user.Password, param.Password)
	if !isPasswordMatch {
		logger.GetLogger(ctx).Error("Wrong password")
		return "", errors.New("wrong password")
	}

	user.Password = ""
	// Create a new token with HS256 signing method and custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "dating-app",
		"sub":  user.Id,
		"user": *user,
	})

	jwtSecret := []byte(app.Config().JWTSecret)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) GetProfile(ctx context.Context) (*entity.User, error) {
	ctx, span := app.Tracer().Start(ctx, "AuthService/GetProfileRequest")
	defer span.End()

	user := auth.GetUser(ctx)

	user, err := s.rUser.GetUserProfile(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserByEmail err: ", err)
		return user, err
	}

	return user, nil

}
