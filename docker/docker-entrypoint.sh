#!/bin/sh

set -e;
export LAB_USER=labapp;
export LAB_GROUP=labapp;

prepare_directories() {
  local log_path="/var/log/labapp";

  if [ "${LAB_LOG_PATH}" != "" ]; then
    log_path="${LAB_LOG_PATH}";
  fi

  if [ ! -d "${log_path}" ]; then
    mkdir ${log_path};
    chown ${LAB_USER}:${LAB_GROUP} -R ${log_path};
  fi
}

prepare_directories;

conf_path_param="";

if [ "${LAB_CONF_PATH}" != "" ]; then
  conf_path_param="--conf-path=${LAB_CONF_PATH}";
fi

if [ $# -gt 0 ]; then
    exec "$@"
else
    exec su-exec ${LAB_USER} /usr/local/bin/labapp/lab server run ${conf_path_param};
fi
