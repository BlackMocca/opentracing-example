
### Running jaeger

```script
docker run \
  --rm \
  --name jaeger-ui \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 16686:16686 \
  -d \
  jaegertracing/all-in-one:1.7@sha256:146de3a8c00e7ce536734d96627a71047d82481c4862ef79560a72dba1b4099a \
  --log-level=debug
```

OPEN Jaeger UI #### http://127.0.0.1:16686

### Running Application

Running application
```
go mod tidy
go run main.go
```

#### Running On docker compose
```
docker-compose up
```

### Test By Curl
```curl
curl -X GET http://127.0.0.1:3000/users
```

![Work flow Tracing](https://github.com/BlackMocca/opentracing-example/blob/master/assets/github/workflow.png?raw=true)


#### Appendix
Jaeger Sampling https://www.jaegertracing.io/docs/1.26/sampling/