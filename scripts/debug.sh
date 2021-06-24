#!/bin/bash

# export NAAS_CONFIG=configs/example_config.yaml
export NAAS_CONFIG=configs/config.yaml
export GRPC_ENABLE=true
export GRPC_GATEWAY_ENABLE=true
export HTTP_ENABLE=true
go run cmd/main.go