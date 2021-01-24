#!/bin/sh
CUR_DIR=$(dirname "$0");

if [ -x "${CUR_DIR}/custom-backend-pre-setup.sh" ]; then
  "${CUR_DIR}"/custom-backend-pre-setup.sh
fi
