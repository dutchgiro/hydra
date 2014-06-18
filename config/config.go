package config

import (
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/innotech/hydra/log"

	"github.com/innotech/hydra/vendors/github.com/BurntSushi/toml"
	etcdConfig "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/discovery"
)

const (
	defaulApplicationFilePath = "/etc/hydra/apps.json"
	defaultBalanceTimeout     = 2500
	defaultConfigFilePath     = "/etc/hydra/hydra.conf"
	// TODO: Set default from name
	defaultDataDirectory           = "/tmp/hydra"
	defaultEtcdAddress             = "127.0.0.1:7401"
	defaultInstanceExpirtationTime = 30
	defaultLoadBalancerAddress     = "*:7777"
	defaultPeerAddress             = "127.0.0.1:7701"
	defaultPeerHeartbeatTiemout    = 50
	defaultPeerElectionTimeout     = 200
	defaultPrivateAPIAddress       = "127.0.0.1:7771"
	defaultPublicAPIAddress        = "127.0.0.1:7772"
	defaultSnapshot                = true
	defaultSnapshotCount           = 20000
	defaultVerbose                 = false
)

type Config struct {
	EtcdConf *etcdConfig.Config
	AppsConf *ApplicationsConfig

	AppsFilePath           string `toml:"apps_file"`
	BalanceTimeout         int    `toml:"balance_timeout"`
	BindAddr               string `toml:"bind_addr"`
	CAFile                 string `toml:"ca_file"`
	CertFile               string `toml:"cert_file"`
	ConfigFilePath         string
	DataDir                string `toml:"data_dir"`
	Discovery              string `toml:"discovery"`
	EtcdAddr               string `toml:"addr"`
	Force                  bool
	InstanceExpirationTime int      `toml:"instance_expiration_time"`
	KeyFile                string   `toml:"key_file"`
	LoadBalancerAddr       string   `toml:"load_balancer_addr"`
	Name                   string   `toml:"name"`
	Peers                  []string `toml:"peers"`
	PrivateAPIAddr         string   `toml:"private_api_addr"`
	PublicAPIAddr          string   `toml:"public_api_addr"`
	Snapshot               bool     `toml:"snapshot"`
	SnapshotCount          int      `toml:"snapshot_count"`
	Verbose                bool     `toml:"verbose"`
	Peer                   struct {
		Addr             string `toml:"addr"`
		BindAddr         string `toml:"bind_addr"`
		CAFile           string `toml:"ca_file"`
		CertFile         string `toml:"cert_file"`
		KeyFile          string `toml:"key_file"`
		HeartbeatTimeout int    `toml:"heartbeat_timeout"`
		ElectionTimeout  int    `toml:"election_timeout"`
	}
}

func New() *Config {
	c := new(Config)
	c.EtcdConf = etcdConfig.New()
	c.AppsFilePath = defaulApplicationFilePath
	c.BalanceTimeout = defaultBalanceTimeout
	c.ConfigFilePath = defaultConfigFilePath
	c.DataDir = defaultDataDirectory
	c.InstanceExpirationTime = defaultInstanceExpirtationTime
	c.LoadBalancerAddr = defaultLoadBalancerAddress
	c.PrivateAPIAddr = defaultPrivateAPIAddress
	c.PublicAPIAddr = defaultPublicAPIAddress
	c.Snapshot = defaultSnapshot
	c.SnapshotCount = defaultSnapshotCount
	c.Verbose = defaultVerbose
	c.Peer.Addr = defaultPeerAddress
	c.Peer.HeartbeatTimeout = defaultPeerHeartbeatTiemout
	c.Peer.ElectionTimeout = defaultPeerElectionTimeout

	return c
}

