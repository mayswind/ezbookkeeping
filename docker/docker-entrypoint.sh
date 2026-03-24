#!/bin/sh

set -e;

conf_path_param="";

if [ "${OSCAR_CONF_PATH}" != "" ]; then
  conf_path_param="--conf-path=${OSCAR_CONF_PATH}";
fi

if [ $# -gt 0 ]; then
    exec "$@"
else
    exec /oscar/oscar server run ${conf_path_param};
fi
