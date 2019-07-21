#!/bin/bash
set -euo pipefail

: ${PYTHON_VERSION:=3.7.4}
: ${CONTAINER_NAME:=slack-wifi-updater}
CONTAINER_IMAGE=docker://docker.io/library/python:${PYTHON_VERSION}-alpine

function cleanup() {
    [[ -n "${container+x}" ]] && buildah rm $container
}

trap 'exit 1' INT HUP QUIT TERM ALRM USR1
trap cleanup 0

container=$(buildah from $CONTAINER_IMAGE)

: ${USER_UID:=1000}
buildah run $container -- sh -c "addgroup -S slack -g $USER_UID && adduser -S slack -G slack -u $USER_UID"
buildah run $container -- apk --no-cache add wireless-tools
buildah run $container -- pip3 install requests
buildah copy $container ./check_wifi.py /usr/bin/check_wifi

buildah config --user=slack --cmd=/usr/bin/check_wifi $container
buildah commit $container $CONTAINER_NAME
