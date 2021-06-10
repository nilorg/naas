#!/bin/bash

export OAUTH2_SERVER=http://localhost:8080/oauth2 \
export OAUTH2_CLIENT_ID=1000 \
export OAUTH2_CLIENT_SECRET=99799a6b-a289-4099-b4ad-b42603c17ffc \
export OAUTH2_REDIRECT_URI=http://localhost:8000/auth/callback

go run cmd/token-server/main.go