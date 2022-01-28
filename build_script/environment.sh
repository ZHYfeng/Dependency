#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
echo $SCRIPT_DIR

PATH_PROJECT=$SCRIPT_DIR/..
PATH_BUILD=$PATH_PROJECT/build
PATH_INSTALL=$PATH_PROJECT/install
export GOROOT=$PATH_BUILD/goroot
export PATH=$GOROOT/bin:$PATH
export GOPATH=$PATH_BUILD/gopath
export PATH=$GOPATH/bin:$PATH
export PATH=$PATH_INSTALL/bin:$PATH
export PKG_CONFIG_PATH=$PATH_INSTALL/lib/pkgconfig:$PKG_CONFIG_PATH
export PATH_PROJECT=$PATH_PROJECT