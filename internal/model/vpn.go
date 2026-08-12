package model

import "net/url"

type Protocol string

const (
	ProtocolVLESS     Protocol = "vless"
	ProtocolVMess     Protocol = "vmess"
	ProtocolTrojan    Protocol = "trojan"
	ProtocolHysteria2 Protocol = "hysteria2"
)

type VPNConfig struct {
	Raw      string   `json:"raw"`
	Protocol Protocol `json:"protocol"`

	Source string `json:"source,omitempty"`

	Address string `json:"address"`
	Port    int    `json:"port"`

	UUID     string `json:"uuid,omitempty"`
	Password string `json:"password,omitempty"`

	Name string `json:"name,omitempty"`

	TLS       TLSConfig       `json:"tls"`
	Transport TransportConfig `json:"transport"`

	UpMbps   int `json:"up_mbps,omitempty"`
	DownMbps int `json:"down_mbps,omitempty"`
}

type TLSConfig struct {
	Enabled     bool   `json:"enabled"`
	ServerName  string `json:"server_name,omitempty"`
	Insecure    bool   `json:"insecure,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	PublicKey   string `json:"public_key,omitempty"`
	ShortID     string `json:"short_id,omitempty"`
}

type TransportConfig struct {
	Type    string `json:"type,omitempty"`
	Path    string `json:"path,omitempty"`
	Host    string `json:"host,omitempty"`
	Service string `json:"service,omitempty"`
}

func (c *VPNConfig) URL() (*url.URL, error) {
	return url.Parse(c.Raw)
}
