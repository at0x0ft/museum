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

make_devcontainer_directory_if_not_exists() {
  if [ ! -d "${1}" ]; then
    mkdir -p "${1}"
  fi
  return 0
}

# If you install 'yq' command in local, comment out yq function below.
yq() {
  local readonly DOCKER_YQ_HELPER_SCRIPT_PATH="${SCRIPT_ROOT}/lib/docker_yq.sh"
  ${DOCKER_YQ_HELPER_SCRIPT_PATH} "${@}"
  return 0
}

generate_devcontainer_json() {
  local readonly YQ_DEVCONTAINER_JSON_PATH='.configs.vscode_devcontainer'
  local readonly YQ_JSON_INDENTATION_SPACES='4'

  yq -o=json -I="${YQ_JSON_INDENTATION_SPACES}" "${YQ_DEVCONTAINER_JSON_PATH}" "${1}" >"${2}"
  return 0
}

generate() {
  # START: temporary initial setup
  local readonly CONFIG_YAML_INPUT_PATH="${SCRIPT_ROOT}/config.yml"
  local readonly PROJECT_PATH="${SCRIPT_ROOT}/test_project"
  # Works like argument
  set -- "${CONFIG_YAML_INPUT_PATH}" "${PROJECT_PATH}"
  # END

  local readonly DEVCONTAINER_DIRNAME='.devcontainer'

  local readonly devcontainer_directory_path="${2}/${DEVCONTAINER_DIRNAME}"
  make_devcontainer_directory_if_not_exists "${devcontainer_directory_path}"
  local readonly devcontainer_json_output_path="${devcontainer_directory_path}/devcontainer.json"

  generate_devcontainer_json "${1}" "${devcontainer_json_output_path}"

  return 0
}
generate
