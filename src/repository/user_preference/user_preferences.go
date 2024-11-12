package user_preference

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

func (r *UserPreferenceRepository) CreateNewUserPreference(ctx context.Context, param entity.UserPreference) (int, error) {
	_, span := app.Tracer().Start(ctx, "UserPreferenceRepository/CreateNewUserPreference")
	defer span.End()

	var newUserPreference entity.UserPreference
	namedStmt, err := r.getNamedStatement(ctx, createNewUserPreference)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &newUserPreference, param); err != nil {
		logger.GetLogger(ctx).Error("CreateNewUserPreference  err: ", err)
		return 0, err
	}

	return newUserPreference.Id, nil
}

func (r *UserPreferenceRepository) UpdateUserPreference(ctx context.Context, param entity.UserPreference) error {
	_, span := app.Tracer().Start(ctx, "UserPreferenceRepository/UpdateUserPreference")
	defer span.End()

	namedStmt, err := r.getNamedStatement(ctx, updateUserPreference)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return err
	}

	_, err = namedStmt.ExecContext(ctx, param)
	if err != nil {
		logger.GetLogger(ctx).Errorf("UpdateUserPreference id:%d, err: %v", param.Id, err)
		return err
	}

	return nil
}

func (r *UserPreferenceRepository) GetUserPreference(ctx context.Context, userId int) (*entity.UserPreference, error) {
	ctx, span := app.Tracer().Start(ctx, "UserPreferenceRepository/GetUserPreference")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":user:" + strconv.Itoa(userId) + ":preference"

	preferenceRedis := r.redisClient.Get(ctx, redisPath).Val()
	if preferenceRedis != "" {
		var unmarshalUserPreference entity.UserPreference
		err := json.Unmarshal([]byte(preferenceRedis), &unmarshalUserPreference)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return nil, err
		}
		return &unmarshalUserPreference, nil
	}

	var userPreference entity.UserPreference
	err := r.masterStmts[getUserPreference].GetContext(ctx, &userPreference, userId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.GetLogger(ctx).Error("GetUserPreference err: ", err)
		}
		return &userPreference, err
	}

	marshaledUserPreference, errMarshal := json.Marshal(&userPreference)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("GetUserPreference json marshal err:%v\n", errMarshal)
	}

	errRedis := r.redisClient.Set(ctx, redisPath, marshaledUserPreference, r.redisConfig.InvalidateTime).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("GetUserPreference redis set err:%v\n", errRedis)
	}

	return &userPreference, nil
}

func (r *UserPreferenceRepository) GetAnotherUserPreferenceByPreference(ctx context.Context, notIncludedUserId []int, minAge int, maxAge int, preferredGender string) ([]*entity.User, error) {
	ctx, span := app.Tracer().Start(ctx, "UserPreferenceRepository/GetAnotherUserPreferenceByPreference")
	defer span.End()

	// Convert each integer to a string
	strIds := make([]string, len(notIncludedUserId))
	for i, id := range notIncludedUserId {
		strIds[i] = strconv.Itoa(id)
	}
	stringIds := strings.Join(strIds, ",")

	query := `SELECT tu.id, tu.is_verified, tu.name, tu.gender, tu.age, tu.bio, tu.photo_url 
			  FROM tbl_user_preferences tup 
			  JOIN tbl_users tu ON tu.id = tup.user_id 
			  WHERE tup.user_id NOT IN (` + stringIds + `)`

	if minAge > 0 {
		query = query + " AND tu.age > " + strconv.Itoa(minAge)
	}

	if maxAge > 0 {
		query = query + " AND tu.age < " + strconv.Itoa(maxAge)
	}

	if preferredGender != "" && preferredGender != "both" {
		query = query + " AND tu.gender = '" + preferredGender + "' "
	}

	logger.GetLogger(ctx).Info("query err: ", query)

	var userPreference []*entity.User
	err := r.db.SelectContext(ctx, &userPreference, query)
	if err != nil {

		logger.GetLogger(ctx).Error("GetAnotherUserPreferenceByPreference err: ", err)

		return userPreference, err
	}

	return userPreference, nil
}
