#.PHONY是一个伪目标，可以防止在Makefile中定义的执行命令的目标和工作目录下的实际文件出现名字冲突，另一种是提高执行makefile时的效率
.PHONY: debug build swagger debug-naas debug-token-server test clean

BIN_PATH = ./bin
BIN_OUTPUT_NAAS_NAME = naas
BIN_OUTPUT_NAAS_TOKEN_SERVER_NAME = token-server

build:
	@go build -ldflags "-w -s" -o $(BIN_PATH)/$(BIN_OUTPUT_NAAS_NAME) ./cmd/main.go
	@go build -ldflags "-w -s" -o $(BIN_PATH)/$(BIN_OUTPUT_NAAS_TOKEN_SERVER_NAME) ./cmd/token-server/main.go

debug-naas: swagger
	@GRPC_ENABLE="true" GRPC_GATEWAY_ENABLE="true" HTTP_ENABLE="true" go run ./cmd/main.go

debug-token-server:
	@export OAUTH2_SERVER="http://naas:8080/oauth2"
	@export OAUTH2_CLIENT_ID="1000"
	@export OAUTH2_CLIENT_SECRET="99799a6b-a289-4099-b4ad-b42603c17ffc"
	@export OAUTH2_REDIRECT_URI="http://naas-admin.nilorg.com/admin/auth/callback"
	@go run ./cmd/main.go

swagger:
	@swag init -g ./internal/server/server.go

test:
	@go test -cover -race ./...

clean:
	@rm -rf $(BIN_PATH)/*
