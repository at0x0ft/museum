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
DEVCONTAINER_DIRNAME='.devcontainer'

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

evaluate_yaml() {
  local readonly DOCKER_COMPOSE_PATH="${SCRIPT_ROOT}/docker-compose.yml"
  local readonly YAML_EVALUATOR_SERVICE_NAME='yaml_evaluator'
  local readonly YAML_EVALUATOR_SERVICE_WORKING_DIRPATH='/workspace'
  local readonly config_yaml_relpath="./$(basename -- "${1}")"
  local readonly YAML_OUTPUT_RELPATH='.'

  cp "${1}" "${2}"

  docker-compose -f "${DOCKER_COMPOSE_PATH}" run --rm \
    --user="$(id -u):$(id -g)" \
    -v "${2}:${YAML_EVALUATOR_SERVICE_WORKING_DIRPATH}" \
    "${YAML_EVALUATOR_SERVICE_NAME}" "${config_yaml_relpath}" "${YAML_OUTPUT_RELPATH}"

  return 0
}

convert_devcontainer_yaml_to_json() {
  local readonly YQ_JSON_INDENTATION_SPACES='4'
  local readonly YQ_EVALUATION_STATEMENT_PATH='.'
  local readonly DEVCONTAINER_JSON_PATH="$(dirname -- "${1}")/devcontainer.json"

  yq -o=json -I="${YQ_JSON_INDENTATION_SPACES}" "${YQ_EVALUATION_STATEMENT_PATH}" "${1}" >"${DEVCONTAINER_JSON_PATH}"
  return 0
}

generate() {
  # TODO: Delete here.
  # START: temporary initial setup
  local readonly PROJECT_PATH="${SCRIPT_ROOT}/test_project"
  local readonly CONFIG_YAML_INPUT_PATH="${SCRIPT_ROOT}/config.yml"
  # Works like argument
  set -- "${PROJECT_PATH}" "${CONFIG_YAML_INPUT_PATH}"
  # END

  local readonly DEVCONTAINER_YAML_NAME='devcontainer.yml'

  local readonly devcontainer_directory_path="${1}/${DEVCONTAINER_DIRNAME}"
  make_devcontainer_directory_if_not_exists "${devcontainer_directory_path}"

  evaluate_yaml "${2}" "${devcontainer_directory_path}"
  convert_devcontainer_yaml_to_json "${devcontainer_directory_path}/${DEVCONTAINER_YAML_NAME}"

  return 0
}
generate
