package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Blackmocca/opentracing-example/models"
	"github.com/Blackmocca/opentracing-example/proto/proto_models"
	"github.com/go-resty/resty/v2"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/spf13/cast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getCover(ctx context.Context, users []*models.User) error {
	span := opentracing.SpanFromContext(ctx)
	client := resty.New()
	client.Debug = true
	host := "http://127.0.0.1:3000"

	for index, _ := range users {
		id := users[index].Id
		url := fmt.Sprintf("%s/users/%s/cover", host, id)

		req := client.R().EnableTrace()

		ext.SpanKindRPCClient.Set(span)
		// ext.HTTPUrl.Set(span, url)
		// ext.HTTPMethod.Set(span, "GET")
		span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)

		resp, err := req.SetContext(ctx).Get(url)
		if err != nil {
			ext.LogError(span, err)
			panic(err)
		}

		if resp.StatusCode() == http.StatusOK {
			var m = map[string]interface{}{}
			json.Unmarshal(resp.Body(), &m)

			users[index].Cover = cast.ToString(m["cover"])
		}
	}
	return nil
}

func getAddress(ctx context.Context, users []*models.User) error {
	address := "127.0.0.1:3100"

	for index, _ := range users {
		id := users[index].Id

		conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(
				otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
			),
			grpc.WithStreamInterceptor(
				otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer()),
			),
		)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := proto_models.NewUserClient(conn)

		req := &proto_models.FetchUserAddressRequest{
			Id: id,
		}
		resp, err := client.FetchUserAddress(ctx, req)
		if err != nil {
			var grpcCode = status.Code(err)
			var status = http.StatusInternalServerError

			switch grpcCode {
			case codes.InvalidArgument:
				status = http.StatusBadRequest
			case codes.NotFound:
				status = http.StatusNotFound
			case codes.Unavailable:
				status = http.StatusServiceUnavailable
			case codes.Unauthenticated:
				status = http.StatusNetworkAuthenticationRequired
			case codes.Unimplemented:
				status = http.StatusNotImplemented
			}
			log.Println(status, err)
			return err
		}

		if resp != nil {
			users[index].Address = resp.GetAddress().Address
		}
	}
	return nil
}
