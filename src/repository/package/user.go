package user

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

func (r *PackageRepository) GetPackageById(ctx context.Context, packageId int) (*entity.Package, error) {
	ctx, span := app.Tracer().Start(ctx, "PackageRepository/GetPackageById")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":package:detail:" + strconv.Itoa(packageId)

	packageRedis := r.redisClient.Get(ctx, redisPath).Val()
	if packageRedis != "" {
		var unmarshalPackage entity.Package
		err := json.Unmarshal([]byte(packageRedis), &unmarshalPackage)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return nil, err
		}
		return &unmarshalPackage, nil
	}

	var packageById entity.Package
	err := r.masterStmts[getPackageById].GetContext(ctx, &packageById, packageId)
	if err != nil {
		logger.GetLogger(ctx).Error("GetPackageById err: ", err)
		return nil, err
	}

	marshaledPackage, errMarshal := json.Marshal(&packageById)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("GetPackageById json marshal err:%v\n", errMarshal)
	}

	errRedis := r.redisClient.Set(ctx, redisPath, marshaledPackage, r.redisConfig.InvalidateTime).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("GetPackageById redis set err:%v\n", errRedis)
	}

	return &packageById, nil
}

func (r *PackageRepository) GetPackages(ctx context.Context) ([]*entity.Package, error) {
	ctx, span := app.Tracer().Start(ctx, "PackageRepository/GetPackages")
	defer span.End()

	redisPath := r.redisConfig.DBPrefix + ":package:list"

	packagesRedis := r.redisClient.Get(ctx, redisPath).Val()
	if packagesRedis != "" {
		var unmarshalPackages []*entity.Package
		err := json.Unmarshal([]byte(packagesRedis), &unmarshalPackages)
		if err != nil {
			logger.GetLogger(ctx).Error("Unmarshal error: ", err)
			return nil, err
		}
		return unmarshalPackages, nil
	}

	var packages []*entity.Package
	err := r.masterStmts[getPackages].SelectContext(ctx, &packages)
	if err != nil {
		logger.GetLogger(ctx).Error("GetPackageById err: ", err)
		return nil, err
	}

	marshaledPackages, errMarshal := json.Marshal(&packages)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("GetPackageById json marshal err:%v\n", errMarshal)
	}

	errRedis := r.redisClient.Set(ctx, redisPath, marshaledPackages, r.redisConfig.InvalidateTime).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("GetPackageById redis set err:%v\n", errRedis)
	}

	return packages, nil
}
