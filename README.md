# Honey

> [Russian version](README_RU.md)

**Honey** is a multi-protocol honeypot designed to analyze external bots that scan, probe, and attack servers. It emulates several common network services with realistic banners and login flows, lures automated attack tools into interacting with it, and silently records every credential they attempt.

## Description

Honey runs as a single Go binary and simultaneously listens on five standard service ports:

All services share a single on-disk virtual filesystem rooted at `root/`. Each service is jailed to a sub-path of that tree (configurable via the `*_path` parameters):

| Service | Port | Emulated target | Emulation root (jailed) |
|---------|------|-----------------|-------------------------|
| Telnet  | 23   | BCM96318 broadband router | `root/` shell (`telnet_path=/`) |
| SSH     | 22   | dropbear 0.53.1 (embedded device) | `root/` shell (`ssh_path=/`) |
| HTTP    | 80   | WordPress 6.8.2 website, or Basic-auth service | `root/home/public_html` site (`http_path`) |
| HTTPS   | 443  | WordPress 6.8.2 website, or Basic-auth service | `root/home/public_html` site (`https_path`) |
| FTP     | 21   | File server | `root/home` file root (`ftp_path`) |

Every service can be individually enabled or disabled through a JSON configuration file.

The emulated targets are deliberately chosen to resemble unattended embedded devices (routers, gateways) and a default WordPress website - the most frequent targets of mass-scanning bots. When an external bot connects and attempts to authenticate, Honey captures the full set of credentials (login and password), the source IP address, and a timestamp, and writes them to a per-protocol log.

By default the honeypot is fully passive: it never executes attacker commands and never grants access - every interaction is observed, recorded, and rejected. Optionally, per-service **emulation modes** can grant a *jailed* session (a confined shell or file root the attacker cannot escape), which lures interactive attackers into running commands that are all logged. Both modes are safe to expose on a public IP address.

## Features

- **Multi-protocol emulation** - Telnet, SSH, HTTP, HTTPS, and FTP services running concurrently in a single process.
- **Service configuration** - every service is individually enabled or disabled via a JSON configuration file, so the honeypot can be tuned to expose only the services of interest.
- **Realistic banners** - the SSH server presents itself as `SSH-2.0-dropbear_0.53.1` and the Telnet service greets the client with a `BCM96318 Broadband Router` login prompt, mimicking a typical unattended embedded device and maximizing the number of bots that engage.
- **HTTPS with TLS** - the HTTPS service is served over TLS using a local certificate (`localhost.crt` / `localhost.key`).
- **WordPress emulation** - in emulation mode the HTTPS service serves a full static snapshot of a real WordPress 6.8.2 site: all pages (homepage, posts, archives, login page), the complete visual design (theme, fonts, styles, scripts, images), and the standard root endpoints (`index.php`, `readme.html`, `license.txt`, `wp-login.php`, `xmlrpc.php`, `wp-*.php`,...). Unknown URLs receive the site's own styled 404 page. Every request is logged with a timestamp, source IP, HTTP method, and URL, luring CMS-targeting bots and scanners into revealing their probing patterns.
- **Credential harvesting** - every login/password pair attempted by attackers is captured across all protocols (Telnet login sequence, SSH password authentication, HTTP/HTTPS Basic header, FTP login).
- **Full audit trail** - each captured attempt is logged with a timestamp and the source IP address, separately per protocol, which allows correlating a single bot's activity across services.
- **Passive or jailed, per service** - in capture-only mode no commands are executed and no access is granted; in emulation mode the attacker gets a jailed shell/file-root that is strictly confined (path traversal is denied) and fully logged, so the honeypot cannot be compromised through the emulated services.
- **Unified virtual filesystem** - every service is jailed to a sub-path of a single on-disk tree rooted at `root/`. The Telnet/SSH shells layer a realistic busybox-based embedded-Linux filesystem (a BCM96318-style router) on top of it: standard top-level directories, busybox applet symlinks, `/etc` configuration files, and `/proc`/`/sys` pseudo-filesystems are emulated in memory, while any real file placed under `root/` overrides the emulated data at the same path. HTTP/HTTPS serve `root/home/public_html` and FTP serves `root/home` from the same tree.
- **Single static binary** - built with `CGO_ENABLED=0`, the result is a self-contained executable with no runtime dependencies.
- **Cross-compilation** - one source tree builds for both Linux and Windows.

## Technologies

