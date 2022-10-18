#!/bin/bash

set -e

# 参数1 要构建的目录名 mp2c
echo "build start"
cd $(dirname $(readlink -f "$0"))/../cmd/$1/
echo "build:$1 GOPROXY=https://goproxy.cn,https://goproxy.io,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64"
git config --global url."git@gitlab.miotech.com:".insteadOf "https://gitlab.miotech.com/"
env GOPROXY=https://goproxy.cn,https://goproxy.io,direct GOPRIVATE="gitlab.miotech.com" GONOSUMDB="gitlab.miotech.com" GONOPROXY="gitlab.miotech.com" go mod download
env GOPROXY=https://goproxy.cn,https://goproxy.io,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o $1 .
echo "build end"