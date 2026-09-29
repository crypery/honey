package emsh

func buildDev(root *vfsNode) {
	dev := root.dir("dev")
	dev.emptyFiles("-rw-rw-rw-",
		"null", "zero", "full", "random", "urandom", "tty", "console", "mem",
		"kmem", "port", "ptmx", "stdin", "stdout", "stderr", "ttyS0", "ttyS1",
		"ttyAMA0", "loop0", "loop1", "mtdblock0", "mtdblock1", "mtdblock2",
		"mtdblock3", "mtd0", "mtd1", "mtd2", "mtd3", "eth0", "wlan0", "br0")
	for _, n := range []string{"fd", "pts", "mtd", "shm", "input", "usb", "sound", "video", "net", "disk", "block", "char"} {
		dev.dir(n)
	}
}

func buildVar(root *vfsNode) {
	varr := root.dir("var")
	varlog := varr.dir("log")
	varlog.file("messages", "-rw-r--r--",
		"Sep 06 10:00:01 BCM96318 kernel: system boot complete\n"+
			"Sep 06 10:00:05 BCM96318 syslogd: syslogd started\n"+
			"Sep 06 10:00:10 BCM96318 dropbear[123]: dropbear 0.53.1 ready to start device\n")
	varlog.emptyFiles("-rw-r--r--", "wtmp", "lastlog", "btmp", "syslog", "kern.log", "daemon.log", "mail.log", "cron.log", "secure")
	varrun := varr.dir("run")
	varrun.emptyFiles("-rw-rw-r--", "utmp", "random_seed")
	varrun.dir("pid")
	varrun.dir("lock")
	varrun.dir("local")
	varspool := varr.dir("spool")
	varspool.dir("cron")
	varspool.dir("mail")
	varspool.dir("mqueue")
	varcache := varr.dir("cache")
	varcache.dir("man")
	varcache.dir("ldconfig")
	varcache.dir("fonts")
	varlib := varr.dir("lib")
	for _, n := range []string{"udhcpc", "pcmcia", "firmware", "modules", "misc", "dropbear", "dnsmasq", "pppd", "strongswan", "nfs", "samba", "smbd", "nmbd", "vsftpd", "tftpd", "inetd", "xinetd", "named", "dhcpd", "ntpd", "openvpn", "rp-pppoe", "quagga", "zebra", "bgpd", "radiusd", "syslogd", "rsyslogd", "cron", "atd", "avahi-daemon", "cupsd", "rpcbind", "portmap"} {
		varlib.dir(n)
	}
	varr.dir("tmp")
	varr.dir("lock")
	varr.dir("opt")
	varr.dir("backup")
	varwww := varr.dir("www")
	varwww.file("index.html", "-rw-r--r--", "<html><head><title>BCM96318</title></head><body><h1>BCM96318 Broadband Router</h1></body></html>\n")
	varwww.dir("html")
	varwww.dir("cgi-bin")
}

func buildHome(root *vfsNode) {
	home := root.dir("home")
	admin := home.dir("admin")
	admin.file(".bashrc", "-rw-r--r--", "# ~/.bashrc\nexport PS1='\\w \\$ '\n")
	admin.file(".profile", "-rw-r--r--", "# ~/.profile\nexport PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\n")
	admin.file(".bash_history", "-rw-r--r--", "ls\npwd\ncat /etc/passwd\nifconfig\n")
	admin.dir(".ssh")
}

func buildRootDir(root *vfsNode) {
	rootdir := root.dir("root")
	rootdir.file(".bashrc", "-rw-r--r--", "# ~/.bashrc\nexport PS1='\\w \\$ '\n")
	rootdir.file(".profile", "-rw-r--r--", "# ~/.profile\nexport PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\n")
	rootdir.file(".bash_history", "-rw-r--r--", "ls\npwd\ncat /etc/shadow\nifconfig\nnetstat -an\n")
	rootdir.file(".netrc", "-rw-------", "machine 192.168.1.10\nlogin admin\npassword admin\n")
	rootdir.file(".kshrc", "-rw-r--r--", "# ~/.kshrc\n")
	rootdir.file(".exrc", "-rw-r--r--", "# ~/.exrc\n")
	rootdir.file(".login", "-rw-r--r--", "# ~/.login\n")
	rootdir.file(".logout", "-rw-r--r--", "# ~/.logout\n")
	rootdir.file(".cshrc", "-rw-r--r--", "# ~/.cshrc\n")
	rootdir.file(".tcshrc", "-rw-r--r--", "# ~/.tcshrc\n")
	rootdir.file(".zshrc", "-rw-r--r--", "# ~/.zshrc\n")
	rootdir.file(".vimrc", "-rw-r--r--", "# ~/.vimrc\nset number\n")
	rootdir.file(".emacs", "-rw-r--r--", "# ~/.emacs\n")
	rootdir.file(".inputrc", "-rw-r--r--", "# ~/.inputrc\n")
	rootdir.file(".gitconfig", "-rw-r--r--", "[user]\n\tname = admin\n\temail = admin@BCM96318\n")
	rootdir.file(".gitignore", "-rw-r--r--", "*.o\n*.a\n")
	rootdir.file(".htaccess", "-rw-r--r--", "# ~/.htaccess\n")
	rootdir.file(".htpasswd", "-rw-r--r--", "admin:$1$Ab3Cd4Ef$Gh5Ij6Kl7Mn8Op9Qr0StUv\n")
	rootdir.file(".my.cnf", "-rw-r--r--", "[client]\nuser = root\npassword = admin\n")
	rootdir.file(".pgpass", "-rw-------", "localhost:5432:*:root:admin\n")
	rootdir.file(".mysql_history", "-rw-r--r--", "show databases;\n")
	rootdir.file(".psql_history", "-rw-r--r--", "\\l\n")
	rootdir.file(".rediscli_history", "-rw-r--r--", "keys *\n")
	rootdir.file(".lesshst", "-rw-r--r--", "")
	rootdir.dir(".ssh")
	rootdir.dir(".gnupg")
	rootdir.dir(".config")
	rootdir.dir(".local")
	rootdir.dir(".cache")
}
