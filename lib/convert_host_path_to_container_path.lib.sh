# ${1} = path, ${2} = host mountpoint path, ${3} = container mountpoint path
# Note: you have to give "ABSOLUTE" path for each argument.
convert_host_path_to_container_path() {
  local readonly host_relpath="${1##${2}}"
  if [ "${host_relpath}" != "${1}" ]; then
    printf '%s%s' "${3}" "${host_relpath}"
    return 0
  fi
  printf '%s' "${1}"
  return 0
}

# Usage Example:
# HOST_MOUNTPOINT_PATH='/hoge/fuga'
# CONTAINER_MOUNTPOINT_PATH='/workspace'
# . './convert_host_path_to_container_path.lib.sh'

# processed_arguments=0
# while [ "${processed_arguments}" -lt "${#}" ]; do
#   converted_argument="$(convert_host_path_to_container_path "${1}" "${HOST_MOUNTPOINT_PATH}" "${CONTAINER_MOUNTPOINT_PATH}")"
#   shift
#   set -- "${@}" "${converted_argument}"
#   processed_arguments=$((processed_arguments + 1))
# done
