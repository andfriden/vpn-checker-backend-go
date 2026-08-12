package service

import (
	"path/filepath"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
	"github.com/andfriden/vpn-checker-backend-go/internal/storage"
)

type ResultsService struct {
	storage *storage.FileStorage
}

func NewResultsService(storage *storage.FileStorage) *ResultsService {
	return &ResultsService{
		storage: storage,
	}
}

func (s *ResultsService) All() ([]model.CheckResult, error) {
	return s.storage.LoadResults()
}

type ResultsQuery struct {
	Working  string
	Protocol string
	Sort     string
	Page     int
	Limit    int
}

type ResultsResponse struct {
	Data       []ResultItem `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

type ResultItem struct {
	Protocol string `json:"protocol,omitempty"`
	Address  string `json:"address,omitempty"`
	Port     int    `json:"port,omitempty"`
	Latency  int64  `json:"latency_ms"`
	IP       string `json:"ip,omitempty"`
	Success  bool   `json:"success"`
	Source   string `json:"source,omitempty"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

func (s *ResultsService) List(
	query ResultsQuery,
) (*ResultsResponse, error) {

	results, err := s.storage.LoadResults()
	if err != nil {
		return nil, err
	}

	results = Filter(
		results,
		query.Working,
		query.Protocol,
	)

	Sort(
		results,
		query.Sort,
	)

	total := len(results)

	page := query.Page
	if page <= 0 {
		page = 1
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}

	pages := 0

	if total > 0 {
		pages = (total + limit - 1) / limit
	}

	if page > pages && pages > 0 {
		page = pages
	}

	start := (page - 1) * limit

	if start > total {
		start = total
	}

	end := start + limit

	if end > total {
		end = total
	}

	return &ResultsResponse{
		Data: convertResults(
			results[start:end],
		),

		Pagination: Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: pages,
		},
	}, nil
}

func (s *ResultsService) Best(
	protocol string,
	limit int,
) ([]model.CheckResult, error) {

	results, err := s.storage.LoadResults()
	if err != nil {
		return nil, err
	}

	results = Filter(
		results,
		"true",
		protocol,
	)

	Sort(
		results,
		"latency",
	)

	results = Deduplicate(
		results,
	)

	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}

	return results, nil
}

func (s *ResultsService) Stats() (interface{}, error) {
	return s.storage.LoadStats()
}

func (s *ResultsService) ExportPath(
	filename string,
) string {

	return filepath.Join(
		s.storage.Dir,
		filename,
	)
}
