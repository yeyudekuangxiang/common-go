#!/bin/bash


echo "package start"
cd $(dirname $(readlink -f "$0"))/../
pwd
env
tag=$CI_COMMIT_SHORT_SHA

container_name='mp2c-go'

replicaCount=1
if [ $CI_COMMIT_REF_NAME = "develop" ]; then
  namespace='dev'
  replicaCount=3
  tag=$CI_COMMIT_REF_NAME
elif [ ${CI_COMMIT_REF_NAME:0:1} = "v" ]; then
  namespace='prod'
  replicaCount=3
  tag=$CI_COMMIT_REF_NAME
elif [ ${CI_COMMIT_REF_NAME:0:8}  = "release-" ]; then
  namespace='prod'
  container_name=${container_name}${CI_COMMIT_REF_NAME}
elif [ ${CI_COMMIT_REF_NAME:0:7} = "hotfix-" ]; then
  namespace='prod'
  container_name=${container_name}${CI_COMMIT_REF_NAME}
elif [ ${CI_COMMIT_REF_NAME:0:8} = "feature-" ]; then
  namespace='dev'
  container_name=${container_name}${CI_COMMIT_REF_NAME}
else
  echo "THIS CI_COMMIT_REF_NAME $CI_COMMIT_REF_NAME DOES NOT REQUIRE DEPLOYMENT"
  exit 1
fi
echo "helm list -n ${namespace}"
helm list -n ${namespace}
echo "helm list -n ${namespace} --filter ^${container_name}$ | grep ${container_name}"
helm list -n ${namespace} --filter ^${container_name}$ | grep ${container_name}
remoteAppVersion=`helm list -n ${namespace} --filter ^${container_name}$ | grep ${container_name} | awk '{print $10}'`
echo "remoteAppVersion"
echo $remoteAppVersion
if [ -z $remoteAppVersion ]; then
  echo "chart \"${container_name}\"  not exists. prepare for create"
  helm install --set image.tag=$tag --set image.branch=$CI_COMMIT_REF_NAME --set replicaCount=${replicaCount} $container_name ./build/$1 -n ${namespace}
else
    echo "chart \"${container_name}\"  exists. prepare for upgrade"
  helm upgrade --set image.tag=$tag --set image.branch=$CI_COMMIT_REF_NAME --set replicaCount=${replicaCount} --set restartTime="'$(date +%s)'" $container_name ./build/$1 -n ${namespace}
fi
