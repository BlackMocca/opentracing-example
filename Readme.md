## Example For Opentracing

![Work flow Tracing](https://github.com/BlackMocca/opentracing-example/blob/master/assets/github/workflow.png?raw=true)

#### running jaeger with docker 
```
docker run -d --name jaeger-ui \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 16686:16686 \
  jaegertracing/all-in-one:1.7@sha256:146de3a8c00e7ce536734d96627a71047d82481c4862ef79560a72dba1b4099a
```

#### Running On docker compose
```
docker-compose up
```

### Test By Curl
Create Migration and Seed data on database postgres
```shellscript
make install-migration
make app.migration.up db_url="postgres://postgres:postgres@psql_db:5432/app_example?sslmode=disable" path=migrations/database/postgres
make app.migration.seed db_url="postgres://postgres:postgres@psql_db:5432/app_example?sslmode=disable" path=migrations/database/postgres/seed/master
make app.migration.seed db_url="postgres://postgres:postgres@psql_db:5432/app_example?sslmode=disable" path=migrations/database/postgres/seed/story-001
```

### Test Get User Without Database
```curl
curl -X GET http://127.0.0.1:3000/users
```

### Test Get User With Database
```curl
curl -X GET http://127.0.0.1:3000/users/database/psql
```

### Test Error 
```curl
curl -X GET http://127.0.0.1:3000/internal-error
curl -X GET http://127.0.0.1:3000/conflict
```

#### Appendix
Jaeger Sampling https://www.jaegertracing.io/docs/1.26/sampling/