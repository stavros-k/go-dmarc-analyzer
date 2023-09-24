package database_sqlite

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
		// Convert the ReportRecordModel to a parsers.Record
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
		// Convert the ReportRecordModel to a parsers.Record
		records[idx] = ModelToReportRecord(record)
	}

	return records, nil
}

// Converts a parsers.Record to a ReportRecordModel
func ReportRecordToModel(reportID string, rec *parsers.Record) *ReportRecordModel {
	return &ReportRecordModel{
		ReportID:                   reportID,
		SourceIP:                   rec.Row.SourceIP,
		Count:                      rec.Row.Count,
		PolicyEvaluatedDisposition: rec.Row.PolicyEvaluated.Disposition,
		PolicyEvaluatedDKIM:        rec.Row.PolicyEvaluated.DKIM,
		PolicyEvaluatedSPF:         rec.Row.PolicyEvaluated.SPF,
		IdentifiersHeaderFrom:      rec.Identifiers.HeaderFrom,
		IdentifiersEnvelopeFrom:    rec.Identifiers.EnvelopeFrom,
		IdentifiersEnvelopeTo:      rec.Identifiers.EnvelopeTo,
		AuthResultsDKIMDomain:      rec.AuthResults.DKIM.Domain,
		AuthResultsDKIMResult:      rec.AuthResults.DKIM.Result,
		AuthResultsDKIMSelector:    rec.AuthResults.DKIM.Selector,
		AuthResultsDKIMHumanResult: rec.AuthResults.DKIM.HumanResult,
		AuthResultsSPFDomain:       rec.AuthResults.SPF.Domain,
		AuthResultsSPFResult:       rec.AuthResults.SPF.Result,
		AuthResultsSPFScope:        rec.AuthResults.SPF.Scope,
		AuthResultsSPFHumanResult:  rec.AuthResults.SPF.HumanResult,
	}
}

// Converts a ReportRecordModel to a parsers.Record
func ModelToReportRecord(r *ReportRecordModel) *parsers.Record {
	return &parsers.Record{
		Row: parsers.Row{
			SourceIP: r.SourceIP,
			Count:    r.Count,
			PolicyEvaluated: parsers.PolicyEvaluated{
				Disposition: r.PolicyEvaluatedDisposition,
				DKIM:        r.PolicyEvaluatedDKIM,
				SPF:         r.PolicyEvaluatedSPF,
			},
		},
		Identifiers: parsers.Identifiers{
			HeaderFrom:   r.IdentifiersHeaderFrom,
			EnvelopeFrom: r.IdentifiersEnvelopeFrom,
			EnvelopeTo:   r.IdentifiersEnvelopeTo,
		},
		AuthResults: parsers.AuthResult{
			DKIM: parsers.DKIMAuthResult{
				Domain:      r.AuthResultsDKIMDomain,
				Result:      r.AuthResultsDKIMResult,
				Selector:    r.AuthResultsDKIMSelector,
				HumanResult: r.AuthResultsDKIMHumanResult,
			},
			SPF: parsers.SPFAuthResult{
				Domain:      r.AuthResultsSPFDomain,
				Result:      r.AuthResultsSPFResult,
				Scope:       r.AuthResultsSPFScope,
				HumanResult: r.AuthResultsSPFHumanResult,
			},
		},
	}
}
