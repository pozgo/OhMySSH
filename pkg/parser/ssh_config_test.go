package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSSHConfigLoad(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config")

	configContent := `# SSH Configuration
Host server1
    HostName example.com
    Port 22
    User admin

Host server2
    HostName 192.168.1.100
    Port 2222
    User root
    IdentityFile ~/.ssh/id_rsa
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config := &SSHConfig{Path: configPath}
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	hosts := config.GetHosts()
	if len(hosts) != 2 {
		t.Errorf("Expected 2 hosts, got %d", len(hosts))
	}

	if hosts[0].Name != "server1" {
		t.Errorf("Expected host name 'server1', got '%s'", hosts[0].Name)
	}

	if hosts[0].Hostname != "example.com" {
		t.Errorf("Expected hostname 'example.com', got '%s'", hosts[0].Hostname)
	}

	if hosts[1].Port != "2222" {
		t.Errorf("Expected port '2222', got '%s'", hosts[1].Port)
	}
}

func TestSSHConfigLoadWithFixture(t *testing.T) {
	// Use test fixture to avoid touching real SSH config
	fixturePath := filepath.Join("..", "..", "test", "fixtures", "sample_ssh_config")
	
	config := &SSHConfig{Path: fixturePath}
	err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load fixture config: %v", err)
	}

	hosts := config.GetHosts()
	if len(hosts) != 4 {
		t.Errorf("Expected 4 hosts in fixture, got %d", len(hosts))
	}

	// Test specific hosts from fixture
	expectedHosts := map[string]string{
		"development-server": "dev.example.com",
		"production-server":  "prod.example.com",
		"local-vm":          "192.168.1.100",
		"bastion":           "bastion.company.com",
	}

	for _, host := range hosts {
		expectedHostname, exists := expectedHosts[host.Name]
		if !exists {
			t.Errorf("Unexpected host found: %s", host.Name)
			continue
		}
		if host.Hostname != expectedHostname {
			t.Errorf("Host %s: expected hostname %s, got %s", host.Name, expectedHostname, host.Hostname)
		}
	}
}

func TestSSHConfigLoadMissingFile(t *testing.T) {
	config := &SSHConfig{Path: "/nonexistent/path/config"}
	err := config.Load()
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}
}

func TestNewSSHConfig(t *testing.T) {
	config := NewSSHConfig()
	if config == nil {
		t.Error("NewSSHConfig returned nil")
	}
	if len(config.Hosts) != 0 {
		t.Error("NewSSHConfig should initialize with empty hosts")
	}
	if config.Path == "" {
		t.Error("NewSSHConfig should set default path")
	}
}