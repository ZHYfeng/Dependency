// Copyright 2018 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/ZHYfeng/Dependency/03-syzkaller/pkg/osutil"
)

type akaros struct{}

func (ctx akaros) build(targetArch, vmType, kernelDir, outputDir, compiler, userspaceDir,
	cmdlineFile, sysctlFile string, config []byte) error {
	configFile := filepath.Join(kernelDir, ".config")
	if err := osutil.WriteFile(configFile, config); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	if err := osutil.SandboxChown(configFile); err != nil {
		return err
	}
	sshkey := filepath.Join(kernelDir, "key")
	sshkeyPub := sshkey + ".pub"
	os.Remove(sshkey)
	os.Remove(sshkeyPub)
	if _, err := osutil.RunCmd(10*time.Minute, "", "ssh-keygen", "-t", "rsa", "-b", "2048",
		"-N", "", "-C", "", "-f", sshkey); err != nil {
		return err
	}
	if err := osutil.SandboxChown(sshkeyPub); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, "", "olddefconfig", "ARCH=x86"); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, "", "xcc"); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, "tools/dev-libs/elfutils", "install"); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, "", "apps-install"); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, "", "fill-kfs"); err != nil {
		return err
	}
	targetKey := filepath.Join(kernelDir, "kern", "kfs", ".ssh", "authorized_keys")
	if err := osutil.Rename(sshkeyPub, targetKey); err != nil {
		return err
	}
	const init = `#!/bin/bash
/ifconfig
dropbear -F 2>db_out &
bash
`
	initFile := filepath.Join(kernelDir, "kern", "kfs", "init.sh")
	if err := osutil.WriteFile(initFile, []byte(init)); err != nil {
		return fmt.Errorf("failed to write init script: %v", err)
	}
	if err := osutil.SandboxChown(initFile); err != nil {
		return err
	}
	if err := os.Chmod(initFile, 0770); err != nil {
		return err
	}
	if err := ctx.cmd(kernelDir, "dropbear", "./CONFIGURE_AKAROS"); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, "dropbear/build"); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, "dropbear/build", "install"); err != nil {
		return err
	}
	if err := ctx.make(kernelDir, ""); err != nil {
		return err
	}
	if err := osutil.WriteFile(filepath.Join(outputDir, "image"), nil); err != nil {
		return fmt.Errorf("failed to write image file: %v", err)
	}
	for src, dst := range map[string]string{
		".config":                    "kernel.config",
		"key":                        "key",
		"obj/kern/akaros-kernel":     "kernel",
		"obj/kern/akaros-kernel-64b": "obj/akaros-kernel-64b",
	} {
		fullSrc := filepath.Join(kernelDir, filepath.FromSlash(src))
		fullDst := filepath.Join(outputDir, filepath.FromSlash(dst))
		if err := osutil.CopyFile(fullSrc, fullDst); err != nil {
			return fmt.Errorf("failed to copy %v: %v", src, err)
		}
	}
	return nil
}

func (ctx akaros) clean(kernelDir, targetArch string) error {
	// Note: this does not clean toolchain and elfutils.
	return ctx.make(kernelDir, "", "realclean")
}

func (ctx akaros) make(kernelDir, runDir string, args ...string) error {
	args = append([]string{"-j", strconv.Itoa(runtime.NumCPU())}, args...)
	return ctx.cmd(kernelDir, runDir, "make", args...)
}

func (ctx akaros) cmd(kernelDir, runDir string, bin string, args ...string) error {
	cmd := osutil.Command(bin, args...)
	if err := osutil.Sandbox(cmd, true, false); err != nil {
		return err
	}
	cmd.Dir = kernelDir
	if runDir != "" {
		cmd.Dir = filepath.Join(kernelDir, filepath.FromSlash(runDir))
	}
	cmd.Env = append([]string{}, os.Environ()...)
	cmd.Env = append(cmd.Env, []string{
		"MAKE_JOBS=" + strconv.Itoa(runtime.NumCPU()),
		"AKAROS_ROOT=" + kernelDir,
		"AKAROS_XCC_ROOT=" + filepath.Join(kernelDir, "toolchain", "x86_64-ucb-akaros-gcc"),
		"X86_64_INSTDIR=" + filepath.Join(kernelDir, "toolchain", "x86_64-ucb-akaros-gcc"),
		"PATH=" + filepath.Join(kernelDir, "toolchain", "x86_64-ucb-akaros-gcc", "bin") +
			string(filepath.ListSeparator) + os.Getenv("PATH"),
	}...)
	_, err := osutil.Run(time.Hour, cmd)
	return err
}
