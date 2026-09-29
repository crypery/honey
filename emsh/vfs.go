package emsh

// This file defines the emulated (virtual) filesystem that the telnet/ssh
// shells present to attackers. It models a minimal busybox-based embedded
// Linux system (the kind found on a BCM96318 broadband router): the standard
// top-level directories, the busybox applet symlinks, the /etc configuration
// files, and the /proc and /sys pseudo-filesystems.
//
// The tree is built once, in memory, and is never backed by real files on
// disk. It is immutable and shared by every shell session. When a real file
// with the same path exists inside the shell's jail root (the unified root/
// tree), the real file's data takes precedence over the emulated data.

// vfsNode is a single entry in the emulated filesystem.
type vfsNode struct {
	name     string
	isDir    bool
	perm     string
	link     string // symlink target, "" when not a symlink
	content  []byte // file content (emulated fallback data)
	vsize    int    // displayed size override (0 = len(content))
	children map[string]*vfsNode
}

func vfsDir(name string) *vfsNode {
	return &vfsNode{name: name, isDir: true, perm: "drwxr-xr-x", children: map[string]*vfsNode{}}
}

func (d *vfsNode) file(name, perm, content string) {
	d.children[name] = &vfsNode{name: name, perm: perm, content: []byte(content)}
}

func (d *vfsNode) fileSize(name, perm string, content []byte, size int) {
	d.children[name] = &vfsNode{name: name, perm: perm, content: content, vsize: size}
}

func (d *vfsNode) symlink(name, target string) {
	d.children[name] = &vfsNode{name: name, perm: "lrwxrwxrwx", link: target}
}

func (d *vfsNode) dir(name string) *vfsNode {
	c := vfsDir(name)
	d.children[name] = c
	return c
}

// addLinks registers a batch of busybox applet symlinks pointing at target.
func (d *vfsNode) addLinks(target string, names ...string) {
	for _, n := range names {
		d.symlink(n, target)
	}
}

// emptyFiles registers a batch of zero-length files with the same permission.
func (d *vfsNode) emptyFiles(perm string, names ...string) {
	for _, n := range names {
		d.file(n, perm, "")
	}
}

// fakeELF returns a small, plausible ELF header so that `cat` of a binary
// shows binary-looking data rather than an empty file.
func fakeELF(name string) []byte {
	b := []byte{0x7f, 'E', 'L', 'F', 0x01, 0x01, 0x01, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	b = append(b, []byte("BusyBox v1.36.1 (2024-01-05 12:00:00 UTC) "+name+"\n")...)
	for i := 0; i < 64; i++ {
		b = append(b, byte(i))
	}
	return b
}

// vfsRoot is the shared emulated filesystem, built once at package load.
var vfsRoot = buildVFS()

func buildVFS() *vfsNode {
	root := vfsDir("/")
	buildBin(root)
	buildSbin(root)
	buildUsr(root)
	buildLib(root)
	buildEtc(root)
	buildDev(root)
	buildProc(root)
	buildSys(root)
	buildVar(root)
	buildHome(root)
	buildRootDir(root)
	for _, n := range []string{"tmp", "mnt", "opt", "media", "srv", "boot", "run"} {
		root.dir(n)
	}
	return root
}
