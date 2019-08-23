#! /usr/bin/python3
import os

# path_root = "/home/yuh/data"
path_root = "/home/yhao"
path_git = os.path.join(path_root, "git")
path_git_work = os.path.join(path_git, "work")
path_result = os.path.join(path_git_work, "result")
path_taint = os.path.join(path_git_work, "script/dependency/taint_info")
file_default_json = os.path.join(path_git_work, "script/dependency/con/default.json")
name_run = "run.py"
path_default_run = os.path.join(path_git_work, "script/dependency/con", name_run)
name_with_dra = "result-with-dra"
name_without_dra = "result-without-dra"
path_linux_bc = os.path.join(path_root, "benchmark/linux/16-linux-clang-np-bc-f")
path_linux = os.path.join(path_root, "benchmark/linux/13-linux-clang-np")
path_kernel = os.path.join(path_linux, "arch/x86/boot/bzImage")

path_syzkaller = os.path.join(path_git, "gopath/src/github.com/google/syzkaller")
file_syzkaller = os.path.join(path_syzkaller, "bin/syz-manager")

name_driver = "built-in"
file_taint = name_driver + ".taint"
file_asm = name_driver + ".s"
file_bc = name_driver + ".bc"
file_json = name_driver + ".json"

dev = {
    "dev_ashmem": {
        "enable_syscalls": [
            "openat$ashmem",
            "ioctl$ASHMEM_SET_NAME",
            "ioctl$ASHMEM_GET_NAME",
            "ioctl$ASHMEM_SET_SIZE",
            "ioctl$ASHMEM_GET_SIZE",
            "ioctl$ASHMEM_SET_PROT_MASK",
            "ioctl$ASHMEM_GET_PROT_MASK",
            "ioctl$ASHMEM_GET_PIN_STATUS",
            "ioctl$ASHMEM_PURGE_ALL_CACHES",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/staging/android/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/staging/android"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_ashmem_ioctl_serialize"),
    },
    "dev_binder": {
        "enable_syscalls": [
            "syz_open_dev$binder",
            "mmap$binder",
            "ioctl$BINDER_SET_MAX_THREADS",
            "ioctl$BINDER_SET_CONTEXT_MGR",
            "ioctl$BINDER_THREAD_EXIT",
            "ioctl$BINDER_GET_NODE_DEBUG_INFO",
            "ioctl$BINDER_WRITE_READ",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
    "dev_block": {
        "enable_syscalls": [
            "openat$nullb",
            "openat$md",
            "ioctl$BLKTRACESETUP",
            "ioctl$BLKTRACESTART",
            "ioctl$BLKTRACESTOP",
            "ioctl$BLKTRACETEARDOWN",
            "ioctl$BLKFLSBUF",
            "ioctl$BLKROSET",
            "ioctl$BLKDISCARD",
            "ioctl$BLKSECDISCARD",
            "ioctl$BLKZEROOUT",
            "ioctl$BLKREPORTZONE",
            "ioctl$BLKRESETZONE",
            "ioctl$BLKRAGET",
            "ioctl$BLKROGET",
            "ioctl$BLKBSZGET",
            "ioctl$BLKPBSZGET",
            "ioctl$BLKIOMIN",
            "ioctl$BLKIOOPT",
            "ioctl$BLKALIGNOFF",
            "ioctl$BLKSECTGET",
            "ioctl$BLKROTATIONAL",
            "ioctl$BLKFRASET",
            "ioctl$BLKBSZSET",
            "ioctl$BLKPG",
            "ioctl$BLKRRPART",
            "ioctl$BLKGETSIZE",
            "ioctl$BLKGETSIZE64",
            "ioctl$HDIO_GETGEO",
            "ioctl$IOC_PR_REGISTER",
            "ioctl$IOC_PR_RESERVE",
            "ioctl$IOC_PR_RELEASE",
            "ioctl$IOC_PR_PREEMPT",
            "ioctl$IOC_PR_PREEMPT_ABORT",
            "ioctl$IOC_PR_CLEAR",
        ],
        "file_bc": os.path.join(path_linux_bc, "block/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "block/"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_blkdev_ioctl_serialize"),
    },
    "dev_cdrom": {
        "enable_syscalls": [
            "syz_open_dev$CDROM_DEV_LINK",
            "ioctl$CDROMPAUSE",
            "ioctl$CDROMRESUME",
            "ioctl$CDROMPLAYMSF",
            "ioctl$CDROMPLAYTRKIND",
            "ioctl$CDROMREADTOCHDR",
            "ioctl$CDROMREADTOCENTRY",
            "ioctl$CDROMSTOP",
            "ioctl$CDROMSTART",
            "ioctl$CDROMEJECT",
            "ioctl$CDROMVOLCTRL",
            "ioctl$CDROMSUBCHNL",
            "ioctl$CDROMREADMODE2",
            "ioctl$CDROMREADMODE1",
            "ioctl$CDROMREADAUDIO",
            "ioctl$CDROMEJECT_SW",
            "ioctl$CDROMMULTISESSION",
            "ioctl$CDROM_GET_MCN",
            "ioctl$CDROMRESET",
            "ioctl$CDROMVOLREAD",
            "ioctl$CDROMREADRAW",
            "ioctl$CDROMREADCOOKED",
            "ioctl$CDROMSEEK",
            "ioctl$CDROMPLAYBLK",
            "ioctl$CDROMREADALL",
            "ioctl$CDROMGETSPINDOWN",
            "ioctl$CDROMSETSPINDOWN",
            "ioctl$CDROMCLOSETRAY",
            "ioctl$CDROM_SET_OPTIONS",
            "ioctl$CDROM_CLEAR_OPTIONS",
            "ioctl$CDROM_SELECT_SPEED",
            "ioctl$CDROM_SELECT_DISK",
            "ioctl$CDROM_MEDIA_CHANGED",
            "ioctl$CDROM_DISC_STATUS",
            "ioctl$CDROM_CHANGER_NSLOTS",
            "ioctl$CDROM_LOCKDOOR",
            "ioctl$CDROM_DEBUG",
            "ioctl$CDROM_GET_CAPABILITY",
            "ioctl$CDROMAUDIOBUFSIZ",
            "ioctl$DVD_READ_STRUCT",
            "ioctl$DVD_WRITE_STRUCT",
            "ioctl$DVD_AUTH",
            "ioctl$CDROM_SEND_PACKET",
            "ioctl$CDROM_NEXT_WRITABLE",
            "ioctl$CDROM_LAST_WRITTEN",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/cdrom/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/cdrom"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_cdrom_ioctl_serialize"),
    },
    "dev_dri": {
        "enable_syscalls": [
            "syz_open_dev$dri",
            "syz_open_dev$dricontrol",
            "syz_open_dev$drirender",
            "ioctl$DRM_IOCTL_VERSION",
            "ioctl$DRM_IOCTL_GET_UNIQUE",
            "ioctl$DRM_IOCTL_GET_MAGIC",
            "ioctl$DRM_IOCTL_IRQ_BUSID",
            "ioctl$DRM_IOCTL_GET_MAP",
            "ioctl$DRM_IOCTL_GET_CLIENT",
            "ioctl$DRM_IOCTL_GET_STATS",
            "ioctl$DRM_IOCTL_GET_CAP",
            "ioctl$DRM_IOCTL_SET_CLIENT_CAP",
            "ioctl$DRM_IOCTL_SET_VERSION",
            "ioctl$DRM_IOCTL_SET_UNIQUE",
            "ioctl$DRM_IOCTL_AUTH_MAGIC",
            "ioctl$DRM_IOCTL_ADD_MAP",
            "ioctl$DRM_IOCTL_RM_MAP",
            "ioctl$DRM_IOCTL_SET_SAREA_CTX",
            "ioctl$DRM_IOCTL_GET_SAREA_CTX",
            "ioctl$DRM_IOCTL_SET_MASTER",
            "ioctl$DRM_IOCTL_DROP_MASTER",
            "ioctl$DRM_IOCTL_ADD_CTX",
            "ioctl$DRM_IOCTL_RM_CTX",
            "ioctl$DRM_IOCTL_GET_CTX",
            "ioctl$DRM_IOCTL_SWITCH_CTX",
            "ioctl$DRM_IOCTL_NEW_CTX",
            "ioctl$DRM_IOCTL_RES_CTX",
            "ioctl$DRM_IOCTL_LOCK",
            "ioctl$DRM_IOCTL_UNLOCK",
            "ioctl$DRM_IOCTL_ADD_BUFS",
            "ioctl$DRM_IOCTL_MARK_BUFS",
            "ioctl$DRM_IOCTL_INFO_BUFS",
            "ioctl$DRM_IOCTL_MAP_BUFS",
            "ioctl$DRM_IOCTL_FREE_BUFS",
            "ioctl$DRM_IOCTL_DMA",
            "ioctl$DRM_IOCTL_CONTROL",
            "ioctl$DRM_IOCTL_AGP_ACQUIRE",
            "ioctl$DRM_IOCTL_AGP_RELEASE",
            "ioctl$DRM_IOCTL_AGP_ENABLE",
            "ioctl$DRM_IOCTL_AGP_INFO",
            "ioctl$DRM_IOCTL_AGP_ALLOC",
            "ioctl$DRM_IOCTL_AGP_FREE",
            "ioctl$DRM_IOCTL_AGP_BIND",
            "ioctl$DRM_IOCTL_AGP_UNBIND",
            "ioctl$DRM_IOCTL_SG_ALLOC",
            "ioctl$DRM_IOCTL_SG_FREE",
            "ioctl$DRM_IOCTL_WAIT_VBLANK",
            "ioctl$DRM_IOCTL_MODESET_CTL",
            "ioctl$DRM_IOCTL_GEM_OPEN",
            "ioctl$DRM_IOCTL_GEM_CLOSE",
            "ioctl$DRM_IOCTL_GEM_FLINK",
            "ioctl$DRM_IOCTL_MODE_GETRESOURCES",
            "ioctl$DRM_IOCTL_PRIME_HANDLE_TO_FD",
            "ioctl$DRM_IOCTL_PRIME_FD_TO_HANDLE",
            "ioctl$DRM_IOCTL_MODE_GETPLANERESOURCES",
            "ioctl$DRM_IOCTL_MODE_GETCRTC",
            "ioctl$DRM_IOCTL_MODE_SETCRTC",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
    "dev_floppy": {
        "enable_syscalls": [
            "syz_open_dev$floppy",
            "ioctl$FLOPPY_FDEJECT",
            "ioctl$FLOPPY_FDCLRPRM",
            "ioctl$FLOPPY_FDSETPRM",
            "ioctl$FLOPPY_FDDEFPRM",
            "ioctl$FLOPPY_FDGETPRM",
            "ioctl$FLOPPY_FDMSGON",
            "ioctl$FLOPPY_FDMSGOFF",
            "ioctl$FLOPPY_FDFMTBEG",
            "ioctl$FLOPPY_FDFMTTRK",
            "ioctl$FLOPPY_FDFMTEND",
            "ioctl$FLOPPY_FDFLUSH",
            "ioctl$FLOPPY_FDSETEMSGTRESH",
            "ioctl$FLOPPY_FDGETMAXERRS",
            "ioctl$FLOPPY_FDSETMAXERRS",
            "ioctl$FLOPPY_FDGETDRVTYP",
            "ioctl$FLOPPY_FDSETDRVPRM",
            "ioctl$FLOPPY_FDGETDRVPRM",
            "ioctl$FLOPPY_FDPOLLDRVSTAT",
            "ioctl$FLOPPY_FDGETDRVSTAT",
            "ioctl$FLOPPY_FDRESET",
            "ioctl$FLOPPY_FDGETFDCSTAT",
            "ioctl$FLOPPY_FDWERRORCLR",
            "ioctl$FLOPPY_FDWERRORGET",
            "ioctl$FLOPPY_FDRAWCMD",
            "ioctl$FLOPPY_FDTWADDLE",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
    "dev_i2c": {
        "enable_syscalls": [
            "syz_open_dev$I2C",
            "ioctl$I2C_RETRIES",
            "ioctl$I2C_TIMEOUT",
            "ioctl$I2C_SLAVE",
            "ioctl$I2C_SLAVE_FORCE",
            "ioctl$I2C_TENBIT",
            "ioctl$I2C_PEC",
            "ioctl$I2C_FUNCS",
            "ioctl$I2C_RDWR",
            "ioctl$I2C_SMBUS",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/i2c/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/i2c"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_i2cdev_ioctl_serialize"),
    },
    "dev_input": {
        "enable_syscalls": [
            "syz_open_dev$mouse",
            "syz_open_dev$mice",
            "syz_open_dev$evdev",
            "write$evdev",
            "ioctl$EVIOCGVERSION",
            "ioctl$EVIOCGID",
            "ioctl$EVIOCGREP",
            "ioctl$EVIOCGKEYCODE",
            "ioctl$EVIOCGKEYCODE_V2",
            "ioctl$EVIOCGEFFECTS",
            "ioctl$EVIOCGMASK",
            "ioctl$EVIOCGNAME",
            "ioctl$EVIOCGPHYS",
            "ioctl$EVIOCGUNIQ",
            "ioctl$EVIOCGPROP",
            "ioctl$EVIOCGMTSLOTS",
            "ioctl$EVIOCGKEY",
            "ioctl$EVIOCGLED",
            "ioctl$EVIOCGSND",
            "ioctl$EVIOCGSW",
            "ioctl$EVIOCGBITKEY",
            "ioctl$EVIOCGBITSND",
            "ioctl$EVIOCGBITSW",
            "ioctl$EVIOCGABS0",
            "ioctl$EVIOCGABS20",
            "ioctl$EVIOCGABS2F",
            "ioctl$EVIOCGABS3F",
            "ioctl$EVIOCSREP",
            "ioctl$EVIOCSKEYCODE",
            "ioctl$EVIOCSKEYCODE_V2",
            "ioctl$EVIOCSFF",
            "ioctl$EVIOCRMFF",
            "ioctl$EVIOCGRAB",
            "ioctl$EVIOCREVOKE",
            "ioctl$EVIOCSMASK",
            "ioctl$EVIOCSCLOCKID",
            "ioctl$EVIOCSABS0",
            "ioctl$EVIOCSABS20",
            "ioctl$EVIOCSABS2F",
            "ioctl$EVIOCSABS3F",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/input/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/input"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_evdev_ioctl_handler_serialize"),
    },
    "dev_ion": {
        "enable_syscalls": [
            "openat$ion",
            "ioctl$ION_IOC_ALLOC",
            "ioctl$ION_IOC_HEAP_QUERY",
            "ioctl$DMA_BUF_IOCTL_SYNC",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/staging/android/ion/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/staging/android/ion"),
            os.path.join(path_linux_bc, "drivers/dma-buf"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_ion_ioctl_serialize"),
    },
    "dev_kvm": {
        "enable_syscalls": [
            "openat$kvm",
            "ioctl$KVM_CREATE_VM",
            "ioctl$KVM_GET_MSR_INDEX_LIST",
            "ioctl$KVM_CHECK_EXTENSION",
            "ioctl$KVM_GET_VCPU_MMAP_SIZE",
            "ioctl$KVM_GET_SUPPORTED_CPUID",
            "ioctl$KVM_GET_EMULATED_CPUID",
            "ioctl$KVM_X86_GET_MCE_CAP_SUPPORTED",
            "ioctl$KVM_GET_API_VERSION",
            "ioctl$KVM_CREATE_VCPU",
            "ioctl$KVM_CHECK_EXTENSION_VM",
            "ioctl$KVM_GET_DIRTY_LOG",
            "ioctl$KVM_CREATE_IRQCHIP",
            "ioctl$KVM_IRQ_LINE",
            "ioctl$KVM_IRQ_LINE_STATUS",
            "ioctl$KVM_GET_IRQCHIP",
            "ioctl$KVM_SET_IRQCHIP",
            "ioctl$KVM_XEN_HVM_CONFIG",
            "ioctl$KVM_GET_CLOCK",
            "ioctl$KVM_SET_CLOCK",
            "ioctl$KVM_SET_USER_MEMORY_REGION",
            "ioctl$KVM_SET_TSS_ADDR",
            "ioctl$KVM_ENABLE_CAP",
            "ioctl$KVM_SET_IDENTITY_MAP_ADDR",
            "ioctl$KVM_SET_BOOT_CPU_ID",
            "ioctl$KVM_PPC_GET_PVINFO",
            "ioctl$KVM_ASSIGN_PCI_DEVICE",
            "ioctl$KVM_DEASSIGN_PCI_DEVICE",
            "ioctl$KVM_ASSIGN_DEV_IRQ",
            "ioctl$KVM_DEASSIGN_DEV_IRQ",
            "ioctl$KVM_SET_GSI_ROUTING",
            "ioctl$KVM_ASSIGN_SET_MSIX_NR",
            "ioctl$KVM_ASSIGN_SET_MSIX_ENTRY",
            "ioctl$KVM_IOEVENTFD",
            "ioctl$KVM_ASSIGN_SET_INTX_MASK",
            "ioctl$KVM_SIGNAL_MSI",
            "ioctl$KVM_CREATE_PIT2",
            "ioctl$KVM_GET_PIT",
            "ioctl$KVM_SET_PIT",
            "ioctl$KVM_GET_PIT2",
            "ioctl$KVM_SET_PIT2",
            "ioctl$KVM_PPC_GET_SMMU_INFO",
            "ioctl$KVM_IRQFD",
            "ioctl$KVM_PPC_ALLOCATE_HTAB",
            "ioctl$KVM_CREATE_DEVICE",
            "ioctl$KVM_REGISTER_COALESCED_MMIO",
            "ioctl$KVM_UNREGISTER_COALESCED_MMIO",
            "ioctl$KVM_SET_NR_MMU_PAGES",
            "ioctl$KVM_GET_NR_MMU_PAGES",
            "ioctl$KVM_REINJECT_CONTROL",
            "ioctl$KVM_HYPERV_EVENTFD",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
    "dev_loop": {
        "enable_syscalls": [
            "syz_open_dev$loop",
            "ioctl$LOOP_SET_FD",
            "ioctl$LOOP_CHANGE_FD",
            "ioctl$LOOP_CLR_FD",
            "ioctl$LOOP_SET_STATUS",
            "ioctl$LOOP_SET_STATUS64",
            "ioctl$LOOP_GET_STATUS",
            "ioctl$LOOP_GET_STATUS64",
            "ioctl$LOOP_SET_CAPACITY",
            "ioctl$LOOP_SET_DIRECT_IO",
            "ioctl$LOOP_SET_BLOCK_SIZE",
            "openat$loop_ctrl",
            "ioctl$LOOP_CTL_GET_FREE",
            "ioctl$LOOP_CTL_ADD",
            "ioctl$LOOP_CTL_REMOVE",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/block/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/block"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_lo_ioctl_serialize"),
    },
    "dev_nbd": {
        "enable_syscalls": [
            "syz_open_dev$ndb",
            "ioctl$NBD_DISCONNECT",
            "ioctl$NBD_CLEAR_SOCK",
            "ioctl$NBD_SET_SOCK",
            "ioctl$NBD_SET_BLKSIZE",
            "ioctl$NBD_SET_SIZE",
            "ioctl$NBD_SET_SIZE_BLOCKS",
            "ioctl$NBD_SET_TIMEOUT",
            "ioctl$NBD_SET_FLAGS",
            "ioctl$NBD_DO_IT",
            "ioctl$NBD_CLEAR_QUE",
            "syz_genetlink_get_family_id$nbd",
            "sendmsg$NBD_CMD_CONNECT",
            "sendmsg$NBD_CMD_DISCONNECT",
            "sendmsg$NBD_CMD_RECONFIGURE",
            "sendmsg$NBD_CMD_STATUS",
            "socketpair$nbd",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/block/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/block"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_nbd_ioctl_serialize"),
    },
    "dev_net_tun": {
        "enable_syscalls": [
            "openat$tun",
            "write$tun",
            "ioctl$TUNGETFEATURES",
            "ioctl$TUNSETQUEUE",
            "ioctl$TUNSETIFF",
            "ioctl$TUNSETIFINDEX",
            "ioctl$TUNGETIFF",
            "ioctl$TUNSETNOCSUM",
            "ioctl$TUNSETPERSIST",
            "ioctl$TUNSETOWNER",
            "ioctl$TUNSETGROUP",
            "ioctl$TUNSETLINK",
            "ioctl$TUNSETOFFLOAD",
            "ioctl$TUNSETTXFILTER",
            "ioctl$SIOCGIFHWADDR",
            "ioctl$SIOCSIFHWADDR",
            "ioctl$TUNGETSNDBUF",
            "ioctl$TUNSETSNDBUF",
            "ioctl$TUNGETVNETHDRSZ",
            "ioctl$TUNSETVNETHDRSZ",
            "ioctl$TUNATTACHFILTER",
            "ioctl$TUNDETACHFILTER",
            "ioctl$TUNGETFILTER",
            "ioctl$TUNSETSTEERINGEBPF",
            "ioctl$TUNSETFILTEREBPF",
            "ioctl$TUNSETVNETLE",
            "ioctl$TUNSETVNETBE",
            "syz_open_dev$loop",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/net/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/net"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info___tun_chr_ioctl_serialize"),
    },
    "dev_ppp": {
        "enable_syscalls": [
            "openat$ppp",
            "write$ppp",
            "ioctl$PPPIOCNEWUNIT",
            "ioctl$PPPIOCATTACH",
            "ioctl$PPPIOCATTCHAN",
            "ioctl$PPPIOCCONNECT",
            "ioctl$PPPIOCDISCONN",
            "ioctl$PPPIOCSCOMPRESS",
            "ioctl$PPPIOCGUNIT",
            "ioctl$PPPIOCSDEBUG",
            "ioctl$PPPIOCGDEBUG",
            "ioctl$PPPIOCGIDLE",
            "ioctl$PPPIOCSMAXCID",
            "ioctl$PPPIOCGNPMODE",
            "ioctl$PPPIOCSNPMODE",
            "ioctl$PPPIOCSPASS",
            "ioctl$PPPIOCSACTIVE",
            "ioctl$PPPIOCSMRRU",
            "ioctl$PPPIOCSMRU1",
            "ioctl$PPPIOCSFLAGS1",
            "ioctl$PPPIOCGFLAGS1",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/net/ppp/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/net/ppp"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_ppp_ioctl_serialize"),
    },
    "dev_ptmx": {
        "enable_syscalls": [
            "openat$ptmx",
            "syz_open_pts",
            "ioctl$TCGETS",
            "ioctl$TCSETS",
            "ioctl$TCSETSW",
            "ioctl$TCSETSF",
            "ioctl$TCGETA",
            "ioctl$TCSETA",
            "ioctl$TCSETAW",
            "ioctl$TCSETAF",
            "ioctl$TIOCGLCKTRMIOS",
            "ioctl$TIOCSLCKTRMIOS",
            "ioctl$TIOCGWINSZ",
            "ioctl$TIOCSWINSZ",
            "ioctl$TCSBRK",
            "ioctl$TCSBRKP",
            "ioctl$TIOCSBRK",
            "ioctl$TIOCCBRK",
            "ioctl$TCXONC",
            "ioctl$FIONREAD",
            "ioctl$TIOCOUTQ",
            "ioctl$TCFLSH",
            "ioctl$TIOCSTI",
            "ioctl$TIOCCONS",
            "ioctl$TIOCSCTTY",
            "ioctl$TIOCNOTTY",
            "ioctl$TIOCGPGRP",
            "ioctl$TIOCSPGRP",
            "ioctl$TIOCGSID",
            "ioctl$TIOCEXCL",
            "ioctl$TIOCNXCL",
            "ioctl$TIOCGETD",
            "ioctl$TIOCSETD",
            "ioctl$TIOCPKT",
            "ioctl$TIOCMGET",
            "ioctl$TIOCMSET",
            "ioctl$TIOCMBIC",
            "ioctl$TIOCMBIS",
            "ioctl$TIOCGSOFTCAR",
            "ioctl$TIOCSSOFTCAR",
            "ioctl$KDGETLED",
            "ioctl$KDSETLED",
            "ioctl$KDGKBLED",
            "ioctl$KDSKBLED",
            "ioctl$KDGKBTYPE",
            "ioctl$KDADDIO",
            "ioctl$KDDELIO",
            "ioctl$KDENABIO",
            "ioctl$KDDISABIO",
            "ioctl$KDSETMODE",
            "ioctl$KDGETMODE",
            "ioctl$KDMKTONE",
            "ioctl$KIOCSOUND",
            "ioctl$GIO_CMAP",
            "ioctl$PIO_CMAP",
            "ioctl$GIO_FONT",
            "ioctl$GIO_FONTX",
            "ioctl$PIO_FONT",
            "ioctl$PIO_FONTX",
            "ioctl$PIO_FONTRESET",
            "ioctl$GIO_SCRNMAP",
            "ioctl$GIO_UNISCRNMAP",
            "ioctl$PIO_SCRNMAP",
            "ioctl$PIO_UNISCRNMAP",
            "ioctl$GIO_UNIMAP",
            "ioctl$PIO_UNIMAP",
            "ioctl$PIO_UNIMAPCLR",
            "ioctl$KDGKBMODE",
            "ioctl$KDSKBMODE",
            "ioctl$KDGKBMETA",
            "ioctl$KDSKBMETA",
            "ioctl$KDGKBENT",
            "ioctl$KDGKBSENT",
            "ioctl$KDSKBSENT",
            "ioctl$KDGKBDIACR",
            "ioctl$KDGETKEYCODE",
            "ioctl$KDSETKEYCODE",
            "ioctl$KDSIGACCEPT",
            "ioctl$VT_OPENQRY",
            "ioctl$VT_GETMODE",
            "ioctl$VT_SETMODE",
            "ioctl$VT_GETSTATE",
            "ioctl$VT_RELDISP",
            "ioctl$VT_ACTIVATE",
            "ioctl$VT_WAITACTIVE",
            "ioctl$VT_DISALLOCATE",
            "ioctl$VT_RESIZE",
            "ioctl$VT_RESIZEX",
            "ioctl$TIOCLINUX2",
            "ioctl$TIOCLINUX3",
            "ioctl$TIOCLINUX4",
            "ioctl$TIOCLINUX5",
            "ioctl$TIOCLINUX6",
            "ioctl$TIOCLINUX7",
            "ioctl$TIOCGSERIAL",
            "ioctl$TIOCSSERIAL",
            "ioctl$TCGETS2",
            "ioctl$TCSETS2",
            "ioctl$TIOCSERGETLSR",
            "ioctl$TIOCGRS485",
            "ioctl$TIOCSRS485",
            "ioctl$TIOCGISO7816",
            "ioctl$TIOCSISO7816",
            "ioctl$TIOCSPTLCK",
            "ioctl$TIOCGPTLCK",
            "ioctl$TIOCGPKT",
            "ioctl$TIOCSIG",
            "ioctl$TIOCVHANGUP",
            "ioctl$TIOCGDEV",
            "ioctl$TCGETX",
            "ioctl$TCSETX",
            "ioctl$TCSETXF",
            "ioctl$TCSETXW",
            "ioctl$TIOCMIWAIT",
            "ioctl$TIOCGICOUNT",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/tty/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/tty"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_tty_ioctl_serialize"),
    },
    "dev_random": {
        "enable_syscalls": [
            "openat$random",
            "openat$urandom",
            "ioctl$RNDGETENTCNT",
            "ioctl$RNDADDTOENTCNT",
            "ioctl$RNDADDENTROPY",
            "ioctl$RNDZAPENTCNT",
            "ioctl$RNDCLEARPOOL",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/char/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/char"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_random_ioctl_serialize"),
    },
    "dev_rfkill": {
        "enable_syscalls": [
            "openat$rfkill",
            "write$rfkill",
            "read$rfkill",
            "ioctl$RFKILL_IOCTL_NOINPUT",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
    "dev_rtc": {
        "enable_syscalls": [
            "openat$rtc",
            "syz_open_dev$rtc",
            "ioctl$RTC_AIE_OFF",
            "ioctl$RTC_AIE_ON",
            "ioctl$RTC_PIE_OFF",
            "ioctl$RTC_PIE_ON",
            "ioctl$RTC_UIE_OFF",
            "ioctl$RTC_UIE_ON",
            "ioctl$RTC_WIE_ON",
            "ioctl$RTC_WIE_OFF",
            "ioctl$RTC_ALM_READ",
            "ioctl$RTC_ALM_SET",
            "ioctl$RTC_RD_TIME",
            "ioctl$RTC_SET_TIME",
            "ioctl$RTC_IRQP_READ",
            "ioctl$RTC_IRQP_SET",
            "ioctl$RTC_EPOCH_READ",
            "ioctl$RTC_EPOCH_SET",
            "ioctl$RTC_WKALM_RD",
            "ioctl$RTC_WKALM_SET",
            "ioctl$RTC_PLL_GET",
            "ioctl$RTC_PLL_SET",
            "ioctl$RTC_VL_READ",
            "ioctl$RTC_VL_CLR",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/rtc/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/rtc"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_rtc_dev_ioctl_serialize"),
    },
    "dev_sg": {
        "enable_syscalls": [
            "syz_open_dev$sg",
            "ioctl$SG_IO",
            "ioctl$SG_SET_TIMEOUT",
            "ioctl$SG_GET_TIMEOUT",
            "ioctl$SG_GET_LOW_DMA",
            "ioctl$SG_GET_SCSI_ID",
            "ioctl$SG_SET_FORCE_PACK_ID",
            "ioctl$SG_GET_PACK_ID",
            "ioctl$SG_GET_NUM_WAITING",
            "ioctl$SG_GET_SG_TABLESIZE",
            "ioctl$SG_SET_RESERVED_SIZE",
            "ioctl$SG_GET_RESERVED_SIZE",
            "ioctl$SG_GET_COMMAND_Q",
            "ioctl$SG_GET_KEEP_ORPHAN",
            "ioctl$SG_GET_VERSION_NUM",
            "ioctl$SG_GET_ACCESS_COUNT",
            "ioctl$SG_EMULATED_HOST",
            "ioctl$SG_SET_COMMAND_Q",
            "ioctl$SG_SET_KEEP_ORPHAN",
            "ioctl$SG_NEXT_CMD_LEN",
            "ioctl$SG_SET_DEBUG",
            "ioctl$SG_SCSI_RESET",
            "ioctl$SG_GET_REQUEST_TABLE",
            "ioctl$SCSI_IOCTL_SEND_COMMAND",
            "ioctl$SCSI_IOCTL_TEST_UNIT_READY",
            "ioctl$SCSI_IOCTL_DOORLOCK",
            "ioctl$SCSI_IOCTL_DOORUNLOCK",
            "ioctl$SCSI_IOCTL_START_UNIT",
            "ioctl$SCSI_IOCTL_STOP_UNIT",
            "ioctl$SCSI_IOCTL_SYNC",
            "ioctl$SCSI_IOCTL_BENCHMARK_COMMAND",
            "ioctl$SCSI_IOCTL_GET_BUS_NUMBER",
            "ioctl$SCSI_IOCTL_GET_PCI",
            "ioctl$SCSI_IOCTL_PROBE_HOST",
            "ioctl$SCSI_IOCTL_GET_IDLUN",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/scsi/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/scsi"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_sg_ioctl_serialize"),
    },
    "dev_snd_control": {
        "enable_syscalls": [
            "syz_open_dev$sndctrl",
            "ioctl$SNDRV_CTL_IOCTL_PVERSION",
            "ioctl$SNDRV_CTL_IOCTL_CARD_INFO",
            "ioctl$SNDRV_CTL_IOCTL_HWDEP_INFO",
            "ioctl$SNDRV_CTL_IOCTL_PCM_NEXT_DEVICE",
            "ioctl$SNDRV_CTL_IOCTL_POWER_STATE",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_LIST",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_INFO",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_READ",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_WRITE",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_LOCK",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_UNLOCK",
            "ioctl$SNDRV_CTL_IOCTL_SUBSCRIBE_EVENTS",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_ADD",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_REPLACE",
            "ioctl$SNDRV_CTL_IOCTL_ELEM_REMOVE",
            "ioctl$SNDRV_CTL_IOCTL_TLV_READ",
            "ioctl$SNDRV_CTL_IOCTL_TLV_WRITE",
            "ioctl$SNDRV_CTL_IOCTL_TLV_COMMAND",
            "ioctl$SNDRV_CTL_IOCTL_HWDEP_NEXT_DEVICE",
            "ioctl$SNDRV_CTL_IOCTL_PCM_INFO",
            "ioctl$SNDRV_CTL_IOCTL_PCM_PREFER_SUBDEVICE",
            "ioctl$SNDRV_CTL_IOCTL_RAWMIDI_NEXT_DEVICE",
            "ioctl$SNDRV_CTL_IOCTL_RAWMIDI_INFO",
            "ioctl$SNDRV_CTL_IOCTL_RAWMIDI_PREFER_SUBDEVICE",
        ],
        "file_bc": os.path.join(path_linux_bc, "sound/core/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "sound/core"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_snd_ctl_ioctl_serialize"),
    },
    "dev_snd_seq": {
        "enable_syscalls": [
            "syz_open_dev$sndseq",
            "write$sndseq",
            "ioctl$SNDRV_SEQ_IOCTL_PVERSION",
            "ioctl$SNDRV_SEQ_IOCTL_CLIENT_ID",
            "ioctl$SNDRV_SEQ_IOCTL_SYSTEM_INFO",
            "ioctl$SNDRV_SEQ_IOCTL_RUNNING_MODE",
            "ioctl$SNDRV_SEQ_IOCTL_GET_CLIENT_INFO",
            "ioctl$SNDRV_SEQ_IOCTL_SET_CLIENT_INFO",
            "ioctl$SNDRV_SEQ_IOCTL_CREATE_PORT",
            "ioctl$SNDRV_SEQ_IOCTL_DELETE_PORT",
            "ioctl$SNDRV_SEQ_IOCTL_GET_PORT_INFO",
            "ioctl$SNDRV_SEQ_IOCTL_SET_PORT_INFO",
            "ioctl$SNDRV_SEQ_IOCTL_SUBSCRIBE_PORT",
            "ioctl$SNDRV_SEQ_IOCTL_UNSUBSCRIBE_PORT",
            "ioctl$SNDRV_SEQ_IOCTL_CREATE_QUEUE",
            "ioctl$SNDRV_SEQ_IOCTL_DELETE_QUEUE",
            "ioctl$SNDRV_SEQ_IOCTL_GET_QUEUE_INFO",
            "ioctl$SNDRV_SEQ_IOCTL_SET_QUEUE_INFO",
            "ioctl$SNDRV_SEQ_IOCTL_GET_NAMED_QUEUE",
            "ioctl$SNDRV_SEQ_IOCTL_GET_QUEUE_STATUS",
            "ioctl$SNDRV_SEQ_IOCTL_GET_QUEUE_TEMPO",
            "ioctl$SNDRV_SEQ_IOCTL_SET_QUEUE_TEMPO",
            "ioctl$SNDRV_SEQ_IOCTL_GET_QUEUE_TIMER",
            "ioctl$SNDRV_SEQ_IOCTL_SET_QUEUE_TIMER",
            "ioctl$SNDRV_SEQ_IOCTL_GET_QUEUE_CLIENT",
            "ioctl$SNDRV_SEQ_IOCTL_SET_QUEUE_CLIENT",
            "ioctl$SNDRV_SEQ_IOCTL_GET_CLIENT_POOL",
            "ioctl$SNDRV_SEQ_IOCTL_SET_CLIENT_POOL",
            "ioctl$SNDRV_SEQ_IOCTL_REMOVE_EVENTS",
            "ioctl$SNDRV_SEQ_IOCTL_QUERY_SUBS",
            "ioctl$SNDRV_SEQ_IOCTL_GET_SUBSCRIPTION",
            "ioctl$SNDRV_SEQ_IOCTL_QUERY_NEXT_CLIENT",
            "ioctl$SNDRV_SEQ_IOCTL_QUERY_NEXT_PORT",
        ],
        "file_bc": os.path.join(path_linux_bc, "sound/core/seq/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "sound/core/seq"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_snd_seq_ioctl_serialize"),
    },
    "dev_snd_timer": {
        "enable_syscalls": [
            "syz_open_dev$sndtimer",
            "ioctl$SNDRV_TIMER_IOCTL_PVERSION",
            "ioctl$SNDRV_TIMER_IOCTL_NEXT_DEVICE",
            "ioctl$SNDRV_TIMER_IOCTL_TREAD",
            "ioctl$SNDRV_TIMER_IOCTL_GINFO",
            "ioctl$SNDRV_TIMER_IOCTL_GPARAMS",
            "ioctl$SNDRV_TIMER_IOCTL_GSTATUS",
            "ioctl$SNDRV_TIMER_IOCTL_SELECT",
            "ioctl$SNDRV_TIMER_IOCTL_INFO",
            "ioctl$SNDRV_TIMER_IOCTL_PARAMS",
            "ioctl$SNDRV_TIMER_IOCTL_STATUS",
            "ioctl$SNDRV_TIMER_IOCTL_START",
            "ioctl$SNDRV_TIMER_IOCTL_STOP",
            "ioctl$SNDRV_TIMER_IOCTL_CONTINUE",
            "ioctl$SNDRV_TIMER_IOCTL_PAUSE",
        ],
        "file_bc": os.path.join(path_linux_bc, "sound/core/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "sound/core"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info___snd_timer_user_ioctl_serialize"),
    },
    "dev_sr": {
        "enable_syscalls": [
            "openat$sr",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_uinput_ioctl_handler_serialize"),
    },
    "dev_uhid": {
        "enable_syscalls": [
            "openat$uhid",
            "write$UHID_CREATE",
            "write$UHID_CREATE2",
            "write$UHID_DESTROY",
            "write$UHID_INPUT",
            "write$UHID_INPUT2",
            "write$UHID_GET_REPORT_REPLY",
            "write$UHID_SET_REPORT_REPLY",
        ],
        "file_bc": os.path.join(path_linux_bc, "drivers/input/misc/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "drivers/input/misc"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_uinput_ioctl_handler_serialize"),
    },
    "dev_uinput": {
        "enable_syscalls": [
            "openat$uinput",
            "write$uinput_user_dev",
            "write$input_event",
            "ioctl$UI_GET_VERSION",
            "ioctl$UI_DEV_CREATE",
            "ioctl$UI_DEV_DESTROY",
            "ioctl$UI_DEV_SETUP",
            "ioctl$UI_SET_EVBIT",
            "ioctl$UI_SET_KEYBIT",
            "ioctl$UI_SET_RELBIT",
            "ioctl$UI_SET_MSCBIT",
            "ioctl$UI_SET_ABSBIT",
            "ioctl$UI_SET_LEDBIT",
            "ioctl$UI_SET_SNDBIT",
            "ioctl$UI_SET_FFBIT",
            "ioctl$UI_SET_SWBIT",
            "ioctl$UI_SET_PROPBIT",
            "ioctl$UI_SET_PHYS",
            "ioctl$UI_BEGIN_FF_UPLOAD",
            "ioctl$UI_BEGIN_FF_ERASE",
            "ioctl$UI_END_FF_UPLOAD",
            "ioctl$UI_END_FF_ERASE",
            "ioctl$UI_GET_SYSNAME",
            "ioctl$UI_ABS_SETUP",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
    "dev_userio": {
        "enable_syscalls": [
            "openat$userio",
            "write$USERIO_CMD_REGISTER",
            "write$USERIO_CMD_SET_PORT_TYPE",
            "write$USERIO_CMD_SEND_INTERRUPT",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
    "dev_video4linux": {
        "enable_syscalls": [
            "syz_open_dev$video",
            "syz_open_dev$video4linux",
            "openat$vimc0",
            "openat$vimc1",
            "openat$vimc2",
            "syz_open_dev$vivid",
            "openat$vim2m",
            "openat$vicodec0",
            "openat$vicodec1",
            "syz_open_dev$swradio",
            "syz_open_dev$radio",
            "syz_open_dev$vbi",
            "syz_open_dev$cec",
            "ioctl$VIDIOC_QUERYCAP",
            "ioctl$VIDIOC_RESERVED",
            "ioctl$VIDIOC_ENUM_FMT",
            "ioctl$VIDIOC_G_FMT",
            "ioctl$VIDIOC_S_FMT",
            "ioctl$VIDIOC_REQBUFS",
            "ioctl$VIDIOC_QUERYBUF",
            "ioctl$VIDIOC_G_FBUF",
            "ioctl$VIDIOC_S_FBUF",
            "ioctl$VIDIOC_OVERLAY",
            "ioctl$VIDIOC_QBUF",
            "ioctl$VIDIOC_EXPBUF",
            "ioctl$VIDIOC_DQBUF",
            "ioctl$VIDIOC_STREAMON",
            "ioctl$VIDIOC_STREAMOFF",
            "ioctl$VIDIOC_G_PARM",
            "ioctl$VIDIOC_S_PARM",
            "ioctl$VIDIOC_G_STD",
            "ioctl$VIDIOC_S_STD",
            "ioctl$VIDIOC_ENUMSTD",
            "ioctl$VIDIOC_ENUMINPUT",
            "ioctl$VIDIOC_G_CTRL",
            "ioctl$VIDIOC_S_CTRL",
            "ioctl$VIDIOC_G_TUNER",
            "ioctl$VIDIOC_S_TUNER",
            "ioctl$VIDIOC_G_AUDIO",
            "ioctl$VIDIOC_S_AUDIO",
            "ioctl$VIDIOC_QUERYCTRL",
            "ioctl$VIDIOC_QUERYMENU",
            "ioctl$VIDIOC_G_INPUT",
            "ioctl$VIDIOC_S_INPUT",
            "ioctl$VIDIOC_G_EDID",
            "ioctl$VIDIOC_S_EDID",
            "ioctl$VIDIOC_G_OUTPUT",
            "ioctl$VIDIOC_S_OUTPUT",
            "ioctl$VIDIOC_ENUMOUTPUT",
            "ioctl$VIDIOC_G_AUDOUT",
            "ioctl$VIDIOC_S_AUDOUT",
            "ioctl$VIDIOC_G_MODULATOR",
            "ioctl$VIDIOC_S_MODULATOR",
            "ioctl$VIDIOC_G_FREQUENCY",
            "ioctl$VIDIOC_S_FREQUENCY",
            "ioctl$VIDIOC_CROPCAP",
            "ioctl$VIDIOC_G_CROP",
            "ioctl$VIDIOC_S_CROP",
            "ioctl$VIDIOC_G_JPEGCOMP",
            "ioctl$VIDIOC_S_JPEGCOMP",
            "ioctl$VIDIOC_QUERYSTD",
            "ioctl$VIDIOC_TRY_FMT",
            "ioctl$VIDIOC_ENUMAUDIO",
            "ioctl$VIDIOC_ENUMAUDOUT",
            "ioctl$VIDIOC_G_PRIORITY",
            "ioctl$VIDIOC_S_PRIORITY",
            "ioctl$VIDIOC_G_SLICED_VBI_CAP",
            "ioctl$VIDIOC_LOG_STATUS",
            "ioctl$VIDIOC_G_EXT_CTRLS",
            "ioctl$VIDIOC_S_EXT_CTRLS",
            "ioctl$VIDIOC_TRY_EXT_CTRLS",
            "ioctl$VIDIOC_ENUM_FRAMESIZES",
            "ioctl$VIDIOC_ENUM_FRAMEINTERVALS",
            "ioctl$VIDIOC_G_ENC_INDEX",
            "ioctl$VIDIOC_ENCODER_CMD",
            "ioctl$VIDIOC_TRY_ENCODER_CMD",
            "ioctl$VIDIOC_DBG_S_REGISTER",
            "ioctl$VIDIOC_DBG_G_REGISTER",
            "ioctl$VIDIOC_S_HW_FREQ_SEEK",
            "ioctl$VIDIOC_S_DV_TIMINGS",
            "ioctl$VIDIOC_G_DV_TIMINGS",
            "ioctl$VIDIOC_DQEVENT",
            "ioctl$VIDIOC_SUBSCRIBE_EVENT",
            "ioctl$VIDIOC_UNSUBSCRIBE_EVENT",
            "ioctl$VIDIOC_CREATE_BUFS",
            "ioctl$VIDIOC_PREPARE_BUF",
            "ioctl$VIDIOC_G_SELECTION",
            "ioctl$VIDIOC_S_SELECTION",
            "ioctl$VIDIOC_DECODER_CMD",
            "ioctl$VIDIOC_TRY_DECODER_CMD",
            "ioctl$VIDIOC_ENUM_DV_TIMINGS",
            "ioctl$VIDIOC_QUERY_DV_TIMINGS",
            "ioctl$VIDIOC_DV_TIMINGS_CAP",
            "ioctl$VIDIOC_ENUM_FREQ_BANDS",
            "ioctl$VIDIOC_DBG_G_CHIP_INFO",
            "ioctl$VIDIOC_QUERY_EXT_CTRL",
            "ioctl$VIDIOC_SUBDEV_G_FMT",
            "ioctl$VIDIOC_SUBDEV_S_FMT",
            "ioctl$VIDIOC_SUBDEV_G_FRAME_INTERVAL",
            "ioctl$VIDIOC_SUBDEV_S_FRAME_INTERVAL",
            "ioctl$VIDIOC_SUBDEV_ENUM_MBUS_CODE",
            "ioctl$VIDIOC_SUBDEV_ENUM_FRAME_SIZE",
            "ioctl$VIDIOC_SUBDEV_ENUM_FRAME_INTERVAL",
            "ioctl$VIDIOC_SUBDEV_G_CROP",
            "ioctl$VIDIOC_SUBDEV_S_CROP",
            "ioctl$VIDIOC_SUBDEV_G_SELECTION",
            "ioctl$VIDIOC_SUBDEV_S_SELECTION",
            "ioctl$VIDIOC_SUBDEV_G_EDID",
            "ioctl$VIDIOC_SUBDEV_S_EDID",
            "ioctl$VIDIOC_SUBDEV_S_DV_TIMINGS",
            "ioctl$VIDIOC_SUBDEV_G_DV_TIMINGS",
            "ioctl$VIDIOC_SUBDEV_ENUM_DV_TIMINGS",
            "ioctl$VIDIOC_SUBDEV_QUERY_DV_TIMINGS",
            "ioctl$VIDIOC_SUBDEV_DV_TIMINGS_CAP",
        ],
        "file_bc": os.path.join(path_linux_bc, "arch/x86/kvm/built-in.bc"),
        "path_s": [
            os.path.join(path_linux_bc, "virt/kvm"),
            os.path.join(path_linux_bc, "arch/x86/kvm"),
        ],
        "file_taint": os.path.join(path_taint, "taint_info_kvm_dev_ioctl_serialize"),
    },
}
