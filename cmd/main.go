package main

import (
	"teamwork/customerimporter"
	"teamwork/logger"
)

func main() {
	log := logger.NewLogger()
	domainProcessor := customerimporter.NewDomainProcessor(log, "./data/customers.csv")
	domainProcessor.GetDomains()
}
