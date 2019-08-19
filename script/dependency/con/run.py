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
name_with_dra = "result-with-dra"
name_without_dra = "result-without-dra"

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

path_image = "/home/yuh/data/benchmark/linux/image"
file_image = "stretch.img"
file_ssh_key = "stretch.id_rsa"

file_log_syzkaller = "log_syzkaller.txt"
file_log_dra = "log_dra.txt"

time_run = 1 * 24 * 60 * 60  # second


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
        self.path = ""

    def execute(self, dra=True):
        if dra:
            name = name_with_dra
        else:
            name = name_without_dra

        path = os.path.join(path_root, name, str(self.index))
        while os.path.exists(path):
            self.index = self.index + 1
            path = os.path.join(path_root, name, str(self.index))
        self.path = path
        os.makedirs(self.path)

        print(os.path.join(path, "img"))
        cmd_cp_img = "cp -rf " + path_image + " " + os.path.join(self.path, "img")
        p_cp_img = subprocess.Popen(cmd_cp_img, shell=True, preexec_fn=os.setsid)
        p_cp_img.wait()

        if dra:
            cmd_cp_built_in = "cp ./built-in.* " + self.path
            p_cp_built_in = subprocess.Popen(cmd_cp_built_in, shell=True, preexec_fn=os.setsid)
            p_cp_built_in.wait()

        f = open(os.path.join(path_root, file_json), "r")
        c = json.load(f)
        f.close()

        c["workdir"] = os.path.join(self.path, path_workdir)
        c["image"] = os.path.join(self.path, "img", file_image)
        c["sshkey"] = os.path.join(self.path, "img", file_ssh_key)

        f = open(os.path.join(self.path, file_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        os.chdir(self.path)
        self.execute_syzkaller()
        if dra:
            self.execute_dra()

    def execute_syzkaller(self):
        f = open(os.path.join(self.path, file_json), "r")
        c = json.load(f)
        f.close()
        c["http"] = "127.0.0.1:" + str(get_open_port())
        self.drpc = "127.0.0.1:" + str(get_open_port())
        c["drpc"] = self.drpc
        f = open(os.path.join(self.path, file_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        self.cmd_syzkaller = file_syzkaller + " -config=./" + file_json + " 2>" + file_log_syzkaller + " 1>&2"
        self.t0 = time.time()
        self.p_syzkaller = subprocess.Popen(self.cmd_syzkaller, shell=True, preexec_fn=os.setsid)

    def execute_dra(self):
        self.cmd_dra = path_dra + " -asm=" + file_asm + " -objdump=" + file_vmlinux_objdump \
                       + " -staticRes=" + file_taint + " -port=" + self.drpc \
                       + file_bc + " 1>" + file_log_dra + " 2>&1"
        self.p_dra = subprocess.Popen(self.cmd_dra, shell=True, preexec_fn=os.setsid)

    def close(self, dra=True):
        os.killpg(os.getpgid(self.p_syzkaller.pid), signal.SIGTERM)
        if dra:
            os.killpg(os.getpgid(self.p_dra.pid), signal.SIGTERM)


def main():
    dra = False
    tasks = [Process() for i in range(number_execute)]
    for i in tasks:
        i.execute(dra)

    # time.sleep(time_run)
    time.sleep(30)

    for i in tasks:
        i.close(dra)


if __name__ == "__main__":
    main()
