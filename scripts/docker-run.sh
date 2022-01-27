#!/bin/bash
exist=`docker ps -a -q -f=name=mio`
if [ -n "${exist}" ]; then
    docker rm -f mio > /dev/null
fi
cd ../ && docker build -t mio/mio:latest .
docker run --name=mio -d -p 80:80 mio/mio:latest