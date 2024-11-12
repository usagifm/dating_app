package user_swipe

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/usagifm/dating-app/lib/helper"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

func (r *UserSwipeRepository) CreateUserSwipe(ctx context.Context, param entity.UserSwipe) (int, error) {
	_, span := app.Tracer().Start(ctx, "UserSwipeRepository/CreateUserSwipe")
	defer span.End()

	var created entity.UserSwipe
	namedStmt, err := r.getNamedStatement(ctx, createUserSwipe)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &created, param); err != nil {
		logger.GetLogger(ctx).Error("CreateUserSwipe  err: ", err)
		return 0, err
	}

	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(param.SwiperId) + ":daily-swipes"

	var unmarshalTodaySwipesUserId []int
	todaySwipesUserId := r.redisClient.Get(ctx, redisPath).Val()
	if todaySwipesUserId != "" {
		err := json.Unmarshal([]byte(todaySwipesUserId), &unmarshalTodaySwipesUserId)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return 0, err
		}
	}

	unmarshalTodaySwipesUserId = append(unmarshalTodaySwipesUserId, param.SwipedId)

	marshaledTodaySwipesUserId, errMarshal := json.Marshal(&unmarshalTodaySwipesUserId)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("CreateUserSwipe : marshaledTodaySwipesUserId json marshal err:%v\n", errMarshal)
	}

	_, end := helper.GetTodayStartAndEnd()
	now := time.Now()
	duration := end.Sub(now)
	errRedis := r.redisClient.Set(ctx, redisPath, marshaledTodaySwipesUserId, duration).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("CreateUserSwipe : marshaledTodaySwipesUserId redis set err:%v\n", errRedis)
	}

	return created.Id, nil
}

func (r *UserSwipeRepository) GetUserSwipe(ctx context.Context, userId int) ([]*entity.UserSwipe, error) {
	ctx, span := app.Tracer().Start(ctx, "UserSwipeRepository/GetUserSwipe")
	defer span.End()

	var userSwipes []*entity.UserSwipe
	err := r.masterStmts[getUserSwipe].SelectContext(ctx, &userSwipes, userId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.GetLogger(ctx).Error("GetUserSwipe err: ", err)
		}
		return userSwipes, err
	}

	return userSwipes, nil
}

func (r *UserSwipeRepository) GetMatchedUserSwipe(ctx context.Context, swipedId int, swiperId int) (*entity.UserSwipe, error) {
	ctx, span := app.Tracer().Start(ctx, "UserSwipeRepository/GetMatchedUserSwipe")
	defer span.End()

	var userSwipes entity.UserSwipe
	err := r.masterStmts[getMatchedUserSwipe].GetContext(ctx, &userSwipes, swipedId, swiperId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.GetLogger(ctx).Error("GetMatchedUserSwipe err: ", err)
			return nil, err
		}
		return nil, nil
	}

	return &userSwipes, nil
}

func (r *UserSwipeRepository) GetTodaySwipesUserId(ctx context.Context, userId int) ([]int, error) {
	ctx, span := app.Tracer().Start(ctx, "UserSwipeRepository/GetTodaySwipesUserId")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(userId) + ":daily-swipes"

	todaySwipesUserId := r.redisClient.Get(ctx, redisPath).Val()
	if todaySwipesUserId != "" {

		var unmarshalTodaySwipesUserId []int
		err := json.Unmarshal([]byte(todaySwipesUserId), &unmarshalTodaySwipesUserId)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return nil, err
		}

		return unmarshalTodaySwipesUserId, nil
	}

	start, end := helper.GetTodayStartAndEnd()

	var userSwipes []entity.UserSwipe
	err := r.masterStmts[getUserSwipePerToday].SelectContext(ctx, &userSwipes, userId, start, end)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.GetLogger(ctx).Error("GetTodaySwipesUserId err: ", err)
		}
		return nil, err
	}

	todaySwipesUserIds := []int{}

	for _, userSwipe := range userSwipes {
		todaySwipesUserIds = append(todaySwipesUserIds, userSwipe.SwipedId)
	}

	marshaledTodaySwipesUserId, errMarshal := json.Marshal(&todaySwipesUserIds)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("GetTodaySwipesUserId json marshal err:%v\n", errMarshal)
	}

	now := time.Now()
	duration := end.Sub(now)
	errRedis := r.redisClient.Set(ctx, redisPath, marshaledTodaySwipesUserId, duration).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("GetTodaySwipesUserId redis set err:%v\n", errRedis)
	}

	return todaySwipesUserIds, nil
}

// statement, err := p.getStatement(ctx, GetListByPackages)
// if err != nil {
// 	logger.GetLogger(ctx).Error("GetListActive Statement err: ", err)
// 	return nil, err
// }

// packageIdsStr := strings.Join(packageIds, ",")

// err = p.redis.WithCache(ctx, fmt.Sprintf(GetListByPackageProductRedisKey, packageIdsStr), &products, func() (interface{}, error) {
// 	var data []*entity.Product
// 	err = statement.SelectContext(ctx, &data, fmt.Sprintf("{%s}", packageIdsStr))
// 	return data, err
// })
// if err != nil {
// 	logger.GetLogger(ctx).Error("GetProductList err: ", err)
// 	return nil, err
// }

// func (r *UserSwipeRepository) get(ctx context.Context, txId string) (*entity.Transaction, error) {
// 	ctx, span := app.Tracer().Start(ctx, "TransactionRepository/GetTransactionByTxId")
// 	defer span.End()

// 	var transaction entity.Transaction
// 	err := r.masterStmts[getTransactionByTxId].GetContext(ctx, &transaction, txId)
// 	if err != nil {
// 		if !errors.Is(err, sql.ErrNoRows) {
// 			logger.GetLogger(ctx).Error("GetTransactionByTxId err: ", err)
// 		}
// 		return &transaction, err
// 	}

// 	return &transaction, nil
// }
