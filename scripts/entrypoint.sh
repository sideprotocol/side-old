#!/bin/bash
set -e

# Run setup_node_ubuntu.sh script
/root/setup_node_container.sh -c sidehub_1818-1 -d aside aetc ausdc aeth -o n

# Start sidechaind
exec sidechaind start --api.enable