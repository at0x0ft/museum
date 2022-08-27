#!/usr/bin/env sh
set -eu

HOST_MOUNTPOINT_PATH='/hoge/fuga'
CONTAINER_MOUNTPOINT_PATH='/workspace'
. '../lib/convert_host_path_to_container_path.lib.sh'

set -- '/hoge/fuga/piyo.yml' '/foo/bar/config.txt' '/hoge' '/hoge/fuga'

processed_arguments=0
while [ "${processed_arguments}" -lt "${#}" ]; do
  converted_argument="$(convert_host_path_to_container_path "${1}" "${HOST_MOUNTPOINT_PATH}" "${CONTAINER_MOUNTPOINT_PATH}")"
  shift
  set -- "${@}" "${converted_argument}"
  processed_arguments=$((processed_arguments + 1))
done

echo "${@}"
