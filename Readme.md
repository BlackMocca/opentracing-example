## Example For Opentracing

![Work flow Tracing](https://github.com/BlackMocca/opentracing-example/blob/master/assets/github/workflow.png?raw=true)

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
curl -X GET http://127.0.0.1:3000/users/database
```

### Test Error 
```curl
curl -X GET http://127.0.0.1:3000/internal-error
curl -X GET http://127.0.0.1:3000/conflict
```

#### Appendix
Jaeger Sampling https://www.jaegertracing.io/docs/1.26/sampling/