#!/usr/bin/env bash

mode=$1;
version=`grep "const LAB_VERSION" lab.go | awk -F '=' '{print $2}' | tr -d ' ' | tr -d '"'`;

if [ "$mode" == "--snapshot" ]; then
  version="SNAPSHOT-"`date "+%Y%m%d"`;
fi

echo "Building docker image...";
docker build -t lab:${version} .
