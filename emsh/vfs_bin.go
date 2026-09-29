package emsh

func buildBin(root *vfsNode) {
	bin := root.dir("bin")
	bin.fileSize("busybox", "-rwxr-xr-x", fakeELF("busybox"), 1046317)
	bin.addLinks("busybox",
		"sh", "ash", "ls", "cat", "echo", "grep", "egrep", "fgrep", "sed", "awk",
		"vi", "more", "less", "head", "tail", "wc", "sort", "uniq", "cut", "tr",
		"expr", "date", "cal", "df", "du", "find", "xargs", "touch", "mkdir",
		"rmdir", "rm", "cp", "mv", "ln", "chmod", "chown", "chgrp", "id",
		"whoami", "uname", "hostname", "ps", "top", "kill", "killall", "free",
		"dd", "cmp", "diff", "md5sum", "sha1sum", "sha256sum", "sha512sum",
		"base64", "xxd", "od", "hexdump", "strings", "stty", "clear", "dmesg",
		"sleep", "test", "true", "false", "printenv", "tar", "gzip", "gunzip",
		"bzip2", "bunzip2", "zip", "unzip", "cpio", "ping", "nslookup", "wget",
		"telnet", "tftp", "ftp", "netstat", "ifconfig", "route", "ip", "reboot",
		"halt", "poweroff", "su", "sudo", "passwd", "useradd", "userdel",
		"groupadd", "groupdel", "crontab", "mktemp", "mkfifo", "mknod", "mkswap",
		"fsck", "mount", "umount", "swapon", "swapoff", "login", "getty", "watch",
		"which", "whereis", "who", "w", "uptime", "env", "readlink", "realpath",
		"stat", "file", "yes", "seq", "shuf", "nl", "paste", "column", "comm",
		"split", "sum", "tac", "tee", "tsort", "uudecode", "uuencode", "arch",
		"basename", "chroot", "cksum", "cryptpw", "dc", "dir", "dircolors",
		"dnsdomainname", "dos2unix", "dumpkmap", "ed", "ee", "expand", "factor",
		"fallocate", "fatattr", "fdflush", "fdformat", "fdisk", "findfs",
		"finger", "fixcon", "flashcp", "flash_erase", "flash_lock", "flash_unlock",
		"flashinfo", "fmt", "fold", "fsync", "fuser", "getopt", "groups",
		"hostid", "hwclock", "i2cdetect", "i2cdump", "i2cget", "i2cset",
		"i2ctransfer", "ifdown", "ifup", "inetd", "init", "insmod", "install",
		"ipaddr", "ipcalc", "iplink", "iproute", "iprule", "iptunnel", "kbd_mode",
		"loadfont", "logger", "logname", "look", "lsmod", "lspci", "lsusb",
		"lzma", "lzmad", "makedevs", "mountpoint", "namei", "nanddump", "nandwrite",
		"nc", "nice", "nmeter", "nohup", "nproc", "ntpd", "openvt", "partprobe",
		"ping6", "pipe_progress", "pkill", "pmap", "popmaildir", "raidautorun",
		"rdate", "renice", "reset", "resize", "resume", "rev", "rmmod",
		"run-parts", "runsv", "runsvdir", "script", "setarch", "setconsole",
		"setfont", "setkeycodes", "setlogcons", "setsid", "showkey", "shred",
		"sl", "smemcap", "start-stop-daemon", "stdbuf", "switch_root", "sync",
		"sysctl", "taskset", "tcpsvd", "time", "timeout", "tracert", "truncate",
		"ts", "tty", "ttysize", "tunctl", "unexpand", "unix2dos", "unlink",
		"unlzma", "unshare", "usleep", "vlock", "zcat", "zcmp", "zdiff",
		"zforce", "zgrep", "zless", "zmore", "znew")
}

