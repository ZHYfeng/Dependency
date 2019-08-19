#! /usr/bin/python
import json
import multiprocessing
import os.path
import shutil
import signal
import socket
import subprocess
import time

number_execute = 6
path_root = os.getcwd()

path_syzkaller = "/home/yuh/data/git/gopath/src/github.com/google/syzkaller"
file_syzkaller = os.path.join(path_syzkaller, "bin/syz-manager")

path_dra = "/home/yuh/data/git/2018_dependency/build/tools/DRA/dra"

path_linux = "/home/yuh/data/benchmark/linux/13-linux-clang-np"
path_kernel = os.path.join(path_linux, "arch/x86/boot/bzImage")
file_vmlinux_objdump = os.path.join(path_linux, "vmlinux.objdump")

name_driver = "built-in"
file_taint = name_driver + ".taint"
file_asm = name_driver + ".s"
file_bc = name_driver + ".bc"
file_json = name_driver + ".json"

path_workdir = "workdir"

path_image = "/home/yuh/data/benchmark/linux/img"
file_image = "stretch.img"
file_ssh_key = "stretch.id_rsa"

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
        self.cmd_dra = ""
        self.cmd_syzkaller = ""
        self.drpc = "0"
        self.index = 0
        while os.path.exists(str(self.index)):
            self.index = self.index + 1
        os.makedirs(str(self.index))

        shutil.copytree(path_image, os.path.join(path_root, str(self.index), "img"))

        f = open(os.path.join(path_root, file_json), "r")
        c = json.load(f)
        f.close()

        c["workdir"] = os.path.join(path_root, str(self.index), path_workdir)
        c["image"] = os.path.join(path_root, str(self.index), "img", file_image)
        c["sshkey"] = os.path.join(path_root, str(self.index), "img", file_ssh_key)

        f = open(os.path.join(path_root, str(self.index), file_json), "w")
        json.dump(c, f, indent=4)
        f.close()
        f = open(os.path.join(path_root, file_json), "w")
        json.dump(c, f, indent=4)
        f.close()

    def execute(self):
        self.execute_syzkaller()
        self.execute_dra()

    def execute_syzkaller(self):
        f = open(os.path.join(path_root, str(self.index), file_json), "r")
        c = json.load(f)
        f.close()
        c["http"] = "127.0.0.1:" + get_open_port()
        self.drpc = "127.0.0.1:" + get_open_port()
        c["drpc"] = self.drpc
        f = open(os.path.join(path_root, str(self.index), file_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        self.cmd_syzkaller = file_syzkaller + " -config=./" + file_json + " 2>&1 1>" + file_log_syzkaller
        self.t0 = time.time()
        # self.p_syzkaller = subprocess.Popen(self.cmd_syzkaller, shell=True, preexec_fn=os.setsid)

    def execute_dra(self):
        self.cmd_dra = path_dra + " -asm=" + file_asm + " -objdump=" + file_vmlinux_objdump \
                       + " -staticRes=" + file_taint + " -port=" + self.drpc \
                       + file_bc + " 2>&1 1>" + file_log_dra
        self.p_dra = subprocess.Popen(self.cmd_dra, shell=True, preexec_fn=os.setsid)

    def close(self):
        os.killpg(os.getpgid(self.p_syzkaller.pid), signal.SIGTERM)
        os.killpg(os.getpgid(self.p_dra.pid), signal.SIGTERM)


def main():
    tasks = [Process() for i in range(number_execute)]
    for i in tasks:
        i.execute()

    time.sleep(time)

    for i in tasks:
        i.close()


if __name__ == "__main__":
    main()
