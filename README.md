# 2018-Dependency

## build llvm

```shell
wget http://releases.llvm.org/7.0.0/llvm-7.0.0.src.tar.xz
tar -xf llvm-7.0.0.src.tar.xz
wget http://releases.llvm.org/7.0.0/cfe-7.0.0.src.tar.xz
tar -xf cfe-7.0.0.src.tar.xz
mv cfe-7.0.0.src llvm-7.0.0.src/tools/clang
mkdir llvm
cd llvm
cmake -G "Unix Makefiles" -DLLVM_ENABLE_RTTI=ON -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/home/yhao016/data/18-Dependency/install/ ../llvm-7.0.0.src/
make -j16
make install
```

## Install Protocol Buffers C++

notice: remove the default protobuf and install new!

```shell
sudo apt -y install autoconf automake libtool curl make g++ unzip
git clone https://github.com/protocolbuffers/protobuf
cd protobuf
git checkout tags/v3.11.0
git submodule update --init --recursive
./autogen.sh
./configure --prefix=/home/yhao016/data/18-Dependency/install  --disable-shared
make -j12
make install
sudo ldconfig
```

## Install gRPC C++

```shell
sudo apt -y install build-essential autoconf libtool pkg-config libgflags-dev libgtest-dev libc++-dev
git clone -b v1.25.0 https://github.com/grpc/grpc
cd grpc
git submodule update --init --recursive
make HAS_SYSTEM_PROTOBUF=false -j12
make install prefix=/home/yhao016/data/18-Dependency/install
sudo ldconfig
```

