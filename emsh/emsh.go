package emsh

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Shell is a jailed command-line shell. It presents an emulated busybox-based
// Linux filesystem (see vfs.go) to the attacker, confined to Root. Every path
// operation is confined to Root; attempts to escape the root are denied. When
// a real file with the same path exists inside Root, its data takes precedence
// over the emulated data. It is used to emulate a terminal for the telnet and
// ssh honeypot services.
type Shell struct {
	Root string
	User string
	Host string
	ip   string
	cwd  string // absolute display path, "/" for the root itself
}

// NewShell creates a shell jailed to root, acting as user on host.
func NewShell(root, user, host string) *Shell {
	return &Shell{Root: root, User: user, Host: host, cwd: "/"}
}

// SetIP stores the client address used for command logging.
func (s *Shell) SetIP(ip string) { s.ip = ip }

// IP returns the stored client address.
func (s *Shell) IP() string { return s.ip }

// view is a merged, attacker-visible description of a single filesystem entry:
// the emulated (VFS) entry overridden by the real on-disk file when present.
type view struct {
	name    string
	isDir   bool
	perm    string
	size    int
	content []byte
	link    string
	real    bool // backed by a real on-disk file
}

// toDisplayPath converts a user supplied path (absolute or relative to cwd)
// into a cleaned absolute display path rooted at "/". Because the path is
// cleaned, ".." components cannot escape above the root.
func (s *Shell) toDisplayPath(p string) string {
	if p == "" {
		p = "."
	}
	var dp string
	if strings.HasPrefix(p, "/") {
		dp = p
	} else {
		dp = s.cwd + "/" + p
	}
	return path.Clean(dp)
}

// vfsLookup walks the emulated filesystem and returns the node for the
// absolute display path dp, or nil when it does not exist in the VFS.
func vfsLookup(root *vfsNode, dp string) *vfsNode {
	node := root
	for _, part := range strings.Split(strings.TrimPrefix(dp, "/"), "/") {
		if part == "" || part == "." {
			continue
		}
		if node.children == nil {
			return nil
		}
		node = node.children[part]
		if node == nil {
			return nil
		}
	}
	return node
}

// nodeView renders a VFS node as a view.
func nodeView(n *vfsNode) *view {
	size := len(n.content)
	if n.vsize > 0 {
		size = n.vsize
	}
	return &view{name: n.name, isDir: n.isDir, perm: n.perm, size: size, content: n.content, link: n.link}
}

// realView stats the real on-disk file at Root+dp and returns a view, or nil
// when no such file exists on disk.
func (s *Shell) realView(dp string) *view {
	realPath := filepath.Join(s.Root, filepath.FromSlash(dp))
	info, err := os.Stat(realPath)
	if err != nil {
		return nil
	}
	v := &view{name: filepath.Base(realPath), isDir: info.IsDir(), real: true}
	if v.isDir {
		v.perm = "drwxr-xr-x"
	} else {
		v.perm = "-rw-r--r--"
		v.size = int(info.Size())
	}
	return v
}

// lookup returns the merged view for an absolute display path: the VFS entry
// (if any) overridden by the real on-disk file (if present).
func (s *Shell) lookup(dp string) (*view, bool) {
	if node := vfsLookup(vfsRoot, dp); node != nil {
		v := nodeView(node)
		if !node.isDir {
			if real := s.realView(dp); real != nil {
				v.size = real.size
				v.perm = real.perm
				v.real = true
				v.content = nil
			}
		}
		return v, true
	}
	if real := s.realView(dp); real != nil {
		return real, true
	}
	return nil, false
}

