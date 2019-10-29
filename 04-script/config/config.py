#! /usr/bin/python3
import json
import os
import shutil

from config import default


def read_s(paths):
    ctx_s = ""
    for path in paths:
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.endswith(".s") and file_name != default.file_asm:
                    f = open(os.path.join(dir_path, file_name), "r")
                    ctx_s = ctx_s + "\n" + f.read()
                    f.close()
    f = open("built-in" + ".s", "w")
    f.write(ctx_s)
    f.close()
    # cat `find -name "*.s"` >> built-in.s


def generate_dev_dir():
    for d in default.dev:
        path = os.path.join(default.path_result, d)
        print(path)
        # if os.path.exists(path):
        #     shutil.rmtree(path)
        if not os.path.exists(path):
            os.makedirs(path)
        os.chdir(path)

        # if not os.path.exists("built-in.json"):
        #     df = open(default.file_default_json, "r")
        #     c = json.load(df)
        #     df.close()
        #     c["enable_syscalls"] = default.dev[d]["enable_syscalls"]
        #     c["syzkaller"] = default.path_syzkaller
        #     c["kernel_obj"] = default.path_linux
        #     c["vm"]["kernel"] = default.path_kernel
        #     f = open("built-in" + ".json", "w")
        #     json.dump(c, f, indent=4)
        #     f.close()
        
        # if not os.path.exists(default.file_asm):
        #     read_s(default.dev[d]["path_s"])
        # if not os.path.exists(default.file_bc):
        #     shutil.copy(default.dev[d]["file_bc"], default.file_bc)
        # if not os.path.exists(default.file_taint):
        #     shutil.copy(default.dev[d]["file_taint"], default.file_taint)
        # if not os.path.exists(default.name_with_dra):
        #     os.makedirs(default.name_with_dra)
        # if not os.path.exists(default.name_without_dra):
        #     os.makedirs(default.name_without_dra)
        # if not os.path.exists(default.name_run):
        #     shutil.copy(default.path_default_run, default.name_run)
        # if not os.path.exists(default.name_run_bash):
        #     shutil.copy(default.path_default_run_bash, default.name_run_bash)

        df = open(default.file_default_json, "r")
        c = json.load(df)
        df.close()
        c["enable_syscalls"] = default.dev[d]["enable_syscalls"]
        c["syzkaller"] = default.path_syzkaller
        c["kernel_obj"] = default.path_linux
        c["vm"]["kernel"] = default.path_kernel
        f = open("built-in" + ".json", "w")
        json.dump(c, f, indent=4)
        f.close()

        read_s(default.dev[d]["path_s"])
        shutil.copy(default.dev[d]["file_bc"], default.file_bc)
        shutil.copy(default.dev[d]["file_taint"], default.file_taint)
        shutil.copy(default.path_default_run, default.name_run)
        shutil.copy(default.path_default, default.name_default)
        shutil.copy(default.path_default_run_bash, default.name_run_bash)
        shutil.copy(default.path_default_remove_bash, default.name_remove_bash)
