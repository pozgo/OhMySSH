package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pozgo/OhMySSH/pkg/parser"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modeNormal mode = iota
	modeEditor
)

type vimMode int

const (
	vimNormal vimMode = iota
	vimInsert
	vimCommand
)

type model struct {
	width         int
	height        int
	sshConfig     *parser.SSHConfig
	hosts         []parser.Host
	selectedIdx   int
	configContent string
	currentMode   mode
	err           error
	textarea      textarea.Model
	saved         bool
	vimMode       vimMode
	commandBuffer string
	keySequence   string
	shouldConnect bool
	selectedHost  parser.Host
}

func initialModel() model {
	config := parser.NewSSHConfig()
	err := config.Load()
	var hosts []parser.Host
	var configContent string
	
	if err == nil {
		hosts = config.GetHosts()
		if content, readErr := ioutil.ReadFile(config.Path); readErr == nil {
			configContent = string(content)
		}
	}

	ta := textarea.New()
	ta.SetValue(configContent)
	ta.Focus()

	return model{
		sshConfig:     config,
		hosts:         hosts,
		selectedIdx:   0,
		configContent: configContent,
		currentMode:   modeNormal,
		err:           err,
		textarea:      ta,
		saved:         false,
		vimMode:       vimNormal,
		commandBuffer: "",
		keySequence:   "",
		shouldConnect: false,
		selectedHost:  parser.Host{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case string:
		// Handle SSH connection messages
		if strings.Contains(msg, "SSH connection") {
			// SSH process has completed, we can continue normally
			return m, nil
		}
	case tea.MouseMsg:
		// Handle mouse clicks in normal mode
		if m.currentMode == modeNormal {
			return m.handleMouseClick(msg)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
		m.textarea.SetHeight(msg.Height - 6)
		return m, nil
	case tea.KeyMsg:
		if m.currentMode == modeEditor {
			return m.handleVimKeybindings(msg)
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.selectedIdx > 0 {
					m.selectedIdx--
				}
				return m, nil
			case "down", "j":
				if m.selectedIdx < len(m.hosts)-1 {
					m.selectedIdx++
				}
				return m, nil
			case "enter", " ":
				// Connect to selected server
				if len(m.hosts) > 0 && m.selectedIdx < len(m.hosts) {
					m.shouldConnect = true
					m.selectedHost = m.hosts[m.selectedIdx]
					return m, tea.Quit
				}
				return m, nil
			case "e":
				m.currentMode = modeEditor
				m.vimMode = vimNormal
				m.commandBuffer = ""
				m.keySequence = ""
				m.textarea.Focus()
				m.textarea.SetValue(m.configContent)
				return m, nil
			}
		}
	}
	return m, nil
}

func (m model) handleVimKeybindings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch m.vimMode {
	case vimNormal:
		// Handle key sequences first
		key := msg.String()
		m.keySequence += key
		
		// Check for complete sequences
		if m.keySequence == "dd" {
			// Delete current line
			m.deleteCurrentLine()
			m.keySequence = ""
			return m, nil
		}
		
		if m.keySequence == "ZZ" {
			// Save and quit (Shift+Z+Z)
			err := m.saveConfig()
			if err != nil {
				m.err = err
				m.keySequence = ""
				return m, nil
			}
			m.saved = true
			m.sshConfig.Load()
			m.hosts = m.sshConfig.GetHosts()
			m.currentMode = modeNormal
			m.textarea.Blur()
			m.vimMode = vimNormal
			m.keySequence = ""
			return m, nil
		}
		
		// Reset sequence if it gets too long or doesn't match any pattern
		if len(m.keySequence) > 2 || (!strings.HasPrefix("dd", m.keySequence) && !strings.HasPrefix("ZZ", m.keySequence)) {
			m.keySequence = key // Start new sequence with current key
		}
		
		// Handle single character commands
		switch key {
		case "esc":
			// Exit editor and return to main view
			m.currentMode = modeNormal
			m.textarea.Blur()
			m.vimMode = vimNormal
			m.keySequence = ""
			return m, nil
		case "i":
			// Enter insert mode
			m.vimMode = vimInsert
			m.keySequence = ""
			return m, nil
		case "a":
			// Enter insert mode after cursor
			m.vimMode = vimInsert
			m.keySequence = ""
			// Move cursor right first
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		case "o":
			// Insert new line and enter insert mode
			m.vimMode = vimInsert
			m.keySequence = ""
			// Add newline
			currentValue := m.textarea.Value()
			m.textarea.SetValue(currentValue + "\n")
			return m, nil
		case ":":
			// Enter command mode
			m.vimMode = vimCommand
			m.commandBuffer = ":"
			m.keySequence = ""
			return m, nil
		case "h", "left":
			// Move cursor left
			m.keySequence = ""
			m.textarea, cmd = m.textarea.Update(tea.KeyMsg{Type: tea.KeyLeft})
			return m, cmd
		case "j", "down":
			// Move cursor down
			m.keySequence = ""
			m.textarea, cmd = m.textarea.Update(tea.KeyMsg{Type: tea.KeyDown})
			return m, cmd
		case "k", "up":
			// Move cursor up
			m.keySequence = ""
			m.textarea, cmd = m.textarea.Update(tea.KeyMsg{Type: tea.KeyUp})
			return m, cmd
		case "l", "right":
			// Move cursor right
			m.keySequence = ""
			m.textarea, cmd = m.textarea.Update(tea.KeyMsg{Type: tea.KeyRight})
			return m, cmd
		}
		return m, nil
		
	case vimInsert:
		switch msg.String() {
		case "esc":
			// Return to normal mode
			m.vimMode = vimNormal
			return m, nil
		default:
			// Pass through to textarea in insert mode
			m.textarea, cmd = m.textarea.Update(msg)
			m.saved = false
			return m, cmd
		}
		
	case vimCommand:
		switch msg.String() {
		case "esc":
			// Cancel command mode
			m.vimMode = vimNormal
			m.commandBuffer = ""
			return m, nil
		case "enter":
			// Execute command
			return m.executeVimCommand()
		case "backspace":
			// Remove character from command buffer
			if len(m.commandBuffer) > 1 {
				m.commandBuffer = m.commandBuffer[:len(m.commandBuffer)-1]
			}
			return m, nil
		default:
			// Add character to command buffer
			if len(msg.String()) == 1 {
				m.commandBuffer += msg.String()
			}
			return m, nil
		}
	}
	
	return m, nil
}

