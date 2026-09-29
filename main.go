/*
 * Git: https://github.com/crypery
 * Author: https://crypery.com
 * License: GNU AGPL v3 (Affero GPL)
 */
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"honey/emsh"
	filedriver "honey/goftp/file-driver"
	"honey/goftp/server"

	"honey/gliderlabs/ssh"
	"honey/reiver/go-oi"
	"honey/reiver/go-telnet"
	"honey/reiver/go-telnet/telsh"
)

const timeLayout = "02/01/2006_15:04:05"

type Config struct {
	TelnetUse   bool   `json:"telnet_use"`
	SSHUse      bool   `json:"ssh_use"`
	FTPUse      bool   `json:"ftp_use"`
	HTTPUse     bool   `json:"http_use"`
	HTTPSUse    bool   `json:"https_use"`
	HTTPSEmule  bool   `json:"https_emule"`
	HTTPEmule   bool   `json:"http_emule"`
	FTPEmule    bool   `json:"ftp_emule"`
	TelnetEmule bool   `json:"telnet_emule"`
	SSHEmule    bool   `json:"ssh_emule"`
	TelnetPath  string `json:"telnet_path"`
	SSHPath     string `json:"ssh_path"`
	HTTPPath    string `json:"http_path"`
	HTTPSPath   string `json:"https_path"`
	FTPPath     string `json:"ftp_path"`
	TelnetServer string `json:"telnet_server"`
	SSHServer    string `json:"ssh_server"`
	FTPServer    string `json:"ftp_server"`
	HTTPServer   string `json:"http_server"`
	HTTPSServer  string `json:"https_server"`
	EmuleLogin   string `json:"emule_login"`
	EmulePass    string `json:"emule_pass"`
}

var (
	logMutex = &sync.Mutex{}
)

func exLog(path string, data string) error {
	logMutex.Lock()
	var file, err = os.OpenFile(filepath.FromSlash(path), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		logMutex.Unlock()
	}()
	_, err = file.Write([]byte(data))
	if err != nil {
		return err
	}
	file.Sync()
	return nil
}

func loadConfig(file string) (*Config, error) {
	cfg := &Config{
		TelnetUse:    true,
		SSHUse:       true,
		FTPUse:       true,
		HTTPUse:      true,
		HTTPSUse:     true,
		HTTPSEmule:   true,
		HTTPEmule:    true,
		FTPEmule:     true,
		TelnetEmule:  true,
		SSHEmule:     true,
		TelnetPath:   "/",
		SSHPath:      "/",
		HTTPPath:     "/home/public_html",
		HTTPSPath:    "/home/public_html",
		FTPPath:      "/home",
		TelnetServer: "BCM96318 Broadband Router",
		SSHServer:    "SSH-2.0-dropbear_0.53.1",
		FTPServer:    "Welcome to Pure-FTPd",
		HTTPServer:   "Apache/2.4.41 (Ubuntu)",
		HTTPSServer:  "Apache/2.4.41 (Ubuntu)",
		EmuleLogin:   "admin",
		EmulePass:    "admin",
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// randomToken returns a random hex string, used to set never-matching
// credentials when a service runs in capture-only (non-emulation) mode.
func randomToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// vfsPath resolves a configured virtual-filesystem path (e.g. "/", "/home",
// "/home/public_html") to the on-disk directory inside the unified root/ tree
// that backs every emulated service.
func vfsPath(p string) string {
	if p == "" || p == "/" {
		return "root"
	}
	return filepath.Join("root", filepath.FromSlash(p))
}

// ftpAuth implements the goftp server.Auth interface. In emulation mode it
// accepts only the configured emule_login/emule_pass; otherwise it always fails
// so that credentials are captured without granting access.
type ftpAuth struct {
	emule bool
	login string
	pass  string
}

func (a ftpAuth) CheckPasswd(name string, pass string) (bool, error) {
	if a.emule && name == a.login && pass == a.pass {
		return true, nil
	}
	return false, nil
}

// telnetShells holds the per-connection emulated shell for telnet sessions,
// keyed by the (per-connection) telnet.Context.
var telnetShells sync.Map

// handleWPLogin emulates the WordPress login form. When the submitted
// credentials match the configured emule_login/emule_pass the client is
// redirected to the (fake) admin area; wrong credentials get an error page.
func handleWPLogin(w http.ResponseWriter, r *http.Request, dir, emuleLogin, emulePass string) {
	q := r.URL.Query()
	login := q.Get("log")
	pass := q.Get("pwd")
	if login != "" && pass != "" {
		if login == emuleLogin && pass == emulePass {
			http.Redirect(w, r, "/wp-admin/", http.StatusFound)
			return
		}
		if data, err := os.ReadFile(filepath.Join(dir, "wp-login.php")); err == nil {
			html := string(data)
			errMsg := "<div class=\"alert alert-danger\" role=\"alert\" style=\"margin-bottom:15px;\">Error: The username or password you entered is incorrect.</div>"
			if i := strings.Index(html, "<form"); i >= 0 {
				html = html[:i] + errMsg + html[i:]
			}
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			w.Write([]byte(html))
			return
		}
	}
	http.ServeFile(w, r, filepath.Join(dir, "wp-login.php"))
}

func serverName(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if name != "" {
			w.Header().Set("Server", name)
		}
		next.ServeHTTP(w, r)
	})
}

