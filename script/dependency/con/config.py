#! /usr/bin/python3
import json
import os
import shutil
import sys

sys.path.append(os.getcwd())
from script.dependency.con import devices


def read_s(paths):
    ctx_s = ""
    for path in paths:
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.endswith(".s") and file_name != devices.file_asm:
                    print(os.path.join(dir_path, file_name))
                    f = open(os.path.join(dir_path, file_name), "r")
                    ctx_s = ctx_s + "\n" + f.read()
                    f.close()
    f = open("built-in" + ".s", "w")
    f.write(ctx_s)
    f.close()
    # cat `find -name "*.s"` >> built-in.s


def generate_dev_dir():
    for d in devices.dev:
        path = os.path.join(devices.path_result, d)
        print(path)
        if os.path.exists(path):
            shutil.rmtree(path)
        if not os.path.exists(path):
            os.makedirs(path)
        os.chdir(path)

        if not os.path.exists("built-in.json"):
            df = open(devices.file_default_json, "r")
            c = json.load(df)
            df.close()
            c["enable_syscalls"] = devices.dev[d]["enable_syscalls"]
            f = open("built-in" + ".json", "w")
            json.dump(c, f, indent=4)
            f.close()

        if not os.path.exists(devices.file_asm):
            read_s(devices.dev[d]["path_s"])
        if not os.path.exists(devices.file_bc):
            shutil.copy(devices.dev[d]["file_bc"], devices.file_bc)
        if not os.path.exists(devices.file_taint):
            shutil.copy(devices.dev[d]["file_taint"], devices.file_taint)
        if not os.path.exists(devices.name_with_dra):
            os.makedirs(devices.name_with_dra)
        if not os.path.exists(devices.name_without_dra):
            os.makedirs(devices.name_without_dra)
        if not os.path.exists(devices.name_run):
            shutil.copy(devices.path_default_run, devices.name_run)


if __name__ == "__main__":
    generate_dev_dir()