func (m model) executeVimCommand() (tea.Model, tea.Cmd) {
	command := m.commandBuffer[1:] // Remove the ':'
	m.commandBuffer = ""
	
	switch command {
	case "w", "write":
		// Save file
		err := m.saveConfig()
		if err != nil {
			m.err = err
		} else {
			m.saved = true
			m.sshConfig.Load()
			m.hosts = m.sshConfig.GetHosts()
		}
		m.vimMode = vimNormal
		return m, nil
	case "q", "quit":
		// Exit editor
		m.currentMode = modeNormal
		m.textarea.Blur()
		m.vimMode = vimNormal
		return m, nil
	case "wq", "x":
		// Save and quit
		err := m.saveConfig()
		if err != nil {
			m.err = err
			m.vimMode = vimNormal
			return m, nil
		}
		m.saved = true
		m.sshConfig.Load()
		m.hosts = m.sshConfig.GetHosts()
		m.currentMode = modeNormal
		m.textarea.Blur()
		m.vimMode = vimNormal
		return m, nil
	case "q!":
		// Force quit without saving
		m.currentMode = modeNormal
		m.textarea.Blur()
		m.vimMode = vimNormal
		return m, nil
	}
	
	// Unknown command, return to normal mode
	m.vimMode = vimNormal
	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	if m.currentMode == modeEditor {
		return m.renderEditor()
	}

	return m.renderNormalMode()
}

