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

SKELETON_FILENAME='skeleton.yml'
SERVICE_CONFIG_TEMPLATE_FILENAME='config.yml'

get_service_config_path() {
  local readonly service_dirpath="$(yq "${1}.path" "${2}")"
  printf '%s/%s' "${service_dirpath}" "${SERVICE_CONFIG_TEMPLATE_FILENAME}"
  return 0
}

get_base_shell_config_path() {
  get_service_config_path '.base_shell' "${1}"
  return 0
}

get_services_count() {
  yq '.services | length' "${1}"
  return 0
}

create_yq_evaluate_statement() {
  local evaluate_statement=". *+ load(\"${1}\")"
  shift
  for f in "${@}"; do
    evaluate_statement="${evaluate_statement} | . *+ load("${f}")"
  done
  printf '%s' "${evaluate_statement}"
  return 0
}

merge_service_configs() {
  local base_shell_config_path="${1}"
  shift
  local readonly evaluate_statement="$(create_yq_evaluate_statement "${@}")"
  yq "${evaluate_statement}" "${base_shell_config_path}"
  # echo ${evaluate_statement}
  return 0
}

restore() {
  local readonly devcontainer_path="${1}"
  shift
  local readonly skeleton_path="${devcontainer_path}/${SKELETON_FILENAME}"

  local readonly base_shell_config_path="$(get_base_shell_config_path "${skeleton_path}")"
  set -- "${base_shell_config_path}"

  local service_index=0
  local readonly services="$(get_services_count "${skeleton_path}")"
  while [ "${service_index}" -lt "${services}" ]; do
    local readonly service_config_path="$(get_service_config_path ".services[${service_index}]" "${skeleton_path}")"
    set -- "${@}" "${skeleton_path}"
    service_index=$((service_index + 1))
  done

  local readonly mixed_config_output_path="${devcontainer_path}/${SERVICE_CONFIG_TEMPLATE_FILENAME}"
  if [ "${#}" -eq 1 ]; then
    cp "${base_shell_config_path}" "${mixed_config_output_path}"
  else
    create_yq_evaluate_statement "${@}" > "${mixed_config_output_path}"
  fi

  return 0
}
restore "${1}"
