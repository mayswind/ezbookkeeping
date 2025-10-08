#!/bin/sh

set -e

# Logging functions
ebk_log() {
	local type="$1"; shift
	printf '%s [%s] [Entrypoint]: %s\n' "$(date +"%F %T")" "$type" "$*"
}

ebk_note() {
	ebk_log Note "$@"
}

ebk_warn() {
	ebk_log Warn "$@" >&2
}

ebk_error() {
	ebk_log ERROR "$@" >&2
	exit 1
}

#######################################
# Get values for environment variables from files,
#   e.g. files from Docker's secrets feature in 
#   /run/secrets.
# Usage: 
#   file_env VAR [DEFAULT]
#   i.e. file_env 'XYZ_DB_PASSWORD' 'example'
#   
#   This will allow to fill in the value of
#   "$XYZ_DB_PASSWORD" from file path's content
#   specified in "$XYZ_DB_PASSWORD_FILE")
#######################################
file_env() {
	local var="$1"
	local fileVar="${var}_FILE"
	local def="${2:-}"

	# Equivalent to Bash's variable expansion ${!var} and functional
	# with all POSIX-compatible shells
	eval val_var=\$$var
	eval val_fileVar=\$$fileVar

	if [ -n "$val_var" ] && [ -n "$val_fileVar" ]; then
		ebk_error "Both $var and $fileVar are set (but are exclusive)"
	fi

	# Set default value if var and fileVar are empty
	local val="$def"

	if [ -n "$val_var" ]; then
		val="$val_var"
	elif [ -n "$val_fileVar" ]; then
		val="$(cat "$val_fileVar")"
	fi

	export "$var"="$val"
	unset "$fileVar"
}

# Convert *_FILE's value to corresponding env var
# To support additional env vars, just add 
# another line containing
# file_env 'ENVIRONMENT_VARIABLE'
file_env 'EBK_SECURITY_SECRET_KEY'
file_env 'EBK_DATABASE_NAME'
file_env 'EBK_DATABASE_USER'
file_env 'EBK_DATABASE_PASSWD'

conf_path_param=""

if [ "${EBK_CONF_PATH}" != "" ]; then
  conf_path_param="--conf-path=${EBK_CONF_PATH}"
fi

if [ $# -gt 0 ]; then
    exec "$@"
else
    exec /ezbookkeeping/ezbookkeeping server run ${conf_path_param}
fi
