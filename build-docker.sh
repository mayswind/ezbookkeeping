#!/usr/bin/env bash

mode=$1;
version=`grep '"version": ' package.json | awk -F ':' '{print $2}' | tr -d ' ' | tr -d ',' | tr -d '"'`;

if [ "$mode" == "--snapshot" ]; then
  version="SNAPSHOT-"`date "+%Y%m%d"`;
fi

echo "Building docker image...";
docker build -t ezbookkeeping:${version} .
