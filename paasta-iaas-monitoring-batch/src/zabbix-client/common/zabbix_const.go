package common

const (
	SYSTEM_CPU_UTIL = "system.cpu.util"
	SYSTEM_CPU_NUM = "system.cpu.num"
	VM_MEMORY_UTILIZATION = "vm.memory.utilization"
	SPACE_UTILIZATION = "vfs.fs.size[/,pused]"
	TOTAL_SPACE = "vfs.fs.size[/,total]"
	USED_SPACE = "vfs.fs.size[/,used]"
	NETWORK_INPUT_PACKET = "net.if.in[*"
	NETWORK_OUTPUT_PACKET = "net.if.out[*"
	NETWORK_INPUT_DROPPED_PACKET = "net.if.in.dropped*"
	NETWORK_OUTPUT_DROPPED_PACKET = "net.if.out.dropped*"
	NETWORK_INPUT_ERROR_PACKET = "net.if.in.errors*"
	NETWORK_OUTPUT_ERROR_PACKET = "net.if.out.errors*"
	DISK_READ_RATE = "vfs.dev.read.rate*"
	DISK_WRITE_RATE = "vfs.dev.write.rate*"
)

