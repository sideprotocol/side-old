#!/usr/bin/env bash

set -eo pipefail

# Run ignite generate openapi command
ignite generate openapi -y

# Convert openapi.yml to swagger.yml
awk '
  /^openapi:/ { print "swagger: '\''2.0'\''"; next }
  1
' ./docs/static/openapi.yml > ./client/docs/swagger-ui/swagger.yml

# Update the fields
gsed -i 's/^\(  title: \).*$/\1Sidechain Chain - gRPC Gateway docs/' ./client/docs/swagger-ui/swagger.yml
gsed -i 's/^\(  description: \).*$/\1A REST interface for state queries and transactions/' ./client/docs/swagger-ui/swagger.yml

# Check if version field exists
if ! grep -q "^  version:" ./client/docs/swagger-ui/swagger.yml; then
    gsed -i '/^info:/ a\
\  version: 1.0.0
' ./client/docs/swagger-ui/swagger.yml
else
    gsed -i 's/^\(  version: \).*$/\11.0.0/' ./client/docs/swagger-ui/swagger.yml
fi