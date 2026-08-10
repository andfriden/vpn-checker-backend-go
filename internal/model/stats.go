package model

import "time"

type Stats struct {
	Total   int `json:"total"`
	Checked int `json:"checked"`
	Working int `json:"working"`
	Failed  int `json:"failed"`

	UpdatedAt time.Time `json:"updated_at"`
}