- **Go 1.21+** - the entire honeypot is written in Go and compiles to a single static binary.
- **golang.org/x/crypto/ssh** - the SSH protocol implementation (key exchange, cryptography, authentication) used by the SSH honeypot.
- **github.com/anmitsu/go-shlex** - shell-style word splitting used to parse SSH command strings.
- **Vendored protocol stacks** - the Telnet, FTP, and SSH server frameworks are vendored inside the repository, so the project builds reproducibly with only the two external modules above.

## Methodology

1. **Deployment** - the binary is deployed on a host with a public IP address and started with privileges sufficient to bind ports 21, 22, 23, 80, and 443. The JSON configuration file selects which services are exposed.
2. **Luring** - mass-scanning bots, port scanners, and credential-stuffing tools discover the open ports and interact with the emulated services exactly as they would with a real, unattended device. The embedded-device banners (dropbear, BCM96318 router) and the WordPress website snapshot are specifically attractive to opportunistic bots: CMS scanners probe the WordPress installation for the login page, admin area, standard root endpoints, and known vulnerabilities.
3. **Capture** - every authentication attempt is intercepted at the protocol level and recorded together with the source IP address and a timestamp. In WordPress emulation mode, every single request (GET/POST and URL) is recorded as an access log. Failed logins are answered with realistic rejections so the bot continues its behavior.
4. **Analysis** - the collected logs allow the operator to study bot behavior: which credential lists are used, which tools and client fingerprints appear, which URLs and endpoints are probed, attack frequency, and source networks. Because the same bot typically tries several protocols, per-protocol logs can be correlated by IP to reconstruct a bot's full attack sequence.

Because the honeypot mimics a low-value embedded device and a default WordPress website, it attracts a high volume of opportunistic bot traffic, making it an effective sensor for observing the current state of automated attacks.

## Configuration

All services are controlled by a JSON configuration file read at startup. Each service can be enabled or disabled independently:

| Parameter | Type | Description |
|-----------|------|-------------|
| `telnet_use`  | bool | Enable the Telnet service (port 23). |
| `ssh_use`     | bool | Enable the SSH service (port 22). |
| `ftp_use`     | bool | Enable the FTP service (port 21). |
| `http_use`    | bool | Enable the HTTP service (port 80). |
| `https_use`   | bool | Enable the HTTPS service (port 443). |
| `http_emule`  | bool | If `true`, the HTTP service serves the WordPress snapshot at `http_path`; if `false`, it requires HTTP Basic authentication and captures credentials. |
| `https_emule` | bool | If `true`, the HTTPS service serves the WordPress snapshot at `https_path`; if `false`, it requires HTTP Basic authentication and captures credentials. |
| `ftp_emule`   | bool | If `true`, FTP authenticates with `emule_login`/`emule_pass` and serves the `ftp_path` root; if `false`, every login is rejected (credentials still captured). |
| `telnet_emule`| bool | If `true`, Telnet authenticates with `emule_login`/`emule_pass` and opens a jailed shell rooted at `telnet_path`; if `false`, every login is rejected (credentials still captured). |
| `ssh_emule`   | bool | If `true`, SSH authenticates with `emule_login`/`emule_pass` and opens a jailed shell rooted at `ssh_path`; if `false`, every login is rejected (credentials still captured). |
| `telnet_path` | string | Sub-path of the unified `root/` tree used as the Telnet shell jail root. Default `/`. |
| `ssh_path`    | string | Sub-path of the unified `root/` tree used as the SSH shell jail root. Default `/`. |
| `http_path`   | string | Sub-path of the unified `root/` tree served by the HTTP service. Default `/home/public_html`. |
| `https_path`  | string | Sub-path of the unified `root/` tree served by the HTTPS service. Default `/home/public_html`. |
| `ftp_path`    | string | Sub-path of the unified `root/` tree served by the FTP service. Default `/home`. |
| `emule_login`  | string | Login accepted by the emulation modes (Telnet, SSH, FTP, WordPress `wp-login.php`). Default `admin`. |
| `emule_pass`   | string | Password accepted by the emulation modes. Default `admin`. |
| `http_server`   | string | Value of the `Server` response header for the HTTP service (port 80). |
| `https_server`  | string | Value of the `Server` response header for the HTTPS service (port 443). |
| `ftp_server`    | string | FTP server name, sent as the `220` welcome banner. |
| `telnet_server` | string | Telnet banner shown to the client on connect. |
| `ssh_server`    | string | SSH version string (banner) advertised during the handshake. |

If the configuration file is missing, all services and all emulation modes are enabled by default, with `emule_login`/`emule_pass` set to `admin` and the `*_path` parameters set to `/` (Telnet, SSH), `/home/public_html` (HTTP, HTTPS), and `/home` (FTP).

### Emulation modes

