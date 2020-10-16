#!/bin/bash

export OAUTH2_SERVER=http://localhost:8080/oauth2 \
export OAUTH2_CLIENT_ID=1002 \
export OAUTH2_CLIENT_SECRET=444ssxa-7b9c-4cf0-sdfas-55bf202b99ba \
export OAUTH2_REDIRECT_URI=http://localhost:8000/auth/callback

go run cmd/token-server/main.go