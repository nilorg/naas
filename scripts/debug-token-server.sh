#!/bin/bash

export OAUTH2_SERVER=http://localhost:8080/oauth2 \
export OAUTH2_CLIENT_ID=1000 \
export OAUTH2_CLIENT_SECRET=83312b80-e69c-43f1-bcaa-358a1d1e7830 \
export OAUTH2_REDIRECT_URI=http://localhost:8000/auth/callback

go run cmd/token-server/main.go