#!/usr/bin/env bash

docker login -u="$USER" -p="$APIKEY"

GO111MODULE=on go mod vendor

docker buildx create --use --name multiplatformbuilder
docker buildx build -t $IMAGE --platform=linux/arm,linux/arm64,linux/amd64 . --push
