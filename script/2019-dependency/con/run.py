#! /usr/bin/python

import multiprocessing
import os.path
import shutil
import signal
import socket
import subprocess
import time

total_cpu = multiprocessing.cpu_count() - 6

path_syzkaller = "/home/yuh/data/git/gopath/src/github.com/google/syzkaller"
file_syzkaller = path_syzkaller + "/bin/syz-manager"

dra_path = "/home/yuh/data/git/2018_dependency/build/tools/DRA/dra"

linux_dir = "/home/yuh/data/benchmark/linux/13-linux-clang-np"
kernel_path = linux_dir + "/arch/x86/boot/bzImage"
file_vmlinux_objdump = linux_dir + "/vmlinux.objdump"
driver_name = ""
static_file = driver_name + ".static"
asm_file = driver_name + ".s"
bc_file = driver_name + ".bc"

http_port = ""
drpc_port = ""
path_workdir = "workdir"

path_image = "/home/yuh/data/benchmark/linux/img"
image_file = "/stretch.img"
sshkey_file = "/stretch.id_rsa"

file_log_syzkaller = "log_syzkaller.txt"
file_log_dra = "log_dra.txt"

time = 1 * 24 * 60 * 60  # second


def get_open_port():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("", 0))
    s.listen(1)
    port = s.getsockname()[1]
    s.close()
    return port


class Process:
    def __init__(self):
        self.cmd_dra = dra_path + " -asm=./built-in.s" + " -objdump=" + file_vmlinux_objdump \
                       + " -staticRes=./built-in.taint" + "./built-in.bc" + " 2>&1 1>" + file_log_dra
        self.cmd_syzkaller = file_syzkaller + " -config=./built-in.json" + " 2>&1 1>" + file_log_syzkaller
        self.index = 0
        while os.path.exists(self.index):
            self.index = self.index + 1
        os.mkdir(self.index)
        os.chdir(self.index)
        shutil.copy(path_image, "./img")
        dirpath = os.getcwd()

    def execute(self):
        self.execute_syzkaller()
        self.execute_dra()

    def execute_dra(self):
        self.p_dra = subprocess.Popen(self.cmd_dra, shell=True, preexec_fn=os.setsid)

    def execute_syzkaller(self):
        self.t0 = time.time()
        self.p_syzkaller = subprocess.Popen(self.cmd_syzkaller, shell=True, preexec_fn=os.setsid)

    def close(self):
        os.killpg(os.getpgid(self.syzkaller.pid), signal.SIGTERM)
        os.killpg(os.getpgid(self.syzkaller.pid), signal.SIGTERM)


tasks = [Process() for i in range(total_cpu)]


def main():
    print("")


if __name__ == "__main__":
    main()
