#!/bin/bash

PROJ="rex2ansi"
VERSION="1.0"

# os
go_build() {
    OS=$1
    ARCH=$2
    
    NAME=$PROJ-$OS-$ARCH-$VERSION
    rm -fr $NAME
    mkdir $NAME; cd $NAME
    GOOS=$OS GOARCH=$ARCH go build frogtoss.com/$PROJ
    sleep 1 # "file changed as we read it"
    cd ..
    tar zcvf $NAME.tar.gz $NAME
}

# NAME=$PROJ-$OS-$VERSION
# mkdir $NAME; cd $NAME
# GOOS=darwin GOARCH=amd64 go build frogtoss.com/rex2ansi
# cd ..
# tar zcvf $NAME.tar.gz $NAME


go_build "darwin" "amd64"  
go_build "linux" "amd64"
go_build "linux" "386"
go_build "windows" "amd64"
go_build "windows" "386"
