package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

func (r *UserPackageRepository) GetUserPackage(ctx context.Context, userId int) (*entity.UserPackage, error) {
	ctx, span := app.Tracer().Start(ctx, "UserPackageRepository/GetUserPackage")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(userId) + ":package:detail"
	userPackageRedis := r.redisClient.Get(ctx, redisPath).Val()
	if userPackageRedis != "" {
		var unmarshalUserPackage entity.UserPackage
		err := json.Unmarshal([]byte(userPackageRedis), &unmarshalUserPackage)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return nil, err
		}
		return &unmarshalUserPackage, nil
	}

	var userPackage entity.UserPackage
	err := r.masterStmts[getUserPackage].GetContext(ctx, &userPackage)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.GetLogger(ctx).Error("GetUserPackage err: ", err)
			return nil, err
		}

		marshaledUserPackage, errMarshal := json.Marshal(&userPackage)
		if errMarshal != nil {
			logger.GetLogger(ctx).Errorf("GetUserPackage json marshal err:%v\n", errMarshal)
		}

		errRedis := r.redisClient.Set(ctx, redisPath, marshaledUserPackage, r.redisConfig.InvalidateTime).Err()
		if errRedis != nil {
			logger.GetLogger(ctx).Errorf("GetUserPackage redis set err:%v\n", errRedis)
		}

		return &userPackage, nil

	}

	marshaledUserPackage, errMarshal := json.Marshal(&userPackage)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("GetUserPackage json marshal err:%v\n", errMarshal)
	}

	time := time.Now()

	duration := userPackage.ValidDate.Sub(time)

	errRedis := r.redisClient.Set(ctx, redisPath, marshaledUserPackage, duration).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("GetUserPackage redis set err:%v\n", errRedis)
	}

	return &userPackage, nil
}

func (r *UserPackageRepository) CreateOrUpdateUserPackage(ctx context.Context, param entity.UserPackage) (int, error) {
	_, span := app.Tracer().Start(ctx, "UserPackageRepository/CreateOrUpdateUserPackage")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(param.UserId) + ":package:detail"

	var userPackage entity.UserPackage
	namedStmt, err := r.getNamedStatement(ctx, createOrUpdateUserPackage)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &userPackage, param); err != nil {
		logger.GetLogger(ctx).Error("CreateOrUpdateUserPackage  err: ", err)
		return 0, err
	}

	errRedis := r.redisClient.Del(ctx, redisPath).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("CreateOrUpdateUserPackage redis delete err:%v\n", errRedis)
	}

	return userPackage.Id, nil
}
