package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(
	w http.ResponseWriter,
	value interface{},
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(value)
}
