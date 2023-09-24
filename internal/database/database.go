package database

import (
	"github.com/stavros-k/go-dmarc-analyzer/internal/parsers"
)

type Storage interface {
	Migrate() error
	CreateReport(*parsers.Report) error
	FindReportByReportID(string) (*parsers.Report, error)
	FindReports() ([]*parsers.Report, error)
	CreateReportRecord(string, *parsers.Record) error
	FindRecordsByReportID(string) ([]*parsers.Record, error)
}
