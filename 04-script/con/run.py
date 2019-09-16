#! /usr/bin/python3
import json
import os.path
import signal
import socket
import subprocess
import sys
import time

path_home = os.path.expanduser("~")
path_root = os.path.join(path_home, "data")
number_execute = 6
path_current = os.getcwd()
path_git = os.path.join(path_root, "git")
path_repo = os.path.join(path_git, "gopath/src/github.com/ZHYfeng/2018-dependency")
path_dra = os.path.join(path_repo, "02-dependency/build/tools/DRA/dra")
path_syzkaller = os.path.join(path_repo, "03-syzkaller/bin/syz-manager")
path_linux = os.path.join(path_root, "benchmark/linux/13-linux-clang-np")
path_kernel = os.path.join(path_linux, "arch/x86/boot/bzImage")
file_vmlinux_objdump = os.path.join(path_linux, "vmlinux.objdump")

name_with_dra = "01-result-with-dra"
name_without_dra = "02-result-without-dra"

name_driver = "built-in"
file_taint = name_driver + ".taint"
file_asm = name_driver + ".s"
file_bc = name_driver + ".bc"
file_json = name_driver + ".json"

name_workdir = "workdir"
path_workdir = os.path.join(path_current, name_workdir)

path_image = os.path.join(path_root, "benchmark/linux/image")
file_image = "stretch.img"
file_ssh_key = "stretch.id_rsa"

file_log_run = "log_run.bash"
file_log_syzkaller = "log_syzkaller.log"
file_log_dra = "log_dra.log"
file_run = "run.bash"

# time_run = 1 * 24 * 60 * 60  # second
time_run = 1 * 12 * 60 * 60  # second


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

    def execute(self, run_f, dra=True):
        if dra:
            name = name_with_dra
        else:
            name = name_without_dra

        path = os.path.join(path_current, name, str(self.index))
        while os.path.exists(path):
            self.index = self.index + 1
            path = os.path.join(path_current, name, str(self.index))
        self.path = path
        os.makedirs(self.path)
        print(self.path)

        cmd_cp_img = "cp -rf " + path_image + " " + os.path.join(self.path, "img")
        p_cp_img = subprocess.Popen(cmd_cp_img, shell=True, preexec_fn=os.setsid)
        p_cp_img.wait()

        if os.path.exists(path_workdir):
            cmd_cp_corpus = "cp -rf " + path_workdir + " " + os.path.join(self.path, name_workdir)
            p_cp_corpus = subprocess.Popen(cmd_cp_corpus, shell=True, preexec_fn=os.setsid)
            p_cp_corpus.wait()

        if dra:
            cmd_cp_built_in = "cp ./built-in.* " + self.path
            p_cp_built_in = subprocess.Popen(cmd_cp_built_in, shell=True, preexec_fn=os.setsid)
            p_cp_built_in.wait()

        f = open(os.path.join(path_current, file_json), "r")
        c = json.load(f)
        f.close()

        c["workdir"] = os.path.join(self.path, name_workdir)
        c["image"] = os.path.join(self.path, "img", file_image)
        c["sshkey"] = os.path.join(self.path, "img", file_ssh_key)

        f = open(os.path.join(self.path, file_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        run_f.write("cd " + path + "\n")
        os.chdir(self.path)
        self.execute_syzkaller(run_f)
        if dra:
            self.execute_dra(run_f)

    def execute_syzkaller(self, run_f):
        f = open(os.path.join(self.path, file_json), "r")
        c = json.load(f)
        f.close()
        c["http"] = "127.0.0.1:" + str(get_open_port())
        self.drpc = "127.0.0.1:" + str(get_open_port())
        c["drpc"] = self.drpc
        f = open(os.path.join(self.path, file_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        self.cmd_syzkaller = path_syzkaller + " -config=./" + file_json + " 2>" + file_log_syzkaller + " 1>&2 &"
        self.t0 = time.time()
        self.real_execute(self.cmd_syzkaller, run_f)

    def execute_dra(self, run_f):
        self.cmd_dra = path_dra + " -asm=" + file_asm + " -objdump=" + file_vmlinux_objdump \
                       + " -staticRes=" + file_taint + " -port=" + self.drpc \
                       + " " + file_bc + " 1>" + file_log_dra + " 2>&1 &"
        self.real_execute(self.cmd_dra, run_f)

    def real_execute(self, cmd, run_f):
        p = subprocess.Popen(cmd, shell=True)
        # p = subprocess.Popen(cmd, shell=True, start_new_session=True)
        self.processes.append(p)
        run_f.write(cmd + "\n")
        run_f.write("PID+=(\"$!\")\n")
        f = open(os.path.join(self.path, file_log_run), "a")
        f.write("#!/bin/bash\n\n")
        f.write(cmd + "\n")
        f.write("#dra pid : " + str(p.pid) + "\n")
        f.close()

    def close(self):
        for p in self.processes:
            os.killpg(os.getpgid(p.pid), signal.SIGTERM)
        os.chdir(self.path)
        cmd_rm_img = "rm -rf img " + file_taint + " " + file_asm + " " + file_bc
        p_rm_img = subprocess.Popen(cmd_rm_img, shell=True, preexec_fn=os.setsid)
        p_rm_img.wait()


def main():
    dra = True
    if len(sys.argv) > 1:
        dra = False

    if dra:
        path_run = os.path.join(path_current, name_with_dra, file_run)
    else:
        path_run = os.path.join(path_current, name_without_dra, file_run)
    run_f = open(path_run, "a")
    run_f.write("#!/bin/bash\n\n")
    run_f.write("PID=()\n")

    tasks = [Process() for i in range(number_execute)]
    for i in tasks:
        i.execute(run_f, dra)

    run_f.write("sleep " + str(time_run) + "\n")
    run_f.write("kill -SIGKILL ${PID[@]}\n")
    run_f.close()

    cmd_ch = "chmod a+x " + path_run
    p_ch = subprocess.Popen(cmd_ch, shell=True, preexec_fn=os.setsid)
    p_ch.wait()

    time.sleep(time_run)
    # time.sleep(30)

    for i in tasks:
        i.close()


if __name__ == "__main__":
    main()
