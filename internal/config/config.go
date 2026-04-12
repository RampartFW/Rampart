package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Backend   BackendConfig   `yaml:"backend"`
	Cluster   ClusterConfig   `yaml:"cluster"`
	Snapshots SnapshotConfig  `yaml:"snapshots"`
	Audit     AuditConfig     `yaml:"audit"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	API       APIConfig       `yaml:"api"`
	WebUI     WebUIConfig     `yaml:"webui"`
	MCP       MCPConfig       `yaml:"mcp"`
	Logging   LoggingConfig   `yaml:"logging"`
	Metrics   MetricsConfig   `yaml:"metrics"`

	loadedFrom string
}

type ServerConfig struct {
	Listen     string    `yaml:"listen"`
	UnixSocket string    `yaml:"unixSocket"`
	TLS        TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
	CA   string `yaml:"ca"`
}

type BackendConfig struct {
	Type     string          `yaml:"type"`
	Nftables NftablesConfig  `yaml:"nftables"`
	Iptables IptablesConfig  `yaml:"iptables"`
	Ebpf     EbpfConfig      `yaml:"ebpf"`
	AWS      AWSConfig       `yaml:"aws"`
	Hybrid   HybridConfig    `yaml:"hybrid"`
}

type NftablesConfig struct {
	TableName string `yaml:"tableName"`
	Binary    string `yaml:"binary"`
}

type IptablesConfig struct {
	ChainPrefix string `yaml:"chainPrefix"`
	Binary      string `yaml:"binary"`
}

type EbpfConfig struct {
	XDPMode   string `yaml:"xdpMode"`
	Interface string `yaml:"interface"`
}

type AWSConfig struct {
	Region          string `yaml:"region"`
	SecurityGroupId string `yaml:"securityGroupId"`
}

type HybridConfig struct {
	FastPath string `yaml:"fastPath"`
	SlowPath string `yaml:"slowPath"`
}

type ClusterConfig struct {
	Enabled   bool      `yaml:"enabled"`
	NodeID    string    `yaml:"nodeId"`
	Listen    string    `yaml:"listen"`
	Advertise string    `yaml:"advertise"`
	Peers     []string  `yaml:"peers"`
	TLS       TLSConfig `yaml:"tls"`
}

type SnapshotConfig struct {
	Directory    string                `yaml:"directory"`
	Retention    SnapshotRetention     `yaml:"retention"`
	AutoSnapshot AutoSnapshotConfig    `yaml:"autoSnapshot"`
}

type SnapshotRetention struct {
	MaxCount int           `yaml:"maxCount"`
	MaxAge   time.Duration `yaml:"maxAge"`
}

type AutoSnapshotConfig struct {
	Interval time.Duration `yaml:"interval"`
	PreApply bool          `yaml:"preApply"`
}

type AuditConfig struct {
	Directory string        `yaml:"directory"`
	Retention time.Duration `yaml:"retention"`
	Compress  bool          `yaml:"compress"`
}

type SchedulerConfig struct {
	CheckInterval time.Duration `yaml:"checkInterval"`
}

type APIConfig struct {
	Keys []APIKey `yaml:"keys"`
}

type APIKey struct {
	Name        string   `yaml:"name"`
	Key         string   `yaml:"key"`
	Permissions []string `yaml:"permissions"`
}

type WebUIConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type MCPConfig struct {
	Enabled bool   `yaml:"enabled"`
	Listen  string `yaml:"listen"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
	File   string `yaml:"file"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Listen:     "0.0.0.0:9443",
			UnixSocket: "/var/run/rampart.sock",
		},
		Backend: BackendConfig{
			Type: "auto",
			Nftables: NftablesConfig{
				TableName: "rampart",
				Binary:    "/usr/sbin/nft",
			},
			Iptables: IptablesConfig{
				ChainPrefix: "RAMPART",
				Binary:      "/usr/sbin/iptables",
			},
		},
		Snapshots: SnapshotConfig{
			Directory: "/var/lib/rampart/snapshots",
			Retention: SnapshotRetention{
				MaxCount: 100,
				MaxAge:   720 * time.Hour,
			},
			AutoSnapshot: AutoSnapshotConfig{
				Interval: 6 * time.Hour,
				PreApply: true,
			},
		},
		Audit: AuditConfig{
			Directory: "/var/lib/rampart/audit",
			Retention: 2160 * time.Hour,
			Compress:  true,
		},
		Scheduler: SchedulerConfig{
			CheckInterval: 30 * time.Second,
		},
		WebUI: WebUIConfig{
			Enabled: true,
			Path:    "/ui",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stderr",
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Path:    "/metrics",
		},
	}
}

func LoadConfig(paths ...string) (*Config, error) {
	cfg := DefaultConfig()

	searchPaths := paths
	if len(searchPaths) == 0 {
		home, _ := os.UserHomeDir()
		searchPaths = []string{
			"rampart.yaml",
			"/etc/rampart/rampart.yaml",
			filepath.Join(home, ".config", "rampart", "rampart.yaml"),
		}
	}

	for _, path := range searchPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config %s: %w", path, err)
		}
		cfg.loadedFrom = path
		break
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("RAMPART_LISTEN"); v != "" {
		cfg.Server.Listen = v
	}
	if v := os.Getenv("RAMPART_BACKEND_TYPE"); v != "" {
		cfg.Backend.Type = v
	}
	// Add more env overrides as needed
}

func (cfg *Config) LoadedFrom() string {
	return cfg.loadedFrom
}