func (m model) renderEditor() string {
	editorStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1)

	// Status indicators
	status := ""
	if m.saved {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(" [SAVED]")
	} else {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Render(" [MODIFIED]")
	}

	// Vim mode indicator
	var vimModeStr string
	var modeColor lipgloss.Color
	switch m.vimMode {
	case vimNormal:
		vimModeStr = "NORMAL"
		modeColor = "39"  // Blue
	case vimInsert:
		vimModeStr = "INSERT"
		modeColor = "46"  // Green
	case vimCommand:
		vimModeStr = "COMMAND"
		modeColor = "214" // Orange
	}
	
	vimModeDisplay := lipgloss.NewStyle().
		Foreground(modeColor).
		Bold(true).
		Render(fmt.Sprintf(" [%s]", vimModeStr))

	// Command buffer and key sequence display
	commandDisplay := ""
	if m.vimMode == vimCommand {
		commandDisplay = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Render(m.commandBuffer)
	} else if m.keySequence != "" && m.vimMode == vimNormal {
		commandDisplay = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Render(m.keySequence)
	}

	// Help text based on mode
	var helpText string
	switch m.vimMode {
	case vimNormal:
		helpText = "i: insert | o: new line | dd: delete line | ZZ: save+quit | :: command | hjkl: move | ESC: exit"
	case vimInsert:
		helpText = "ESC: normal mode | Type to edit"
	case vimCommand:
		helpText = ":w save | :q quit | :wq save+quit | :q! force quit | ESC: cancel"
	}

	header := fmt.Sprintf("SSH Config Editor (Vim Mode)%s%s\n\n%s%s\n\n", 
		status, vimModeDisplay, helpText, commandDisplay)
	
	return editorStyle.Render(header + m.textarea.View())
}

func (m model) renderNormalMode() string {
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196"))

		return errorStyle.Render(fmt.Sprintf("Error loading SSH config: %v", m.err))
	}

	// Calculate panel dimensions (30/70 split)
	leftWidth := int(float64(m.width) * 0.3)
	rightWidth := m.width - leftWidth
	topHeight := (m.height - 2) / 2  // Subtract 2 for status bar
	bottomHeight := (m.height - 2) - topHeight

	// Server list panel (top-left)
	serverList := m.renderServerList(leftWidth, topHeight)
	
	// Server metadata panel (bottom-left)
	metadata := m.renderServerMetadata(leftWidth, bottomHeight)
	
	// Config preview panel (right)
	configPreview := m.renderConfigPreview(rightWidth, m.height-2)

	// Combine left panels
	leftPanel := lipgloss.JoinVertical(lipgloss.Left, serverList, metadata)
	
	// Combine all panels
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, configPreview)
	
	// Add status bar
	statusBar := m.renderStatusBar()
	
	return lipgloss.JoinVertical(lipgloss.Left, mainView, statusBar)
}

func (m model) renderServerList(width, height int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1)

	title := lipgloss.NewStyle().Bold(true).Render("SSH Servers")
	
	if len(m.hosts) == 0 {
		return style.Render(title + "\n\nNo servers found in SSH config")
	}

	var serverItems []string
	for i, host := range m.hosts {
		if i == m.selectedIdx {
			// Selected server - bold with highlighted background
			selectedStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("15")).  // Bright white text
				Background(lipgloss.Color("62")).  // Blue background
				Padding(0, 1)
			serverItems = append(serverItems, selectedStyle.Render("▶ "+host.Name))
		} else {
			// Unselected servers
			serverItems = append(serverItems, "  "+host.Name)
		}
	}

	content := title + "\n\n" + strings.Join(serverItems, "\n")
	return style.Render(content)
}

