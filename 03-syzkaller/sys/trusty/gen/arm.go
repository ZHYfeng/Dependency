// AUTOGENERATED FILE
// +build !syz_target syz_target,syz_os_trusty,syz_arch_arm

package gen

import . "github.com/ZHYfeng/2018-Dependency/03-syzkaller/prog"
import . "github.com/ZHYfeng/2018-Dependency/03-syzkaller/sys/trusty"

func init() {
	RegisterTarget(&Target{OS: "trusty", Arch: "arm", Revision: revision_arm, PtrSize: 4, PageSize: 4096, NumPages: 4096, DataOffset: 536870912, Syscalls: syscalls_arm, Resources: resources_arm, Structs: structDescs_arm, Consts: consts_arm}, InitTarget)
}

var resources_arm = []*ResourceDesc(nil)

var structDescs_arm = []*KeyedStruct{
	{Key: StructKey{Name: "dma_pmem"}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "dma_pmem", TypeSize: 4}, Fields: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "todo", TypeSize: 4}}},
	}}},
	{Key: StructKey{Name: "ipc_msg"}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "ipc_msg", TypeSize: 4}, Fields: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "todo", TypeSize: 4}}},
	}}},
	{Key: StructKey{Name: "ipc_msg", Dir: 1}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "ipc_msg", TypeSize: 4, ArgDir: 1}, Fields: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "todo", TypeSize: 4, ArgDir: 1}}},
	}}},
	{Key: StructKey{Name: "ipc_msg_info"}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "ipc_msg_info", TypeSize: 4}, Fields: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "todo", TypeSize: 4}}},
	}}},
	{Key: StructKey{Name: "uevent"}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "uevent", TypeSize: 4}, Fields: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "todo", TypeSize: 4}}},
	}}},
	{Key: StructKey{Name: "uevent", Dir: 1}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "uevent", TypeSize: 4, ArgDir: 1}, Fields: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "todo", TypeSize: 4, ArgDir: 1}}},
	}}},
	{Key: StructKey{Name: "uuid", Dir: 1}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "uuid", TypeSize: 4, ArgDir: 1}, Fields: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "todo", TypeSize: 4, ArgDir: 1}}},
	}}},
}

