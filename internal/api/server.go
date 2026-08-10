package api

import (
	"fmt"
	"net/http"
)

func Start(
	host string,
	port int,
	handler http.Handler,
) error {

	addr := fmt.Sprintf(
		"%s:%d",
		host,
		port,
	)

	return http.ListenAndServe(
		addr,
		handler,
	)
}
