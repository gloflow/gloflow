

run DBs needed for tests (mongo, postgresql)
```
docker run -p 27017:27017 mongo

docker run --name gf-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=gf_tests -e POSTGRES_USER=gf -d postgres

```