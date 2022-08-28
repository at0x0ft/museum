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
CONFIG_FILENAME='config.yml'

make_devcontainer_directory_if_not_exists() {
  if [ ! -d "${1}" ]; then
    mkdir -p "${1}"
  fi

  local readonly host_mountpoint_path="$(yq '.variables.arguments.base_shell.host_mountpoint_path' "${2}")"
  local readonly container_terminal_cwd="${1}/${host_mountpoint_path}/$(yq '.variables.arguments.base_shell.container_terminal_cwd' "${2}")"
  local readonly absolute_container_terminal_cwd="$(readlinkf "${container_terminal_cwd}")"
  if [ ! -d "${absolute_container_terminal_cwd}" ]; then
    mkdir -p "${absolute_container_terminal_cwd}"
  fi

  return 0
}

evaluate_yaml() {
  local readonly DOCKER_COMPOSE_PATH="${SCRIPT_ROOT}/docker-compose.yml"
  local readonly YAML_EVALUATOR_SERVICE_NAME='yaml_evaluator'
  local readonly YAML_EVALUATOR_SERVICE_WORKING_DIRPATH='/workspace'
  local readonly CONFIG_YAML_RELPATH="./${CONFIG_FILENAME}"
  local readonly OUTPUT_DIRECTORY_RELPATH='.'
  local readonly host_mountpoint_abspath="$(readlinkf "${1}")"

  docker-compose -f "${DOCKER_COMPOSE_PATH}" run --rm \
    --user="$(id -u):$(id -g)" \
    -v "${host_mountpoint_abspath}:${YAML_EVALUATOR_SERVICE_WORKING_DIRPATH}" \
    "${YAML_EVALUATOR_SERVICE_NAME}" "${CONFIG_YAML_RELPATH}" "${OUTPUT_DIRECTORY_RELPATH}"

  return 0
}

convert_devcontainer_yaml_to_json() {
  local readonly YQ_JSON_INDENTATION_SPACES='4'
  local readonly YQ_EVALUATION_STATEMENT_PATH='.'
  local readonly yaml_path="${1}/devcontainer.yml"
  local readonly json_path="${1}/devcontainer.json"

  yq -o=json -I="${YQ_JSON_INDENTATION_SPACES}" "${YQ_EVALUATION_STATEMENT_PATH}" "${yaml_path}" >"${json_path}"
  return 0
}

deploy_docker_config() {
  local readonly yq_service_docker_context_path=".services.${1}.build.context"
  local readonly docker_compose_path="${2}/docker-compose.yml"
  local readonly service_docker_context_relpath="$(yq "${yq_service_docker_context_path}" "${docker_compose_path}")"

  if [ "${service_docker_context_relpath}" != 'null' ]; then
    local readonly service_docker_srcpath="./services/${1}/docker"
    local readonly service_docker_dstpath="${2}/${service_docker_context_relpath}"
    mkdir -p "$(dirname -- "${service_docker_dstpath}")"
    if [ -d "${service_docker_dstpath}" ]; then
      rm -rf "${service_docker_dstpath}"
    fi
    cp -r "${service_docker_srcpath}" "${service_docker_dstpath}"
  fi

  return 0
}

deploy_service_configs() {
  deploy_docker_config 'base_shell' "${1}"
  return 0
}

deploy() {
  evaluate_yaml "${1}"
  convert_devcontainer_yaml_to_json "${1}"
  deploy_service_configs "${1}"
  return 0
}
deploy "${1}"
