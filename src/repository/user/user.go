package user

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

func (r *UserRepository) CreateNewUser(ctx context.Context, param entity.User) (int, error) {
	_, span := app.Tracer().Start(ctx, "UserRepository/CreateNewUser")
	defer span.End()

	var newUser entity.User
	namedStmt, err := r.getNamedStatement(ctx, createNewUser)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &newUser, param); err != nil {
		logger.GetLogger(ctx).Error("CreateSubscriptionLog  err: ", err)
		return 0, err
	}

	return newUser.Id, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, param entity.User) error {
	_, span := app.Tracer().Start(ctx, "UserRepository/UpdateUser")
	defer span.End()

	namedStmt, err := r.getNamedStatement(ctx, updateUser)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return err
	}

	_, err = namedStmt.ExecContext(ctx, param)
	if err != nil {
		logger.GetLogger(ctx).Errorf("UpdateUser id:%d, err: %v", param.Id, err)
		return err
	}

	return nil
}

func (r *UserRepository) GetUserProfile(ctx context.Context, userId int) (*entity.User, error) {
	ctx, span := app.Tracer().Start(ctx, "UserRepository/GetUserProfile")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(userId) + ":profile"

	profileRedis := r.redisClient.Get(ctx, redisPath).Val()
	if profileRedis != "" {
		var unmarshalUserProfile entity.User
		err := json.Unmarshal([]byte(profileRedis), &unmarshalUserProfile)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return nil, err
		}
		return &unmarshalUserProfile, nil
	}

	var userProfile entity.User
	err := r.masterStmts[getUserProfile].GetContext(ctx, &userProfile, userId)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserProfile err: ", err)
		return nil, err
	}

	marshaledUserProfile, errMarshal := json.Marshal(&userProfile)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("GetUserProfile json marshal err:%v\n", errMarshal)
	}

	errRedis := r.redisClient.Set(ctx, redisPath, marshaledUserProfile, r.redisConfig.InvalidateTime).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("GetUserProfile redis set err:%v\n", errRedis)
	}

	return &userProfile, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, span := app.Tracer().Start(ctx, "UserRepository/GetUserByEmail")
	defer span.End()

	var userProfile entity.User
	err := r.masterStmts[getUserByEmail].GetContext(ctx, &userProfile, email)
	if err != nil {
		logger.GetLogger(ctx).Error("GetUserByEmail err: ", err)
		return nil, err
	}

	return &userProfile, nil
}
