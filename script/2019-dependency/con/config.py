import os
import shutil
import json

from con import dev


def read_s(paths):
    ctx_s = ""
    for path in paths:
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.endswith(".s"):
                    f = open(os.path.join(dir_path, file_name), "r")
                    ctx_s = ctx_s + "\n" + f.read()
    f = open("built-in" + ".s", "w")
    f.write(ctx_s)
    f.close()


def generate_dev_dir():
    for d in dev.dev:
        path = dev.path_result + "/" + d
        print(path)
        if os.path.exists(path):
            shutil.rmtree(path)
        if not os.path.exists(path):
            os.makedirs(path)
            os.chdir(path)

            if not os.path.exists(d + ".json"):
                df = open(dev.path_default_json, "r")
                c = json.load(df)
                df.close()
                c["enable_syscalls"] = dev.dev[d]["enable_syscalls"]
                f = open("built-in" + ".json", "w")
                json.dump(c, f, indent=4)
                f.close()

            if not os.path.exists("./built-in.s"):
                read_s(dev.dev[d]["path_s"])
            if not os.path.exists("./built-in.bc"):
                shutil.copy(dev.dev[d]["file_bc"], "./built-in.bc")
            if not os.path.exists("./built-in.taint"):
                shutil.copy(dev.dev[d]["file_taint"], "./built-in.taint")


if __name__ == "__main__":
    generate_dev_dir()
