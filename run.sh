#!/usr/bin/env bash

tar -xzvf dev.cloud.360baige.com.tar.gz

killall dev.cloud.360baige.com

nohup ./dev.cloud.360baige.com  >/dev/null 2>error.log &
