#!/bin/bash

echo 'hack 4 good... _\,,/'
echo $GOPATH

SRC=src/$REPOSITORY/$GROUP/$PROJECT/

rm -rf $GOPATH/$SRC/
echo 'removed: '$GOPATH/$SRC

shopt -s extglob

mkdir -p $GOPATH/$SRC && mv !(hack*) $GOPATH/$SRC

shopt -u extglob

GOTO=$PWD

cd $GOPATH/$SRC

go get -v -t ./...

echo 'calling make: '$COMMAND
make $COMMAND

mv -f $BINARY $GOTO/$BINARY
