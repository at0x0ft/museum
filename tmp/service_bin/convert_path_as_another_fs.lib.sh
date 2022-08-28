# args:
#   ${1} = source file system path,
#   ${2} = source file system base path,
#   ${3} = destination file system base path
# Note: you have to give "ABSOLUTE" path for all arguments.
convert_path_as_another_fs() {
  local readonly relpath_from_source_fs="${1##${2}}"
  if [ "${relpath_from_source_fs}" != "${1}" ]; then
    printf '%s%s' "${3}" "${relpath_from_source_fs}"
    return 0
  fi
  printf '%s' "${1}"
  return 0
}

# Usage Example:
# HOST_MOUNTPOINT_PATH='/hoge/fuga'
# CONTAINER_MOUNTPOINT_PATH='/workspace'
# . './convert_path_as_another_fs.lib.sh'

# processed_arguments=0
# while [ "${processed_arguments}" -lt "${#}" ]; do
#   converted_argument="$(convert_path_as_another_fs "${1}" "${HOST_MOUNTPOINT_PATH}" "${CONTAINER_MOUNTPOINT_PATH}")"
#   shift
#   set -- "${@}" "${converted_argument}"
#   processed_arguments=$((processed_arguments + 1))
# done
