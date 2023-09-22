package database

import (
	"github.com/stavros-k/go-dmarc-analyzer/internal/types"
	"gorm.io/gorm"
)

type Storage interface {
	CreateReport(*types.Report) error
	FindReportByReportID(string) (*types.Report, error)
	FindReports() ([]*types.Report, error)
	CreateReportRecord(string, *types.Record) error
	FindRecordsByReportID(string) ([]*types.Record, error)
}

type SqliteStorage struct {
	db *gorm.DB
}

func NewSqliteStorage(db *gorm.DB) *SqliteStorage {
	return &SqliteStorage{
		db: db,
	}
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
