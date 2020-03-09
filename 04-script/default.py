#! /usr/bin/python3
import os

encoding = 'utf-8'

path_home = os.path.expanduser("~")
path_root = os.path.join(path_home, "data")
path_git = os.path.join(path_root, "git")
path_git_repo = os.path.join(path_git, "gopath/src/github.com/ZHYfeng/2018_dependency")

path_result = os.path.join(path_git_repo, "06-result")
path_git_script = os.path.join(path_git_repo, "04-script")
path_taint = os.path.join(path_git_script, "taint_info")
file_syzkaller_json = os.path.join(path_git_script, "config/syzkaller.json")
file_default_json = os.path.join(path_git_script, "config/default.json")
name_default = "default.py"
path_default = os.path.join(path_git_script, "default", name_default)
name_run = "run.py"
path_default_run = os.path.join(path_git_script, "default", name_run)
name_run_bash = "run.bash"
path_default_run_bash = os.path.join(path_git_script, "default", name_run_bash)
name_remove_bash = "remove.bash"
path_default_remove_bash = os.path.join(path_git_script, "default", name_remove_bash)

name_with_dra = "01-result-with-dra"
name_without_dra = "02-result-without-dra"
path_linux_bc = os.path.join(
    path_root, "benchmark/linux/16-linux-clang-np-bc-f")
path_linux = os.path.join(path_root, "benchmark/linux/13-linux-clang-np")
path_kernel = os.path.join(path_linux, "arch/x86/boot/bzImage")

path_syzkaller = os.path.join(path_git_repo, "03-syzkaller")
file_syzkaller = os.path.join(path_syzkaller, "bin/syz-manager")

name_driver = "built-in"
name_asm = name_driver + ".s"
name_bc = name_driver + ".bc"
name_syzkaller_json = "syzkaller.json"
name_dra_json = "dra.json"

number_execute = 1
path_current = os.getcwd()
path_git = os.path.join(path_root, "git")
path_repo = os.path.join(path_git, "gopath/src/github.com/ZHYfeng/2018_dependency")
path_dra = os.path.join(path_repo, "02-dependency/cmake-build-debug/tools/DRA/dra")
path_syzkaller_manager = os.path.join(path_repo, "03-syzkaller/bin/syz-manager")
file_vmlinux_objdump = os.path.join(path_linux, "vmlinux.objdump")

name_workdir = "workdir"
path_workdir = os.path.join(path_current, name_workdir)

path_image = os.path.join(path_root, "benchmark/linux/image")
file_image = "stretch.img"
file_ssh_key = "stretch.id_rsa"

file_log_run = "log_run.bash"
file_log_syzkaller = "log_syzkaller.log"
file_log_dra = "log_dra.log"
file_run = "run.bash"

length = 1 * 24 * 60
time_run = length * 60  # second

do_figure = True
confidence = 0.95
name_dev = "dev_"
name_stat = "statistics.bin"
name_data = "data.bin"
name_data_result = "data.txt"
name_base = "base"
path_home = os.path.expanduser("~")
path_root = os.path.join(path_home, "data")
path_git = os.path.join(path_root, "git")
path_repo = os.path.join(path_git, "gopath/src/github.com/ZHYfeng/2018_dependency")
path_a2i = os.path.join(path_repo, "02-dependency/cmake-build-debug/tools/A2I/a2i")
