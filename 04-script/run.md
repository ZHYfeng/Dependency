# All commemd

## run dra

```shell script
/data/yhao/git/2018_dependency/build/tools/DRA/dra -asm=/data/yhao/benchmark/linux/16-linux-clang-np-bc-f/built-in.s -objdump=/data/yhao/benchmark/linux/13-linux-clang-np/vmlinux.objdump -staticRes=/data/yhao/benchmark/linux/result/taint_info/taint_info_tty_ioctl_serialize /data/yhao/benchmark/linux/16-linux-clang-np-bc-f/built-in.bc 2>&1 | tee ./result-cpp.log
```

## test

```shell script
/data/yhao/git/2018_dependency/build/tools/DRA/dra -asm=/home/yuh/data/benchmark/linux/16-linux-clang-np-bc-f/drivers/tty/built-in.s -objdump=/data/yhao/benchmark/linux/13-linux-clang-np/vmlinux.objdump -staticRes=/data/yhao/git/work/result/taint_info/taint_info_tty_ioctl_serialize -port=127.0.0.1:22223 /home/yuh/data/benchmark/linux/16-linux-clang-np-bc-f/drivers/tty/built-in.bc 2>&1 | tee ./result-cpp.log
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
sudo qemu-system-x86_64   -kernel /home/yhao/benchmark/linux/13-linux-clang-np/arch/x86/boot/bzImage   -append "console=ttyS0 root=/dev/sda debug earlyprintk=serial slub_debug=QUZ"  -hda ./img/stretch.img   -net user,hostfwd=tcp::10021-:22 -net nic   -enable-kvm   -nographic   -m 2G   -smp 2   -pidfile vm.pid   2>&1 | tee vm.log
```

## ssh to qemu

```shell script
ssh -i /data/yhao/benchmark/linux/img/stretch.id_rsa -p 10021 -o "StrictHostKeyChecking no" root@localhost
```

## scp to qemu

```shell script
scp -P 1569 -F "/dev/null" -o "UserKnownHostsFile=/dev/null" -o "BatchMode=yes" -o "IdentitiesOnly=yes" -o "StrictHostKeyChecking=no" -o "ConnectTimeout=10" -i /data/yhao/benchmark/linux/img/stretch.id_rsa -v /data/yhao/git/gopath/src/github.com/google/syzkaller/bin/linux_amd64/syz-executor root@localhost:/syz-executor
```

2019/05/22 12:58:44 running command: qemu-system-x86_64 []string{"-m", "2048", "-smp", "2", "-net", "nic,model=e1000", "-net", "user,host=10.0.2.10,hostfwd=tcp::1569-:22", "-display", "none", "-serial", "stdio", "-no-reboot", "-enable-kvm", "-cpu", "host", "-hda", "/data/yhao/benchmark/linux/img/stretch.img", "-snapshot", "-kernel", "/data/yhao/benchmark/linux/12-linux-clang-np/arch/x86/boot/bzImage", "-append", "earlyprintk=serial oops=panic nmi_watchdog=panic panic=1 ftrace_dump_on_oops=orig_cpu rodata=n vsyscall=native net.ifnames=0 biosdevname=0 root=/dev/sda console=ttyS0 kvm-intel.nested=1 kvm-intel.unrestricted_guest=1 kvm-intel.vmm_exclusive=1 kvm-intel.fasteoi=1 kvm-intel.ept=1 kvm-intel.flexpriority=1 kvm-intel.vpid=1 kvm-intel.emulate_invalid_guest_state=1 kvm-intel.eptad=1 kvm-intel.enable_shadow_vmcs=1 kvm-intel.pml=1 kvm-intel.enable_apicv=1 "}

## remove log file

```shell script
sudo rm -fr ./cover_uncover.txt dependency.log result-cpp.log result-syzkaller.log data.txt data.bin coverage.bin ./workdir/
```

## run a2i

```shell script
/home/yuh/data/git/gopath/src/github.com/ZHYfeng/2018_dependency/02-dependency/build/tools/A2I/a2i -asm=built-in.s -objdump=/home/yuh/data/benchmark/linux/13-linux-clang-np/vmlinux.objdump -staticRes=built-in.taint built-in.bc
```