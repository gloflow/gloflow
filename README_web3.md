# gloflow-web3-monitor
monitoring of the Web3 networks (Ethereum, IPFS, etc.)





The goal is to create a monitoring tool for the Ethereum mainnet.  
Main gf-web3-monitor server manages the spawning of modified GF go-ethereum nodes that connect into the mainnet. These nodes are modified to contain custom instrumentation, and are sending data to a common queue system that is shared with the gf-web3-monitor server. The main API server is a consumer, and the GF go-ethereum node is a producer.  

pre-built container is available in a Dockerhub repo - glofloworg/gf_web3_monitor  
`docker run glofloworg/gf_web3_monitor`  
ENV vars for the container are:  
- `GF_PORT`
- `GF_PORT_METRICS`
- `GF_MONGODB_HOST`
- `GF_INFLUXDB_HOST`
- `GF_AWS_SQS_QUEUE`
- `GF_WORKERS_AWS_DISCOVERY`
- `GF_WORKERS_HOSTS`
- `GF_SENTRY_ENDPOINT`
- `GF_EVENTS_CONSUME`
- `GF_PY_PLUGINS_BASE_DIR_PATH`




WORKER_INSPECTOR
Agent usually running on the same host as an Ethereum node (geth).
- mainly using eth-rpc API to communicate with Eth node.
- uses some Geth specific API's.
- joins several datastructures on blocks and tx's.
- provides REST API.
- in the future will contain functions as well that will assume that they're running on the same host as the Eth node.
- ideally it will run and query a full archive geth node to get all of the expected data (tx traces, acc balances, etc.)





GO TESTS:
```bash

# working dir - ./py/ops/
# ENV VARS:
# - GF_GETH_HOST=127.0.0.1
# - GF_SENTRY_ENDPOINT=...
$ python3 gf_builder_cli.py -run=test_go


# working dir - ./go/gf_eth_monitor_core/
# ENV VARS:
# - GF_TEST_WORKER_INSPECTOR_HOST_PORT=127.0.0.1:9000
$ go test -v

#--------------
# working dir - ./go/gf_eth_indexer
# ENV VARS:
# - GF_TEST_WORKER_INSPECTOR_HOST_PORT=127.0.0.1:9000
$ go test -v -run Test__indexer_core

# ENV VARS:
# - AWS_REGION=us-east-1
# - AWS_ACCESS_KEY_ID=...
# - AWS_SECRET_ACCESS_KEY=...
# - GF_TEST_WORKER_INSPECTOR_HOST_PORT=127.0.0.1:9000
$ go test -v -run Test__indexer_http

#--------------
# working dir - ./go/gf_eth_blocks
# ENV VARS:
# - GF_TEST_WORKER_INSPECTOR_HOST_PORT=127.0.0.1:9000
$ go test -v -run Test__blocks

```



PY TESTS:
```bash

# working dir - ./py/ops/
# ENV VARS:
# - GF_GETH_HOST=...
# - GF_SENTRY_ENDPOINT=...
$ python3 gf_builder_cli.py -run=test_py

```

