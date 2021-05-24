#!/bin/bash

if [[ ! -d "/root/.bhcd/config" ]]; then
  echo "Initialize /root/.bhcd"
  cp -r /go/initial-node/* /root/.bhcd/
fi

cd /go; ./bhcd start $@
