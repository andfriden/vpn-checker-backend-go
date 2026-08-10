package api

import (
	"net/http"
	"path/filepath"
)

func (h *Handler) Export(
	w http.ResponseWriter,
	r *http.Request,
) {

	file := filepath.Join(
		h.storage.Dir,
		"all-working.txt",
	)

	w.Header().Set(
		"Content-Type",
		"text/plain; charset=utf-8",
	)

	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=all-working.txt",
	)

	http.ServeFile(
		w,
		r,
		file,
	)
}
