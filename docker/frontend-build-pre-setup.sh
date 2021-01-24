#!/bin/sh
CUR_DIR=$(dirname "$0");

if [ -x "${CUR_DIR}/custom-frontend-pre-setup.sh" ]; then
  "${CUR_DIR}"/custom-frontend-pre-setup.sh
fi
