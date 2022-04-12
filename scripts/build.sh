#!/bin/bash
# 参数1 要构建的目录名 mp2c
echo "build start"
cd $(dirname $(readlink -f "$0"))/../cmd/$1/
pwd
env GOPROXY=https://goproxy.cn,https://goproxy.io,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o $1 .
echo "build end"