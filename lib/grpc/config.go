package grpc

import (
	"context"

	"github.com/usagifm/dating-app/lib/logger"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	GRPCPaymentServiceURL string
}

func ConnectToGRPC(ctx context.Context, config *GRPCConfig) (*grpc.ClientConn, error) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	GRPCConn, err := grpc.Dial(config.GRPCPaymentServiceURL, opts...)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to ping the grpc service, err:%v", err)
		return nil, err
	}

	return GRPCConn, nil

}