func (m model) renderServerMetadata(width, height int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("105")).
		Padding(1)

	title := lipgloss.NewStyle().Bold(true).Render("Server Details")
	
	if len(m.hosts) == 0 || m.selectedIdx >= len(m.hosts) {
		return style.Render(title + "\n\nNo server selected")
	}

	selected := m.hosts[m.selectedIdx]
	
	var details []string
	details = append(details, fmt.Sprintf("Host: %s", selected.Name))
	
	if selected.Hostname != "" {
		details = append(details, fmt.Sprintf("Hostname: %s", selected.Hostname))
	}
	
	if selected.User != "" {
		details = append(details, fmt.Sprintf("User: %s", selected.User))
	}
	
	if selected.Port != "" {
		details = append(details, fmt.Sprintf("Port: %s", selected.Port))
	}
	
	for key, value := range selected.Options {
		details = append(details, fmt.Sprintf("%s: %s", strings.Title(key), value))
	}

	content := title + "\n\n" + strings.Join(details, "\n")
	return style.Render(content)
}

func (m model) renderConfigPreview(width, height int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")).
		Padding(1)

	title := lipgloss.NewStyle().Bold(true).Render("SSH Config Preview")
	
	// Use highlighted content that shows the selected server
	highlightedContent := m.highlightSelectedServerInConfig()
	content := title + "\n\n" + highlightedContent
	
	// Note: We don't truncate highlighted content as it contains escape sequences
	// that would break the highlighting if cut off
	
	return style.Render(content)
}

func (m model) renderStatusBar() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("255")).
		Padding(0, 1)

	var statusItems []string
	statusItems = append(statusItems, "q: quit")
	statusItems = append(statusItems, "↑↓: navigate")
	statusItems = append(statusItems, "enter: connect")
	statusItems = append(statusItems, "e: edit config")
	statusItems = append(statusItems, "click: select server/edit config")

	return style.Render(strings.Join(statusItems, "  |  "))
}

func (m *model) saveConfig() error {
	content := m.textarea.Value()
	
	// Create backup before saving
	err := m.createBackup()
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}
	
	// Get current file info to preserve permissions
	fileInfo, err := os.Stat(m.sshConfig.Path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	
	// Write file with proper permissions (600 for SSH config)
	err = ioutil.WriteFile(m.sshConfig.Path, []byte(content), 0600)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}
	
	// Restore original file permissions if they were different
	if fileInfo.Mode() != 0600 {
		err = os.Chmod(m.sshConfig.Path, fileInfo.Mode())
		if err != nil {
			return fmt.Errorf("failed to restore file permissions: %v", err)
		}
	}
	
	// Update the config content
	m.configContent = content
	
	return nil
}

func (m *model) createBackup() error {
	// Create backup with timestamp
	backupPath := m.sshConfig.Path + ".backup." + fmt.Sprintf("%d", os.Getpid())
	
	// Read original file
	originalContent, err := ioutil.ReadFile(m.sshConfig.Path)
	if err != nil {
		return err
	}
	
	// Write backup
	err = ioutil.WriteFile(backupPath, originalContent, 0600)
	if err != nil {
		return err
	}
	
	return nil
}

func (m *model) deleteCurrentLine() {
	content := m.textarea.Value()
	lines := strings.Split(content, "\n")
	
	if len(lines) == 0 {
		return
	}
	
	// Find current cursor position to determine which line to delete
	// For simplicity, we'll delete the first line if we can't determine cursor position
	// In a full vim implementation, we'd track cursor line position
	
	// Remove the first line (this is a simplified approach)
	if len(lines) > 1 {
		newContent := strings.Join(lines[1:], "\n")
		m.textarea.SetValue(newContent)
	} else {
		// If only one line, clear it
		m.textarea.SetValue("")
	}
	
	m.saved = false
}

