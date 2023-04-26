#!/bin/bash
set -e

# Run setup_node_ubuntu.sh script
/root/setup_node_container.sh

# Start sidechaind
exec sidechaind start