// listDir returns the merged, sorted listing of an absolute display path that
// must be a directory: the VFS children plus any real on-disk children.
func (s *Shell) listDir(dp string) []view {
	seen := map[string]bool{}
	var out []view
	if node := vfsLookup(vfsRoot, dp); node != nil && node.isDir {
		names := make([]string, 0, len(node.children))
		for name := range node.children {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			child := node.children[name]
			v := nodeView(child)
			if !child.isDir {
				if real := s.realView(path.Join(dp, name)); real != nil {
					v.size = real.size
					v.perm = real.perm
					v.real = true
				}
			}
			out = append(out, *v)
			seen[name] = true
		}
	}
	realDir := filepath.Join(s.Root, filepath.FromSlash(dp))
	if entries, err := os.ReadDir(realDir); err == nil {
		for _, e := range entries {
			if seen[e.Name()] {
				continue
			}
			v := view{name: e.Name(), isDir: e.IsDir(), real: true}
			if e.IsDir() {
				v.perm = "drwxr-xr-x"
			} else {
				v.perm = "-rw-r--r--"
				if info, err := e.Info(); err == nil {
					v.size = int(info.Size())
				}
			}
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

// Command returns the first word of a command line.
func Command(line string) string {
	f := strings.Fields(line)
	if len(f) == 0 {
		return ""
	}
	return f[0]
}

// Exec runs a single command line and returns its output. The second return
// value is true when the session should terminate (exit/logout).
func (s *Shell) Exec(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", false
	}
	fields := strings.Fields(line)
	cmd := fields[0]
	args := fields[1:]
	switch cmd {
	case "exit", "logout":
		return "logout\r\n", true
	case "help":
		return s.help(), false
	case "pwd":
		return s.cwd + "\r\n", false
	case "whoami":
		return s.User + "\r\n", false
	case "id":
		return fmt.Sprintf("uid=0(%s) gid=0(root) groups=0(root)\r\n", s.User), false
	case "uname":
		return s.uname(args), false
	case "hostname":
		return s.Host + "\r\n", false
	case "ls":
		return s.ls(args), false
	case "cd":
		return s.cd(args), false
	case "cat":
		return s.cat(args), false
	case "echo":
		return strings.Join(args, " ") + "\r\n", false
	case "clear":
		return "\033[2J\033[H", false
	default:
		return cmd + ": command not found\r\n", false
	}
}

func (s *Shell) help() string {
	return "Available commands:\r\n" +
		"  ls [-la] [path]  - list directory\r\n" +
		"  cd [path]        - change directory\r\n" +
		"  cat <file>       - print file\r\n" +
		"  pwd              - print working directory\r\n" +
		"  whoami           - print current user\r\n" +
		"  id               - print user id\r\n" +
		"  uname [-a]       - system information\r\n" +
		"  hostname         - print hostname\r\n" +
		"  echo <text>      - print text\r\n" +
		"  clear            - clear the screen\r\n" +
		"  help             - this help\r\n" +
		"  exit             - leave the session\r\n"
}

func (s *Shell) uname(args []string) string {
	for _, a := range args {
		if strings.Contains(a, "a") {
			return fmt.Sprintf("Linux %s 3.10.100 #1 SMP Fri Jan 5 12:00:00 UTC 2024 mips GNU/Linux\r\n", s.Host)
		}
	}
	return "Linux\r\n"
}

func (s *Shell) ls(args []string) string {
	var target string
	long := false
	all := false
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			if strings.Contains(a, "l") {
				long = true
			}
			if strings.Contains(a, "a") {
				all = true
			}
			continue
		}
		target = a
	}
	if target == "" {
		target = "."
	}
	dp := s.toDisplayPath(target)
	v, ok := s.lookup(dp)
	if !ok {
		return "ls: " + target + ": No such file or directory\r\n"
	}
	if !v.isDir {
		if long {
			return s.longEntry(*v)+ "\r\n"
		}
		return v.name + "\r\n"
	}
	var b strings.Builder
	for _, e := range s.listDir(dp) {
		if !all && strings.HasPrefix(e.name, ".") {
			continue
		}
		if long {
			b.WriteString(s.longEntry(e))
		} else {
			name := e.name
			if e.isDir {
				name += "/"
			}
			b.WriteString(name + "  ")
		}
	}
	return b.String()
}

// longEntry renders a single entry in "ls -l" style.
func (s *Shell) longEntry(e view) string {
	perm := e.perm
	if e.link != "" {
		perm = "lrwxrwxrwx"
	}
	size := e.size
	if e.isDir {
		size = 4096
	}
	name := e.name
	if e.isDir {
		name += "/"
	}
	if e.link != "" {
		return fmt.Sprintf("%s  1 root root     %d Sep  6 10:00 %s -> %s\r\n", perm, size, name, e.link)
	}
	return fmt.Sprintf("%s  1 root root     %d Sep  6 10:00 %s\r\n", perm, size, name)
}

