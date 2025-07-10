#!/bin/bash

# ip=`ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|grep -v 172.\*.\*.\*|awk '{print $2}'|tr -d "addr:"`
# echo $ip
#
# cd configs
# sed -i -e "s/127.0.0.1/$ip/g" config.yaml
# cd ..
docker stop student
docker rm student
docker rmi student

docker build --progress=plain -t student .
docker run --name student -d -p 8600:8600 student
