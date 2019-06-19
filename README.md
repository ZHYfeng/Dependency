# 2018_dependency

## Install Protocol Buffers C++

notice: remove the default protobuf and install new!

```shell
git clone https://github.com/protocolbuffers/protobuf.git
cd protobuf
git submodule update --init --recursive
./autogen.sh
./configure
make
make check
sudo make install
sudo ldconfig
```

## Install gRPC C++

```shell
sudo apt-get install build-essential autoconf libtool pkg-config libgflags-dev libgtest-dev libc++-dev
git clone -b $(curl -L https://grpc.io/release) https://github.com/grpc/grpc
cd grpc
git submodule update --init
// remove the  -Werror in Makefile line 356
make
sudo make install
sudo ldconfig
```

## build 2018_dependency

```shell
git clone git@github.com:ZHYfeng/2018_dependency.git
cd 2018_dependency
mkdir build
cd build
cmake ..
make
```

or

```shell
git clone git@github.com:ZHYfeng/2018_dependency.git
cd 2018_dependency
./build.sh
```
notice: change LLVM_DIR to your path.

The binary is in build/tools/DRA/dra.

If meet error `./tools/DRA/dra: error while loading shared libraries: libgrpc++.so.1: cannot open shared object file: No such file or directory`.

run `export LD_LIBRARY_PATH="/usr/local/lib"`

Because gRPC will install in `/usr/local/lib` but default it only search in `/usr/lib`.

run `build/tools/DRA/dra --help` get the usage.

`input bitcode` is necessary.

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

The path of all result: `/data/yhao/benchmark/linux/result` in terran.
There are some config json of syzkaller and static analysis result.

start dra: 
```
/data/yhao/git/2018_dependency/build/tools/DRA/dra -asm=/data/yhao/benchmark/linux/12-linux-clang-np/built-in.s -objdump=/data/yhao/benchmark/linux/12-linux-clang-np/vmlinux.objdump -staticRes=/data/yhao/benchmark/linux/result/taint_info_cdrom_ioctl_serialize /data/yhao/benchmark/linux/12-linux-clang-np/built-in.bc 2>&1 | tee /data/yhao/benchmark/linux/result/result-cpp.log
```
start syz-manager: 
```
sudo /data/yhao/git/gopath/src/github.com/google/syzkaller/bin/syz-manager -config=my_clang_cdrom.cfg -debug 2>&1 | tee result-syzkaller.log
```
Those path is in terran, change the path if necessary.


The path of Linux kernel: `/data/yhao/benchmark/linux`.

`12-linux-clang-np/` is the Linux kernel built with clang.

`15-linux-clang-np-bc-f/` is `12-linux-clang-np/` with bc file.
now we can use this.

`13-linux-clang-np/` is the Linux kernel with new driver.

`16-linux-clang-np-bc-f/` is `13-linux-clang-np/` with bc file.

I will deal with it and get some files in the future.

