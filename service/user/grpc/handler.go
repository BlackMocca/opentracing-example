package user

import (
	"context"
	"time"

	"github.com/Blackmocca/opentracing-example/proto/proto_models"
	"github.com/Blackmocca/opentracing-example/service/user"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type grpcHandler struct {
	userUs user.UserUsecase
}

func NewGRPCHandler(userUs user.UserUsecase) proto_models.UserServer {
	return &grpcHandler{
		userUs: userUs,
	}
}

func (g grpcHandler) FetchUserAddress(ctx context.Context, req *proto_models.FetchUserAddressRequest) (*proto_models.FetchUserAddressResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.LogFields(
		log.String("id", req.GetId()),
	)

	userId := req.GetId()
	time.Sleep(time.Duration(2 * time.Second))

	address := map[string]string{
		"1": "02/223 building ABC Bangkok",
		"2": "11/223 building XYZ Bangkok",
	}

	resp := &proto_models.FetchUserAddressResponse{
		Address: &proto_models.UserAddress{
			Address: address[userId],
		},
	}
	return resp, nil
}
