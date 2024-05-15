#!/usr/bin/env bash

shopt -s expand_aliases
set -xe

PWD=$(pwd)
if [[ -z "$IMAGE" ]]; then
    IMAGE=$(cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1|tr 'A-Z' 'a-z')
fi

docker login

docker build -t $IMAGE -f ${PWD}/Dockerfile ${PWD}/release

docker tag $IMAGE covsj/gotool

docker push covsj/gotool


shopt -u expand_aliases
