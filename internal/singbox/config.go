package singbox

type Config struct {
	Log       LogConfig        `json:"log"`
	DNS       *DNSConfig       `json:"dns,omitempty"`
	Inbounds  []InboundConfig  `json:"inbounds"`
	Outbounds []OutboundConfig `json:"outbounds"`
	Route     RouteConfig      `json:"route"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type DNSConfig struct {
	Servers []DNSServer `json:"servers"`
}

type DNSServer struct {
	Type   string `json:"type"`
	Tag    string `json:"tag"`
	Server string `json:"server,omitempty"`
	Detour string `json:"detour,omitempty"`
}

type InboundConfig struct {
	Type       string `json:"type"`
	Tag        string `json:"tag"`
	Listen     string `json:"listen"`
	ListenPort int    `json:"listen_port"`
}

type OutboundConfig struct {
	Type string `json:"type"`
	Tag  string `json:"tag"`

	Server     string `json:"server,omitempty"`
	ServerPort int    `json:"server_port,omitempty"`

	UUID     string `json:"uuid,omitempty"`
	Password string `json:"password,omitempty"`

	UpMbps   int `json:"up_mbps,omitempty"`
	DownMbps int `json:"down_mbps,omitempty"`

	TLS *TLSConfig `json:"tls,omitempty"`

	Transport *TransportConfig `json:"transport,omitempty"`
}

type TLSConfig struct {
	Enabled bool `json:"enabled"`

	ServerName string `json:"server_name,omitempty"`

	Insecure bool `json:"insecure,omitempty"`

	UTLS *UTLSConfig `json:"utls,omitempty"`

	Reality *RealityConfig `json:"reality,omitempty"`
}

type UTLSConfig struct {
	Enabled bool `json:"enabled"`

	Fingerprint string `json:"fingerprint,omitempty"`
}

type RealityConfig struct {
	Enabled bool `json:"enabled"`

	PublicKey string `json:"public_key,omitempty"`

	ShortID string `json:"short_id,omitempty"`
}

type TransportConfig struct {
	Type string `json:"type,omitempty"`

	Path string `json:"path,omitempty"`

	Headers map[string]string `json:"headers,omitempty"`

	ServiceName string `json:"service_name,omitempty"`
}

type RouteConfig struct {
	Final string `json:"final"`

	DefaultDomainResolver string `json:"default_domain_resolver,omitempty"`
}