var syscalls_arm = []*Syscall{
	{NR: 18, Name: "accept", CallName: "accept", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle_id", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "peer_uuid", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "uuid", Dir: 1}}},
	}},
	{NR: 2, Name: "brk", CallName: "brk", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "brk", TypeSize: 4}}},
	}},
	{NR: 19, Name: "close", CallName: "close", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle_id", TypeSize: 4}}},
	}},
	{NR: 17, Name: "connect", CallName: "connect", Args: []Type{
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "path", TypeSize: 4}, Type: &BufferType{TypeCommon: TypeCommon{TypeName: "string", IsVarlen: true}, Kind: 2}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
	}},
	{NR: 3, Name: "exit_etc", CallName: "exit_etc", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "status", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
	}},
	{NR: 11, Name: "finish_dma", CallName: "finish_dma", Args: []Type{
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "uaddr", TypeSize: 4}, Type: &BufferType{TypeCommon: TypeCommon{TypeName: "array", ArgDir: 1, IsVarlen: true}}},
		&LenType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "len", FldName: "size", TypeSize: 4}}, Buf: "uaddr"},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
	}},
	{NR: 32, Name: "get_msg", CallName: "get_msg", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "msg_info", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "ipc_msg_info"}}},
	}},
	{NR: 7, Name: "gettime", CallName: "gettime", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "clock_id", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "time", TypeSize: 4}, Type: &IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int64", TypeSize: 8, ArgDir: 1}}}},
	}},
	{NR: 21, Name: "handle_set_create", CallName: "handle_set_create"},
	{NR: 22, Name: "handle_set_ctrl", CallName: "handle_set_ctrl", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "cmd", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "evt", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "uevent"}}},
	}},
	{NR: 5, Name: "ioctl", CallName: "ioctl", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "fd", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "req", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "buf", TypeSize: 4}, Type: &BufferType{TypeCommon: TypeCommon{TypeName: "array", IsVarlen: true}}},
	}},
	{NR: 8, Name: "mmap", CallName: "mmap", Args: []Type{
		&VmaType{TypeCommon: TypeCommon{TypeName: "vma", FldName: "uaddr", TypeSize: 4}},
		&LenType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "len", FldName: "size", TypeSize: 4}}, Buf: "uaddr"},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle", TypeSize: 4}}},
	}},
	{NR: 9, Name: "munmap", CallName: "munmap", Args: []Type{
		&VmaType{TypeCommon: TypeCommon{TypeName: "vma", FldName: "uaddr", TypeSize: 4}},
		&LenType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "len", FldName: "size", TypeSize: 4}}, Buf: "uaddr"},
	}},
	{NR: 6, Name: "nanosleep", CallName: "nanosleep", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "clock_id", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int64", FldName: "sleep_time", TypeSize: 8}}},
	}},
	{NR: 16, Name: "port_create", CallName: "port_create", Args: []Type{
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "path", TypeSize: 4}, Type: &BufferType{TypeCommon: TypeCommon{TypeName: "string", IsVarlen: true}, Kind: 2}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "num_recv_bufs", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "recv_buf_size", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
	}},
	{NR: 10, Name: "prepare_dma", CallName: "prepare_dma", Args: []Type{
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "uaddr", TypeSize: 4}, Type: &BufferType{TypeCommon: TypeCommon{TypeName: "array", ArgDir: 1, IsVarlen: true}}},
		&LenType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "len", FldName: "size", TypeSize: 4}}, Buf: "uaddr"},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "flags", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "pmem", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "dma_pmem"}}},
	}},
	{NR: 34, Name: "put_msg", CallName: "put_msg", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "msg_id", TypeSize: 4}}},
	}},
	{NR: 4, Name: "read", CallName: "read", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "fd", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "msg", TypeSize: 4}, Type: &BufferType{TypeCommon: TypeCommon{TypeName: "array", ArgDir: 1, IsVarlen: true}}},
		&LenType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "len", FldName: "size", TypeSize: 4}}, Buf: "msg"},
	}},
	{NR: 33, Name: "read_msg", CallName: "read_msg", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "msg_id", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "offset", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "msg", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "ipc_msg", Dir: 1}}},
	}},
	{NR: 35, Name: "send_msg", CallName: "send_msg", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "msg", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "ipc_msg"}}},
	}},
	{NR: 20, Name: "set_cookie", CallName: "set_cookie", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle", TypeSize: 4}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "intptr", FldName: "cookie", TypeSize: 4}}},
	}},
	{NR: 24, Name: "wait", CallName: "wait", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "handle_id", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "event", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "uevent"}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "timeout_msecs", TypeSize: 4}}},
	}},
	{NR: 25, Name: "wait_any", CallName: "wait_any", Args: []Type{
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "event", TypeSize: 4}, Type: &StructType{Key: StructKey{Name: "uevent", Dir: 1}}},
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "timeout_msecs", TypeSize: 4}}},
	}},
	{NR: 1, Name: "write", CallName: "write", Args: []Type{
		&IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", FldName: "fd", TypeSize: 4}}},
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "msg", TypeSize: 4}, Type: &BufferType{TypeCommon: TypeCommon{TypeName: "array", IsVarlen: true}}},
		&LenType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "len", FldName: "size", TypeSize: 4}}, Buf: "msg"},
	}},
}

var consts_arm = []ConstValue{
	{Name: "__NR_accept", Value: 18},
	{Name: "__NR_brk", Value: 2},
	{Name: "__NR_close", Value: 19},
	{Name: "__NR_connect", Value: 17},
	{Name: "__NR_exit_etc", Value: 3},
	{Name: "__NR_finish_dma", Value: 11},
	{Name: "__NR_get_msg", Value: 32},
	{Name: "__NR_gettime", Value: 7},
	{Name: "__NR_handle_set_create", Value: 21},
	{Name: "__NR_handle_set_ctrl", Value: 22},
	{Name: "__NR_ioctl", Value: 5},
	{Name: "__NR_mmap", Value: 8},
	{Name: "__NR_munmap", Value: 9},
	{Name: "__NR_nanosleep", Value: 6},
	{Name: "__NR_port_create", Value: 16},
	{Name: "__NR_prepare_dma", Value: 10},
	{Name: "__NR_put_msg", Value: 34},
	{Name: "__NR_read", Value: 4},
	{Name: "__NR_read_msg", Value: 33},
	{Name: "__NR_send_msg", Value: 35},
	{Name: "__NR_set_cookie", Value: 20},
	{Name: "__NR_wait", Value: 24},
	{Name: "__NR_wait_any", Value: 25},
	{Name: "__NR_write", Value: 1},
}

const revision_arm = "6c90c8f358f48e5e7692725606e77e45f66eb6e8"