## build 2018-Dependency
```shell
git clone git@github.com:ZHYfeng/2018-Dependency.git
cd 2018-Dependency
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

~/data/git/gopath/src/github.com/ZHYfeng/2018-Dependency/02-dependency/cmake-build-debug/tools/A2I/a2i -objdump=/home/yhao016/data/2018-Dependency/13-linux-clang-np/vmlinux.objdump


https://zenodo.org/record/5348989/files/Build_Linux_Kernel_Into_LLVM_Bitcode.zip


# result

## requirement

```shell script
sudo apt install python3-tk
python3 -m pip install --user matplotlib
python3 -m pip install --user python3-protobuf
python3 -m pip install --user protobuf
python3 -m pip install --user scipy
```

## statistically result

+ mean
+ median
+ max/min
+ confidence intervals
+ Mann Whitney U-test
+ Vargha and Delaney's A12 statistics

# All commemd

## run dra

```shell script
/data/yhao/git/2018-Dependency/build/tools/DRA/dra -asm=/data/yhao/benchmark/linux/16-linux-clang-np-bc-f/built-in.s -objdump=/data/yhao/benchmark/linux/13-linux-clang-np/vmlinux.objdump -staticRes=/data/yhao/benchmark/linux/result/taint_info/taint_info_tty_ioctl_serialize /data/yhao/benchmark/linux/16-linux-clang-np-bc-f/built-in.bc 2>&1 | tee ./result-cpp.log
```

## test

```shell script
/data/yhao/git/2018-Dependency/build/tools/DRA/dra -asm=/home/yuh/data/benchmark/linux/16-linux-clang-np-bc-f/drivers/tty/built-in.s -objdump=/data/yhao/benchmark/linux/13-linux-clang-np/vmlinux.objdump -staticRes=/data/yhao/git/work/result/taint_info/taint_info_tty_ioctl_serialize -port=127.0.0.1:22223 /home/yuh/data/benchmark/linux/16-linux-clang-np-bc-f/drivers/tty/built-in.bc 2>&1 | tee ./result-cpp.log
```

## run syzkaller

```shell script
sudo /data/yhao/git/gopath/src/github.com/google/syzkaller/bin/syz-manager -config=./json/my_clang_ptmx.json -debug 2>&1 | tee ./result-syzkaller.log
```

## get addr2line

```shell script
cat ./address.txt | addr2line -afi -e /data/yhao/benchmark/linux/12-linux-clang-np/vmlinux > address-result.txt
```

## start qemu

```shell script
qemu-system-x86_64   -kernel ~/data/benchmark/linux/13-linux-clang-np/arch/x86/boot/bzImage   -append "console=ttyS0 root=/dev/sda debug earlyprintk=serial slub_debug=QUZ"  -hda ./image/stretch.img   -net user,hostfwd=tcp::10021-:22 -net nic   -enable-kvm   -nographic   -m 2G   -smp 2 -s  -pidfile vm.pid   2>&1 | tee vm.log
```

## ssh to qemu

```shell script
ssh -i ~/data/benchmark/linux/image/stretch.id_rsa -p 10021 -o "StrictHostKeyChecking no" root@localhost
```

## scp to qemu

```shell script
scp -P 10021 -F "/dev/null" -o "UserKnownHostsFile=/dev/null" -o "BatchMode=yes" -o "IdentitiesOnly=yes" -o "StrictHostKeyChecking=no" -o "ConnectTimeout=10" -i ~/data/benchmark/linux/image/stretch.id_rsa -v root@localhost:/syz-executor .
```

2019/05/22 12:58:44 running command: qemu-system-x86_64 []string{"-m", "2048", "-smp", "2", "-net", "nic,model=e1000", "-net", "user,host=10.0.2.10,hostfwd=tcp::1569-:22", "-display", "none", "-serial", "stdio", "-no-reboot", "-enable-kvm", "-cpu", "host", "-hda", "/data/yhao/benchmark/linux/img/stretch.img", "-snapshot", "-kernel", "/data/yhao/benchmark/linux/12-linux-clang-np/arch/x86/boot/bzImage", "-append", "earlyprintk=serial oops=panic nmi_watchdog=panic panic=1 ftrace_dump_on_oops=orig_cpu rodata=n vsyscall=native net.ifnames=0 biosdevname=0 root=/dev/sda console=ttyS0 kvm-intel.nested=1 kvm-intel.unrestricted_guest=1 kvm-intel.vmm_exclusive=1 kvm-intel.fasteoi=1 kvm-intel.ept=1 kvm-intel.flexpriority=1 kvm-intel.vpid=1 kvm-intel.emulate_invalid_guest_state=1 kvm-intel.eptad=1 kvm-intel.enable_shadow_vmcs=1 kvm-intel.pml=1 kvm-intel.enable_apicv=1 "}

## remove log file

```shell script
sudo rm -fr ./cover_uncover.txt dependency.log result-cpp.log result-syzkaller.log data.txt data.bin coverage.bin ./workdir/
```

## run a2i

```shell script
/home/yuh/data/git/gopath/src/github.com/ZHYfeng/2018-Dependency/02-dependency/build/tools/A2I/a2i -asm=built-in.s -objdump=/home/yuh/data/benchmark/linux/13-linux-clang-np/vmlinux.objdump -staticRes=built-in.taint built-in.bc
```

## prepare kernel
```shell script

llvm-link -o built-in.bc arch/x86/kernel/head64.bc arch/x86/kernel/ebda.bc arch/x86/kernel/platform-quirks.bc init/built-in.bc usr/built-in.bc arch/x86/built-in.bc kernel/built-in.bc certs/built-in.bc mm/built-in.bc fs/built-in.bc ipc/built-in.bc security/built-in.bc crypto/built-in.bc block/built-in.bc lib/built-in.bc arch/x86/lib/built-in.bc drivers/built-in.bc sound/built-in.bc firmware/built-in.bc arch/x86/pci/built-in.bc arch/x86/power/built-in.bc arch/x86/video/built-in.bc net/built-in.bc virt/built-in.bc
llvm-dis built-in.bc
rm -rf built-in.s
cat `find -name "*.s"` >> built-in.s
/home/yuh/data/git/gopath/src/github.com/ZHYfeng/2018-Dependency/02-dependency/build/tools/A2L/a2l -objdump=vmlinux.objdump
```

```shell script
python3 ./main.py generate ~/data/work
python3 ./main.py read ~/data/work
```

go run 
