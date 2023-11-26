// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email customerimporter along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"teamwork/logger"
)

type DomainProcessor struct {
	sync.WaitGroup
	sync.Mutex
	filePath string
	log      logger.Logger
}

func NewDomainProcessor(logger *logger.Logger, filePath string) *DomainProcessor {
	return &DomainProcessor{
		log:      *logger,
		filePath: filePath,
	}
}

// extractDomain extract domain from en email address
func extractDomain(email string) string {
	getDomain := strings.Split(email, "@")

	if len(getDomain) == 2 {
		return getDomain[1]
	}

	return ""
}

// checkHeader check headers for en email header
func (d *DomainProcessor) checkHeaders() (int, [][]string, error) {
	emailIndex := -1

	f, err := os.Open(d.filePath)
	if err != nil {
		d.log.Errorf("unable to read file: %v", err)
		return emailIndex, nil, errors.New("unable to read file")
	}
	defer f.Close()

	row, err := csv.NewReader(f).Read()
	if err != nil {
		d.log.Errorf("unable to parse file as CSV: %v", err)
		return emailIndex, nil, errors.New("unable to parse file as CSV")
	}

	for i, header := range row {
		if header == "email" {
			emailIndex = i
			break
		}
	}

	if emailIndex == -1 {
		d.log.Errorf("email column not found in the CSV file: %v", err)
		return emailIndex, nil, errors.New("email column not found in the CSV file")
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		d.log.Errorf("unable to reset file pointe: %v", err)
		return emailIndex, nil, errors.New("unable to reset file pointer")
	}

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		d.log.Errorf("unable to parse file as CSV: %v", err)
		return emailIndex, nil, errors.New("unable to parse file as CSV")
	}

	return emailIndex, records, nil
}

func processRows(emailIndex int, row []string) string {
	if len(row) > emailIndex {
		return extractDomain(row[emailIndex])
	}

	return ""
}

func (d *DomainProcessor) GetDomains() {
	emailIndex, records, err := d.checkHeaders()
	if err != nil {
		d.log.Errorf("error while checking file headers: %v", err)
		return
	}

	totalWorker := int64(10)

	wp := NewWorkerPool(totalWorker)
	wp.Run()

	results := make(chan string, 100)
	reportingDone := make(chan struct{})

	report := NewReport(d.log)

	go func() {
		for domain := range results {
			report.Add(domain)
		}

		reportingDone <- struct{}{}
	}()

	for _, row := range records[1:] {
		wp.AddTask(func() {
			results <- processRows(emailIndex, row)
		})
	}

	wp.Close()
	<-wp.Done()
	close(results)
	<-reportingDone

	report.Print()
}
