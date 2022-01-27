#!/bin/bash

PATH_PROJECT=$HOME/Dependency
PATH_BUILD=$PATH_PROJECT/build
export GOROOT=$PATH_BUILD/goroot
export PATH=$GOROOT/bin:$PATH
export GOPATH=$PATH_BUILD/gopath
export PATH=$GOPATH/bin:$PATH
export PATH=$PATH_BUILD/install/bin:$PATH
export PKG_CONFIG_PATH=$PATH_BUILD/install/lib/pkgconfig:$PKG_CONFIG_PATH