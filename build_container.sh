#!/bin/bash
set -euo pipefail

container_image=docker://docker.io/library/alpine:3.10

function cleanup() {
    [[ -n "${container+x}" ]] && buildah rm $container
}

trap cleanup 0

container_name=${1:-wireless-tools}
container=$(buildah from $container_image)

buildah run $container -- sh -c "addgroup -S slack -g 1000 && adduser -S slack -G slack -u 1000"
buildah run $container -- apk --no-cache add python3 wireless-tools
buildah run $container -- pip3 install requests
buildah copy $container ./check_wifi.py /usr/bin/check_wifi

buildah config --user=slack --cmd=/usr/bin/check_wifi $container
buildah commit $container $container_name
