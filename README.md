# 2018_dependency

## Install Protocol Buffers C++

notice: remove the default protobuf and install new!

```shell
git clone https://github.com/protocolbuffers/protobuf.git
cd protobuf
git submodule update --init --recursive
./autogen.sh
./configure --prefix=/usr
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

The binary is in build/tools/DRA/dra.

If meet error `./tools/DRA/dra: error while loading shared libraries: libgrpc++.so.1: cannot open shared object file: No such file or directory`.

run `export LD_LIBRARY_PATH="/usr/local/lib"`

Because gRPC will install in `/usr/local/lib` but default it only search in `/usr/lib`.

run `build/tools/DRA/dra --help` get the usage.

`input bitcode` is necessary.