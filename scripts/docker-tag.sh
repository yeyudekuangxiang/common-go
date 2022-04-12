#!/bin/bash
function getTag() {
    if [ $1 = "develop" ]; then
      return $1
    elif [ ${1:0:1} = "v" ]; then
      return $1
    elif [ ${1:0:8}  = "release-" ]; then
      tag =  `$1 | md5sum`
      return ${tag:0:8}
    elif [ ${1:0:7} = "hotfix-" ]; then
      tag =  `$1 | md5sum`
        return ${tag:0:8}
    elif [ ${1:0:7} = "feature-" ]; then
      tag =  `$1 | md5sum`
        return ${tag:0:8}
    else
      return ""
    fi
}
