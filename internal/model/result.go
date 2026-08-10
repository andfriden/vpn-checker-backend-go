package model

import "time"

type CheckResult struct {
	Config *VPNConfig `json:"config"`

	Success bool `json:"success"`

	IP string `json:"ip,omitempty"`

	Latency time.Duration `json:"latency"`

	Error string `json:"error,omitempty"`
}
