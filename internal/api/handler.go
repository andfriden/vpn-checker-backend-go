package api

import (
	"github.com/andfriden/vpn-checker-backend-go/internal/app"
	"github.com/andfriden/vpn-checker-backend-go/internal/service"
)

type Handler struct {
	runner *app.Runner

	resultsService *service.ResultsService
}

func NewHandler(
	runner *app.Runner,
	resultsService *service.ResultsService,
) *Handler {

	return &Handler{
		runner:         runner,
		resultsService: resultsService,
	}
}
