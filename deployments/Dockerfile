#build stage
FROM nilorg/golang:1.16.4 AS builder
WORKDIR /src
COPY . .
RUN go build -mod=vendor -o ./bin/app ./cmd/main.go

#final stage
FROM nilorg/alpine:latest
WORKDIR /workspace
VOLUME [ "/workspace" ]
COPY configs configs
COPY web web
COPY --from=builder /src/bin/app .
ENV HTTP_ENABLE=true
ENV GRPC_ENABLE=true
ENV GRPC_GATEWAY_ENABLE=true
ENTRYPOINT ./app
EXPOSE 8080 5000 9000