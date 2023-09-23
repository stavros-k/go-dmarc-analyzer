package database

import (
	"time"

	"github.com/stavros-k/go-dmarc-analyzer/internal/parsers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage interface {
	CreateReport(*parsers.Report) error
	FindReportByReportID(string) (*parsers.Report, error)
	FindReports() ([]*parsers.Report, error)
	CreateReportRecord(string, *parsers.Record) error
	FindRecordsByReportID(string) ([]*parsers.Record, error)
}

type SqliteStorage struct {
	db *gorm.DB
}

func NewSqliteStorage(dbPath string) (*SqliteStorage, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
		NowFunc:        func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return nil, err
	}

	return &SqliteStorage{
		db: db,
	}, nil
}
func (s *SqliteStorage) Migrate() error {
	models := []interface{}{
		&ReportModel{},
		&ReportRecordModel{},
		&AddressModel{},
	}

	for _, model := range models {
		if err := s.db.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}