func basicAuth(logPath string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			td := time.Now().Format(timeLayout)
			exLog(logPath, td+"\t"+ip+"\t"+username+"\t"+password+"\r\n")
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func accessLog(logPath string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		td := time.Now().Format(timeLayout)
		exLog(logPath, td+"\t"+ip+"\t"+r.Method+"\t"+r.URL.String()+"\r\n")
		next.ServeHTTP(w, r)
	})
}

func wordpressHandler(dir, emuleLogin, emulePass string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path
		if upath == "/wp-login.php" {
			handleWPLogin(w, r, dir, emuleLogin, emulePass)
			return
		}
		if strings.HasSuffix(upath, "/") {
			d := filepath.Join(dir, filepath.FromSlash(path.Clean(upath)))
			for _, idx := range []string{"index.html", "index.php"} {
				p := filepath.Join(d, idx)
				if _, err := os.Stat(p); err == nil {
					http.ServeFile(w, r, p)
					return
				}
			}
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, filepath.Join(dir, "404.html"))
			return
		}
		p := filepath.Join(dir, filepath.FromSlash(path.Clean(upath)))
		if _, err := os.Stat(p); os.IsNotExist(err) {
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, filepath.Join(dir, "404.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
}

// sshSessionHandler returns the SSH session handler. In emulation mode it runs
// a jailed shell rooted at the unified root/ tree and logs every command. A
// one-shot exec request (e.g. plink ... "ls") is run directly; otherwise an
// interactive shell is started.
func sshSessionHandler(cfg *Config) ssh.Handler {
	return func(s ssh.Session) {
		if !cfg.SSHEmule {
			<-s.Context().Done()
			return
		}
		ip := strings.Split(s.RemoteAddr().String(), ":")[0]
		shell := emsh.NewShell(vfsPath(cfg.SSHPath), s.User(), "BCM96318")
		logCmd := func(cmd, line string) {
			td := time.Now().Format(timeLayout)
			exLog("log/ssh.txt", td+"\t"+ip+"\t"+cmd+"\t"+line+"\r\n")
		}
		if cmd := s.RawCommand(); cmd != "" {
			// one-shot exec request (e.g. plink ... "ls"): the command arrives
			// in the exec payload, not on stdin, so run it directly and let the
			// library close the session.
			logCmd(emsh.Command(cmd), cmd)
			out, _ := shell.Exec(cmd)
			if out != "" {
				s.Write([]byte(out))
			}
			return
		}
		prompt := "\r\n>"
		shell.Run(s, s, prompt, logCmd)
	}
}

func main() {

	mime.AddExtensionType(".php", "text/html")

	fLog, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		log.SetOutput(fLog)
	}

	cfg, err := loadConfig("config.json")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("MAIN error loading config: %s\n", err)
			return
		}
		log.Printf("MAIN warning: config.json not found, using defaults\n")
	}

	err = os.MkdirAll(path.Join("log"), 0777)
	if err != nil {
		log.Printf("MAIN error creating dir: %s\n", err)
		return
	}

	for _, d := range []string{"root", "root/etc", "root/etc/ssh", "root/root", "root/root/.ssh", "root/var", "root/var/log", "root/home", "root/home/backup", "root/home/public_html"} {
		if err := os.MkdirAll(path.Join(d), 0777); err != nil {
			log.Printf("MAIN error creating dir: %s\n", err)
			return
		}
	}

	if cfg.TelnetUse {
		log.Println("telnet port: 23")
		go func() {
			shellHandler := telsh.NewShellHandler()
			shellHandler.WelcomeMessage = cfg.TelnetServer + "\r\n"
			shellHandler.Prompt = "\r\n>"
			if cfg.TelnetEmule {
				log.Println("telnet emule: " + cfg.TelnetPath)
				shellHandler.Login = cfg.EmuleLogin
				shellHandler.Password = cfg.EmulePass
			} else {
				// capture-only: never-matching credentials, login always fails
				shellHandler.Login = "honey_" + randomToken()
				shellHandler.Password = "honey_" + randomToken()
			}

			// "auth" producer: logs the credential attempt and, in emulation
			// mode with matching credentials, creates the jailed shell.
			// NOTE: the work is done in the ProducerFunc (called at Produce time),
			// NOT the HandlerFunc — the telsh auth path never calls handler.Run().
			authProducer := telsh.ProducerFunc(func(ctx telnet.Context, name string, args ...string) telsh.Handler {
				if name == "auth" && len(args) >= 4 {
					td, ip, login, password := args[0], args[1], args[2], args[3]
					exLog("log/telnet.txt", td+"\t"+ip+"\t"+login+"\t"+password+"\r\n")
					if cfg.TelnetEmule && login == cfg.EmuleLogin && password == cfg.EmulePass {
						shell := emsh.NewShell(vfsPath(cfg.TelnetPath), login, "BCM96318")
						shell.SetIP(ip)
						telnetShells.Store(ctx, shell)
					}
				}
				return telsh.PromoteHandlerFunc(func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, a ...string) error {
					return nil
				})
			})
			shellHandler.Register("auth", authProducer)

			// else producer: handles every command, logs it and runs it in the
			// jailed shell.
			cmdProducer := telsh.ProducerFunc(func(ctx telnet.Context, name string, args ...string) telsh.Handler {
				return telsh.PromoteHandlerFunc(func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, a ...string) error {
					val, ok := telnetShells.Load(ctx)
					if !ok {
						oi.LongWriteString(stdout, name+": command not found\r\n")
						return nil
					}
					shell := val.(*emsh.Shell)
					line := name
					if len(a) > 0 {
						line += " " + strings.Join(a, " ")
					}
					td := time.Now().Format(timeLayout)
					exLog("log/telnet.txt", td+"\t"+shell.IP()+"\t"+name+"\t"+line+"\r\n")
					out, _ := shell.Exec(line)
					if out != "" {
						oi.LongWriteString(stdout, out)
					}
					return nil
				}, args...)
			})
			shellHandler.RegisterElse(cmdProducer)

			addr := ":23"
			if err := telnet.ListenAndServe(addr, shellHandler); nil != err {
				panic(err)
			}
		}()
	}

	if cfg.SSHUse {
		log.Println("ssh port: 22")
		go func() {
			// Password auth: always logs the attempt. In emulation mode it
			// succeeds only for the configured emule_login/emule_pass.
			passwordHandler := func(ctx ssh.Context, td string, ip string, login string, password string) bool {
				exLog("log/ssh.txt", td+"\t"+ip+"\t"+login+"\t"+password+"\r\n")
				if cfg.SSHEmule && login == cfg.EmuleLogin && password == cfg.EmulePass {
					return true
				}
				return false
			}

			// Session handler: in emulation mode runs a jailed shell rooted at
			// the unified root/ tree and logs every command.
			ssh.Handle(sshSessionHandler(cfg))

			idle := 3 * time.Second
			max := 3 * time.Second
			if cfg.SSHEmule {
				log.Println("ssh emule: " + cfg.SSHPath)
				idle = 5 * time.Minute
				max = 30 * time.Minute
			}
			serverSSH := &ssh.Server{
				Addr:            ":22",
				MaxTimeout:      max,
				IdleTimeout:     idle,
				PasswordHandler: passwordHandler,
				Version:         cfg.SSHServer,
			}
			if err := serverSSH.ListenAndServe(); nil != err {
				panic(err)
			}
		}()
	}
	if cfg.HTTPUse {
		log.Println("http port: 80")
		go func() {
			var handler http.Handler
			if cfg.HTTPEmule {
				log.Println("http emule: " + cfg.HTTPPath)
				handler = serverName(cfg.HTTPServer, accessLog("log/http.txt", wordpressHandler(vfsPath(cfg.HTTPPath), cfg.EmuleLogin, cfg.EmulePass)))
			} else {
				handler = serverName(cfg.HTTPServer, basicAuth("log/http.txt"))
			}
			srv := &http.Server{
				Addr:         ":80",
				Handler:      handler,
				IdleTimeout:  3 * time.Second,
				ReadTimeout:  3 * time.Second,
				WriteTimeout: 3 * time.Second,
			}
			if err := srv.ListenAndServe(); nil != err {
				panic(err)
			}
		}()
	}
	if cfg.HTTPSUse {
		log.Println("https port: 443")
		go func() {
			var handler http.Handler
			if cfg.HTTPSEmule {
				log.Println("https emule: " + cfg.HTTPSPath)
				handler = serverName(cfg.HTTPSServer, accessLog("log/https.txt", wordpressHandler(vfsPath(cfg.HTTPSPath), cfg.EmuleLogin, cfg.EmulePass)))
			} else {
				handler = serverName(cfg.HTTPSServer, basicAuth("log/https.txt"))
			}
			srv := &http.Server{
				Addr:         ":443",
				Handler:      handler,
				IdleTimeout:  3 * time.Second,
				ReadTimeout:  3 * time.Second,
				WriteTimeout: 3 * time.Second,
			}
			if err := srv.ListenAndServeTLS("./localhost.crt", "./localhost.key"); nil != err {
				panic(err)
			}
		}()
	}
	if cfg.FTPUse {
		log.Println("ftp port: 21")
		rootPath := vfsPath(cfg.FTPPath)
		var auth server.Auth = ftpAuth{}
		if cfg.FTPEmule {
			log.Println("ftp emule: " + cfg.FTPPath)
			auth = ftpAuth{emule: true, login: cfg.EmuleLogin, pass: cfg.EmulePass}
		}
		factory := &filedriver.FileDriverFactory{
			RootPath: rootPath,
			Perm:     server.NewSimplePerm("user", "group"),
		}
		opts := &server.ServerOpts{
			Factory:        factory,
			Port:           21,
			Hostname:       "0.0.0.0",
			Auth:           auth,
			WelcomeMessage: cfg.FTPServer,
		}
		server.AuthInfo = func(td string, ip string, login string, password string) {
			exLog("log/ftp.txt", td+"\t"+ip+"\t"+login+"\t"+password+"\r\n")
			return
		}
		ftp := server.NewServer(opts)
		if err := ftp.ListenAndServe(); nil != err {
			panic(err)
		}
	}

	select {}
}
