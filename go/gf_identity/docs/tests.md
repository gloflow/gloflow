



for gf_identity tests the following ENV vars are used when run in github actions:
```
env:
    GF_LOG_LEVEL: debug
    GF_ALCHEMY_SERVICE_ACC__API_KEY: ${{ secrets.GF_ALCHEMY_SERVICE_ACC__API_KEY }}

    # currently these are used by web3 tests
    GF_TEST_MONGODB_HOST_PORT: mongo
    GF_TEST_SQL_HOST_PORT: postgres

    AUTH0_DOMAIN: ${{ secrets.GF_AUTH0_DOMAIN }}

```
