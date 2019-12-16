#! /usr/bin/python3
import json
import os
import shutil

from default import default


def read_s(path, files):
    ctx_s = ""
    for file in files:
        for (dir_path, dir_names, file_names) in os.walk(os.path.join(path, file)):
            for file_name in file_names:
                if file_name.endswith(".s") and file_name != default.file_asm:
                    f = open(os.path.join(dir_path, file_name), "r")
                    ctx_s = ctx_s + "\n" + f.read()
                    f.close()
    f = open(default.file_asm, "w")
    f.write(ctx_s)
    f.close()
    # cat `find -name "*.s"` >> built-in.s


def generate_dev_dir():
    file = open(default.file_default_json, "r")
    default_json = json.load(file)
    for device in default_json:
        path = os.path.join(default.path_result, device)
        print(path)

        # if os.path.exists(path):
        #     shutil.rmtree(path)

        if not os.path.exists(path):
            os.makedirs(path)
        os.chdir(path)

        # if not os.path.exists("built-in.json"):
        #     df = open(default.file_syzkaller_json, "r")
        #     c = json.load(df)
        #     df.close()
        #     c["enable_syscalls"] = default_json[device]["enable_syscalls"]
        #     c["syzkaller"] = default.path_syzkaller
        #     c["kernel_obj"] = default.path_linux
        #     c["vm"]["kernel"] = default.path_kernel
        #     f = open("built-in" + ".json", "w")
        #     json.dump(c, f, indent=4)
        #     f.close()
        #
        # if not os.path.exists(default.file_asm):
        #     read_s(default.path_linux_bc, default_json[device]["path_s"])
        # if not os.path.exists(default.file_bc):
        #     shutil.copy(os.path.join(default.path_linux_bc, default_json[device]["file_bc"]), default.file_bc)
        # if not os.path.exists(default.file_taint):
        #     shutil.copy(os.path.join(default.path_taint, default_json[device]["file_taint"]), default.file_taint)
        # if not os.path.exists(default.name_with_dra):
        #     os.makedirs(default.name_with_dra)
        # if not os.path.exists(default.name_without_dra):
        #     os.makedirs(default.name_without_dra)
        # if not os.path.exists(default.name_run):
        #     shutil.copy(default.path_default_run, default.name_run)
        # if not os.path.exists(default.name_run_bash):
        #     shutil.copy(default.path_default_run_bash, default.name_run_bash)

        df = open(default.file_syzkaller_json, "r")
        c = json.load(df)
        df.close()
        c["enable_syscalls"] = default_json[device]["enable_syscalls"]
        c["syzkaller"] = default.path_syzkaller
        c["kernel_obj"] = default.path_linux
        c["vm"]["kernel"] = default.path_kernel
        f = open("built-in" + ".json", "w")
        json.dump(c, f, indent=4)
        f.close()

        read_s(default.path_linux_bc, default_json[device]["path_s"])
        shutil.copy(os.path.join(default.path_linux_bc, default_json[device]["file_bc"]), default.file_bc)
        shutil.copy(os.path.join(default.path_taint, default_json[device]["file_taint"]), default.file_taint)
        shutil.copy(default.path_default_run, default.name_run)
        shutil.copy(default.path_default, default.name_default)
        shutil.copy(default.path_default_run_bash, default.name_run_bash)
        shutil.copy(default.path_default_remove_bash, default.name_remove_bash)

        ff = open(default.file_function, "w")
        if "function" in default_json[device]:
            json.dump(default_json[device]["function"], ff, indent=4, sort_keys=True)
        ff.close()
