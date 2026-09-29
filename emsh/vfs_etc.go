package emsh

func buildEtc(root *vfsNode) {
	etc := root.dir("etc")
	etc.file("passwd", "-rw-r--r--",
		"root:x:0:0:root:/root:/bin/sh\n"+
			"daemon:x:1:1:daemon:/usr/sbin:/bin/sh\n"+
			"bin:x:2:2:bin:/bin:/bin/false\n"+
			"sys:x:3:3:sys:/dev:/bin/false\n"+
			"sync:x:4:65534:sync:/bin:/bin/sync\n"+
			"games:x:5:60:games:/usr/games:/bin/false\n"+
			"man:x:6:12:man:/var/cache/man:/bin/false\n"+
			"lp:x:7:7:lp:/var/spool/lpd:/bin/false\n"+
			"mail:x:8:8:mail:/var/mail:/bin/false\n"+
			"news:x:9:9:news:/var/spool/news:/bin/false\n"+
			"uucp:x:10:10:uucp:/var/spool/uucp:/bin/false\n"+
			"proxy:x:13:13:proxy:/bin:/bin/false\n"+
			"www-data:x:33:33:www-data:/var/www:/bin/false\n"+
			"backup:x:34:34:backup:/var/backups:/bin/false\n"+
			"list:x:38:38:Mailing List Manager:/var/list:/bin/false\n"+
			"irc:x:39:39:ircd:/var/run/ircd:/bin/false\n"+
			"gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/bin/false\n"+
			"nobody:x:65534:65534:nobody:/nonexistent:/bin/false\n"+
			"admin:x:1000:1000:admin:/home/admin:/bin/sh\n")
	etc.file("shadow", "-rw-r-----",
		"root:$1$Xq7Kp2mZ$9vN3rT8wY1uI5oP0aS4dF6gH:19000:0:99999:7:::\n"+
			"daemon:*:19000:0:99999:7:::\n"+
			"bin:*:19000:0:99999:7:::\n"+
			"sys:*:19000:0:99999:7:::\n"+
			"sync:*:19000:0:99999:7:::\n"+
			"games:*:19000:0:99999:7:::\n"+
			"man:*:19000:0:99999:7:::\n"+
			"lp:*:19000:0:99999:7:::\n"+
			"mail:*:19000:0:99999:7:::\n"+
			"news:*:19000:0:99999:7:::\n"+
			"uucp:*:19000:0:99999:7:::\n"+
			"proxy:*:19000:0:99999:7:::\n"+
			"www-data:*:19000:0:99999:7:::\n"+
			"backup:*:19000:0:99999:7:::\n"+
			"list:*:19000:0:99999:7:::\n"+
			"irc:*:19000:0:99999:7:::\n"+
			"gnats:*:19000:0:99999:7:::\n"+
			"nobody:*:19000:0:99999:7:::\n"+
			"admin:$1$Ab3Cd4Ef$Gh5Ij6Kl7Mn8Op9Qr0StUv:19000:0:99999:7:::\n")
	etc.file("group", "-rw-r--r--",
		"root:x:0:\ndaemon:x:1:\nbin:x:2:\nsys:x:3:\nadm:x:4:\ntty:x:5:\ndisk:x:6:\n"+
			"lp:x:7:\nmail:x:8:\nnews:x:9:\nuucp:x:10:\nman:x:12:\nproxy:x:13:\nkmem:x:15:\n"+
			"dialout:x:20:\nfax:x:21:\nvoice:x:22:\ncdrom:x:24:\nfloppy:x:25:\ntape:x:26:\n"+
			"sudo:x:27:\naudio:x:29:\ndip:x:30:\nwww-data:x:33:\nbackup:x:34:\noperator:x:37:\n"+
			"list:x:38:\nirc:x:39:\nsrc:x:40:\ngnats:x:41:\nshadow:x:42:\nutmp:x:43:\n"+
			"verbose:x:44:\nkiosk:x:45:\nplugdev:x:46:\nstaff:x:50:\ngames:x:60:\nusers:x:100:\n"+
			"nogroup:x:65534:\nadmin:x:1000:\n")
	etc.file("hostname", "-rw-r--r--", "BCM96318\n")
	etc.file("hosts", "-rw-r--r--",
		"127.0.0.1       localhost\n127.0.1.1       BCM96318\n\n"+
			"# The following lines are desirable for IPv6 capable hosts\n"+
			"::1             ip6-localhost ip6-loopback\nfe00::0         ip6-localnet\n"+
			"ff00::0         ip6-mcastprefix\nff02::1         ip6-allnodes\n"+
			"ff02::2         ip6-allrouters\n192.168.1.10    BCM96318\n")
	etc.file("resolv.conf", "-rw-r--r--", "# dyn:\nnameserver 4.2.2.1\nnameserver 4.2.2.2\n")
	etc.file("profile", "-rw-r--r--",
		"# /etc/profile\n#\n# This file is used by the shell /bin/sh\n#\n"+
			"export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\n"+
			"export PS1='\\w \\$ '\nexport PS2='> '\nexport PS3='# '\nexport PS4='+ '\n"+
			"export PWD=/\nexport HOME=/\nexport HOSTNAME=`/bin/hostname`\n"+
			"export HOSTTYPE=mips\nexport OS=Linux\numask 022\n")
	etc.file("inittab", "-rw-r--r--",
		"# /etc/inittab init(8) configuration for BusyBox\n#\n"+
			"# Format for each entry: <id>:<runlevels>:<action>:<process>\n#\n"+
			"::sysinit:/etc/init.d/rcS\n::askfirst:-/bin/sh\ntty2::askfirst:-/bin/sh\n"+
			"tty3::askfirst:-/bin/sh\ntty4::askfirst:-/bin/sh\n"+
			"tty4::respawn:/sbin/getty 38400 tty5\ntty5::respawn:/sbin/getty 38400 tty6\n"+
			"::restart:/sbin/init\n::ctrlaltdel:/sbin/reboot\n::shutdown:/bin/umount -a -r\n"+
			"::shutdown:/sbin/swapoff -a\n")
	etc.file("issue", "-rw-r--r--", "BCM96318 Broadband Router login: \n")
	etc.file("motd", "-rw-r--r--", "Welcome to the BCM96318 Broadband Router\nLast login: Sat Sep  6 10:00:01 on ttyS0\n")
	etc.file("services", "-rw-r--r--",
		"tcpmux\t\t1/tcp\necho\t\t7/tcp\necho\t\t7/udp\ndiscard\t\t9/tcp\n"+
			"discard\t\t9/udp\nsystat\t\t11/tcp\nsystat\t\t11/udp\ndaytime\t\t13/tcp\n"+
			"daytime\t\t13/udp\nftp\t\t21/tcp\nftp-data\t\t20/tcp\ntelnet\t\t23/tcp\n"+
			"smtp\t\t25/tcp\ntime\t\t37/tcp\ntime\t\t37/udp\ndomain\t\t53/tcp\n"+
			"domain\t\t53/udp\nshell\t\t51/tcp\nshell\t\t51/udp\nlogin\t\t52/tcp\n"+
			"ntp\t\t123/udp\nssh\t\t22/tcp\nhttp\t\t80/tcp\nhttps\t\t443/tcp\n")
	etc.file("protocols", "-rw-r--r--",
		"ip\t\t0\tIP\tinternet protocol\nhop\t\t1\tHOP\tInternet level relay\n"+
			"icmp\t\t1\tICMP\tinternet control message\nigmp\t\t2\tIGMP\tinternet gateway message\n"+
			"ggp\t\t3\tGGP\tgeneric gateway-support\nipinip\t\t4\tIP-IN-IP\tip transit network\n"+
			"st\t\t5\tST\tsecurity support\ntcp\t\t6\tTCP\ttransmission control\n"+
			"cbt\t\t7\tCBT\t\negp\t\t8\tEGP\texterior gateway protocol\n"+
			"igg\t\t9\tIGG\tinternet group management\npup\t\t12\tPUP\tPARC universal packet protocol\n"+
			"udp\t\t17\tUDP\tuser datagram protocol\nidpr\t\t22\tIDPR\tinternet date packet format\n"+
			"idpr-cmtp\t23\tIDPR-CMTP\tCMTP\nvrrp\t\t112\tVRRP\tVirtual Router Redundancy Protocol\n")
	etc.file("shells", "-rw-r--r--", "/bin/sh\n/bin/ash\n/bin/bash\n")
	etc.file("nsswitch.conf", "-rw-r--r--",
		"passwd:\t\tfiles\ngroup:\t\tfiles\nshadow:\t\tfiles\nhosts:\t\tfiles dns\n"+
			"networks:\t\tfiles\nservices:\t\tfiles\nprotocols:\t\tfiles\nethers:\t\tfiles\nrpc:\t\tfiles\n")
	etc.file("inputrc", "-rw-r--r--", "# /etc/inputrc\n")
	etc.file("fstab", "-rw-r--r--",
		"# device\tmount-point\ttype\toptions\t\tdump\tpass\n"+
			"/dev/root\t/\t\tromfs\tro\t\t\t0\t0\nproc\t\t/proc\t\tproc\t\trw\t\t\t0\t0\n"+
			"sysfs\t\t/sys\t\tsysfs\t\trw\t\t\t0\t0\ntmpfs\t\t/tmp\t\ttmpfs\t\trw\t\t\t0\t0\n"+
			"tmpfs\t\t/var/run\ttmpfs\t\trw\t\t\t0\t0\n/dev/mtdblock3\t/var\text2\t\trw\t\t\t0\t0\n")
	etc.file("mtab", "-rw-r--r--",
		"/dev/root\t/\tromfs\tro\t0\t0\nproc\t/proc\tproc\trw\t0\t0\n"+
			"sysfs\t/sys\tsysfs\trw\t0\t0\ntmpfs\t/tmp\ttmpfs\trw\t0\t0\n"+
			"tmpfs\t/var/run\ttmpfs\trw\t0\t0\n/dev/mtdblock3\t/var\text2\trw\t0\t0\n")
	etc.file("adjtime", "-rw-r--r--", "3 0 0\nF\n")
	etc.file("localtime", "-rw-r--r--", "")
	etc.file("magic", "-rw-r--r--", "")
	etc.file("magic.mime", "-rw-r--r--", "")
	etc.file("mailcap", "-rw-r--r--", "")
	etc.file("mime.types", "-rw-r--r--",
		"text/html\thtml htm\ntext/plain\ttxt\ntext/css\tcss\n"+
			"application/javascript\tjs\nimage/jpeg\tjpg jpeg\nimage/png\tpng\nimage/gif\tgif\n")
	etc.file("mailrc", "-rw-r--r--", "")
	etc.file("manpath", "-rw-r--r--", "")
	etc.file("networks", "-rw-r--r--", "loopback\t\t127\ndefault\t\t0\nbroadcast\t\t255\nlocalnet\t\t192.168.1\n")
	etc.file("passwd-", "-rw-r--r--", "")
	etc.file("securetty", "-rw-r--r--", "console\ntty1\ntty2\ntty3\ntty4\n")
	etc.file("selinux", "-rw-r--r--", "")
	etc.file("sysctl.conf", "-rw-r--r--",
		"# /etc/sysctl.conf\nkernel.domainname = (none)\nnet.ipv4.ip_forward = 1\nnet.ipv4.conf.all.rp_filter = 1\n")
	etc.file("syslog.conf", "-rw-r--r--",
		"# /etc/syslog.conf\n*.info;auth.notice\t/dev/console\nauthpriv\t\t/private\n*.warn\t\t*;\n*.err\t\t*;\n")
	etc.file("termcap", "-rw-r--r--", "")
	etc.file("timezones", "-rw-r--r--", "")
	etc.file("ttys", "-rw-r--r--", "")
	etc.file("xdg", "-rw-r--r--", "")

	initd := etc.dir("init.d")
	initd.file("rcS", "-rwxr-xr-x",
		"#!/bin/sh\n#\n# System initialization\n#\numask 022\n"+
			"/bin/mount -t proc proc /proc\n/bin/mount -t sysfs sysfs /sys\n"+
			"/bin/mount -t tmpfs tmpfs /tmp\n/bin/mount -t tmpfs tmpfs /var/run\n"+
			"hostname BCM96318\n/bin/echo \"Initializing system...\"\n"+
			"# start network\nifconfig lo 127.0.0.1 up\n"+
			"ifconfig eth0 192.168.1.10 netmask 255.255.255.0 up\nroute add default gw 192.168.1.1\n"+
			"# start services\n/usr/sbin/dropbear -p 22\n/usr/sbin/httpd -p 80\n/usr/sbin/inetd\n"+
			"/bin/echo \"System startup complete\"\n")
	initd.file("rc.local", "-rwxr-xr-x", "#!/bin/sh\n#\n# rc.local - executed at the end of each multiuser runlevel\n#\nexit 0\n")
	initd.file("S01logging", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -S -x /sbin/syslogd\n")
	initd.file("S02network", "-rwxr-xr-x", "#!/bin/sh\nifconfig lo 127.0.0.1 up\n")
	initd.file("S03dropbear", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -S -x /usr/sbin/dropbear\n")
	initd.file("S04httpd", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -S -x /usr/sbin/httpd\n")
	initd.file("S05inetd", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -S -x /usr/sbin/inetd\n")
	initd.file("S06crond", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -S -x /usr/sbin/crond\n")
	initd.file("K01crond", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -K -x /usr/sbin/crond\n")
	initd.file("K02inetd", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -K -x /usr/sbin/inetd\n")
	initd.file("K03httpd", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -K -x /usr/sbin/httpd\n")
	initd.file("K04dropbear", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -K -x /usr/sbin/dropbear\n")
	initd.file("K05network", "-rwxr-xr-x", "#!/bin/sh\nifconfig lo down\n")
	initd.file("K06logging", "-rwxr-xr-x", "#!/bin/sh\nstart-stop-daemon -K -x /sbin/syslogd\n")

	network := etc.dir("network")
	network.file("ifconfig", "-rw-r--r--",
		"auto lo\niface lo inet loopback\n\nauto eth0\niface eth0 inet static\n    address 192.168.1.10\n    netmask 255.255.255.0\n    gateway 192.168.1.1\n")
	network.file("interfaces", "-rw-r--r--", "auto lo\niface lo inet loopback\n")
	network.file("resolv.conf", "-rw-r--r--", "nameserver 4.2.2.1\n")
	network.file("networks", "-rw-r--r--", "loopback 127\n")

	dropbear := etc.dir("dropbear")
	dropbear.file("id_rsa", "-rw-r--r--",
		"-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA7v5Xc0m8lQ2pR4tW6yZ8aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL7mN9oP1qR3sT5uV7wX9yZ1aB3cD5eF7gH9iJ1kL3mN5oP7qR9sT1uV3wX5yZ7aB9cD1eF3gH5iJ7kL9mN1oP3qR5sT7uV9wX1yZ3aB5cD7eF9gH1iJ3kL5mN7oP9qR1sT3uV5wX7yZ9aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL\n-----END RSA PRIVATE KEY-----\n")
	dropbear.file("id_rsa.pub", "-rw-r--r--",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDu/ldzSbyVjQ2pR4tW6yZ8aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL7mN9oP1qR3sT5uV7wX9yZ1aB3cD5eF7gH9iJ1kL3mN5oP7qR9sT1uV3wX5yZ7aB9cD1eF3gH5iJ7kL9mN1oP3qR5sT7uV9wX1yZ3aB5cD7eF9gH1iJ3kL5mN7oP9qR1sT3uV5wX7yZ9aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL root@BCM96318\n")
	dropbear.file("id_dsa", "-rw-r--r--",
		"-----BEGIN DSA PRIVATE KEY-----\nMIIBvwIBAAKBgQDu/ldzSbyVjQ2pR4tW6yZ8aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL7mN9oP1qR3sT5uV7wX9yZ1aB3cD5eF7gH9iJ1kL3mN5oP7qR9sT1uV3wX5yZ7aB9cD1eF3gH5iJ7kL9mN1oP3qR5sT7uV9wX1yZ3aB5cD7eF9gH1iJ3kL5mN7oP9qR1sT3uV5wX7yZ9aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL\n-----END DSA PRIVATE KEY-----\n")
	dropbear.file("id_dsa.pub", "-rw-r--r--",
		"ssh-dss AAAAB3NzaC1lZWRlcnNhaWEAAAADAQABAAABAQDu/ldzSbyVjQ2pR4tW6yZ8aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL7mN9oP1qR3sT5uV7wX9yZ1aB3cD5eF7gH9iJ1kL3mN5oP7qR9sT1uV3wX5yZ7aB9cD1eF3gH5iJ7kL9mN1oP3qR5sT7uV9wX1yZ3aB5cD7eF9gH1iJ3kL5mN7oP9qR1sT3uV5wX7yZ9aB1cD3eF5gH7iJ9kL1mN3oP5qR7sT9uV1wX3yZ5aB7cD9eF1gH3iJ5kL root@BCM96318\n")

	etc.dir("httpd")
	etc.dir("firewall")
	crontabs := etc.dir("crontabs")
	crontabs.file("root", "-rw-r--r--", "# m h dom mon dow command\n*/5 * * * * /usr/bin/backup.sh\n")
	ssl := etc.dir("ssl")
	ssl.dir("certs")
	ssl.dir("private")
	ssl.dir("newcerts")
	ssl.dir("misc")
	udhcpcd := etc.dir("udhcpc")
	udhcpcd.file("default.script", "-rw-r--r--", "#!/bin/sh\nexit 0\n")
	for _, n := range []string{"rc.d", "opt", "X11", "pki", "security", "skel", "systemd", "default", "logrotate.d", "logcheck", "update-motd.d", "bash_completion.d"} {
		etc.dir(n)
	}
	profiled := etc.dir("profile.d")
	profiled.file("sh-local.sh", "-rw-r--r--", "# /etc/profile.d/sh-local.sh\n")
	profiled.file("vmware.sh", "-rw-r--r--", "# /etc/profile.d/vmware.sh\n")
}
