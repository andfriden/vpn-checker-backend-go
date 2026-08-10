package worker

import (
	"sync"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)


type Job struct {
	Config model.VPNConfig
}


type Result struct {
	Config model.VPNConfig
	Result model.CheckResult
}


type Pool struct {
	Workers int
}


func New(workers int) *Pool {

	return &Pool{
		Workers: workers,
	}
}



func (p *Pool) Run(
	configs []model.VPNConfig,
	check func(model.VPNConfig) model.CheckResult,
) []model.CheckResult {


	jobs := make(chan model.VPNConfig)

	results := make(chan model.CheckResult)


	var wg sync.WaitGroup


	for i := 0; i < p.Workers; i++ {

		wg.Add(1)


		go func() {

			defer wg.Done()


			for cfg := range jobs {

				result := check(cfg)

				results <- result
			}


		}()

	}



	go func(){

		for _,cfg := range configs {

			jobs <- cfg

		}


		close(jobs)

	}()



	go func(){

		wg.Wait()

		close(results)

	}()



	var output []model.CheckResult


	for r := range results {

		output = append(output,r)

	}


	return output
}
