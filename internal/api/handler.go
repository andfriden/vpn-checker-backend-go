package api

import (
	"github.com/andfriden/vpn-checker-backend-go/internal/app"
	"github.com/andfriden/vpn-checker-backend-go/internal/service"
)

type Handler struct {
	runner         *app.Runner
	resultsService *service.ResultsService

	blackURL string
	whiteURL string
}

func NewHandler(
	runner *app.Runner,
	resultsService *service.ResultsService,
	blackURL string,
	whiteURL string,
) *Handler {

	return &Handler{
		runner:         runner,
		resultsService: resultsService,
		blackURL:       blackURL,
		whiteURL:       whiteURL,
	}
}
