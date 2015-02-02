#!/bin/sh

usage() {
  cat <<-USAGE
Usage: ${0##*/} [<option>...] [<path>]
Install clownfish to <path>. Default:
  * determined interactively if possible
  * /usr/local/bin if run as root
  * /usr/local/bin if it is writable
  * $HOME/bin otherwise
Options:
  -v <version>           Install version <version>.      Default: $DEFAULT_VERSION
  -o <operating system>  Install for <operating system>. Default: $DEFAULT_OS
  -a <architecture>      Install for <architecture>.     Default: $DEFAULT_ARCH
USAGE
}

warn() {
  echo "${0##*/}: $*" >&2
}

fail() {
  warn "$*"
  exit 1
}

setup_downloader() {
  if command -v curl >/dev/null; then
    DOWNLOADER=curl
  elif command -v wget >/dev/null; then
    DOWNLOADER=wget
  else
    fail "missing dependencies, please install one of: curl, wget"
  fi
}

setup_defaults() {
  DEFAULT_DIR="$(default_dir)"
  DEFAULT_OS="$(default_os)"
  DEFAULT_ARCH="$(default_arch)"
  DEFAULT_VERSION="$(latest_version)" || exit $?
}

default_dir() {
  local DIR=/usr/local/bin
  if [ "$(id -u)" -eq 0 ]; then
    echo "$DIR"
  elif [ -d "$DIR" -a -w "$DIR" ]; then
    echo "$DIR"
  else
    echo "$HOME/bin"
  fi
}

default_os() {
  case "$(uname -s)" in
    Darwin) echo "mac";;
    *)      echo "linux";;
  esac
}

default_arch() {
  echo "$(getconf LONG_BIT)bit"
}

latest_version() {
  follow https://github.com/clownfish/cli/releases/latest | filter_version_tag ||
  fail 'Failed to get latest version. Check your internet connection.'
}

follow() {
  case "$DOWNLOADER" in
    curl) curl --head --silent "$@";;
    wget) wget --server-response --quiet --output-document=/dev/null "$@" 2>&1;;
  esac | tr -d '\r'
}

filter_version_tag() {
  awk -v FS=/ '/Location:/{print $NF}'
}

parse_options() {
  while getopts o:a:v:h OPTNAME; do
    case $OPTNAME in
      o) OS=$OPTARG;;
      a) ARCH=$OPTARG;;
      v) VERSION=$OPTARG;;
      h) usage; exit;;
      *) usage >&2; exit 64;;
    esac
  done
  shift $(expr $OPTIND - 1)
  case $# in
    0) ;;
    1) DIR=$1;;
    *) usage >&2; exit 64;;
  esac
}

use_defaults() {
  : ${DIR:="$(read_directory)"}
  : ${DIR:=$DEFAULT_DIR}
  : ${OS:=$DEFAULT_OS}
  : ${ARCH:=$DEFAULT_ARCH}
  : ${VERSION:=$DEFAULT_VERSION}
}

read_directory() {
  if [ -t 0 ]; then
    read -p "install into: [$DEFAULT_DIR] " REPLY
    if [ -n "$REPLY" ]; then
      echo "$REPLY"
    else
      echo "$DEFAULT_DIR"
    fi
  fi
}

check_directory() {
  check_directory_exists
  check_directory_writable
  check_directory_in_path
}

check_directory_exists() {
  test -d "$DIR" ||
  mkdir -p "$DIR" ||
  fail "unable to create installation directory $DIR"
}

check_directory_writable() {
  test -w "$DIR" ||
  fail "unable to write installation directory $DIR"
}

check_directory_in_path() {
  echo "$PATH" | tr : "\n" | grep -q "^$DIR$" ||
  warn "installation directory $DIR not in PATH"
}

install() {
  CLOWNFISH_FILES="https://github.com/keighl/clownfish/releases/download/$VERSION/clownfish-$OS-$ARCH.tgz"
  echo "\nDownloading $CLOWNFISH_FILES...\n"
  download $CLOWNFISH_FILES | extract "$DIR"
  [ $? -eq 0 ] && echo "Installed clownfish to $DIR"
}

download() {
  case "$DOWNLOADER" in
    curl) curl -s --location "$@";;
    wget) wget --output-document=- "$@";;
  esac ||
  fail "failed to download $*"
}

extract() {
  tar xzf "$1" clownfish ||
  fail "failed to extract clownfish"
}

setup_downloader
setup_defaults
parse_options "$@"
use_defaults
check_directory
install