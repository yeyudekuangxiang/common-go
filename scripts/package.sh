#!/bin/bash

source ./docker-tag.sh

echo "package start"
cd $(dirname $(readlink -f "$0"))/../
pwd
# 参数1要打包的目录名 mp2c 分支名 v1.1.1 develop release-
tag = getTag $CI_COMMIT_REF_NAME
IMAGE_TAG=${CI_BUILD_DOCKER_URL}mioyingtech/mio-go:$tag
docker login -u $CI_BUILD_DOCKER_USER -p $CI_BUILD_TOKEN $CI_BUILD_DOCKER_URL
docker build -t $IMAGE_TAG -f ./build/$1/Dockerfile ./cmd/$1
docker push $IMAGE_TAG
echo "package end"