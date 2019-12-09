# 2018_dependency

## Install Protocol Buffers C++

notice: remove the default protobuf and install new!

```shell
sudo apt -y install autoconf automake libtool curl make g++ unzip
git clone https://github.com/protocolbuffers/protobuf.git
cd protobuf
git submodule update --init --recursive
./autogen.sh
./configure --prefix=~/data/build
make
make check
make install
sudo ldconfig
```

## Install gRPC C++

```shell
sudo apt -y install build-essential autoconf libtool pkg-config libgflags-dev libgtest-dev libc++-dev
git clone -b $(curl -L https://grpc.io/release) https://github.com/grpc/grpc
cd grpc
git submodule update --init
make HAS_SYSTEM_PROTOBUF=false
make install prefix=~/data/build
sudo ldconfig
```

## build 2018_dependency
```shell
git clone git@github.com:ZHYfeng/2018_dependency.git
cd 2018_dependency
bash build.bash
```

notice: change LLVM_DIR to your path.
The binary is in build/tools/DRA/dra.

run `build/tools/DRA/dra --help` get the usage.

## build syzkaller

### install go
```
wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
tar -xf go1.12.6.linux-amd64.tar.gz
mv go goroot
export GOROOT=`pwd`/goroot
export PATH=$GOROOT/bin:$PATH
mkdir gopath
export GOPATH=`pwd`/gopath
```
### install protobuf and grpc
```
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u google.golang.org/grpc
```
### install syzkaller
```
git clone git@github.com:ZHYfeng/2019_syzkaller.git
mv 2019_syzkaller gopath/src/github.com/google/syzkaller/
cd gopath/src/github.com/google/syzkaller/
mkdir workdir
cd ./pkg/dra
./build-protoc.sh
cd ../..
make
cd workdir
```

## how to use

The path of Linux kernel: `/home/yuh/data/benchmark/linux`.

`12-=linux-clang-np/` is the Linux kernel built with clang.
`15-linux-clang-np-bc-f/` is `12-linux-clang-np/` with bc file.

`13-linux-clang-np/` is the Linux kernel with new driver.
`16-linux-clang-np-bc-f/` is `13-linux-clang-np/` with bc file.