func (s *Shell) cd(args []string) string {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}
	dp := s.toDisplayPath(target)
	v, ok := s.lookup(dp)
	if !ok {
		return "cd: " + target + ": No such file or directory\r\n"
	}
	if !v.isDir {
		return "cd: " + target + ": Not a directory\r\n"
	}
	s.cwd = dp
	return ""
}

func (s *Shell) cat(args []string) string {
	if len(args) == 0 {
		return "cat: missing operand\r\n"
	}
	var b strings.Builder
	for _, a := range args {
		dp := s.toDisplayPath(a)
		v, ok := s.lookup(dp)
		if !ok {
			b.WriteString("cat: " + a + ": No such file or directory\r\n")
			continue
		}
		if v.isDir {
			b.WriteString("cat: " + a + ": Is a directory\r\n")
			continue
		}
		data := v.content
		if v.real {
			d, err := os.ReadFile(filepath.Join(s.Root, filepath.FromSlash(dp)))
			if err != nil {
				b.WriteString("cat: " + a + ": No such file or directory\r\n")
				continue
			}
			data = d
		}
		b.Write(data)
		if len(data) == 0 || data[len(data)-1] != '\n' {
			b.WriteString("\r\n")
		}
	}
	return b.String()
}

// Run executes an interactive shell loop: it reads command lines from r, writes
// a prompt and the command output to w, and calls logCmd(cmd, line) for every
// executed command. It returns when the input is closed or exit is issued.
//
// Both "\r" and "\n" are treated as line terminators (a real terminal sends
// "\r" on Enter, while scripted clients often send "\n" or "\r\n"), so the
// shell works with PuTTY, plink, ssh and the like.
func (s *Shell) Run(r io.Reader, w io.Writer, prompt string, logCmd func(cmd, line string)) {
	log.Printf("emsh: Run start user=%s cwd=%s", s.User, s.cwd)
	w.Write([]byte(prompt))
	var line []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			log.Printf("emsh: Run read %d bytes: %q", n, buf[:n])
			// Echo the input back to the terminal, like a real PTY does, so the
			// user sees what they type. Convert \r to \r\n so the cursor moves
			// to the next line; a \r\n pair is echoed only once.
			var echoBuf []byte
			for i := 0; i < n; i++ {
				b := buf[i]
				if b == '\r' {
					echoBuf = append(echoBuf, '\r', '\n')
					if i+1 < n && buf[i+1] == '\n' {
						i++ // the \n belongs to the \r\n pair
					}
				} else {
					echoBuf = append(echoBuf, b)
				}
			}
			if len(echoBuf) > 0 {
				w.Write(echoBuf)
			}
			line = append(line, buf[:n]...)
			for {
				// Find the earliest line terminator (\r or \n).
				crIdx := bytes.IndexByte(line, '\r')
				lfIdx := bytes.IndexByte(line, '\n')
				var idx int
				var isCR bool
				switch {
				case crIdx >= 0 && lfIdx >= 0:
					if crIdx < lfIdx {
						idx, isCR = crIdx, true
					} else {
						idx, isCR = lfIdx, false
					}
				case crIdx >= 0:
					idx, isCR = crIdx, true
				case lfIdx >= 0:
					idx, isCR = lfIdx, false
				default:
					idx = -1
				}
				if idx < 0 {
					break
				}
				cmdLine := strings.TrimRight(string(line[:idx]), "\r")
				line = line[idx+1:]
				// A "\r\n" pair is a single terminator: drop the trailing \n.
				if isCR && len(line) > 0 && line[0] == '\n' {
					line = line[1:]
				}
				log.Printf("emsh: Run terminator (isCR=%v) cmdLine=%q rest=%q", isCR, cmdLine, line)
				if cmdLine != "" {
					cmd := Command(cmdLine)
					if logCmd != nil {
						logCmd(cmd, cmdLine)
					}
					out, exit := s.Exec(cmdLine)
					log.Printf("emsh: Run exec %q -> out=%q exit=%v", cmdLine, out, exit)
					if out != "" {
						w.Write([]byte(out))
					}
					if exit {
						log.Printf("emsh: Run exit")
						return
					}
				}
				w.Write([]byte(prompt))
			}
		}
		if err != nil {
			log.Printf("emsh: Run read err=%v, returning", err)
			return
		}
	}
}
