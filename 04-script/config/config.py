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
                if file_name.endswith(".s") and file_name != default.name_asm:
                    f = open(os.path.join(dir_path, file_name), "r")
                    ctx_s = ctx_s + "\n" + f.read()
                    f.close()
    f = open(default.name_asm, "w")
    f.write(ctx_s)
    f.close()
    # cat `find -name "*.s"` >> built-in.s


def read_file_syzkaller_json():
    df = open(default.file_syzkaller_json, "r")
    c = json.load(df)
    df.close()
    c["syzkaller"] = default.path_syzkaller
    c["kernel_obj"] = default.path_linux
    c["vm"]["kernel"] = default.path_kernel
    return c


def copy_files():
    shutil.copy(default.path_default_run, default.name_run)
    shutil.copy(default.path_default, default.name_default)
    shutil.copy(default.path_default_run_bash, default.name_run_bash)
    shutil.copy(default.path_default_remove_bash, default.name_remove_bash)


def generate_dev_dir():
    file = open(default.file_default_json, "r")
    default_json = json.load(file)
    for device in default_json:
        path = os.path.join(default.path_result, device)
        print(path)

        if os.path.exists(path):
            shutil.rmtree(path)

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

        os.makedirs(default.name_with_dra)
        os.makedirs(default.name_without_dra)

        syzkaller_json = read_file_syzkaller_json()
        syzkaller_json["enable_syscalls"] = default_json[device]["enable_syscalls"]
        f = open(default.name_syzkaller_json, "w")
        json.dump(syzkaller_json, f, indent=4)
        f.close()

        dra_json = {}
        dra_json[device] = default_json[device]
        f = open(default.name_dra_json, "w")
        json.dump(dra_json, f, indent=4)

        read_s(default.path_linux_bc, default_json[device]["path_s"])
        shutil.copy(os.path.join(default.path_linux_bc, default_json[device]["file_bc"]), default.name_bc)
        shutil.copy(os.path.join(default.path_taint, default_json[device]["file_taint"]),
                    default_json[device]["file_taint"])
        copy_files()

        # ff = open(default.name_dra_json, "w")
        # if "function" in default_json[device]:
        #     json.dump(default_json[device]["function"], ff, indent=4, sort_keys=True)
        # ff.close()

    overall = ["dev_cdrom", "dev_kvm", "dev_ptmx", "dev_snd_seq", ]

    path = os.path.join(default.path_result, "overall")
    print(path)

    if os.path.exists(path):
        shutil.rmtree(path)

    if not os.path.exists(path):
        os.makedirs(path)
    os.chdir(path)

    os.makedirs(default.name_with_dra)
    os.makedirs(default.name_without_dra)

    syzkaller_json = read_file_syzkaller_json()
    for d in overall:
        syzkaller_json["enable_syscalls"] += default_json[d]["enable_syscalls"]
    f = open(default.name_syzkaller_json, "w")
    json.dump(syzkaller_json, f, indent=4)
    f.close()

    dra_json = {}
    for d in overall:
        dra_json[d] = default_json[d]
    f = open(default.name_dra_json, "w")
    json.dump(dra_json, f, indent=4)

    s_files = []
    for d in overall:
        shutil.copy(os.path.join(default.path_taint, default_json[d]["file_taint"]),
                    default_json[d]["file_taint"])
        for s in default_json[d]["path_s"]:
            s_files.append(s)
    read_s(default.path_linux_bc, s_files)
    shutil.copy(os.path.join(default.path_linux_bc, default.name_bc), default.name_bc)

    copy_files()
