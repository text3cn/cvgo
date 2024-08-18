package types

// Redis 配置
type RedisConfig struct {
	DefaultConn string `mapstructure:"default_conn"`
	Host        string
	Port        int
	Auth        string
	Db          int
}

// swagger 配置
type SwaggerConfig struct {
	FilePath string `mapstructure:"filepath"`
}

// discovery
type EtcdConfig struct {
	DiscoveryIntervalSeconds int `mapstructure:"discovery_interval_seconds"`
	Server                   DiscoveryServer
	Client                   DiscoveryClient
}

type DiscoveryServer struct {
	Endpoints         []string `mapstructure:"endpoints"`
	DialTimeoutSecods int      `mapstructure:"dial_timeout_secods"`
}

type DiscoveryClient struct {
	ServiceName string `mapstructure:"service_name"`
	ServiceAddr string `mapstructure:"service_addr"`
}
