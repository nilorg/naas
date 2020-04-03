#!/bin/bash

# 输出目录
GO_PUT_PATH='./'

protoc -I/usr/local/protoc/3.7.1/include -I. --go_out=paths=source_relative,plugins=grpc:${GO_PUT_PATH} ./pkg/proto/*.proto
protoc -I/usr/local/protoc/3.7.1/include -I. --grpc-gateway_out=paths=source_relative,logtostderr=true:${GO_PUT_PATH} ./pkg/proto/*.proto