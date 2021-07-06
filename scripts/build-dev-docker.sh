#!/bin/bash
docker build -f deployments/Dockerfile -t nilorg/naas:dev .
docker push nilorg/naas:dev

# docker build -f deployments/Dockerfile.token-server -t nilorg/naas-token-server:dev .
# docker push nilorg/naas-token-server:dev