func buildSbin(root *vfsNode) {
	sbin := root.dir("sbin")
	sbin.addLinks("../bin/busybox",
		"init", "halt", "poweroff", "reboot", "ifconfig", "route", "netstat",
		"ping", "nslookup", "wget", "telnetd", "ftpd", "tftpd", "httpd",
		"syslogd", "klogd", "getty", "inetd", "crond", "crontab", "fdisk",
		"mkfs", "mkswap", "fsck", "swapon", "swapoff", "mount", "umount",
		"losetup", "mdev", "modprobe", "insmod", "rmmod", "lsmod", "depmod",
		"brctl", "vconfig", "dnsmasq", "udhcpd", "udhcpc", "dropbear",
		"dropbearkey", "dropbearconvert", "scp", "sftp", "sshd", "busybox",
		"ifup", "ifdown", "ifstatus", "ifrename", "iwconfig", "iwlist",
		"iwpriv", "iwscan", "wpa_supplicant", "hostapd", "minissdpd",
		"mini_httpd", "thttpd", "lighttpd", "apache2", "nginx", "named",
		"bind", "dhcpd", "dhcrelay", "ntpd", "chronyd", "openvpn", "pppd",
		"rp-pppoe", "quagga", "zebra", "bgpd", "ospfd", "ripd", "strongswan",
		"charon", "radiusd", "squid", "memcached", "redis-server", "mysql",
		"mysqld", "postgres", "mongod", "dockerd", "containerd", "runc",
		"podman", "lxc", "lxd", "syslog-ng", "rsyslogd", "cron", "atd",
		"anacron", "bluetoothd", "udisksd", "avahi-daemon", "cupsd", "smbd",
		"nmbd", "winbind", "samba", "nfsd", "mount.nfs", "rpc.mountd",
		"rpc.statd", "rpc.idmapd", "portmap", "rpcbind", "xinetd", "tcpd",
		"tcpdump", "wireshark", "tshark", "ethertypes", "edit", "nano",
		"emacs", "joe", "vim", "ex", "ed", "sedit", "pico", "screen", "tmux",
		"expect", "curl", "ssh", "rsync", "rsh", "rlogin", "rcp", "ncat",
		"socat")
}

