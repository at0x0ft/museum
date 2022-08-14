#!/usr/bin/env sh
set -eu

# ref: https://github.com/ko1nksm/readlinkf/blob/master/readlinkf.sh
readlinkf() {
  [ "${1:-}" ] || return 1
  max_symlinks=40
  CDPATH='' # to avoid changing to an unexpected directory

  target=$1
  [ -e "${target%/}" ] || target=${1%"${1##*[!/]}"} # trim trailing slashes
  [ -d "${target:-/}" ] && target="$target/"

  cd -P . 2>/dev/null || return 1
  while [ "$max_symlinks" -ge 0 ] && max_symlinks=$((max_symlinks - 1)); do
    if [ ! "$target" = "${target%/*}" ]; then
      case $target in
        /*) cd -P "${target%/*}/" 2>/dev/null || break ;;
        *) cd -P "./${target%/*}" 2>/dev/null || break ;;
      esac
      target=${target##*/}
    fi

    if [ ! -L "$target" ]; then
      target="${PWD%/}${target:+/}${target}"
      printf '%s\n' "${target:-/}"
      return 0
    fi

    # `ls -dl` format: "%s %u %s %s %u %s %s -> %s\n",
    #   <file mode>, <number of links>, <owner name>, <group name>,
    #   <size>, <date and time>, <pathname of link>, <contents of link>
    # https://pubs.opengroup.org/onlinepubs/9699919799/utilities/ls.html
    link=$(ls -dl -- "$target" 2>/dev/null) || break
    target=${link#*" $target -> "}
  done
  return 1
}
SCRIPT_PATH="$(readlinkf "${0}")"
SCRIPT_ROOT="$(dirname -- "${SCRIPT_PATH}")"
REPOSITORY_ROOT=$(readlinkf "${SCRIPT_ROOT}/..")

docker_yq() {
  local readonly DOCKER_COMPOSE_PATH="${REPOSITORY_ROOT}/docker-compose.yml"
  local readonly DOCKER_COMPOSE_HOST_MOUNTPOINT_PATH="${REPOSITORY_ROOT}"
  local readonly CONTAINER_WORKING_DIRECTORY_PATH='/workspace'
  local readonly PATH_CONVERT_LIBRARY_PATH="${SCRIPT_ROOT}/convert_host_path_to_container_path.lib.sh"
  . "${PATH_CONVERT_LIBRARY_PATH}"

  local processed_arguments=0
  while [ "${processed_arguments}" -lt "${#}" ]; do
    converted_argument="$(convert_host_path_to_container_path "${1}" "${DOCKER_COMPOSE_HOST_MOUNTPOINT_PATH}" "${CONTAINER_WORKING_DIRECTORY_PATH}")"
    shift
    set -- "${@}" "${converted_argument}"
    processed_arguments=$((processed_arguments + 1))
  done

  docker-compose -f "${DOCKER_COMPOSE_PATH}" run --rm -T yq "${@}"
  return 0
}
docker_yq "${@}"