When an emulation mode is enabled and the attacker supplies the configured `emule_login`/`emule_pass`, Honey grants a **jailed** emulated session instead of rejecting it:

- **Telnet / SSH** - an interactive shell (`ls`, `cd`, `cat`, `pwd`, `whoami`, `id`, `uname`, `hostname`, `echo`, `clear`, `help`, `exit`) that presents an emulated busybox-based Linux filesystem (a BCM96318-style embedded router: `/bin` busybox applets, `/etc` configuration, `/proc` and `/sys` pseudo-filesystems, `/dev`, `/var`, `/root`, …). The emulated tree is built in memory and layered over the unified `root/` directory (the `telnet_path`/`ssh_path` sub-path, `/` by default); attempts to escape (`cd../../etc`, …) are denied. Any file or directory that actually exists under `root/` overrides the emulated data at the same path, so the operator can drop real files (e.g. a custom `/etc/passwd` or `/var/log/messages`) into `root/` to control exactly what the attacker sees. Every executed command is logged as `timestamp<TAB>ip<TAB>command<TAB>full-line`.
- **FTP** - the attacker browses and downloads (and uploads) files inside the `ftp_path` sub-path of `root/` (`/home` by default) only; path traversal is blocked.
- **WordPress** - submitting the credentials to `wp-login.php` redirects to a fake `wp-admin/` dashboard.

When an emulation mode is disabled, the service still captures and logs every credential attempt but never authorizes it. This lets the operator run a fully passive (capture-only) honeypot or a fully interactive (jailed) one, per service.

### Unified virtual filesystem

All services share a single on-disk tree rooted at `root/`. The Telnet and SSH jailed shells do not expose the host's real filesystem: they layer a virtual filesystem (VFS) that models a minimal busybox-based embedded Linux system - the kind found on a BCM96318 broadband router - over the `root/` directory. The emulated tree is built once, in memory, at startup and is shared by every session:

- **`/bin`, `/sbin`, `/usr/bin`, `/usr/sbin`** - the `busybox` binary plus symlinks for the standard busybox applets (`ls`, `cat`, `ps`, `ifconfig`, `netstat`, `ping`, `grep`, `sed`, `awk`, …).
- **`/etc`** - configuration files (`passwd`, `shadow`, `group`, `hosts`, `hostname`, `profile`, `inittab`, `init.d/`, `network/`, `dropbear/`, `ssh/`, `ssl/`, …).
- **`/proc`, `/sys`** - pseudo-filesystems (`cpuinfo`, `meminfo`, `net/`, `version`, …).
- **`/dev`, `/var`, `/root`, `/home`, `/tmp`, …`** - device nodes, logs, and home directories.

The same `root/` tree also backs the other services: HTTP/HTTPS serve `root/home/public_html` (the WordPress snapshot) and FTP serves `root/home`.

Two rules govern what the Telnet/SSH attacker sees:

1. **Emulated by default** - any path that exists in the VFS is shown with its emulated content, even though no such file exists on disk.
2. **Real files win** - if a real file or directory with the same path exists under `root/`, its data (content, size, permissions) overrides the emulated entry. This lets the operator customize exactly what is exposed by dropping real files into `root/`, without touching the emulated tree.

Because the VFS is immutable and every service is confined to its sub-path of `root/`, the honeypot never reads or writes outside that tree, and no attacker command can reach the host filesystem.

## Modules

| Module | Role |
|--------|------|
| `golang.org/x/crypto/ssh` | Core SSH protocol implementation (key exchange, cryptography, authentication) used by the SSH honeypot. |
| `github.com/anmitsu/go-shlex` | Splits SSH command strings into shell-style words for parsing. |
| `gliderlabs/ssh` (vendored) | High-level SSH server framework: session handling, password-authentication callbacks, TCP/IP forwarding, keep-alive. |
| `reiver/go-telnet` (vendored) | Telnet server: protocol negotiation, connection handling, optional TLS. |
| `reiver/go-telnet/telsh` (vendored) | Telnet shell layer: login/password prompt, command registration and dispatch. |
| `reiver/go-oi` (vendored) | I/O helper primitives (long writes, seekable readers) used by the Telnet stack. |
| `goftp` (vendored) | FTP server: command processing, authentication, passive/active data connections, plus a file-system driver. |

## Build

Build scripts are provided for both targets:

```bat
build_windows.bat:: builds honey.exe (Windows amd64)
build_linux.bat:: builds honey (Linux amd64)
```

Or manually:

```
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o honey.exe.
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o honey.
```

## License

GNU AGPL v3 (Affero GPL)
