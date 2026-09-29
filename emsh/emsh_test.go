package emsh

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newTestShell(t *testing.T) *Shell {
	t.Helper()
	root := t.TempDir()
	return NewShell(root, "root", "router")
}

func TestPwd(t *testing.T) {
	s := newTestShell(t)
	out, _ := s.Exec("pwd")
	if out != "/\r\n" {
		t.Fatalf("pwd = %q, want /", out)
	}
}

func TestLsRoot(t *testing.T) {
	s := newTestShell(t)
	out, _ := s.Exec("ls")
	for _, want := range []string{"bin", "etc", "proc", "sys", "dev", "var", "root", "home"} {
		if !strings.Contains(out, want) {
			t.Errorf("ls root missing %q in %q", want, out)
		}
	}
}

func TestCdAndPwd(t *testing.T) {
	s := newTestShell(t)
	if out, _ := s.Exec("cd /etc"); out != "" {
		t.Fatalf("cd /etc returned %q", out)
	}
	out, _ := s.Exec("pwd")
	if out != "/etc\r\n" {
		t.Fatalf("pwd = %q, want /etc", out)
	}
}

func TestCdInvalid(t *testing.T) {
	s := newTestShell(t)
	out, _ := s.Exec("cd /nope")
	if !strings.Contains(out, "No such file or directory") {
		t.Fatalf("cd /nope = %q", out)
	}
}

func TestCatVFS(t *testing.T) {
	s := newTestShell(t)
	out, _ := s.Exec("cat /etc/hostname")
	if !strings.Contains(out, "BCM96318") {
		t.Fatalf("cat /etc/hostname = %q, want to contain BCM96318", out)
	}
}

func TestCatMissing(t *testing.T) {
	s := newTestShell(t)
	out, _ := s.Exec("cat /etc/does-not-exist")
	if !strings.Contains(out, "No such file or directory") {
		t.Fatalf("cat missing = %q", out)
	}
}

func TestRealFileOverridesVFS(t *testing.T) {
	s := newTestShell(t)
	// /etc/hostname exists in the VFS; drop a real file at the same path.
	if err := os.MkdirAll(filepath.Join(s.Root, "etc"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s.Root, "etc", "hostname"), []byte("real-host\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, _ := s.Exec("cat /etc/hostname")
	if !strings.Contains(out, "real-host") {
		t.Fatalf("real file should override VFS, got %q", out)
	}
}

func TestRealDirAppearsInLs(t *testing.T) {
	s := newTestShell(t)
	if err := os.MkdirAll(filepath.Join(s.Root, "custom"), 0o755); err != nil {
		t.Fatal(err)
	}
	out, _ := s.Exec("ls /")
	if !strings.Contains(out, "custom") {
		t.Fatalf("real dir should appear in ls, got %q", out)
	}
}

func TestDotfileHiding(t *testing.T) {
	s := newTestShell(t)
	plain, _ := s.Exec("ls /root")
	all, _ := s.Exec("ls -a /root")
	if strings.Contains(plain, ".bashrc") && !strings.Contains(all, ".bashrc") {
		t.Fatalf("dotfile handling inconsistent: plain=%q all=%q", plain, all)
	}
	if !strings.Contains(all, ".bashrc") {
		t.Fatalf("ls -a /root should show .bashrc, got %q", all)
	}
}

func TestLsLong(t *testing.T) {
	s := newTestShell(t)
	out, _ := s.Exec("ls -l /etc")
	if !strings.Contains(out, "root root") {
		t.Fatalf("ls -l should show owner, got %q", out)
	}
}

// TestRunCRLineEnding checks that Shell.Run treats a bare "\r" (what PuTTY and
// other real terminals send on Enter) as a line terminator, not just "\n".
func TestRunCRLineEnding(t *testing.T) {
	s := newTestShell(t)
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()

	done := make(chan struct{})
	copied := make(chan struct{})
	var outBuf bytes.Buffer
	go func() {
		io.Copy(&outBuf, outR)
		close(copied)
	}()
	go func() {
		s.Run(inR, outW, "\r\n>", nil)
		outW.Close()
		close(done)
	}()

	// Send "ls" terminated by a bare \r, like PuTTY does on Enter.
	if _, err := inW.Write([]byte("ls\r")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	// Send "pwd" terminated by \r\n, like a scripted client does.
	if _, err := inW.Write([]byte("pwd\r\n")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	inW.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after input close")
	}
	select {
	case <-copied:
	case <-time.After(2 * time.Second):
		t.Fatal("output not fully read")
	}
	data := outBuf.String()
	t.Logf("run output: %q", data)
	if !strings.Contains(data, "bin") {
		t.Fatalf("ls (\\r-terminated) missing in output: %q", data)
	}
	if !strings.Contains(data, "/\r\n") {
		t.Fatalf("pwd (\\r\\n-terminated) missing in output: %q", data)
	}
}
