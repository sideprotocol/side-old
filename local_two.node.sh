# #!/bin/bash

# # 1. Initialize the nodes and get their IDs
# docker-compose up -d
# sleep 20  # Wait for the nodes to fully start before proceeding
# docker run --rm -v $(pwd)/data1:/root/.sidechaind sidechain/node sidechaind tendermint show-node-id > node1.id
# docker run --rm -v $(pwd)/data2:/root/.sidechaind sidechain/node sidechaind tendermint show-node-id > node2.id
# docker-compose down

# # # 2. Replace the seed nodes in the docker-compose file with the actual IDs
# #sed -i "s/tcp:\/\/node2:26656/tcp:\/\/$(cat node2.id):26656/" docker-compose.yml
# #sed -i "s/tcp:\/\/node1:26656/tcp:\/\/$(cat node1.id):26656/" docker-compose.yml

# # # 3. Start the nodes using the modified docker-compose file
# # docker-compose up


#!/bin/bash

# Bring up the nodes initially to generate the node IDs.
docker-compose up -d

# Wait for the nodes to fully start and generate their node IDs.
sleep 30

# Fetch the node IDs.
# Fetch the node IDs.
NODE_ID1=$(docker exec -it sidechain-node1-1 sidechaind tendermint show-node-id | tr -d '\n')
NODE_ID2=$(docker exec -it sidechain-node2-1 sidechaind tendermint show-node-id | tr -d '\n')
echo $NODE_ID1
echo $NODE_ID2

# Replace the seeds in the docker-compose file with the actual node IDs.
sed -i.bak -e "s|node1:26656|${NODE_ID1}@node1:26656|g" docker-compose.yml
sed -i.bak -e "s|node2:26656|${NODE_ID2}@node2:26656|g" docker-compose.yml

cat docker-compose.yml
# Restart the nodes with the updated configuration.
docker-compose down
#docker-compose up
