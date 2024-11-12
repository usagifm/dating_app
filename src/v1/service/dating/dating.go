package dating

import (
	"context"
	"errors"
	"time"

	"github.com/usagifm/dating-app/lib/helper"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
	"github.com/usagifm/dating-app/src/middleware/auth"
	"github.com/usagifm/dating-app/src/v1/contract"
)

func (s *DatingService) GetUserPreference(ctx context.Context) (*entity.UserPreference, error) {
	ctx, span := app.Tracer().Start(ctx, "DatingService/GetUserPreference")
	defer span.End()

	user := auth.GetUser(ctx)

	userPreference, err := s.rUserPreference.GetUserPreference(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserPreference err: ", err)
		return nil, err
	}

	return userPreference, nil
}

func (s *DatingService) UpdateUserPreference(ctx context.Context, param contract.UpdateUserPreferenceRequest) error {
	ctx, span := app.Tracer().Start(ctx, "DatingService/UpdateUserPreference")
	defer span.End()

	user := auth.GetUser(ctx)

	userPreference, err := s.rUserPreference.GetUserPreference(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserPreference err: ", err)
		return err
	}

	userPreference.MaxAge = param.MaxAge
	userPreference.MinAge = param.MinAge
	userPreference.PreferredGender = param.PreferredGender

	err = s.rUserPreference.UpdateUserPreference(ctx, *userPreference)
	if err != nil {
		logger.GetLogger(ctx).Error("UpdateUserPreference err: ", err)
		return err
	}

	return nil
}

func (s *DatingService) GetProfilesByPreference(ctx context.Context) ([]*entity.User, error) {
	ctx, span := app.Tracer().Start(ctx, "DatingService/GetProfilesByPreference")
	defer span.End()

	user := auth.GetUser(ctx)

	userPreference, err := s.rUserPreference.GetUserPreference(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserPreference err: ", err)
		return nil, err
	}

	notIncludedUserId := []int{user.Id}

	todaySwipedUserId, err := s.rUserSwipe.GetTodaySwipesUserId(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetTodaySwipesUserId err: ", err)
		return nil, err
	}
	notIncludedUserId = append(notIncludedUserId, todaySwipedUserId...)

	matchedUserId, err := s.rUserMatch.GetUserMatchesUserId(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserMatchesUserId err: ", err)
		return nil, err
	}
	notIncludedUserId = append(notIncludedUserId, matchedUserId...)

	userProfiles, err := s.rUserPreference.GetAnotherUserPreferenceByPreference(ctx, notIncludedUserId, userPreference.MinAge, userPreference.MaxAge, userPreference.PreferredGender)
	if err != nil {
		logger.GetLogger(ctx).Error("GetAnotherUserPreferenceByPreference err: ", err)
		return nil, err
	}

	return userProfiles, nil
}

func (s *DatingService) Swipe(ctx context.Context, param contract.SwipeRequest) (bool, error) {
	ctx, span := app.Tracer().Start(ctx, "DatingService/Swipe")
	defer span.End()

	user := auth.GetUser(ctx)
	var isHavePackage bool
	userPackage, err := s.rUserPackage.GetUserPackage(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserPackage err: ", err)
		return false, err
	}

	if userPackage.Id == 0 {
		isHavePackage = false
	} else {
		now := time.Now()
		if userPackage.ValidDate.Unix() > now.Unix() {
			isHavePackage = true
		}
	}

	todaySwipedUserId, err := s.rUserSwipe.GetTodaySwipesUserId(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetTodaySwipesUserId err: ", err)
		return false, err
	}

	if !isHavePackage {
		if len(todaySwipedUserId) > 9 {
			logger.GetLogger(ctx).Error("already reach daily quota ")
			return false, errors.New("already reach daily quota ")
		}
	}

	isExist := helper.Contains(todaySwipedUserId, param.SwipedId)
	if isExist {
		logger.GetLogger(ctx).Error("already swiped this user today")
		return false, errors.New("already swiped this user today")
	}

	matchedUserId, err := s.rUserMatch.GetUserMatchesUserId(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserMatchesUserId err: ", err)
		return false, err
	}

	isExist = helper.Contains(matchedUserId, param.SwipedId)
	if isExist {
		logger.GetLogger(ctx).Error("already matched with this user")
		return false, err
	}

	userSwipe := entity.UserSwipe{
		SwiperId:  user.Id,
		SwipedId:  param.SwipedId,
		SwipeType: param.SwipeType,
	}

	_, errCreate := s.rUserSwipe.CreateUserSwipe(ctx, userSwipe)
	if errCreate != nil {
		logger.GetLogger(ctx).Error("CreateUserSwipe err: ", errCreate)
		return false, errCreate
	}

	isMatched, err := s.rUserSwipe.GetMatchedUserSwipe(ctx, user.Id, param.SwipedId)
	if err != nil {
		logger.GetLogger(ctx).Error("GetMatchedUserSwipe err: ", err)
		return false, err
	}

	if isMatched != nil {

		userMatch := entity.UserMatch{
			UserId1: user.Id,
			UserId2: param.SwipedId,
		}
		_, err := s.rUserMatch.CreateUserMatch(ctx, userMatch)
		if err != nil {
			logger.GetLogger(ctx).Error("CreateUserMatch err: ", err)
			return false, err
		}

		_ = s.rUserMatch.InvalidateUserMatchesRedis(ctx, user.Id)
		go s.rUserMatch.GetUserMatchesUserId(ctx, user.Id)

		_ = s.rUserMatch.InvalidateUserMatchesRedis(ctx, param.SwipedId)
		go s.rUserMatch.GetUserMatchesUserId(ctx, param.SwipedId)

		return true, nil
	}

	return false, nil
}

