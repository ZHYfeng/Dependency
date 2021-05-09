// Copyright 2018 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package linux_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ZHYfeng/2018-Dependency/03-syzkaller/prog"
	_ "github.com/ZHYfeng/2018-Dependency/03-syzkaller/sys/linux/gen"
)

func TestSanitize(t *testing.T) {
	target, err := prog.GetTarget("linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		input  string
		output string
	}{
		{
			`syslog(0x10000000006, 0x0, 0x0)`,
			`syslog(0x9, 0x0, 0x0)`,
		},
		{
			`syslog(0x10000000007, 0x0, 0x0)`,
			`syslog(0x9, 0x0, 0x0)`,
		},
		{
			`syslog(0x1, 0x0, 0x0)`,
			`syslog(0x1, 0x0, 0x0)`,
		},

		{
			`ptrace(0xf000000000, 0x0)`,
			`ptrace(0xffffffffffffffff, 0x0)`,
		},
		{
			`ptrace$peek(0x0, 0x0, &(0x7f0000000000))`,
			`ptrace$peek(0xffffffffffffffff, 0x0, &(0x7f0000000000))`,
		},
		{
			`ptrace(0x1, 0x0)`,
			`ptrace(0x1, 0x0)`,
		},
		{
			`arch_prctl$ARCH_SET_GS(0xf00000001002, 0x0)`,
			`arch_prctl$ARCH_SET_GS(0x1001, 0x0)`,
		},
		{
			`arch_prctl$ARCH_SET_GS(0x1003, 0x0)`,
			`arch_prctl$ARCH_SET_GS(0x1003, 0x0)`,
		},
		{
			`ioctl(0x0, 0x200000c0045877, 0x0)`,
			`ioctl(0x0, 0xc0045878, 0x0)`,
		},
		{
			`ioctl$int_in(0x0, 0x2000008004587d, 0x0)`,
			`ioctl$int_in(0x0, 0x6609, 0x0)`,
		},
		{
			`fanotify_mark(0x1, 0x2, 0x407fe029, 0x3, 0x0)`,
			`fanotify_mark(0x1, 0x2, 0x4078e029, 0x3, 0x0)`,
		},
		{
			`fanotify_mark(0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffff8ffff, 0xffffffffffffffff, 0x0)`,
			`fanotify_mark(0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffff8ffff, 0xffffffffffffffff, 0x0)`,
		},
		{
			`syz_init_net_socket$bt_hci(0x1, 0x0, 0x0)`,
			`syz_init_net_socket$bt_hci(0xffffffffffffffff, 0x0, 0x0)`,
		},
		{
			`syz_init_net_socket$bt_hci(0x27, 0x0, 0x0)`,
			`syz_init_net_socket$bt_hci(0x27, 0x0, 0x0)`,
		},
		{
			`syz_init_net_socket$bt_hci(0x1a, 0x0, 0x0)`,
			`syz_init_net_socket$bt_hci(0x1a, 0x0, 0x0)`,
		},
		{
			`syz_init_net_socket$bt_hci(0x1f, 0x0, 0x0)`,
			`syz_init_net_socket$bt_hci(0x1f, 0x0, 0x0)`,
		},
		{
			`mmap(0x0, 0x0, 0x0, 0x0, 0x0, 0x0)`,
			`mmap(0x0, 0x0, 0x0, 0x10, 0x0, 0x0)`,
		},
		{
			`mremap(0x0, 0x0, 0x0, 0xcc, 0x0)`,
			`mremap(0x0, 0x0, 0x0, 0xcc, 0x0)`,
		},
		{
			`mremap(0x0, 0x0, 0x0, 0xcd, 0x0)`,
			`mremap(0x0, 0x0, 0x0, 0xcf, 0x0)`,
		},
		{
			`
mknod(0x0, 0x1000, 0x0)
mknod(0x0, 0x8000, 0x0)
mknod(0x0, 0xc000, 0x0)
mknod(0x0, 0x2000, 0x0)
mknod(0x0, 0x6000, 0x0)
mknod(0x0, 0x6000, 0x700)
`,
			`
mknod(0x0, 0x1000, 0x0)
mknod(0x0, 0x8000, 0x0)
mknod(0x0, 0xc000, 0x0)
mknod(0x0, 0x8000, 0x0)
mknod(0x0, 0x8000, 0x0)
mknod(0x0, 0x6000, 0x700)
`,
		},
		{
			`
exit(0x3)
exit(0x43)
exit(0xc3)
exit(0xc3)
exit_group(0x5a)
exit_group(0x43)
exit_group(0x443)
`,
			`
exit(0x3)
exit(0x1)
exit(0x1)
exit(0x1)
exit_group(0x5a)
exit_group(0x1)
exit_group(0x1)
`,
		},
		{
			`
syz_open_procfs(0x0, &(0x7f0000000000)='io')
syz_open_procfs(0x0, &(0x7f0000000000)='exe')
syz_open_procfs(0x0, &(0x7f0000000000)='exe\x00')
syz_open_procfs(0x0, &(0x7f0000000000)='/exe')
syz_open_procfs(0x0, &(0x7f0000000000)='./exe\x00')
`,
			`
syz_open_procfs(0x0, &(0x7f0000000000)='io')
syz_open_procfs(0x0, &(0x7f0000000000)='net\x00')
syz_open_procfs(0x0, &(0x7f0000000000)='net\x00')
syz_open_procfs(0x0, &(0x7f0000000000)='net\x00')
syz_open_procfs(0x0, &(0x7f0000000000)='net\x00')
			`,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			p, err := target.Deserialize([]byte(test.input), prog.Strict)
			if err != nil {
				t.Fatal(err)
			}
			got := strings.TrimSpace(string(p.Serialize()))
			want := strings.TrimSpace(test.output)
			if got != want {
				t.Fatalf("input:\n%v\ngot:\n%v\nwant:\n%s", test.input, got, want)
			}
		})
	}
}
