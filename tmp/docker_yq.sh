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

PATH_CONVERT_LIBRARY_PATH="${SCRIPT_ROOT}/convert_path_as_another_fs.lib.sh"
. "${PATH_CONVERT_LIBRARY_PATH}"

convert_host_path_to_container_path() {
  echo "[DEBUG] ${1}" 1>&2
  if [ ! -e "${1}" ]; then
    printf '%s' "${1}"
    return 0
  fi

  local readonly host_abspath="$(readlinkf "${1}")"
  local readonly DOCKER_COMPOSE_HOST_MOUNTPOINT_PATH="${REPOSITORY_ROOT}"
  local readonly CONTAINER_WORKING_DIRECTORY_PATH='/workspace'
  printf "[DEBUG]: 1 = %s, 2 = %s, 3 = %s\n" ${host_abspath} ${DOCKER_COMPOSE_HOST_MOUNTPOINT_PATH} ${CONTAINER_WORKING_DIRECTORY_PATH} 1>&2
  convert_path_as_another_fs "${host_abspath}" "${DOCKER_COMPOSE_HOST_MOUNTPOINT_PATH}" "${CONTAINER_WORKING_DIRECTORY_PATH}"
  return 0
}

docker_yq() {
  local readonly DOCKER_COMPOSE_PATH="${REPOSITORY_ROOT}/docker-compose.yml"

  local processed_arguments=0
  while [ "${processed_arguments}" -lt "${#}" ]; do
    converted_argument="$(convert_host_path_to_container_path "${1}")"
    shift
    set -- "${@}" "${converted_argument}"
    processed_arguments=$((processed_arguments + 1))
  done

  docker-compose -f "${DOCKER_COMPOSE_PATH}" run --rm -T yq "${@}"
  return 0
}
docker_yq "${@}"
