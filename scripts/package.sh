#!/bin/bash


echo "package start"
cd $(dirname $(readlink -f "$0"))/../
pwd
# 参数1要打包的目录名 mp2c 分支名 v1.1.1 develop release-
env
tag=$CI_COMMIT_SHORT_SHA
if [ $CI_COMMIT_REF_NAME = "develop" ]; then
  namespace='dev'
  tag=$CI_COMMIT_REF_NAME
elif [ ${CI_COMMIT_REF_NAME:0:1} = "v" ]; then
  namespace='prod'
  tag=$CI_COMMIT_REF_NAME
fi
IMAGE_TAG=${CI_BUILD_DOCKER_URL}mioyingtech/mio-go:$tag
echo $IMAGE_TAG
docker login -u $CI_BUILD_DOCKER_USER -p $CI_BUILD_TOKEN $CI_BUILD_DOCKER_URL
docker build -t $IMAGE_TAG -f ./build/package/$1/Dockerfile ./cmd/$1
docker push $IMAGE_TAG
echo "package end"