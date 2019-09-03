#!/bin/bash

export GOPATH=$PWD
basePath=$PWD
golint src/core/*.go 
golint src/logger/*.go
go vet src/core/*.go 
go vet src/logger/*.go


GOOS=linux GOARCH=amd64 go build -o bin/lotterySystem core
if [ $? -ne 0 ];then
echo "build is failed"
exit 1
fi

rm -rf deploy-package

export PACKAGEPATH=deploy-package/LotterySystem

mkdir -p $PACKAGEPATH
if [ $? -ne 0 ];then
echo "creating directory deploy-package is failed"
exit 1
fi

Srcfiles="bin config"
cp -avdrf $Srcfiles $PACKAGEPATH
if [ $? -ne 0 ];then
echo "copying files is failed"
exit 1
fi
