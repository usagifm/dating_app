package user_match

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

func (r *UserMatchRepository) CreateUserMatch(ctx context.Context, param entity.UserMatch) (int, error) {
	_, span := app.Tracer().Start(ctx, "UserMatchRepository/CreateUserMatch")
	defer span.End()

	var created entity.UserMatch
	namedStmt, err := r.getNamedStatement(ctx, createUserMatch)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return 0, err
	}

	if errCreated := namedStmt.GetContext(ctx, &created, param); errCreated != nil {
		logger.GetLogger(ctx).Error("CreateUserMatch  err: ", errCreated)
		return 0, errCreated
	}

	return created.Id, nil
}

func (r *UserMatchRepository) GetUserMatchesUserId(ctx context.Context, userId int) ([]int, error) {
	ctx, span := app.Tracer().Start(ctx, "UserMatchRepository/GetUserMatches")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(userId) + ":matches"

	userMatchesUserIdRedis := r.redisClient.Get(ctx, redisPath).Val()
	if userMatchesUserIdRedis != "" {
		var unmarshalUserMatchesUserIdRedis []int
		err := json.Unmarshal([]byte(userMatchesUserIdRedis), &unmarshalUserMatchesUserIdRedis)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return nil, err
		}

		return unmarshalUserMatchesUserIdRedis, nil
	}

	var userMatches []*entity.UserMatch
	err := r.masterStmts[getUserMatches].SelectContext(ctx, &userMatches, userId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.GetLogger(ctx).Error("GetUserSwipe err: ", err)
		}
		return nil, err
	}

	userMatchesUserIds := []int{}

	for _, userMatch := range userMatches {

		if userMatch.UserId1 != userId {
			userMatchesUserIds = append(userMatchesUserIds, userMatch.UserId1)
		} else if userMatch.UserId2 != userId {
			userMatchesUserIds = append(userMatchesUserIds, userMatch.UserId2)
		}

	}

	marshaledUserMatchUserId, errMarshal := json.Marshal(&userMatchesUserIds)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("GetUserMatches json marshal err:%v\n", errMarshal)
	}

	errRedis := r.redisClient.Set(ctx, redisPath, marshaledUserMatchUserId, r.redisConfig.InvalidateTime).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("GetUserMatches redis set err:%v\n", errRedis)
	}

	return userMatchesUserIds, nil
}

func (r *UserMatchRepository) InvalidateUserMatchesRedis(ctx context.Context, userId int) error {
	ctx, span := app.Tracer().Start(ctx, "UserMatchRepository/InvalidateUserMatched")
	defer span.End()
	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(userId) + ":matches"

	errRedis := r.redisClient.Del(ctx, redisPath).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("InvalidateUserMatched redis delete err:%v\n", errRedis)
		return errRedis
	}

	return nil
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
