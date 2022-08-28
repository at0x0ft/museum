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

PATH_CONVERT_LIBRARY_PATH="${SCRIPT_ROOT}/convert_path_as_another_fs.lib.sh"
. "${PATH_CONVERT_LIBRARY_PATH}"

convert_base_shell_service_path_to_host_path() {
  convert_path_as_another_fs "${1}" "${CONTAINER_WORKSPACE_FOLDER}" "${LOCAL_WORKSPACE_FOLDER}"
  return 0
}

convert_host_path_to_dst_service_path() {
  convert_path_as_another_fs "${1}" "${LOCAL_WORKSPACE_FOLDER}" "${2}"
  return 0
}

# TODO: if union above functions, the function will be declared below.
convert_base_shell_service_path_to_dst() {
  convert_path_as_another_fs "${1}" "${CONTAINER_WORKSPACE_FOLDER}" "${2}"
  return 0
}

convert_path() {
  if [ ! -e "${1}" ]; then
    printf '%s' "${1}"
    return 0
  fi

  local readonly abspath=$(readlinkf "${1}")

  is_execute_wrapper_script() {
    local readonly split_path="${path#"${EXECUTE_WRAPPER_SCRIPT_DIRECTORY}"}"
    if [ "${path}" != "${split_path}" ]; then
      return 0
    fi
    return 1
  }

  if is_execute_wrapper_script "${argument}"; then
    local readonly wrapper_script_absolute_path=$(readlinkf "${argument}")
    local readonly split_execute_command="${wrapper_script_absolute_path#"${EXECUTE_WRAPPER_SCRIPT_DIRECTORY}"}"
    printf '%s' "${split_execute_command#?}"
  else
    printf '%s' $(convert_devcontainer_filepath_to_runtime_container_filepath "${argument}")
  fi

  return 0
}

run_docker_compose() {
  local readonly DOCKER_EXEC_OPTIONS="${1}"
  shift
  local readonly EXECUTE_COMMAND="${1}"
  shift
  local readonly DEVCONTAINER_ROOT="${CONTAINER_WORKSPACE_FOLDER}/.devcontainer"
  local readonly ENV_FILE="${CONTAINER_WORKSPACE_FOLDER}/.env"
  local readonly DEVCONTAINER_CONFIG_FILE="${DEVCONTAINER_ROOT}/config"
  local readonly EXECUTE_WRAPPER_SCRIPT_DIRECTORY="${DEVCONTAINER_ROOT}/bin"
  local line
  for line in $(cat "${ENV_FILE}"); do
    eval "local readonly ${line}"
  done
  for line in $(cat "${DEVCONTAINER_CONFIG_FILE}"); do
    eval "local readonly ${line}"
  done

  get_runtime_container_name() {
    printf '%s_%s_1' "${COMPOSE_PROJECT_NAME}" "${RUNTIME_CONTAINER_SERVICE_NAME}"
    return 0
  }

  convert_devcontainer_filepath_to_runtime_container_filepath() {
    local readonly file_absolute_path=$(readlinkf "${1}")
    local readonly docker_host_filepath="${LOCAL_WORKSPACE_FOLDER}${file_absolute_path#$CONTAINER_WORKSPACE_FOLDER}"
    printf '%s%s\n' "${RUNTIME_CONTAINER_WORKING_DIRECTORY}" "${docker_host_filepath#$LOCAL_WORKSPACE_FOLDER}"
    return 0
  }

  convert_path_in_arguments() {

    local result=''
    local argument_index=1
    for argument in "${@}"; do
      result="${result}"$(convert_path "${argument}")
      if [ "${argument_index}" -lt "${#}" ]; then
        result="${result} "
      fi
      argument_index=$((argument_index+1))
    done
    printf '%s' "${result}"
    return 0
  }

  local readonly runtime_container_name=$(get_runtime_container_name)
  local readonly runtime_container_current_working_dirpath=$(convert_devcontainer_filepath_to_runtime_container_filepath $(pwd))
  local readonly converted_execute_arguments=$(convert_path_in_arguments "${@}")
  local readonly change_directory_and_execute_command="cd ${runtime_container_current_working_dirpath} && ${EXECUTE_COMMAND} ${converted_execute_arguments}"
  docker container exec "${DOCKER_EXEC_OPTIONS}" \
    "${runtime_container_name}" \
    sh -c "${change_directory_and_execute_command}"

  local readonly DST_SERVICE_NAME="${1}"
  shift
  local readonly DST_MOUNTPOINT_PATH="${}"
  local readonly exec_command=...
  convert...

  local readonly dst_current_working_directory="$(convert_base_shell_service_path_to_dst "$(pwd)")"

  docker-compose run --rm -T --entrypoint='sh' "${DST_SERVICE_NAME}" -c \
    "cd ${dst_current_working_directory} && ${exec_command} ${@}"
  return 0
}
