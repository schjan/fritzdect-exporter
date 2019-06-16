#!/usr/bin/env bash

docker login -u="$QUAY_USER" -p="$QUAY_APIKEY" quay.io;

GO111MODULE=on go mod vendor;
docker build . -t $IMAGE:amd64 -t $IMAGE:latest   --build-arg opts="CGO_ENABLED=0 GOOS=linux GOARCH=amd64";
docker build . -t $IMAGE:arm32v6 --build-arg opts="CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6";

docker push $IMAGE;

docker manifest create $IMAGE $IMAGE:amd64 $IMAGE:arm32v6;
docker manifest annotate $IMAGE $IMAGE:arm32v6 --os linux --arch arm;
docker manifest push $IMAGE;
