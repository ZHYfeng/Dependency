#! /usr/bin/python3
import json
import os.path
import signal
import socket
import subprocess
import sys
import time

import default


# from default import default


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
            name = default.name_with_dra
        else:
            name = default.name_without_dra

        path = os.path.join(default.path_current, name, str(self.index))
        while os.path.exists(path):
            self.index = self.index + 1
            path = os.path.join(default.path_current, name, str(self.index))
        self.path = path
        os.makedirs(self.path)
        print(self.path)

        cmd_cp_img = "cp -rf " + default.path_image + " " + os.path.join(self.path, "img")
        p_cp_img = subprocess.Popen(cmd_cp_img, shell=True, preexec_fn=os.setsid)
        p_cp_img.wait()

        if os.path.exists(default.path_workdir):
            cmd_cp_corpus = "cp -rf " + default.path_workdir + " " + os.path.join(self.path, default.name_workdir)
            p_cp_corpus = subprocess.Popen(cmd_cp_corpus, shell=True, preexec_fn=os.setsid)
            p_cp_corpus.wait()

        cmd_cp_built_in = "cp ./built-in.* " + self.path
        p_cp_built_in = subprocess.Popen(cmd_cp_built_in, shell=True, preexec_fn=os.setsid)
        p_cp_built_in.wait()

        cmd_cp_built_in = "cp ./*_serialize " + self.path
        p_cp_built_in = subprocess.Popen(cmd_cp_built_in, shell=True, preexec_fn=os.setsid)
        p_cp_built_in.wait()

        cmd_cp_built_in = "cp " + default.name_dra_json + " " + self.path
        p_cp_built_in = subprocess.Popen(cmd_cp_built_in, shell=True, preexec_fn=os.setsid)
        p_cp_built_in.wait()

        f = open(os.path.join(default.path_current, default.name_syzkaller_json), "r")
        c = json.load(f)
        f.close()

        c["workdir"] = os.path.join(self.path, default.name_workdir)
        c["image"] = os.path.join(self.path, "img", default.file_image)
        c["sshkey"] = os.path.join(self.path, "img", default.file_ssh_key)

        f = open(os.path.join(self.path, default.name_syzkaller_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        run_f.write("cd " + path + "\n")
        os.chdir(self.path)
        self.execute_syzkaller(run_f, dependency_priority=dra)
        self.execute_dra(run_f)

    def execute_syzkaller(self, run_f, dependency_task=False, dependency_priority=False):
        f = open(os.path.join(self.path, default.name_syzkaller_json), "r")
        c = json.load(f)
        f.close()
        c["http"] = "127.0.0.1:" + str(get_open_port())
        self.drpc = "127.0.0.1:" + str(get_open_port())
        c["drpc"] = self.drpc
        c["dependency_task"] = dependency_task
        c["dependency_priority"] = dependency_priority
        f = open(os.path.join(self.path, default.name_syzkaller_json), "w")
        json.dump(c, f, indent=4)
        f.close()

        self.cmd_syzkaller = default.path_syzkaller_manager + " -config=./" + default.name_syzkaller_json + " 2>" + default.file_log_syzkaller + " 1>&2 &"
        self.t0 = time.time()
        self.real_execute(self.cmd_syzkaller, run_f)

    def execute_dra(self, run_f):
        self.cmd_dra = default.path_dra + " -asm=" + default.name_asm + " -objdump=" + default.file_vmlinux_objdump \
                       + " -bc=" + default.name_bc + " -port=" + self.drpc \
                       + " " + default.name_dra_json + " 1>" + default.file_log_dra + " 2>&1 &"
        self.real_execute(self.cmd_dra, run_f)

    def real_execute(self, cmd, run_f):
        p = subprocess.Popen(cmd, shell=True)
        # p = subprocess.Popen(cmd, shell=True, start_new_session=True)
        self.processes.append(p)
        run_f.write(cmd + "\n")
        run_f.write("PID+=(\"$!\")\n")
        f = open(os.path.join(self.path, default.file_log_run), "a")
        f.write("#!/bin/bash\n\n")
        f.write(cmd + "\n")
        f.write("#dra pid : " + str(p.pid) + "\n")
        f.close()

    def close(self):
        for p in self.processes:
            os.killpg(os.getpgid(p.pid), signal.SIGTERM)

    def remove(self):
        os.chdir(self.path)
        cmd_rm_img = "rm -rf img " + " " + default.name_asm + " " + default.name_bc
        p_rm_img = subprocess.Popen(cmd_rm_img, shell=True, preexec_fn=os.setsid)
        p_rm_img.wait()


def main():
    dra = True
    if len(sys.argv) > 1:
        dra = False

    if dra:
        path_run = os.path.join(default.path_current, default.name_with_dra, default.file_run)
    else:
        path_run = os.path.join(default.path_current, default.name_without_dra, default.file_run)
    run_f = open(path_run, "a")
    run_f.write("#!/bin/bash\n\n")
    run_f.write("PID=()\n")

    tasks = [Process() for i in range(default.number_execute)]
    for i in tasks:
        i.execute(run_f, dra)

    run_f.write("sleep " + str(default.time_run) + "\n")
    run_f.write("kill -SIGKILL ${PID[@]}\n")
    run_f.close()

    cmd_ch = "chmod a+x " + path_run
    p_ch = subprocess.Popen(cmd_ch, shell=True, preexec_fn=os.setsid)
    p_ch.wait()

    time.sleep(default.time_run)
    # time.sleep(30)

    for i in tasks:
        i.remove()

    for i in tasks:
        i.close()


if __name__ == "__main__":
    main()
