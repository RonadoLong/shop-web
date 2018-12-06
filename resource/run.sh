#!/bin/bash
DOCKER_IMG=client_img
DOCKER_CONTAIN=client_containt
APP_PORT=9527
HOST_PORT=9527
#APP_LOG_DIR=/var/log/
#HOST_LOG_DIR=/var/log/
docker stop ${DOCKER_CONTAIN}
docker rm  ${DOCKER_CONTAIN}
docker rmi ${DOCKER_IMG}
docker build -t ${DOCKER_IMG} .
#if [ -n "$5" -a -n "$6" ];then
#docker run --name ${DOCKER_CONTAIN} -p ${HOST_PROT}:${APP_PORT} -v ${HOST_LOG_DIR}:${APP_LOG_DIR} -v /etc/localtime:/etc/localtime  -d ${DOCKER_IMG}
#else
docker run --privileged=true --name ${DOCKER_CONTAIN} -p ${HOST_PORT}:${APP_PORT} -d ${DOCKER_IMG}
#fi