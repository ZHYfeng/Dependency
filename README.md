# 1. what the artifact does

This artifact is for paper "Demystifying the Dependency Challenge in Kernel Fuzzing". 
Fuzz testing operating system kernels remains a daunting task to date. 
One known challenge is that much of the kernel code is locked under specific kernel states and current kernel fuzzers are not effective in exploring such an enormous state space. 
We refer to this problem as the dependency challenge. 
Though there are some efforts trying to address the dependency challenge, the prevalence and categorization of dependencies have never been studied. 
Most prior work simply attempted to recover dependencies opportunistically whenever they are relatively easy to recognize. 
We undertake a substantial measurement study to systematically understand the real challenge behind dependencies.
In one word, the artifact is to help researchers to understand the dependency challenge in kernel fuzzing.

# 2. where it can be obtained
We have made everything ready in the virtual machine (username & password: icse22ae):  
Ubuntu 20.04 64 ICSE22AE 1.ova in `https://drive.google.com/drive/folders/1Ts4P4iC2PHihtBviSXMUkn3My0PLkowN?usp=sharing`

You can also get the source code: `https://github.com/ZHYfeng/Dependency`

# 3. how to repeat/replicate/reproduce the results presented in the paper

## build our tools (skip this step if using virtual machine)
```
sudo apt install -y git
git clone https://github.com/ZHYfeng/Dependency.git
cd Dependency
bash build_script/build.bash
```

## prepare kernel and image (skip this step if using virtual machine)
1. configure the kernel and image based on the requirement of syzkaller, mv image to `path-of-Dependency/workdir/image`
    > doc of syzkaller： https://github.com/google/syzkaller/blob/master/docs/linux/setup_ubuntu-host_qemu-vm_x86-64-kernel.md  
    > the image we build: image.tar.gz in `https://drive.google.com/drive/folders/1Ts4P4iC2PHihtBviSXMUkn3My0PLkowN?usp=sharing`
2. add `-fsanitize-coverage=no-prune` to `CFLAGS_KCOV` in kernel config
3. build kernel using clang and mv it to `path-of-Dependency/workdir/13-linux-clang-np`
   > the kernel we build: linux-clang-np.tar.gz in `https://drive.google.com/drive/folders/1Ts4P4iC2PHihtBviSXMUkn3My0PLkowN?usp=sharing`
4. copy the kernel and generate bitcode of kernel using `-fembed-bitcode -save-temps=obj`
    > https://github.com/ZHYfeng/Build_Linux_Kernel_Into_LLVM_Bitcode/tree/master/01-change_makefile  
    > the bitcode we build:  linux-clang-np-bc-f.tar.gz in `https://drive.google.com/drive/folders/1Ts4P4iC2PHihtBviSXMUkn3My0PLkowN?usp=sharing`
5. preprocess kernel in order to save time
   ```
   cd path-of-Dependency/workdir/13-linux-clang-np
   objdump -d vmlinux > vmlinux.objdump
   a2l -objdump=vmlinux.objdump
   ```

## prepare workdir (skip this step if using virtual machine)
> the workdir we prepare: workdir.tar.gz in `https://drive.google.com/drive/folders/1Ts4P4iC2PHihtBviSXMUkn3My0PLkowN?usp=sharing`
1. make a directory called `dev_xxx` in `path-of-Dependency/workdir`
2. copy the bitcode(.bc) and assembly code(.s) to the directory and rename it to `built-in.bc` and `built-in.s`
3. copy the configuration files `path-of-Dependency/04-experiment_script/json/dra.json` and `path-of-Dependency/04-experiment_script/json/syzkaller.json`.
   > change the value of `file_bc` in `dra.json` to the relative path for the bitcode of device driver you test  
   > change the value of `path_s` in `dra.json` to the relative path of device driver you test  
4. copy the run script `path-of-Dependency/04-experiment_script/python/run.py`
5. generate static analysis results based on the static-taint-analysis-component `https://zenodo.org/record/5348989/files/static-taint-analysis-component.zip`

## running the fuzzing
(the path based on virtual machine)
1. active the environment
    ```
    source /home/icse22ae/Dependency/environment.sh
    ```
2. pick one device driver in `/home/icse22ae/Dependency/workdir/workdir`, for example`cdrom`:
    ```
    cd /home/icse22ae/Dependency/workdir/workdir/dev_cdrom
    ```
3. configure the run script
    > time_run: the second of fuzzing time.  
    > number_execute: the number of fuzzing runs.  
    > number_vm_count: the number of vm in each fuzzing.  

    In our paper, `time_run` is at least 48 hours, `number_execute` is 3 and `number_vm_count` is 32.  
    For artifact evaluation, `number_execute` and `number_vm_count` could be 1.  
    `time_run` should be at least 5 mins(20 mins for device driver kvm)
4. run our tool using script
    It will automatically stop after `time_run`.
    ```
    python3 run.py
    ```
5. read the results  
    still in the same environment in step 1 and the same path in step 2.
    ```
    go run /home/icse22ae/Dependency/03-syzkaller/tools/read_result/ -a2i
    ```
    Based on the different fuzzing configuration and device driver, the time would be differnet.  
    For cdrom, it should be several mins. For kvm, it needs several hours.

## understand the results
You can find the results used in our paper in `/home/icse22ae/Dependency/workdir/data`.  
### results after step 4 run our tool using script
1. The `dataDependency.bin`, `dataResult.bin`, `dataRunTime.bin`, `statistics.bin` in `./0` or `./1` or `./2` are the resutls in protobuf format.
    > The protobuf files are in `/home/icse22ae/Dependency/05-proto`
### results after step 5 read the results
2. `0_coverage.txt` is the coverage of the fuzzing in `./0`. `coverage.txt` is the average coverage of all runs.Each line is `time@number-of-edge`.
3. `conditionD.txt` lists all unresolved condition related to dependency.
4. `conditionND.txt` lists all unresolved condition not related to dependency.
5. `conditionDN.txt` lists all unresolved condition related to dependency but our static analysis can not find their write statements.
6. `intersection.txt` is the intersection coverage of all runs and `union_coverage.txt` is the union coverage of all runs. Each line is the address of the edge.
7. `OutsideFunctions.txt` is the `Unreachable Functions Elimination` mentioned in our paper.
8. `statistic.txt` is the statistic used in our paper.
9. `uncovered.txt` lists all uncovered edge and its unresovled conditions, and `uncovered_more.txt` lists more details about them.