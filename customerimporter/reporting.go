package customerimporter

import (
	"sort"
	"teamwork/logger"
)

type DomainCounter struct {
	Domain string
	Users  int
}

type Report struct {
	domainCustomerCount map[string]int
	RowsProcessed       int
	logger              logger.Logger
}

func NewReport(logger logger.Logger) *Report {
	return &Report{
		domainCustomerCount: make(map[string]int),
		RowsProcessed:       0,
		logger:              logger,
	}
}

func (r *Report) Add(domain string) {
	r.RowsProcessed++
	r.domainCustomerCount[domain]++
}

func (r *Report) Sort() []DomainCounter {
	var results []DomainCounter
	var domains []string

	for domain := range r.domainCustomerCount {
		domains = append(domains, domain)
	}

	sort.Strings(domains)

	for _, domain := range domains {
		users := r.domainCustomerCount[domain]
		results = append(results, DomainCounter{
			Domain: domain,
			Users:  users,
		})
	}

	return results
}

func (r *Report) Print() {
	r.logger.Infof("Rows processed: %v", r.RowsProcessed)
	results := r.Sort()
	r.logger.Infof("Finished processing. Domain counts: %+v", results)
}
