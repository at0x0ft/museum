#!/usr/bin/env sh
set -eu

DOTFILES_URL='https://github.com/at0x0ft/dotfiles.git'
if [ "${USER_NAME}" =  'root' ]; then
  USER_HOME_DIRECTORY='/root'
else
  USER_HOME_DIRECTORY="/home/${USER_NAME}"
fi
DOTFILES_REPOSITORY_DSTPATH="${USER_HOME_DIRECTORY}/.dotfiles"
DOTFILES_INSTALL_COMMAND="${DOTFILES_REPOSITORY_DSTPATH}/src/bin/install.sh -econtainer"

apt-get update
DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
  libncurses-dev \
  unzip \
  curl \
  file \
  jq \
  make \
  autoconf
chsh -s $(which zsh) "${USER_NAME}"

su "${USER_NAME}" -c "git clone '${DOTFILES_URL}' ${DOTFILES_REPOSITORY_DSTPATH} && ${DOTFILES_INSTALL_COMMAND}"