// Load configures hydra, it can be loaded from both system file,
// custom file or command line arguments and the values extracted from
// files they can be overriden with the command line arguments.
func (c *Config) Load(arguments []string) error {
	var path string
	f := flag.NewFlagSet("hydra", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&path, "config", "", "path to config file")
	f.Parse(arguments)

	if path != "" {
		// Load from config file specified in arguments.
		if err := c.LoadFile(path); err != nil {
			return err
		}
	} else {
		// Load from system file.
		if err := c.LoadSystemFile(); err != nil {
			return err
		}

	}

	// Load from command line flags.
	if err := c.LoadFlags(arguments); err != nil {
		return err
	}

	// Complete configuration.
	if c.Name == "" {
		c.NameFromHostname()
	}

	// Attempt cluster discovery.
	if c.Discovery != "" {
		if err := c.handleDiscovery(); err != nil {
			return err
		}
	}

	// Force remove server configuration if specified.
	if c.Force {
		if err := c.Reset(); err != nil {
			return err
		}
	}

	return nil
}

// LoadSystemFile loads from the default hydra configuration file path if it exists.
func (c *Config) LoadSystemFile() error {
	if _, err := os.Stat(c.ConfigFilePath); os.IsNotExist(err) {
		return nil
	}
	return c.LoadFile(c.ConfigFilePath)
}

// LoadFile loads configuration from a file.
func (c *Config) LoadFile(path string) error {
	_, err := toml.DecodeFile(path, &c)
	return err
}

// LoadFlags loads configuration from command line flags.
func (c *Config) LoadFlags(arguments []string) error {
	var peers, ignoredString string

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&c.EtcdAddr, "addr", c.EtcdAddr, "")
	f.StringVar(&c.AppsFilePath, "apps-file", c.AppsFilePath, "")
	f.IntVar(&c.BalanceTimeout, "balance-timeout", c.BalanceTimeout, "")
	f.StringVar(&c.BindAddr, "bind-addr", c.BindAddr, "")
	f.StringVar(&c.CAFile, "ca-file", c.CAFile, "")
	f.StringVar(&c.CertFile, "cert-file", c.CertFile, "")
	f.StringVar(&c.DataDir, "data-dir", c.DataDir, "")
	f.StringVar(&c.Discovery, "discovery", c.Discovery, "")
	f.BoolVar(&c.Force, "f", false, "")
	f.BoolVar(&c.Force, "force", false, "")
	f.IntVar(&c.InstanceExpirationTime, "instance-expiration-time", c.InstanceExpirationTime, "")
	f.StringVar(&c.KeyFile, "key-file", c.KeyFile, "")
	f.StringVar(&c.LoadBalancerAddr, "load-balancer-addr", c.LoadBalancerAddr, "")
	f.StringVar(&c.Name, "name", c.Name, "")
	f.StringVar(&peers, "peers", "", "")
	f.StringVar(&c.PrivateAPIAddr, "private-api-addr", c.PrivateAPIAddr, "")
	f.StringVar(&c.PublicAPIAddr, "public-api-addr", c.PublicAPIAddr, "")
	f.BoolVar(&c.Snapshot, "snapshot", true, "")
	f.IntVar(&c.SnapshotCount, "snapshot-count", c.SnapshotCount, "")
	f.BoolVar(&c.Verbose, "v", c.Verbose, "")
	f.BoolVar(&c.Verbose, "verbose", c.Verbose, "")

	f.StringVar(&c.Peer.Addr, "peer-addr", c.Peer.Addr, "")
	f.StringVar(&c.Peer.BindAddr, "peer-bind-addr", c.Peer.BindAddr, "")
	f.StringVar(&c.Peer.CAFile, "peer-ca-file", c.Peer.CAFile, "")
	f.StringVar(&c.Peer.CertFile, "peer-cert-file", c.Peer.CertFile, "")
	f.StringVar(&c.Peer.KeyFile, "peer-key-file", c.Peer.KeyFile, "")
	f.IntVar(&c.Peer.HeartbeatTimeout, "peer-heartbeat-timeout", c.Peer.HeartbeatTimeout, "")
	f.IntVar(&c.Peer.ElectionTimeout, "peer-election-timeout", c.Peer.ElectionTimeout, "")

	// BEGIN IGNORED FLAGS
	f.StringVar(&ignoredString, "config", "", "")
	// BEGIN IGNORED FLAGS

	if err := f.Parse(arguments); err != nil {
		return err
	}

	// Convert some parameters to lists.
	if peers != "" {
		c.Peers = strings.Split(peers, ",")
	}

	return nil
}

