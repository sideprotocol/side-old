# Join a network

## Initialize node configuration files
```shell
export TESTNET_PATH=${PWD}/testnet
cd build/
./sidechaind init outsider --home ${TESTNET_PATH}/outsider
```

## Get the genesis file for the existing network

```shell
cp ${TESTNET_PATH}/node0/config/genesis.json ${TESTNET_PATH}/outsider/config/
```

## Add at least one seed peer

### Add seed peer (you need a seed peer online for this to work)

```shell
# export SEEDER_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps | grep docker_seeder | awk '{print $1}'))
sed -i "s/seeds = \"\"/seeds = \"6a4fa5865eda7c4f3e9d257190f4635008a3c8d6@13.212.61.41:26656\"/g" ${TESTNET_PATH}/outsider/config/config.toml
```

### Set node to look for peers constantly

```shell
sed -i 's/seed_mode = false/seed_mode = true/g' ${TESTNET_PATH}/outsider/config/config.toml
```

### On local networks, change addr_book_strict flag to false

```shell
sed -i 's/addr_book_strict = true/addr_book_strict = false/g' ${TESTNET_PATH}/outsider/config/config.toml
```

## Build docker image

```shell
cd ../..
docker build -f .docker/Dockerfile-outsider . -t docker_outsider
```

## Start the node

```shell
docker run --name=outsider --network=docker_default -d docker_outsider
```