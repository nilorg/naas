#build stage
FROM nilorg/golang:latest AS builder
WORKDIR /src
COPY . .
RUN go build -o ./bin/app -i ./cmd/main.go

#final stage
FROM nilorg/alpine:latest
WORKDIR /workspace
COPY --from=builder /src/bin/app .
ENTRYPOINT ./app
EXPOSE 8080 5000