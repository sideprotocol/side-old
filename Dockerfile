FROM golang:1.19.4-bullseye AS build-env

WORKDIR /go/src/github.com/sideprotocol/sidchain

RUN apt-get update -y
RUN apt-get install git -y

COPY . .

RUN make build

FROM golang:1.19.4-bullseye

RUN apt-get update -y
RUN apt-get install ca-certificates jq -y

WORKDIR /root

COPY --from=build-env /go/src/github.com/sideprotocol/sidchain/build/sidechaind /usr/bin/sidechaind

# Copy the setup_node_ubuntu.sh script to the container
COPY --from=build-env /go/src/github.com/sideprotocol/sidchain/scripts/setup_node_container.sh /root/setup_node_container.sh

# Copy the entrypoint.sh script to the container
COPY --from=build-env /go/src/github.com/sideprotocol/sidchain/scripts/entrypoint.sh /root/entrypoint.sh
RUN  chmod +x /root/setup_node_container.sh /root/entrypoint.sh

EXPOSE 26656 26657 1317 9090 8545 8546

ENTRYPOINT ["/root/entrypoint.sh"]
