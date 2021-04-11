#!/bin/sh

set -e;

conf_path_param="";

if [ "${EBK_CONF_PATH}" != "" ]; then
  conf_path_param="--conf-path=${EBK_CONF_PATH}";
fi

if [ $# -gt 0 ]; then
    exec "$@"
else
    exec /usr/local/bin/ezbookkeeping/ezbookkeeping server run ${conf_path_param};
fi
