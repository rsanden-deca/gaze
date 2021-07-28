package common

import (
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "time"
)

var config *Config

type Config struct {
    Logfpath              string
    Hostname              string
    Service               string
    TTL                   time.Duration
    RiemannEndpoints      []string
    RiemannConnectTimeout time.Duration
    RiemannRetryDelay     time.Duration
    RiemannRetryAttempts  uint
}

func UsageExit() {
    fmt.Printf("Usage: %s --hostname myhostname --service myservice --ttl 10 --riemann 127.0.0.1:5555 --riemann 127.0.0.2:5555 --riemann 127.0.0.3:5555 /path/to/file.log\n", os.Args[0])
    os.Exit(1)
}

func (self *Config) Check() bool {
    if self.Logfpath == "" { return false }
    if self.Hostname == "" { return false }
    if self.Service == "" { return false }
    if self.TTL == 0 { return false }
    if self.RiemannConnectTimeout == 0 { return false }

    if len(self.RiemannEndpoints) == 0 { return false }
    for _, endpoint := range self.RiemannEndpoints {
        pieces := strings.Split(endpoint, ":")
        if len(pieces) != 2 { return false }
        _, err := strconv.ParseInt(pieces[1], 10, 64)
        if err != nil { return false }
    }

    return true
}

func GetHostname() string {
    hostname, err := os.Hostname()
    if err != nil {
        log.Printf("ERROR : Failed to determine hostname : %s\n", err.Error())
        os.Exit(1)
    }
    return hostname
}

type arrayFlags []string

func (self *arrayFlags) String() string {
    return strings.Join(*self, " ")
}

func (self *arrayFlags) Set(value string) error {
    *self = append(*self, value)
    return nil
}

func GetNewConfig() *Config {
    cfg := &Config{}
    hostnamePtr := flag.String("hostname", "", "Hostname to use for :host in riemann events")
    servicePtr := flag.String("service", "", "Service to use for :service in riemann events")
    ttlPtr := flag.Uint("ttl", 60, "TTL to use (in seconds) for :ttl in riemann events")

    var riemannEndpoints arrayFlags
    flag.Var(&riemannEndpoints, "riemann", "Riemann endpoint host:port")

    flag.Parse()

    if flag.NArg() != 1 {
        UsageExit()
    }
    cfg.Logfpath = flag.Arg(0)

    cfg.Hostname = *hostnamePtr
    cfg.Service = *servicePtr
    cfg.TTL = time.Duration(uint(*ttlPtr) * uint(time.Second))
    cfg.RiemannConnectTimeout = time.Duration(RIEMANN_CONNECT_TIMEOUT_MILLIS * uint(time.Millisecond))
    cfg.RiemannRetryDelay = time.Duration(RIEMANN_RETRY_DELAY_MILLIS * uint(time.Millisecond))
    cfg.RiemannRetryAttempts = RIEMANN_RETRY_ATTEMPTS
    cfg.RiemannEndpoints = []string(riemannEndpoints)

    if cfg.Hostname == "" {
        cfg.Hostname = GetHostname()
    }

    if !cfg.Check() {
        UsageExit()
    }
    return cfg
}

func GetConfig() *Config {
    if config == nil {
        config = GetNewConfig()
    }
    return config
}
