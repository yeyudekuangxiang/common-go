#!/bin/bash

set -e

echo "remove start"
container_name='mp2c-go'

if [ ${CI_COMMIT_REF_NAME:0:8}  = "release-" ]; then
  namespace='mio-prod'
  container_name=${container_name}${CI_COMMIT_REF_NAME}
elif [ ${CI_COMMIT_REF_NAME:0:7} = "hotfix-" ]; then
  namespace='mio-prod'
  container_name=${container_name}${CI_COMMIT_REF_NAME}
elif [ ${CI_COMMIT_REF_NAME:0:8} = "feature-" ]; then
  namespace='mio-dev'
  container_name=${container_name}${CI_COMMIT_REF_NAME}
else
  echo "仅支持 release-(.+) hotfix-(.+) feature-(.+) 自动删除"
  exit 1
fi
remoteAppVersion=`helm list -n ${namespace} --filter ^${container_name}$ | grep ${container_name} | awk '{print $10}'`
echo "remove namespace:$namespace  container_name:$container_name  remoteAppVersion:$remoteAppVersion"

if [ -z $remoteAppVersion ]; then
  echo "release \"${container_name}\"  not exists. need not remove"
else
    echo "release \"${container_name}\"  exists. prepare for remove"
  helm uninstall $container_name  -n ${namespace}
fi
echo "remove end"
