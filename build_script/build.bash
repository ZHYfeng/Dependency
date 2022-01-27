#!/bin/bash

export PATH_PROJECT=$PWD

# apt install
sudo apt install -y llvm-7 clang-7
sudo apt install -y make gcc flex bison libncurses-dev libelf-dev libssl-dev
sudo apt install -y qemu-system-x86
sudo apt install -y git cmake build-essential autoconf libtool pkg-config
sudo usermod -a -G kvm `whoami`

# make dir for project
## build: goroot, gopath and grpc
## install: INSTALL_PREFIX of grpc and https://github.com/ZHYfeng/Dependency.git
PATH_BUILD=$PATH_PROJECT/build
PATH_INSTALL=$PATH_PROJECT/install
if [ -d $PATH_BUILD ]
then
    echo "[*] $PATH_BUILD exist"
else 
    mkdir $PATH_BUILD
fi
if [ -d $PATH_INSTALL ]
then
    echo "[*] $PATH_INSTALL exist"
else 
    mkdir $PATH_INSTALL
fi

# get golang
cd $PATH_BUILD
wget https://go.dev/dl/go1.17.linux-amd64.tar.gz
tar -xf go1.17.linux-amd64.tar.gz
mv go goroot
mkdir gopath
export GOROOT=$PATH_BUILD/goroot
export PATH=$GOROOT/bin:$PATH
export GOPATH=$PATH_BUILD/gopath
export PATH=$GOPATH/bin:$PATH

# build grpc
cd $PATH_BUILD
git clone --recurse-submodules -b v1.42.0 https://github.com/grpc/grpc
mkdir -p grpc-build
pushd grpc-build
cmake -DgRPC_INSTALL=ON \
    -DgRPC_BUILD_TESTS=OFF \
    -DCMAKE_INSTALL_PREFIX=$PATH_INSTALL \
    ../grpc
make -j
make install
popd
sudo ldconfig
export PATH=$PATH_INSTALL/bin:$PATH
export PKG_CONFIG_PATH=$PATH_INSTALL/lib/pkgconfig:$PKG_CONFIG_PATH
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

# update proto code
cd $PATH_PROJECT/05-proto
bash ./build.bash

# build DRA
cd $PATH_PROJECT/02-dependency
bash ./build.bash

# build syzkaller
cd $PATH_PROJECT/03-syzkaller
bash ./build.bash

# generate environment.sh
cd $PATH_PROJECT
echo "export GOROOT=$PATH_BUILD/goroot" >> environment.sh
echo "export PATH=\$GOROOT/bin:\$PATH" >> environment.sh
echo "export GOPATH=$PATH_BUILD/gopath" >> environment.sh
echo "export PATH=\$GOPATH/bin:\$PATH" >> environment.sh
echo "export PATH=$PATH_INSTALL/bin:\$PATH" >> environment.sh
echo "export PKG_CONFIG_PATH=$PATH_INSTALL/lib/pkgconfig:\$PKG_CONFIG_PATH" >> environment.sh