func (s *DatingService) GetUserMatches(ctx context.Context) ([]*entity.User, error) {
	ctx, span := app.Tracer().Start(ctx, "DatingService/GetUserPreference")
	defer span.End()

	user := auth.GetUser(ctx)
	matchesUser := []*entity.User{}

	matchesUserId, err := s.rUserMatch.GetUserMatchesUserId(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserMatchesUserId err: ", err)
		return nil, err
	}

	for _, userId := range matchesUserId {

		user, err := s.rUser.GetUserProfile(ctx, userId)
		if err != nil {
			logger.GetLogger(ctx).Error("GetUserProfile err: ", err)
			return nil, err
		}

		matchesUser = append(matchesUser, user)

	}

	return matchesUser, nil
}

func (s *DatingService) GetPackages(ctx context.Context) ([]*entity.Package, error) {
	ctx, span := app.Tracer().Start(ctx, "DatingService/GetPackages")
	defer span.End()

	packages, err := s.rPackage.GetPackages(ctx)
	if err != nil {
		logger.GetLogger(ctx).Error("GetPackages err: ", err)
		return nil, err
	}

	return packages, nil
}

func (s *DatingService) BuyPackage(ctx context.Context, param contract.BuyPackageRequest) error {
	ctx, span := app.Tracer().Start(ctx, "DatingService/GetPackages")
	defer span.End()

	user := auth.GetUser(ctx)

	packageDetail, err := s.rPackage.GetPackageById(ctx, param.PackageId)
	if err != nil {
		logger.GetLogger(ctx).Error("GetPackages err: ", err)
		return err
	}

	userPackage, err := s.rUserPackage.GetUserPackage(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserPackage err: ", err)
		return err
	}

	logger.GetLogger(ctx).Info("validDate err: ", userPackage)

	now := time.Now()
	var validDate time.Time
	if userPackage.Id == 0 {
		validDate = now.Add(time.Duration(packageDetail.Periode*24) * time.Hour)
		logger.GetLogger(ctx).Info("validDate dalem err: ", validDate)
	} else {
		now := time.Now()
		if userPackage.ValidDate.Unix() > now.Unix() {
			validDate = userPackage.ValidDate.Add(time.Duration(packageDetail.Periode*24) * time.Hour)
		}
	}

	logger.GetLogger(ctx).Info("validDate err: ", validDate)

	newUserPackage := entity.UserPackage{
		UserId:    user.Id,
		PackageId: packageDetail.Id,
		ValidDate: validDate,
	}
	_, errCreateOrUpdate := s.rUserPackage.CreateOrUpdateUserPackage(ctx, newUserPackage)
	if errCreateOrUpdate != nil {
		logger.GetLogger(ctx).Error("CreateOrUpdateUserPackage err: ", errCreateOrUpdate)
		return err
	}

	userProfile, err := s.rUser.GetUserProfile(ctx, user.Id)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserByEmail err: ", err)
		return err
	}

	userProfile.IsVerified = true
	errUpdateUser := s.rUser.UpdateUser(ctx, *userProfile)
	if errUpdateUser != nil {
		logger.GetLogger(ctx).Error("errUpdateUser err: ", errUpdateUser)
		return err
	}

	return nil
}
