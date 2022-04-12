#!/bin/bash

$1 chart名称

source ./docker-tag.sh

echo "package start"
cd $(dirname $(readlink -f "$0"))
pwd
tag = getTag $CI_COMMIT_REF_NAME

container_name='mp2c-go'

replicaCount = 1
if [ $CI_COMMIT_REF_NAME = "develop" ]; then
  namespace='dev'
  replicaCount = 3
elif [ ${CI_COMMIT_REF_NAME:0:1} = "v" ]; then
  namespace='prod'
  replicaCount = 3
elif [ ${CI_COMMIT_REF_NAME:0:8}  = "release-" ]; then
  namespace='prod'
  container_name=${container_name}${1}
elif [ ${CI_COMMIT_REF_NAME:0:7} = "hotfix-" ]; then
  namespace='prod'
  container_name=${container_name}${1}
elif [ ${CI_COMMIT_REF_NAME:0:7} = "feature-" ]; then
  namespace='dev'
  container_name=${container_name}${1}
else
  echo "THIS CI_COMMIT_REF_NAME $CI_COMMIT_REF_NAME DOES NOT REQUIRE DEPLOYMENT"
  exit 0
fi
remoteAppVersion=`helm list -n ${namespace} --filter \^${container_name}\$ | grep ${container_name} | awk '{print ${10}}'`
echo $remoteAppVersion
if [ -z remoteAppVersion ]; then
  echo "deployment \"${container_name}\"  not exists. prepare for create"
  helm install --set image.tag=$tag --set image.branch=$CI_COMMIT_REF_NAME --set replicaCount=${replicaCount} $container_name ./$1 -n ${namespace}
else
  helm upgrade --set image.tag=$2 --set image.branch=$CI_COMMIT_REF_NAME --set replicaCount=${replicaCount} $container_name ./$1 -n ${namespace}
fi
