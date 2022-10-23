package report

import (
	"balance_api/internal/entity"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// BalanceReport keeps a report dir name
type BalanceReport struct {
	reportDir string
}

// New creates new dir and constructs BalanceReport
func New(d string) (*BalanceReport, error) {
	err := os.Mkdir(strings.Trim(d, "/"), 0750)
	if err != nil {
		return nil, err
	}
	return &BalanceReport{
		reportDir: d,
	}, nil
}

// GetDir is a getter for reportDir field of BalanceReport
func (r *BalanceReport) GetDir() string {
	return r.reportDir
}

// Create writes entity.Report to a csv file
func (r *BalanceReport) Create(ctx context.Context, name string, report entity.Report) (string, error) {
	name = name + ".csv"
	file, err := os.Create(r.reportDir + name)
	if err != nil {
		return "", fmt.Errorf("ReportFile - Create: %w", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	for _, v := range report.Sums {
		line := []string{v.Name, v.Sum}
		err = w.Write(line)
		if err != nil {
			return "", fmt.Errorf("ReportFile - Create: %w", err)
		}
		w.Flush()
	}
	return name, nil
}
