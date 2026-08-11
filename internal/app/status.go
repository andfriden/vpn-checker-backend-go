package app

import "time"

type Status struct {
	Status string `json:"status"`

	Total   int `json:"total"`
	Checked int `json:"checked"`

	Working int `json:"working"`
	Failed  int `json:"failed"`

	Progress int `json:"progress"`

	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`

	ElapsedSeconds float64 `json:"elapsed_seconds"`

	CurrentSpeed float64 `json:"current_speed"`

	EstimatedSecondsLeft float64 `json:"estimated_seconds_left"`
}
