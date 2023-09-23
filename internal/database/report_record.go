package database

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/parsers"
	"gorm.io/gorm"
)

type ReportRecordModel struct {
	ID                         uint   `gorm:"primaryKey"`
	CreatedAt                  int64  `gorm:"autoCreateTime"`
	ReportID                   string `gorm:"foreignKey:ReportID"`
	SourceIP                   string
	Count                      int
	PolicyEvaluatedDisposition string
	PolicyEvaluatedDKIM        string
	PolicyEvaluatedSPF         string
	IdentifiersHeaderFrom      string
	IdentifiersEnvelopeFrom    string
	IdentifiersEnvelopeTo      string
	AuthResultsDKIMDomain      string
	AuthResultsDKIMResult      string
	AuthResultsDKIMSelector    string
	AuthResultsDKIMHumanResult string
	AuthResultsSPFDomain       string
	AuthResultsSPFResult       string
	AuthResultsSPFScope        string
	AuthResultsSPFHumanResult  string
}

func (s *SqliteStorage) CreateReportRecord(reportID string, record *parsers.Record) error {
	r := ReportRecordToModel(reportID, record)
	err := s.db.Create(r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Infof("Record with source IP %s already exists, skipping", record.Row.SourceIP)
			return nil
		}

		return err
	}

	return nil
}

func (s *SqliteStorage) FindRecordsByReportID(reportID string) ([]*parsers.Record, error) {
	reportRecordModels := []*ReportRecordModel{}

	if err := s.db.Where("report_id = ?", reportID).Find(&reportRecordModels).Error; err != nil {
		return nil, err
	}

	records := make([]*parsers.Record, len(reportRecordModels))
	for idx, record := range reportRecordModels {
		records[idx] = ModelToReportRecord(record)
	}

	return records, nil
}

func (s *SqliteStorage) FindRecords() ([]*parsers.Record, error) {
	reportRecordModels := []*ReportRecordModel{}

	if err := s.db.Find(&reportRecordModels).Error; err != nil {
		return nil, err
	}

	records := make([]*parsers.Record, len(reportRecordModels))
	for idx, record := range reportRecordModels {
		records[idx] = ModelToReportRecord(record)
	}

	return records, nil
}

// Converts a parsers.Record to a ReportRecordModel
func ReportRecordToModel(reportID string, record *parsers.Record) *ReportRecordModel {
	return &ReportRecordModel{
		ReportID:                   reportID,
		SourceIP:                   record.Row.SourceIP,
		Count:                      record.Row.Count,
		PolicyEvaluatedDisposition: record.Row.PolicyEvaluated.Disposition,
		PolicyEvaluatedDKIM:        record.Row.PolicyEvaluated.DKIM,
		PolicyEvaluatedSPF:         record.Row.PolicyEvaluated.SPF,
		IdentifiersHeaderFrom:      record.Identifiers.HeaderFrom,
		IdentifiersEnvelopeFrom:    record.Identifiers.EnvelopeFrom,
		IdentifiersEnvelopeTo:      record.Identifiers.EnvelopeTo,
		AuthResultsDKIMDomain:      record.AuthResults.DKIM.Domain,
		AuthResultsDKIMResult:      record.AuthResults.DKIM.Result,
		AuthResultsDKIMSelector:    record.AuthResults.DKIM.Selector,
		AuthResultsDKIMHumanResult: record.AuthResults.DKIM.HumanResult,
		AuthResultsSPFDomain:       record.AuthResults.SPF.Domain,
		AuthResultsSPFResult:       record.AuthResults.SPF.Result,
		AuthResultsSPFScope:        record.AuthResults.SPF.Scope,
		AuthResultsSPFHumanResult:  record.AuthResults.SPF.HumanResult,
	}
}

// Converts a ReportRecordModel to a parsers.Record
func ModelToReportRecord(reportRecordModel *ReportRecordModel) *parsers.Record {
	return &parsers.Record{
		Row: parsers.Row{
			SourceIP: reportRecordModel.SourceIP,
			Count:    reportRecordModel.Count,
			PolicyEvaluated: parsers.PolicyEvaluated{
				Disposition: reportRecordModel.PolicyEvaluatedDisposition,
				DKIM:        reportRecordModel.PolicyEvaluatedDKIM,
				SPF:         reportRecordModel.PolicyEvaluatedSPF,
			},
		},
		Identifiers: parsers.Identifiers{
			HeaderFrom:   reportRecordModel.IdentifiersHeaderFrom,
			EnvelopeFrom: reportRecordModel.IdentifiersEnvelopeFrom,
			EnvelopeTo:   reportRecordModel.IdentifiersEnvelopeTo,
		},
		AuthResults: parsers.AuthResult{
			DKIM: parsers.DKIMAuthResult{
				Domain:      reportRecordModel.AuthResultsDKIMDomain,
				Result:      reportRecordModel.AuthResultsDKIMResult,
				Selector:    reportRecordModel.AuthResultsDKIMSelector,
				HumanResult: reportRecordModel.AuthResultsDKIMHumanResult,
			},
			SPF: parsers.SPFAuthResult{
				Domain:      reportRecordModel.AuthResultsSPFDomain,
				Result:      reportRecordModel.AuthResultsSPFResult,
				Scope:       reportRecordModel.AuthResultsSPFScope,
				HumanResult: reportRecordModel.AuthResultsSPFHumanResult,
			},
		},
	}
}