func (m model) highlightSelectedServerInConfig() string {
	if len(m.hosts) == 0 || m.selectedIdx >= len(m.hosts) {
		return m.configContent
	}
	
	selectedHostName := m.hosts[m.selectedIdx].Name
	lines := strings.Split(m.configContent, "\n")
	var highlightedLines []string
	
	inSelectedBlock := false
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Check if this line starts a new host block
		if strings.HasPrefix(strings.ToLower(trimmedLine), "host ") {
			// Extract host name from line
			parts := strings.Fields(trimmedLine)
			if len(parts) >= 2 {
				hostName := parts[1]
				if hostName == selectedHostName {
					inSelectedBlock = true
					// Highlight the Host line
					highlightedLine := lipgloss.NewStyle().
						Background(lipgloss.Color("62")).
						Foreground(lipgloss.Color("15")).
						Bold(true).
						Render(line)
					highlightedLines = append(highlightedLines, highlightedLine)
					continue
				} else {
					inSelectedBlock = false
				}
			}
		}
		
		// If we're in the selected block and this is an indented line (configuration)
		if inSelectedBlock && strings.HasPrefix(line, "    ") {
			// Highlight configuration lines
			highlightedLine := lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("15")).
				Render(line)
			highlightedLines = append(highlightedLines, highlightedLine)
		} else if inSelectedBlock && trimmedLine == "" {
			// Keep empty lines in the block but don't highlight
			highlightedLines = append(highlightedLines, line)
		} else {
			// Regular line or end of selected block
			if inSelectedBlock && trimmedLine != "" && !strings.HasPrefix(line, "    ") && !strings.HasPrefix(strings.ToLower(trimmedLine), "host ") {
				inSelectedBlock = false
			}
			highlightedLines = append(highlightedLines, line)
		}
	}
	
	return strings.Join(highlightedLines, "\n")
}

func (m model) handleMouseClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only handle left clicks
	if msg.Type != tea.MouseLeft {
		return m, nil
	}
	
	// Calculate panel dimensions (30/70 split) - same as in renderNormalMode
	leftWidth := int(float64(m.width) * 0.3)
	
	// Check if click is in the right panel (config preview)
	if msg.X >= leftWidth {
		// Click is in the config preview panel - open editor
		m.currentMode = modeEditor
		m.vimMode = vimNormal
		m.commandBuffer = ""
		m.keySequence = ""
		m.textarea.Focus()
		m.textarea.SetValue(m.configContent)
		return m, nil
	}
	
	// Check if click is in the left panel (server list)
	if msg.X < leftWidth {
		topHeight := (m.height - 2) / 2  // Subtract 2 for status bar
		
		// Check if click is in the server list (top-left panel)
		if msg.Y < topHeight {
			// Calculate which server was clicked based on Y position
			// Account for border and padding (2 lines for title + spacing)
			serverLineOffset := 3
			if msg.Y >= serverLineOffset && len(m.hosts) > 0 {
				clickedServerIdx := msg.Y - serverLineOffset
				if clickedServerIdx >= 0 && clickedServerIdx < len(m.hosts) {
					m.selectedIdx = clickedServerIdx
				}
			}
		}
	}
	
	return m, nil
}

func connectToServer(host parser.Host) {
	// Print beautiful connection info
	fmt.Printf("\n")
	fmt.Printf("🚀 Connecting to server via OhMySSH...\n")
	fmt.Printf("┌─────────────────────────────────────────┐\n")
	fmt.Printf("│ Server: %-31s │\n", host.Name)
	if host.Hostname != "" {
		fmt.Printf("│ Host:   %-31s │\n", host.Hostname)
	}
	if host.User != "" {
		fmt.Printf("│ User:   %-31s │\n", host.User)
	}
	if host.Port != "" {
		fmt.Printf("│ Port:   %-31s │\n", host.Port)
	}
	fmt.Printf("└─────────────────────────────────────────┘\n")
	fmt.Printf("Command: ssh %s\n", host.Name)
	fmt.Printf("\n")

	// Execute the SSH command through the user's shell to preserve wrappers and environment
	// This ensures all custom SSH wrappers, aliases, and PATH modifications are respected
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash" // fallback to bash
	}
	
	// Use -i flag to make it an interactive shell so aliases and functions are loaded
	cmd := exec.Command(shell, "-i", "-c", "ssh "+host.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		fmt.Printf("SSH connection failed: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	finalModel, err := p.Run()
	if err != nil {
		log.Printf("Error starting application: %v", err)
		os.Exit(1)
	}
	
	// Check if we should connect to a server
	if m, ok := finalModel.(model); ok && m.shouldConnect {
		connectToServer(m.selectedHost)
	}
}