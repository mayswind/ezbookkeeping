#!/bin/sh

set -e;

conf_path_param="";

if [ "${LAB_CONF_PATH}" != "" ]; then
  conf_path_param="--conf-path=${LAB_CONF_PATH}";
fi

if [ $# -gt 0 ]; then
    exec "$@"
else
    exec /usr/local/bin/labapp/lab server run ${conf_path_param};
fi
