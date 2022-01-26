// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

//go:generate bash -c "echo -en '// AUTOGENERATED FILE\n\n' > linux_generated.go"
//go:generate bash -c "echo -en 'package build\n\n' >> linux_generated.go"
//go:generate bash -c "echo -en 'const createImageScript = `#!/bin/bash\n' >> linux_generated.go"
//go:generate bash -c "cat ../../tools/create-gce-image.sh | grep -v '#' >> linux_generated.go"
//go:generate bash -c "echo -en '`\n\n' >> linux_generated.go"

package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/ZHYfeng/Dependency/03-syzkaller/pkg/osutil"
)

type linux struct{}

func (linux linux) build(targetArch, vmType, kernelDir, outputDir, compiler, userspaceDir,
	cmdlineFile, sysctlFile string, config []byte) error {
	if err := linux.buildKernel(targetArch, kernelDir, outputDir, compiler, config); err != nil {
		return err
	}
	if err := linux.createImage(targetArch, vmType, kernelDir, outputDir, userspaceDir, cmdlineFile,
		sysctlFile); err != nil {
		return err
	}
	return nil
}

func (linux) buildKernel(targetArch, kernelDir, outputDir, compiler string, config []byte) error {
	configFile := filepath.Join(kernelDir, ".config")
	if err := osutil.WriteFile(configFile, config); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	if err := osutil.SandboxChown(configFile); err != nil {
		return err
	}
	// One would expect olddefconfig here, but olddefconfig is not present in v3.6 and below.
	// oldconfig is the same as olddefconfig if stdin is not set.
	// Note: passing in compiler is important since 4.17 (at the very least it's noted in the config).
	cmd := osutil.Command("make", "oldconfig", "CC="+compiler)
	if err := osutil.Sandbox(cmd, true, true); err != nil {
		return err
	}
	cmd.Dir = kernelDir
	if _, err := osutil.Run(10*time.Minute, cmd); err != nil {
		return err
	}
	// Write updated kernel config early, so that it's captured on build failures.
	outputConfig := filepath.Join(outputDir, "kernel.config")
	if err := osutil.CopyFile(configFile, outputConfig); err != nil {
		return err
	}
	// We build only zImage/bzImage as we currently don't use modules.
	cpu := strconv.Itoa(runtime.NumCPU())

	var target string
	switch targetArch {
	case "386", "amd64":
		target = "bzImage"
	case "ppc64le":
		target = "zImage"
	}
	cmd = osutil.Command("make", target, "-j", cpu, "CC="+compiler)

	if err := osutil.Sandbox(cmd, true, true); err != nil {
		return err
	}
	cmd.Dir = kernelDir
	if _, err := osutil.Run(time.Hour, cmd); err != nil {
		return err
	}
	vmlinux := filepath.Join(kernelDir, "vmlinux")
	outputVmlinux := filepath.Join(outputDir, "obj", "vmlinux")
	if err := osutil.Rename(vmlinux, outputVmlinux); err != nil {
		return fmt.Errorf("failed to rename vmlinux: %v", err)
	}
	return nil
}

func (linux) createImage(targetArch, vmType, kernelDir, outputDir, userspaceDir, cmdlineFile, sysctlFile string) error {
	tempDir, err := ioutil.TempDir("", "syz-build")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	scriptFile := filepath.Join(tempDir, "create.sh")
	if err := osutil.WriteExecFile(scriptFile, []byte(createImageScript)); err != nil {
		return fmt.Errorf("failed to write script file: %v", err)
	}

	var kernelImage string
	switch targetArch {
	case "386", "amd64":
		kernelImage = "arch/x86/boot/bzImage"
	case "ppc64le":
		kernelImage = "arch/powerpc/boot/zImage.pseries"
	}
	kernelImagePath := filepath.Join(kernelDir, filepath.FromSlash(kernelImage))
	cmd := osutil.Command(scriptFile, userspaceDir, kernelImagePath, targetArch)
	cmd.Dir = tempDir
	cmd.Env = append([]string{}, os.Environ()...)
	cmd.Env = append(cmd.Env,
		"SYZ_VM_TYPE="+vmType,
		"SYZ_CMDLINE_FILE="+osutil.Abs(cmdlineFile),
		"SYZ_SYSCTL_FILE="+osutil.Abs(sysctlFile),
	)
	if _, err = osutil.Run(time.Hour, cmd); err != nil {
		return fmt.Errorf("image build failed: %v", err)
	}
	// Note: we use CopyFile instead of Rename because src and dst can be on different filesystems.
	imageFile := filepath.Join(outputDir, "image")
	if err := osutil.CopyFile(filepath.Join(tempDir, "disk.raw"), imageFile); err != nil {
		return err
	}
	keyFile := filepath.Join(outputDir, "key")
	if err := osutil.CopyFile(filepath.Join(tempDir, "key"), keyFile); err != nil {
		return err
	}
	if err := os.Chmod(keyFile, 0600); err != nil {
		return err
	}
	return nil
}

func (linux) clean(kernelDir, targetArch string) error {
	cpu := strconv.Itoa(runtime.NumCPU())
	cmd := osutil.Command("make", "distclean", "-j", cpu)
	if err := osutil.Sandbox(cmd, true, true); err != nil {
		return err
	}
	cmd.Dir = kernelDir
	_, err := osutil.Run(10*time.Minute, cmd)
	return err
}
