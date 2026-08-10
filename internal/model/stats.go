package model

import "time"

type ProtocolStats struct {
	Total   int `json:"total"`
	Working int `json:"working"`
	Failed  int `json:"failed"`
}

type Stats struct {
	Total   int `json:"total"`
	Checked int `json:"checked"`
	Working int `json:"working"`
	Failed  int `json:"failed"`

	Protocols map[string]ProtocolStats `json:"protocols"`

	UpdatedAt time.Time `json:"updated_at"`
}