func buildUsr(root *vfsNode) {
	us := root.dir("usr")
	usrbin := us.dir("bin")
	usrbin.addLinks("../../bin/busybox",
		"python", "python2", "perl", "ruby", "php", "node", "make", "gcc",
		"cc", "g++", "c++", "ar", "as", "ld", "nm", "objdump", "strip", "cpp",
		"bison", "flex", "autoconf", "automake", "pkg-config", "cmake", "git",
		"svn", "cvs", "hg", "mercurial", "curl", "ssh", "scp", "sftp", "rsync",
		"rsh", "rlogin", "rcp", "telnet", "ftp", "tftp", "nc", "ncat", "socat",
		"tcpdump", "wireshark", "tshark", "ethertypes", "edit", "nano", "emacs",
		"joe", "vim", "vi", "ex", "ed", "sedit", "pico", "screen", "tmux",
		"expect", "less", "more", "man", "info", "help", "apropos", "whatis",
		"whereis", "which", "locate", "updatedb", "find", "beep", "biff",
		"bunzip2", "bzip2", "bzip2recover", "cal", "chattr", "chroot", "cksum",
		"clear", "cmp", "comm", "cp", "cpio", "csplit", "cut", "date", "dd",
		"ddrescue", "df", "diff", "diff3", "dir", "dircolors", "dmesg", "du",
		"echo", "ed", "env", "expand", "expr", "factor", "false", "fdisk",
		"fgrep", "file", "find", "fmt", "fold", "free", "fsck", "fsync",
		"fuser", "groups", "gunzip", "gzip", "halt", "head", "hexdump",
		"hostid", "hostname", "id", "ifconfig", "install", "ip", "kill",
		"killall", "ln", "logger", "logname", "look", "ls", "lsattr", "lsmod",
		"lspci", "lsusb", "lzma", "md5sum", "mkdir", "mkfifo", "mknod",
		"mktemp", "more", "mount", "mv", "nice", "nl", "nohup", "nproc",
		"nslookup", "od", "paste", "pathchk", "ping", "ping6", "pr",
		"printenv", "printf", "ps", "pwd", "readlink", "realpath", "reboot",
		"renice", "reset", "rev", "rm", "rmdir", "route", "rpcinfo", "rpcping",
		"rpcwait", "seq", "setarch", "setpriv", "setsid", "sha1sum",
		"sha224sum", "sha256sum", "sha384sum", "sha512sum", "shred", "shuf",
		"sleep", "sort", "split", "stat", "stdbuf", "stty", "sum", "sync",
		"tac", "tail", "tar", "tee", "telnet", "test", "time", "timeout",
		"touch", "tr", "true", "truncate", "ts", "tty", "uname", "unexpand",
		"uniq", "units", "unlink", "unshare", "uptime", "users", "uuencode",
		"uudecode", "vdir", "vi", "watch", "wc", "wget", "which", "who",
		"whoami", "whois", "xargs", "xxd", "yes", "zcat", "zcmp", "zdiff",
		"zforce", "zgrep", "znew")
	usrsbin := us.dir("sbin")
	usrsbin.addLinks("../../bin/busybox",
		"sshd", "dropbear", "httpd", "apache2", "nginx", "lighttpd", "thttpd",
		"mini_httpd", "ftpd", "vsftpd", "proftpd", "pure-ftpd", "sftp-server",
		"tftpd", "inetd", "xinetd", "named", "bind", "dnsmasq", "dhcpd",
		"dhcrelay", "ntpd", "chronyd", "openvpn", "vpnd", "pppd", "rp-pppoe",
		"quagga", "zebra", "bgpd", "ospfd", "ripd", "isisd", "bird",
		"strongswan", "swanctl", "charon", "radiusd", "freeradius", "squid",
		"varnish", "memcached", "redis-server", "mysql", "mysqld", "postgres",
		"postgresql", "mongod", "couchdb", "rabbitmq-server", "elasticsearch",
		"kibana", "grafana", "docker", "dockerd", "containerd", "runc",
		"podman", "lxc", "lxd", "syslogd", "syslog-ng", "rsyslogd", "cron",
		"atd", "anacron", "bluetoothd", "udisksd", "avahi-daemon", "cupsd",
		"smbd", "nmbd", "winbind", "samba", "nfsd", "mount.nfs", "rpc.mountd",
		"rpc.statd", "rpc.idmapd", "portmap", "rpcbind", "tcpdump", "wireshark",
		"tshark", "ethertypes", "edit", "nano", "emacs", "joe", "vim", "ex",
		"ed", "sedit", "pico", "screen", "tmux", "expect", "curl", "ssh",
		"scp", "sftp", "rsync", "rsh", "rlogin", "rcp", "nc", "ncat", "socat")
	usrlib := us.dir("lib")
	usrlib.emptyFiles("-rw-r--r--", "libstdc++.so.6", "libgcc_s.so.1", "libm.so.6", "libpthread.so.0")
	usrshare := us.dir("share")
	usrshare.dir("misc")
	udhcpc := usrshare.dir("udhcpc")
	udhcpc.file("default.script", "-rw-r--r--", "#!/bin/sh\n[ -n \"$1\" ] && echo \"udhcpc: $1\"\nexit 0\n")
	usrshare.dir("terminfo")
	usrshare.dir("man")
	usrshare.dir("doc")
	us.dir("local")
	us.dir("include")
}

func buildLib(root *vfsNode) {
	lib := root.dir("lib")
	lib.symlink("ld-uClibc.so.0", "ld-uClibc.so.0.9.33")
	lib.symlink("libc.so.0", "libc.so.0.9.33")
	lib.emptyFiles("-rw-r--r--",
		"ld-uClibc.so.0.9.33", "libc.so.0.9.33", "libpthread.so.0",
		"libdl.so.0", "libm.so.0", "libnsl.so.0", "librt.so.0", "libutil.so.0",
		"libgcc_s.so.1", "libstdc++.so.6")
	lib.dir("modules")
	lib.dir("firmware")
}
