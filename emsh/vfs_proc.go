package emsh

func buildProc(root *vfsNode) {
	proc := root.dir("proc")
	proc.file("version", "-r--r--r--", "Linux version 2.6.32.41 (root@buildhost) (gcc version 4.3.3) #1 SMP Fri Jan 5 12:00:00 UTC 2024\n")
	proc.file("cpuinfo", "-r--r--r--",
		"system type\t\t: BCM96318\nprocessor\t\t: 0\ncpu model\t\t: Broadcom BCM3384 V2.1\n"+
			"BogoMIPS\t\t: 210.00\nwait and kill\t\t: yes\nvector unit\t\t: VU0\n"+
			"secondary CPU\t\t: none\ncpu clock frequency\t: 210 MHz\nmemory frequency\t: 105 MHz\n"+
			"BCL2 controllers\t: 0\nISA DMA channels\t: 0\n")
	proc.file("meminfo", "-r--r--r--",
		"MemTotal:\t\t16384 kB\nMemFree:\t\t8192 kB\nBuffers:\t\t512 kB\nCached:\t\t2048 kB\n"+
			"SwapTotal:\t\t0 kB\nSwapFree:\t\t0 kB\nActive:\t\t3072 kB\nInactive:\t\t1024 kB\n"+
			"Dirty:\t\t0 kB\nWriteback:\t\t0 kB\nAnonPages:\t\t1536 kB\nMapped:\t\t512 kB\n"+
			"Slab:\t\t1024 kB\nSReclaimable:\t512 kB\nSUnreclaim:\t512 kB\nPageTables:\t128 kB\n"+
			"NFS_Unstable:\t0 kB\nBounce:\t\t0 kB\nCommitLimit:\t8192 kB\nCommitted_AS:\t2048 kB\n"+
			"VmallocTotal:\t1048576 kB\nVmallocUsed:\t1024 kB\nVmallocChunk:\t1047552 kB\n")
	proc.file("uptime", "-r--r--r--", "12345.67 23456.78\n")
	proc.file("loadavg", "-r--r--r--", "0.05 0.01 0.00 1/42 1234\n")
	proc.file("cmdline", "-r--r--r--", "brcmboard=0x1310,b00  console=ttyS0,115200 root=/dev/mtdblock2 rootfstype=squashfs\n")
	proc.file("mounts", "-r--r--r--",
		"rootfs / rootfs rw 0 0\n/dev/root / romfs ro 0 0\nproc /proc proc rw 0 0\n"+
			"sysfs /sys sysfs rw 0 0\ntmpfs /tmp tmpfs rw 0 0\ntmpfs /var/run tmpfs rw 0 0\n/dev/mtdblock3 /var ext2 rw 0 0\n")
	proc.file("mountinfo", "-r--r--r--",
		"36 35 0:23 / / rw relatime - romfs /dev/root ro\n37 0 0:24 / /proc rw - proc proc rw\n"+
			"38 0 0:25 / /sys rw - sysfs sysfs rw\n39 0 0:26 / /tmp rw - tmpfs tmpfs rw\n")
	proc.file("devices", "-r--r--r--",
		"block:\t7:0 /dev/root\nblock:\t31:0 /dev/mtdblock0\nblock:\t31:1 /dev/mtdblock1\n"+
			"block:\t31:2 /dev/mtdblock2\nblock:\t31:3 /dev/mtdblock3\ncharacter:\t5:0 /dev/tty\n"+
			"character:\t5:1 /dev/console\ncharacter:\t1:3 /dev/null\ncharacter:\t1:5 /dev/zero\n"+
			"character:\t1:7 /dev/full\ncharacter:\t1:8 /dev/random\ncharacter:\t1:9 /dev/urandom\n")
	proc.file("filesystems", "-r--r--r--",
		"Nodefs:\next2\next3\nvfat\nmsdos\nromfs\nsquashfs\nproc\nsysfs\ntmpfs\ndevtmpfs\ndebugfs\nsecurityfs\nrpc_pipefs\n")
	proc.file("interrupts", "-r--r--r--",
		"\t\t\tCPU0\n  0:          0          0 XT-PIC  timer\n  1:          0          0 XT-PIC  cascade\n"+
			"  2:          0          0 XT-PIC\n  3:          0          0 XT-PIC\n  4:          0          0 XT-PIC\n"+
			"  5:          0          0 XT-PIC  eth0\n  6:          0          0 XT-PIC  ttyS0\n  7:          0          0 XT-PIC\n")
	proc.file("iomem", "-r--r--r--",
		"00000000-0009f7ff : System RAM\n00100000-07ffffff : System RAM\n1fc00000-1fc00000 : bcm63xx_enet\n")
	proc.file("partitions", "-r--r--r--",
		"major minor  #blocks  name\n31:0    16384  mtdblock0\n31:1    25600  mtdblock1\n31:2   262144  mtdblock2\n31:3    65536  mtdblock3\n")
	proc.file("stat", "-r--r--r--",
		"cpu  12345 678 91011 1234567 890 0 1234 0 0 0\ncpu0 12345 678 91011 1234567 890 0 1234 0 0 0\n"+
			"intr 1234567 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n"+
			"ctxt 1234567\nbtime 1704350400\nprocesses 1234\nprocs_running 1\nprocs_blocked 0\n")
	proc.file("swaps", "-r--r--r--", "Filename\t\t\t\tType\t\tSize\tUsed\tPriority\n")
	proc.file("vmstat", "-r--r--r--",
		"nr_free_pages 2048\nnr_inactive_anon 512\nnr_active_anon 1024\nnr_inactive_file 512\n"+
			"nr_active_file 1536\nnr_dirty 0\nnr_writeback 0\nnr_slab_reclaimable 512\n"+
			"nr_slab_unreclaimable 512\nnr_page_table_pages 128\n")
	proc.file("diskstats", "-r--r--r--",
		" 31   0 mtdblock0 0 0 0 0 0 0 0 0 0 0 0 0 0\n 31   1 mtdblock1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n"+
			" 31   2 mtdblock2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n 31   3 mtdblock3 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n")
	proc.file("kmsg", "-rw-r--r--", "")
	proc.file("sysrq-trigger", "-rw-r--r--", "")
	proc.file("timer_list", "-r--r--r--", "")
	proc.file("slabinfo", "-r--r--r--", "")
	proc.file("buddyinfo", "-r--r--r--", "")
	proc.file("page", "-r--r--r--", "")
	proc.file("pagetypeinfo", "-r--r--r--", "")
	proc.file("zoneinfo", "-r--r--r--", "")
	proc.file("execdomains", "-r--r--r--", "0\tdefault\t\t\n")
	proc.file("cgroups", "-r--r--r--",
		"12:cpuset /system\n11:blkio /system\n10:freezer /system\n9:devices /system\n8:memory /system\n"+
			"7:hugetlb /system\n6:perf_event /system\n5:net_cls /system\n4:net_prio /system\n3:pid /system\n"+
			"2:cpu,cpuacct /system\n1:blkio /system\n")
	proc.file("cgroup", "-r--r--r--",
		"12:cpuset /\n11:blkio /\n10:freezer /\n9:devices /\n8:memory /\n7:hugetlb /\n6:perf_event /\n"+
			"5:net_cls /\n4:net_prio /\n3:pid /\n2:cpu,cpuacct /\n1:blkio /\n")
	proc.file("keys", "-rw-r--r--", "")
	proc.file("key-users", "-rw-r--r--", "")
	proc.file("sched_debug", "-r--r--r--", "")
	proc.file("mdstat", "-r--r--r--", "No arrays configured\n")
	proc.file("pressure", "-r--r--r--", "")
	proc.dir("1")
	proc.dir("self")
	procnet := proc.dir("net")
	procnet.file("dev", "-r--r--r--",
		"Inter-|   Receive                                                |  Transmit\n"+
			" face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed\n"+
			"    lo: 1234567    8901    0    0    0     0          0         0  1234567    8901    0    0    0     0       0          0\n"+
			"  eth0: 98765432  654321    0    0    0     0          0         0  12345678   98765    0    0    0     0       0          0\n")
	procnet.file("tcp", "-r--r--r--",
		"  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n"+
			"   0: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0\n"+
			"   1: 00000000:0017 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12346 1 0000000000000000 100 0 0 10 0\n"+
			"   2: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000   333        0 12347 1 0000000000000000 100 0 0 10 0\n")
	procnet.file("tcp6", "-r--r--r--", "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n")
	procnet.file("udp", "-r--r--r--",
		"  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops\n"+
			"   0: 00000000:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 12348 2 0000000000000000 0\n")
	procnet.file("udp6", "-r--r--r--", "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops\n")
	procnet.file("arp", "-r--r--r--",
		"IP address              HW type     Flags       HW address            Mask     Device\n"+
			"192.168.1.1           0x1         0x2         aa:bb:cc:dd:ee:ff     *        eth0\n"+
			"192.168.1.50          0x1         0x0         11:22:33:44:55:66     *        eth0\n")
	procnet.file("if_inet6", "-r--r--r--", "6: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 state UNKNOWN qlen 1\n    inet6 ::1/128 scope host \n       valid_lft forever preferred_lft forever\n")
	procnet.file("route", "-r--r--r--",
		"Iface\tDestination\tGateway \tFlags\tRefCnt\tUse\tMetric\tMask\t\tMTU\tWindow\tIRTT\n"+
			"eth0\t00000000\t00000001\t0003\t0\t0\t0\t00000000\t0\t0\t0\neth0\t0001A8C0\t00000000\t0001\t0\t0\t0\t00FFFFFF\t0\t0\t0\n")
	procnet.emptyFiles("-r--r--r--", "snmp", "snmp6", "raw", "igmp", "protocols", "fib_trie", "fib_triestat", "psfilter", "netfilter", "netlink", "packet", "unix", "raw6", "igmp6", "if_status", "dev_mcast", "netstat", "snmp_agent")
	for _, n := range []string{"bus", "driver", "fs", "acpi", "scsi", "ide", "fb", "pci", "usb", "iso_8859_1", "tty", "sysvipc", "cpu", "mem", "irq", "dma", "sys"} {
		proc.dir(n)
	}
}

func buildSys(root *vfsNode) {
	sys := root.dir("sys")
	for _, n := range []string{"class", "devices", "firmware", "kernel", "module", "power", "bus", "hypervisor", "block", "net", "dev", "fs"} {
		sys.dir(n)
	}
}
