#!/usr/bin/env bash

set -eo pipefail

# Delete original document
rm -rf ../docs/static/openapi.yml

# Run ignite generate openapi command
ignite generate openapi -y

# Convert openapi.yml to swagger.yml
awk '
  /^openapi:/ { print "swagger: '\''2.0'\''"; next }
  1
' ./docs/static/openapi.yml > ./client/docs/swagger-ui/swagger.yaml

# Determine the platform and set the appropriate sed command
if [[ "$(uname)" == "Darwin" ]]; then
    SEDCMD="gsed"
else
    SEDCMD="sed"
fi

# Update the fields
$SEDCMD -i 's/^\(  title: \).*$/\1Sidechain Chain - gRPC Gateway docs/' ./client/docs/swagger-ui/swagger.yaml
$SEDCMD -i 's/^\(  description: \).*$/\1A REST interface for state queries and transactions/' ./client/docs/swagger-ui/swagger.yaml

# Remove the existing version field
$SEDCMD -i '/^  version: /d' ./client/docs/swagger-ui/swagger.yaml

# Add the version field after the description field
$SEDCMD -i '/^\(  description: \).*$/a\
  version: 1.0.0' ./client/docs/swagger-ui/swagger.yaml
