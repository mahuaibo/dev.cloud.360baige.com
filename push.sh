#!/usr/bin/env bash

bee pack -be GOOS=linux

scp dev.cloud.360baige.com.tar.gz ma@123.56.6.206:/home/wealth/cloud/dev.cloud.360baige.com/
