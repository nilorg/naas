#!/bin/bash
docker build -f deployments/Dockerfile -t nilorg.azurecr.io/naas:dev .
docker push nilorg.azurecr.io/naas:dev

# docker build -f deployments/Dockerfile.token-server -t nilorg.azurecr.io/naas-token-server:dev .
# docker push nilorg.azurecr.io/naas-token-server:dev