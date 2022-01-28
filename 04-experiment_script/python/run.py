#! /usr/bin/python3
import json
import os
import signal
import socket
import subprocess
import sys
import time


time_run = 2 * 24 * 60 * 60  # second
number_execute = 1
number_vm_count = 2

path_project = os.environ.get('PATH_PROJECT')
path_syzkaller = os.path.join(path_project, "03-syzkaller")
path_image = os.path.join(path_project, "workdir/image")
path_linux = os.path.join(path_project, "workdir/13-linux-clang-np")
path_kernel = os.path.join(path_linux, "arch/x86/boot/bzImage")
path_vmlinux_objdump = os.path.join(path_linux, "vmlinux.objdump")

name_driver = "built-in"
name_asm = name_driver + ".s"
name_bc = name_driver + ".bc"
name_dra_json = "dra.json"
name_syzkaller_json = "syzkaller.json"
file_log_syzkaller = "log_syzkaller.log"
file_log_dra = "log_dra.log"

path_current = os.getcwd()
path_dra = "dra"
path_syzkaller_manager = os.path.join(path_syzkaller, "bin/syz-manager")

name_workdir = "workdir"
file_image = "stretch.img"
file_ssh_key = "stretch.id_rsa"

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
        self.processes = []

    def execute(self, dra=False):
        path = os.path.join(path_current, str(self.index))
        while os.path.exists(path):
            self.index = self.index + 1
            path = os.path.join(path_current, str(self.index))
        self.path = path
        os.makedirs(self.path)
        print(self.path)

        cmd_cp_img = "cp -rf " + path_image + " " + os.path.join(self.path, "img")
        p_cp_img = subprocess.Popen(cmd_cp_img, shell=True, preexec_fn=os.setsid)
        p_cp_img.wait()

        f = open(os.path.join(path_current, name_syzkaller_json), "r")
        c = json.load(f)
        f.close()

        c["workdir"] = os.path.join(self.path, name_workdir)
        c["image"] = os.path.join(self.path, "img", file_image)
        c["sshkey"] = os.path.join(self.path, "img", file_ssh_key)
        c["syzkaller"] = path_syzkaller
        c["kernel_obj"] = path_linux
        c["vm"]["kernel"] = path_kernel
        c["vm"]["count"] = number_vm_count

        f = open(os.path.join(self.path, name_syzkaller_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        os.chdir(self.path)
        self.execute_syzkaller(dependency_task=dra)
        self.execute_dra()

    def execute_syzkaller(self, dependency_task=False, dependency_priority=False):
        f = open(os.path.join(self.path, name_syzkaller_json), "r")
        c = json.load(f)
        f.close()
        c["http"] = "127.0.0.1:" + str(get_open_port())
        self.drpc = "127.0.0.1:" + str(get_open_port())
        c["drpc"] = self.drpc
        c["dependency_task"] = dependency_task
        c["dependency_priority"] = dependency_priority
        f = open(os.path.join(self.path, name_syzkaller_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        self.cmd_syzkaller = path_syzkaller_manager + " -config=./" + name_syzkaller_json + " 2>" + file_log_syzkaller + " 1>&2 &"
        self.t0 = time.time()
        self.real_execute(self.cmd_syzkaller)

    def execute_dra(self):
        self.cmd_dra = path_dra + " -asm=../" + name_asm + " -objdump=" + path_vmlinux_objdump \
                       + " -bc=../" + name_bc + " -port=" + self.drpc \
                       + " ../" + name_dra_json + " 1>" + file_log_dra + " 2>&1 &"
        self.real_execute(self.cmd_dra)

    def real_execute(self, cmd):
        p = subprocess.Popen(cmd, shell=True)
        # p = subprocess.Popen(cmd, shell=True, start_new_session=True)
        self.processes.append(p)

    def close(self):
        for p in self.processes:
            os.killpg(os.getpgid(p.pid), signal.SIGTERM)


def main():
    dra = True
    if len(sys.argv) > 1:
        dra = False

    tasks = [Process() for i in range(number_execute)]
    for i in tasks:
        i.execute(dra)

    time.sleep(time_run)
    # time.sleep(30)

    for i in tasks:
        i.close()


if __name__ == "__main__":
    main()
