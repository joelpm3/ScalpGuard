package config

import (
	"gopkg.in/yaml.v2"
	"os" // Add this line to import the os package
)


// ServerConfig represents the server configuration structure.
type ServerConfig struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Backends []BackendConfig `yaml:"backends"`
}

// BackendConfig represents the backend server configuration.
type BackendConfig struct {
	Name     string   `yaml:"name"`
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	Whitelist Whitelist `yaml:"whitelist"`
	Blacklist Blacklist `yaml:"blacklist"`
}

type Whitelist struct {
	UserAgents            []string `yaml:"user_agents"`
	TLSFingerprintsJA3    []string `yaml:"tls_fingerprints_ja3"`
	TLSFingerprintsJA3NE  []string `yaml:"tls_fingerprints_ja3_no_extension"`
	HTTP2Fingerprints     []string `yaml:"http2_fingerprints"`
	BrowserSettings       []string `yaml:"browser_settings"`
}

type Blacklist struct {
	UserAgents            []string `yaml:"user_agents"`
	TLSFingerprintsJA3    []string `yaml:"tls_fingerprints_ja3"`
	TLSFingerprintsJA3NE  []string `yaml:"tls_fingerprints_ja3_no_extension"`
	HTTP2Fingerprints     []string `yaml:"http2_fingerprints"`
	BrowserSettings       []string `yaml:"browser_settings"`
}

// BrowserConfig represents browser configurations.
type BrowserConfig struct {
	Browsers []Browser `yaml:"browsers"`
}

type Browser struct {
	Name                    string   `yaml:"name"`
	UserAgents              []string `yaml:"user_agents"`
	TLSFingerprintsJA3      []string `yaml:"tls_fingerprints_ja3"`
	TLSFingerprintsJA3NE    []string `yaml:"tls_fingerprints_ja3_no_extension"`
	HTTP2Fingerprints       []string `yaml:"http2_fingerprints"`
}

// LoadServerConfig loads server configuration from a YAML file.
func LoadServerConfig(path string) (ServerConfig, error) {
	var config ServerConfig
	bytes, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(bytes, &config)
	return config, err
}

// LoadBrowserConfig loads browser configuration from a YAML file.
func LoadBrowserConfig(path string) (BrowserConfig, error) {
	var config BrowserConfig
	bytes, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(bytes, &config)
	return config, err
}
