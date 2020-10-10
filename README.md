# gloflow-ethmonitor
monitoring of the Ethereum network





The goal is to create a monitoring tool for the Ethereum mainnet.  
Main gf-ethmonitor server manages the spawning of modified GF go-ethereum nodes that connect into the mainnet. These nodes are customized to contain custom instrumentation, and are sending data to a common queue system that is shared with the gf-ethmonitor server. The main server is a consumer, and the GF go-ethereum node is a producer.  

pre-built container is available in a Dockerhub repo - glofloworg/gf_eth_monitor  
`docker run glofloworg/gf_eth_monitor`  
ENV vars for the container are:  
- `GF_PORT`
- `GF_MONGODB_HOST`