#!/bin/bash

build=""
# 接收是否build参数
if [ "$1" == "build" ]; then
    build="--build"
fi

docker compose down
docker compose up -d $build