// LoadEtcdConfig loads etcd configuration.
func (c *Config) LoadEtcdConfig() error {
	fileContent := c.makeEtcdConfig()
	f, _ := ioutil.TempFile("", "")
	f.WriteString(fileContent)
	f.Close()
	defer os.Remove(f.Name())
	if err := c.EtcdConf.Load([]string{"-config", f.Name()}); err != nil {
		return err
	}

	return nil
}

// NameFromHostname sets the machine name from the hostname. This is to help
// people get started without thinking up a name.
func (c *Config) NameFromHostname() {
	host, err := os.Hostname()
	if err != nil && host == "" {
		log.Fatal("Node name required and hostname not set. e.g. '-name=name'")
	}
	c.Name = host
}

// Reset removes all server configuration files.
func (c *Config) Reset() error {
	if err := os.RemoveAll(filepath.Join(c.DataDir, "log")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(c.DataDir, "conf")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(c.DataDir, "snapshot")); err != nil {
		return err
	}

	return nil
}

// makeEtcdConfig generates a file content from the hydra settings emulating
// an etcd configuration file.
func (c *Config) makeEtcdConfig() string {
	var content string
	addLineToFileContent := func(line string) {
		content = content + line + "\n"
	}
	if c.EtcdAddr == "" {
		addLineToFileContent(`addr = "` + defaultEtcdAddress + `"`)
	} else {
		addLineToFileContent(`addr = "` + c.EtcdAddr + `"`)
	}
	addLineToFileContent(`ca_file = "` + c.CAFile + `"`)
	addLineToFileContent(`cert_file = "` + c.CertFile + `"`)
	addLineToFileContent(`data_dir = "` + c.DataDir + `"`)
	addLineToFileContent(`key_file = "` + c.KeyFile + `"`)
	addLineToFileContent(`name = "` + c.Name + `"`)
	peers := ""
	for i, addr := range c.Peers {
		if i > 0 {
			peers = peers + ", "
		}
		peers = peers + `"` + addr + `"`
	}
	addLineToFileContent(`peers = [` + peers + `]`)
	addLineToFileContent(`snapshot = ` + strconv.FormatBool(c.Snapshot))
	addLineToFileContent(`snapshot_count = ` + strconv.FormatInt(int64(c.SnapshotCount), 10))
	addLineToFileContent(`[peer]`)
	addLineToFileContent(`addr = "` + c.Peer.Addr + `"`)
	addLineToFileContent(`bind_addr = "` + c.Peer.BindAddr + `"`)
	addLineToFileContent(`ca_file = "` + c.Peer.CAFile + `"`)
	addLineToFileContent(`cert_file = "` + c.Peer.CertFile + `"`)
	addLineToFileContent(`key_file = "` + c.Peer.KeyFile + `"`)
	addLineToFileContent(`heartbeat_timeout = ` + strconv.FormatInt(int64(c.Peer.HeartbeatTimeout), 10))
	addLineToFileContent(`election_timeout = ` + strconv.FormatInt(int64(c.Peer.ElectionTimeout), 10))

	return content
}

func (c *Config) handleDiscovery() error {
	p, err := discovery.Do(c.Discovery, c.Name, c.Peer.Addr)

	// This is fatal, discovery encountered an unexpected error
	// and we have no peer list.
	if err != nil && len(c.Peers) == 0 {
		log.Fatalf("Discovery failed and a backup peer list wasn't provided: %v", err)
		return err
	}

	// Warn about errors coming from discovery, this isn't fatal
	// since the user might have provided a peer list elsewhere.
	if err != nil {
		log.Warnf("Discovery encountered an error but a backup peer list (%v) was provided: %v", c.Peers, err)
	}

	for i := range p {
		// Strip the scheme off of the peer if it has one
		// TODO(bp): clean this up!
		purl, err := url.Parse(p[i])
		if err == nil {
			p[i] = purl.Host
		}
	}

	c.Peers = p

	return nil
}
