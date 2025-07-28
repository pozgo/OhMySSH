package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Host struct {
	Name     string
	Hostname string
	Port     string
	User     string
	Options  map[string]string
}

type SSHConfig struct {
	Hosts []Host
	Path  string
}

func NewSSHConfig() *SSHConfig {
	return &SSHConfig{
		Hosts: make([]Host, 0),
		Path:  getDefaultSSHConfigPath(),
	}
}

func getDefaultSSHConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".ssh", "config")
}

func (c *SSHConfig) Load() error {
	file, err := os.Open(c.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Clear existing hosts to prevent duplicates
	c.Hosts = make([]Host, 0)
	
	scanner := bufio.NewScanner(file)
	var currentHost *Host

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch key {
		case "host":
			if currentHost != nil {
				c.Hosts = append(c.Hosts, *currentHost)
			}
			currentHost = &Host{
				Name:    value,
				Options: make(map[string]string),
			}
		case "hostname":
			if currentHost != nil {
				currentHost.Hostname = value
			}
		case "port":
			if currentHost != nil {
				currentHost.Port = value
			}
		case "user":
			if currentHost != nil {
				currentHost.User = value
			}
		default:
			if currentHost != nil {
				currentHost.Options[key] = value
			}
		}
	}

	if currentHost != nil {
		c.Hosts = append(c.Hosts, *currentHost)
	}

	return scanner.Err()
}

func (c *SSHConfig) GetHosts() []Host {
	return c.Hosts
}