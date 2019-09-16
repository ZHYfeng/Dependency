package main

import (
	"bytes"
	"fmt"
	pb "github.com/google/syzkaller/pkg/dra"
	"testing"
)

func TestRemoveSameResource(t *testing.T) {
	t.Parallel()
	tests := [][]byte{
		[]byte(
			"r0 = openat$kvm(0xffffffffffffff9c, &(0x7f0000000140)='/dev/kvm\x00', 0x0, 0x0)\n" +
				"r1 = openat$kvm(0xffffffffffffff9c, &(0x7f0000000140)='/dev/kvm\x00', 0x0, 0x0)\n" +
				"r2 = ioctl$KVM_CREATE_VM(r1, 0xae01, 0x0)\n" +
				"r3 = ioctl$KVM_CREATE_VM(r1, 0xae01, 0x0)\n" +
				"ioctl$KVM_SET_TSS_ADDR(r2, 0xae47, 0x0)\n" +
				"ioctl$KVM_SET_TSS_ADDR(r3, 0xae47, 0x0)\n" +
				"r4 = ioctl$KVM_CREATE_VM(r0, 0xae01, 0x0)\n" +
				"ioctl$KVM_SET_IDENTITY_MAP_ADDR(r4, 0x4008ae48, &(0x7f0000000100)=0x2000)\n" +
				"ioctl$KVM_CREATE_VCPU(r4, 0xae41, 0x0)\n"),
		[]byte(
			"r0 = openat$ptmx(0xffffffffffffff9c, &(0x7f0000000000)='/dev/ptmx\x00', 0x802, 0x0)\n" +
				"ioctl$TCSETS2(r0, 0x402c542b, &(0x7f0000000040)={0x0, 0x0, 0x0, 0x0, 0x0, \"00a38726bcf0c5eeaa0df10eadbc6230241cdf\"})\n" +
				"r1 = openat$ptmx(0xffffffffffffff9c, &(0x7f0000000000)='/dev/ptmx\x00', 0x802, 0x0)\n" +
				"ioctl$TCSETS2(r1, 0x402c542b, &(0x7f0000000040)={0x0, 0x0, 0x0, 0x0, 0x0, \"00a38726bcf0c5eeaa0df10eadbc6230241cdf\"})\n"),
		[]byte(
			"r0 = openat$kvm(0xffffffffffffff9c, &(0x7f0000000400)='/dev/kvm\x00', 0x0, 0x0)\n" +
				"r1 = ioctl$KVM_CREATE_VM(r0, 0xae01, 0x0)\n" +
				"r2 = openat$kvm(0xffffffffffffff9c, &(0x7f0000000040)='/dev/kvm\x00', 0x0, 0x0)\n" +
				"r3 = ioctl$KVM_CREATE_VM(r2, 0xae01, 0x0)\n" +
				"ioctl$KVM_CREATE_IRQCHIP(r3, 0xae60)\n" +
				"ioctl$KVM_CREATE_VCPU(r3, 0xae41, 0x0)\n" +
				"ioctl$KVM_CREATE_PIT2(r1, 0x4040ae77, &(0x7f0000000100))\n" +
				"ioctl$KVM_SET_PIT2(r1, 0xae71, &(0x7f0000000080))\n" +
				"ioctl$KVM_SET_USER_MEMORY_REGION(r1, 0x4020ae46, &(0x7f0000000000)={0x0, 0x0, 0x0, 0x2000, &(0x7f0000ffb000/0x2000)=nil})\n"),
		[]byte(
			"r0 = openat$ptmx(0xffffffffffffff9c, &(0x7f0000000000)='/dev/ptmx\x00', 0x802, 0x0)\n" +
				"ioctl$TIOCSETD(r0, 0x5423, &(0x7f0000000040))\n" +
				"r1 = openat$ptmx(0xffffffffffffff9c, &(0x7f0000000000)='/dev/ptmx\x00', 0x802, 0x0)\n" +
				"ioctl$TCFLSH(r1, 0x40045436, 0x2)\n"),
	}
	results := [][]byte{
		[]byte(
			"r0 = openat$kvm(0xffffffffffffff9c, &(0x7f0000000140)='/dev/kvm\x00', 0x0, 0x0)\n" +
				"r2 = ioctl$KVM_CREATE_VM(r0, 0xae01, 0x0)\n" +
				"ioctl$KVM_SET_TSS_ADDR(r2, 0xae47, 0x0)\n" +
				"ioctl$KVM_SET_TSS_ADDR(r2, 0xae47, 0x0)\n" +
				"ioctl$KVM_SET_IDENTITY_MAP_ADDR(r2, 0x4008ae48, &(0x7f0000000100)=0x2000)\n" +
				"ioctl$KVM_CREATE_VCPU(r2, 0xae41, 0x0)\n"),
		[]byte(
			"r0 = openat$ptmx(0xffffffffffffff9c, &(0x7f0000000000)='/dev/ptmx\x00', 0x802, 0x0)\n" +
				"ioctl$TCSETS2(r0, 0x402c542b, &(0x7f0000000040)={0x0, 0x0, 0x0, 0x0, 0x0, \"00a38726bcf0c5eeaa0df10eadbc6230241cdf\"})\n" +
				"ioctl$TCSETS2(r0, 0x402c542b, &(0x7f0000000040)={0x0, 0x0, 0x0, 0x0, 0x0, \"00a38726bcf0c5eeaa0df10eadbc6230241cdf\"})\n"),
		[]byte(
			"r0 = openat$kvm(0xffffffffffffff9c, &(0x7f0000000400)='/dev/kvm\x00', 0x0, 0x0)\n" +
				"r1 = ioctl$KVM_CREATE_VM(r0, 0xae01, 0x0)\n" +
				"ioctl$KVM_CREATE_IRQCHIP(r1, 0xae60)\n" +
				"ioctl$KVM_CREATE_VCPU(r1, 0xae41, 0x0)\n" +
				"ioctl$KVM_CREATE_PIT2(r1, 0x4040ae77, &(0x7f0000000100))\n" +
				"ioctl$KVM_SET_PIT2(r1, 0xae71, &(0x7f0000000080))\n" +
				"ioctl$KVM_SET_USER_MEMORY_REGION(r1, 0x4020ae46, &(0x7f0000000000)={0x0, 0x0, 0x0, 0x2000, &(0x7f0000ffb000/0x2000)=nil})\n"),
		[]byte(
			"r0 = openat$ptmx(0xffffffffffffff9c, &(0x7f0000000000)='/dev/ptmx\x00', 0x802, 0x0)\n" +
				"ioctl$TIOCSETD(r0, 0x5423, &(0x7f0000000040))\n" +
				"ioctl$TCFLSH(r0, 0x40045436, 0x2)\n"),
	}
	idxs := [][]int{
		{0, 0, 1, 1, 2, 3, 1, 4, 5},
		{0, 1, 0, 2},
		{0, 1, 0, 1, 2, 3, 4, 5, 6},
		{0, 1, 0, 2},
	}

	for ti, test := range tests {
		test := test
		t.Run(fmt.Sprint(ti), func(t *testing.T) {
			result, idx := removeSameResource(test)
			if !bytes.Equal(result, results[ti]) {
				t.Fatalf("result: \n%s", result)
			}
			for i, ii := range idxs[ti] {
				if ii != idx[i] {
					t.Fatalf("idx: %v", idx)
				}
			}
		})
	}
}
func TestCheckPath(t *testing.T) {
	t.Parallel()
	tests := [][][]uint32{
		{{0x44f9dade, 0xe7a93944, 0x6e1608f, 0x2099698c, 0xe432dbda, 0xb7491474, 0x282f1b47, 0x144f7b47, 0x81006cf3},
			{0x81006cf3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
	}
	results := []int{0}

	for ti, test := range tests {
		t.Run(fmt.Sprint(ti), func(t *testing.T) {
			_, _, idx := pb.CheckPath(test[0], test[1])
			if idx != results[ti] {
				t.Fatalf("idx: %v", idx)
			}
		})
	}
}
