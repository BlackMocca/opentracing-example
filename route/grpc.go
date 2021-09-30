package route

import (
	"github.com/Blackmocca/opentracing-example/proto/proto_models"
	"google.golang.org/grpc"
)

type GRPCRoute struct {
	server *grpc.Server
}

func NewGRPCRoute(server *grpc.Server) *GRPCRoute {
	return &GRPCRoute{server}
}

func (g *GRPCRoute) RegisterUserHandler(handler proto_models.UserServer) {
	proto_models.RegisterUserServer(g.server, handler)
}
