<div align="center">

# 🚀 OhMySSH

*A modern, interactive SSH connection manager with a beautiful terminal UI*

[![Version](https://img.shields.io/github/v/release/pozgo/OhMySSH?style=for-the-badge)](https://github.com/pozgo/OhMySSH/releases)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-00ADD8.svg?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/github/workflow/status/pozgo/OhMySSH/CI?style=for-the-badge)](https://github.com/pozgo/OhMySSH/actions)

**Built with Go • Inspired by lazygit • Powered by Bubble Tea**

</div>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎯 **Core Features**
- 🚀 **Interactive SSH Manager** - Beautiful TUI interface
- 📝 **Vim-style Editor** - Modal editing with syntax highlighting  
- 🖱️ **Mouse Support** - Click to select and interact
- 🔧 **Smart Parser** - Automatic SSH config detection
- 🛡️ **Safe Editing** - Automatic backups & validation

</td>
<td width="50%">

### ⚡ **Performance**
- 🌐 **Cross-platform** - Linux, macOS, Windows
- ⚡ **Lightning Fast** - Built with Go
- 🪶 **Lightweight** - Single binary, no dependencies
- 🔄 **Real-time** - Live config preview
- 🎨 **Beautiful** - Stunning terminal UI

</td>
</tr>
</table>

---

## 🎬 Demo

<div align="center">

### 🖥️ **Beautiful Terminal Interface**

</div>

```text
┌─────────────────────┬─────────────────────────────────┐
│ 🖥️  SSH Servers     │                                 │
│ ▶ 🌐 web-server     │        📄 SSH Config Preview    │
│   🗄️ database-srv   │                                 │
│   🔧 dev-server     │  Host web-server                │
│   🚀 production     │      HostName example.com       │
├─────────────────────┤      User admin                 │
│ 📊 Server Details   │      Port 22                    │
│ 🏠 Host: web-server │      IdentityFile ~/.ssh/key    │
│ 🌍 IP: example.com  │                                 │
│ 👤 User: admin      │  Host database-srv              │
│ 🔌 Port: 22         │      HostName db.example.com    │
└─────────────────────┴─────────────────────────────────┘
  q: quit | ↑↓: navigate | ⏎: connect | e: edit | 🖱️: click
```

<div align="center">

*Three-panel layout with server list, details, and live config preview*

</div>

---

## 🚀 Quick Start

## 🛠️ Build from Source

**Prerequisites:**
- Go 1.19 or higher
- Git

**Step-by-step build instructions:**

```bash
# 1. Clone the repository
git clone https://github.com/pozgo/OhMySSH.git
cd OhMySSH

# 2. Download and resolve Go module dependencies
go mod tidy

# 3. Build for current platform
make build

# 4. Run the application
./OhMySSH
```

**Alternative build commands:**
```bash
# Build for all platforms (Linux, macOS, Windows)
make build-all

# Run tests
make test

# Clean build artifacts
make clean

# Install to system (requires sudo)
make install
```

**If you encounter dependency issues:**
```bash
# Force download specific modules
go mod download github.com/charmbracelet/bubbles
go mod download github.com/charmbracelet/bubbletea  
go mod download github.com/charmbracelet/lipgloss

# Then retry the build
make build
```

### 🎯 Launch & Connect

```bash
# 1. Launch OhMySSH
./OhMySSH

# 2. Navigate with ↑↓ or mouse
# 3. Press Enter to connect
# 4. Press 'e' to edit configs
```

---

## 🎮 Controls & Navigation

<div align="center">

### ⌨️ **Keyboard Shortcuts**

</div>

<table>
<tr>
<td width="33%">

#### 🧭 **Navigation**
| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate servers |
| `k` `j` | Vim-style navigation |
| `⏎` `Space` | Connect to server |
| `Tab` | Switch panels |

</td>
<td width="33%">

#### ⚙️ **Actions**
| Key | Action |
|-----|--------|
| `e` | Edit SSH config |
| `r` | Refresh server list |
| `f` | Search/filter |
| `?` | Show help |

</td>
<td width="33%">

#### 🚪 **Exit**
| Key | Action |
|-----|--------|
| `q` | Quit application |
| `Ctrl+C` | Force quit |
| `Esc` | Cancel/back |

</td>
</tr>
</table>

<div align="center">

### 🖱️ **Mouse Controls**
*Click anywhere to interact!*

</div>

<table>
<tr>
<td align="center" width="33%">

🖱️ **Click Server**<br/>
*Select from list*

</td>
<td align="center" width="33%">

🖱️ **Click Config**<br/>
*Open editor*

</td>
<td align="center" width="33%">

🖱️ **Right Click**<br/>
*Context menu*

</td>
</tr>
</table>

---

## 📝 Vim-Style Editor

<div align="center">

### 🎯 **Modal Editing Experience**

</div>

<table>
<tr>
<td width="33%">

#### 🎯 **Normal Mode**
```
┌─────────────────┐
│ Press 'i' to    │
│ enter Insert    │
│ mode            │
│                 │
│ dd - Delete line│
│ ZZ - Save & exit│
│ :  - Command    │
└─────────────────┘
```

</td>
<td width="33%">

#### ✏️ **Insert Mode**
```
┌─────────────────┐
│ Type to edit    │
│ your SSH config │
│                 │
│ ESC - Return to │
│ Normal mode     │
│                 │
│ Auto-complete   │
└─────────────────┘
```

</td>
<td width="33%">

#### 💻 **Command Mode**
```
┌─────────────────┐
│ :w   - Save     │
│ :q   - Quit     │
│ :wq  - Save+Quit│
│ :q!  - Force    │
│                 │
│ Full vim power! │
└─────────────────┘
```

</td>
</tr>
</table>

---

## 🔧 SSH Configuration

<div align="center">

### 📋 **Example Configuration**

</div>

```ssh
# 🌐 Production Web Server
Host web-server
    HostName production.example.com
    User admin
    Port 22
    IdentityFile ~/.ssh/production_key
    ServerAliveInterval 60
    ForwardAgent yes

# 🗄️ Database Server with Jump Host  
Host db-server
    HostName 10.0.1.100
    User dbadmin
    Port 5432
    ProxyJump bastion.example.com
    LocalForward 5432 localhost:5432

# 🔧 Development Environment
Host dev-server
    HostName dev.example.com
    User developer
    Port 2222
    IdentityFile ~/.ssh/dev_key
    ForwardX11 yes
    StrictHostKeyChecking no

# 🚀 Kubernetes Cluster
Host k8s-master
    HostName k8s.example.com
    User k8s-admin
    IdentityFile ~/.ssh/k8s_key
    DynamicForward 8080
```

<details>
<summary>🔍 <strong>Supported SSH Options</strong></summary>

| Option | Description | Example |
|--------|-------------|---------|
| `HostName` | Server address | `example.com` |
| `User` | Username | `admin` |
| `Port` | SSH port | `22` |
| `IdentityFile` | Private key | `~/.ssh/id_rsa` |
| `ProxyJump` | Jump host | `bastion.com` |
| `LocalForward` | Port forwarding | `8080:localhost:80` |
| `DynamicForward` | SOCKS proxy | `1080` |
| `ForwardAgent` | SSH agent | `yes` |
| `ForwardX11` | X11 forwarding | `yes` |
| `ServerAliveInterval` | Keep alive | `60` |

</details>

---

## 🔗 Connection Flow

<div align="center">

### 🎯 **Beautiful Connection Experience**

</div>

```console
🚀 Connecting to server via OhMySSH...
┌─────────────────────────────────────────┐
│ 🖥️  Server: web-server                  │
│ 🌍 Host:   production.example.com       │
│ 👤 User:   admin                        │
│ 🔌 Port:   22                           │
│ 🔑 Key:    ~/.ssh/production_key        │
└─────────────────────────────────────────┘
Command: ssh web-server

🔍 Retrieving authentication details...
🔐 Available methods: SSH key, password
🔑 Attempting SSH key authentication...
🚀 Connecting with SSH key authentication...
✅ Connected successfully!
```

<div align="center">

*Clean connection output with server details and authentication status*

</div>

---

## 🛡️ Security & Safety

<div align="center">

### 🔒 **Enterprise-Grade Security**

</div>

<table>
<tr>
<td width="50%">

#### 🛡️ **Safety Features**
- ✅ **Automatic Backups** before editing
- ✅ **Permission Preservation** (600/644)
- ✅ **Config Validation** prevents corruption
- ✅ **Test Isolation** protects real configs
- ✅ **Safe Mode** for critical operations

</td>
<td width="50%">

#### 🔐 **Security**
- ✅ **No Key Storage** - uses system SSH agent
- ✅ **Environment Preservation** - respects SSH wrappers
- ✅ **Audit Trail** - logs all config changes
- ✅ **Read-only Mode** for sensitive environments
- ✅ **Encrypted Backups** optional

</td>
</tr>
</table>

---

## 🏗️ Build & Development

<div align="center">

### 🛠️ **Build System**

</div>

<table>
<tr>
<td width="50%">

#### 🚀 **Quick Commands**
```bash
make build      # Current platform
make build-all  # All platforms  
make test       # Run tests
make install    # System install
make clean      # Clean artifacts
```

</td>
<td width="50%">

#### 🌍 **Cross-Platform**
- 🐧 **Linux** (amd64, arm64)
- 🍎 **macOS** (Intel, Apple Silicon)
- 🪟 **Windows** (amd64, arm64)
- 🐳 **Docker** support included

</td>
</tr>
</table>

<details>
<summary>📁 <strong>Project Structure</strong></summary>

```
🏗️ ohmyssh/
├── 📁 cmd/ohmyssh/         # 🚀 Main application
│   └── 📄 main.go
├── 📁 pkg/parser/          # 🔧 SSH config parser
│   ├── 📄 ssh_config.go
│   └── 📄 ssh_config_test.go
├── 📁 test/               # 🧪 Test fixtures
│   ├── 📁 fixtures/
│   └── 📄 README.md
├── 📁 build/              # 📦 Build outputs
├── 📁 assets/             # 🎨 Images & docs
├── 📄 Makefile           # 🛠️ Build automation
├── 📄 build.sh           # 📦 Cross-platform builds
├── 📄 go.mod             # 📋 Go dependencies
└── 📄 README.md          # 📖 This file
```

</details>

---

## 🐛 Troubleshooting

<div align="center">

### 🔧 **Common Issues & Solutions**

</div>

<details>
<summary>❌ <strong>SSH Config Not Found</strong></summary>

```bash
# Check if config exists
ls -la ~/.ssh/config

# Create if missing
touch ~/.ssh/config
chmod 600 ~/.ssh/config

# Verify permissions
stat -c %a ~/.ssh/config  # Should be 600
```

</details>

<details>
<summary>⌨️ <strong>Editor Not Working</strong></summary>

```bash
# Reset to normal mode
Press: ESC ESC ESC

# Force quit editor
Type: :q!

# Check terminal compatibility
echo $TERM
```

</details>

<details>
<summary>🔗 <strong>Connection Fails</strong></summary>

```bash
# Test SSH config manually
ssh -F ~/.ssh/config -T your-server

# Debug mode
ssh -vvv your-server

# Check OhMySSH debug mode
DEBUG=1 ./ohmyssh
```

</details>

<details>
<summary>🖱️ <strong>Mouse Not Working</strong></summary>

```bash
# Check terminal mouse support
echo $TERM

# Try different terminal
# - iTerm2 (macOS)
# - Windows Terminal
# - GNOME Terminal
# - Alacritty

# Fallback to keyboard navigation
```

</details>

<details>
<summary>🔨 <strong>Build Issues</strong></summary>

```bash
# Missing go.sum entries error
go mod tidy

# Module download failures
go clean -modcache
go mod download

# Build fails with missing dependencies
go mod tidy
go mod download
make build

# Permission denied on executable
chmod +x ./OhMySSH

# Go version compatibility
go version  # Should be 1.19+
```

</details>

---

## 🌟 Contributing

<div align="center">

### 🤝 **Join the Community!**

[![Contributors](https://img.shields.io/github/contributors/pozgo/OhMySSH?style=for-the-badge)](https://github.com/pozgo/OhMySSH/graphs/contributors)
[![Issues](https://img.shields.io/github/issues/pozgo/OhMySSH?style=for-the-badge)](https://github.com/pozgo/OhMySSH/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/pozgo/OhMySSH?style=for-the-badge)](https://github.com/pozgo/OhMySSH/pulls)

</div>

```bash
# 1. 🍴 Fork the repository
# 2. 🌿 Create feature branch
git checkout -b feature/amazing-feature

# 3. ✍️ Make your changes
# 4. ✅ Run tests
make test

# 5. 📤 Push and create PR
git push origin feature/amazing-feature
```

<div align="center">

### 💝 **Ways to Contribute**

</div>

<table>
<tr>
<td align="center" width="25%">

🐛 **Bug Reports**<br/>
*Help us improve*

</td>
<td align="center" width="25%">

✨ **New Features**<br/>
*Add awesome functionality*

</td>
<td align="center" width="25%">

📚 **Documentation**<br/>
*Make it even clearer*

</td>
<td align="center" width="25%">

🎨 **UI/UX**<br/>
*Make it more beautiful*

</td>
</tr>
</table>

---

## 📄 License

<div align="center">

This project is licensed under the **MIT License**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

*See [LICENSE](LICENSE) file for details*

</div>

---

## 🙏 Acknowledgments

<div align="center">

### 💎 **Built with Amazing Tools**

</div>

<table>
<tr>
<td align="center" width="33%">

[![Bubble Tea](https://img.shields.io/badge/Bubble_Tea-FF6B9D?style=for-the-badge&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea)<br/>
*TUI Framework*

</td>
<td align="center" width="33%">

[![Lip Gloss](https://img.shields.io/badge/Lip_Gloss-FFD93D?style=for-the-badge&logo=go&logoColor=black)](https://github.com/charmbracelet/lipgloss)<br/>
*Styling Engine*

</td>
<td align="center" width="33%">

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)<br/>
*Programming Language*

</td>
</tr>
</table>

<div align="center">

### 🌟 **Special Thanks**

*Inspired by [lazygit](https://github.com/jesseduffield/lazygit) • Built for the SSH community*

---

<img src="https://raw.githubusercontent.com/pozgo/OhMySSH/main/assets/footer.png" alt="Footer" width="100%"/>

**Made with ❤️ and lots of ☕ for developers everywhere**

[![GitHub Stars](https://img.shields.io/github/stars/pozgo/OhMySSH?style=social)](https://github.com/pozgo/OhMySSH/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/pozgo/OhMySSH?style=social)](https://github.com/pozgo/OhMySSH/network/members)
[![GitHub Watchers](https://img.shields.io/github/watchers/pozgo/OhMySSH?style=social)](https://github.com/pozgo/OhMySSH/watchers)

</div>