package database_sqlite

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/parsers"
	"gorm.io/gorm"
)

type ReportModel struct {
	ReportID                               string `gorm:"primaryKey"`
	CreatedAt                              int64  `gorm:"autoCreateTime"`
	Version                                string
	ReportMetadataOrgName                  string
	ReportMetadataEmail                    string
	ReportMetadataExtraContactInfo         string
	ReportDateRangeBegin                   time.Time
	ReportDateRangeEnd                     time.Time
	PolicyPublishedDomain                  string
	PolicyPublishedAlignmentModeDKIM       string
	PolicyPublishedAlignmentModeSPF        string
	PolicyPublishedPolicy                  string
	PolicyPublishedSubdomainPolicy         string
	PolicyPublishedPercentage              int
	PolicyPublishedFailureReportingOptions rune
}

func (s *SqliteStorage) CreateReport(report *parsers.Report) error {
	r := ReportToModel(report)

	err := s.db.Create(r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Infof("Report with ID %s already exists, skipping", report.ReportMetadata.ReportID)
			return nil
		}

		return err
	}

	for _, record := range report.Records {
		if err := s.CreateReportRecord(report.ReportMetadata.ReportID, &record); err != nil {
			return err
		}
	}

	return nil
}

func (s *SqliteStorage) FindReportByReportID(reportID string) (*parsers.Report, error) {
	report := &ReportModel{}

	if err := s.db.Where("report_id = ?", reportID).First(report).Error; err != nil {
		return nil, err
	}

	records, err := s.FindRecordsByReportID(reportID)
	if err != nil {
		return nil, err
	}

	// Convert the ReportModel to a parsers.Report
	return ModelToReport(report, records), nil
}

func (s *SqliteStorage) FindReports() ([]*parsers.Report, error) {
	var reports []*parsers.Report

	if err := s.db.Find(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}

// Converts a parsers.Report to a ReportModel
func ReportToModel(r *parsers.Report) *ReportModel {
	return &ReportModel{
		ReportID:                               r.ReportMetadata.ReportID,
		Version:                                r.Version,
		ReportMetadataOrgName:                  r.ReportMetadata.OrgName,
		ReportMetadataEmail:                    r.ReportMetadata.Email,
		ReportMetadataExtraContactInfo:         r.ReportMetadata.ExtraContactInfo,
		ReportDateRangeBegin:                   time.Unix(r.ReportMetadata.DateRange.Begin, 0).UTC(),
		ReportDateRangeEnd:                     time.Unix(r.ReportMetadata.DateRange.End, 0).UTC(),
		PolicyPublishedDomain:                  r.PolicyPublished.Domain,
		PolicyPublishedAlignmentModeDKIM:       r.PolicyPublished.AlignmentModeDKIM,
		PolicyPublishedAlignmentModeSPF:        r.PolicyPublished.AlignmentModeSPF,
		PolicyPublishedPolicy:                  r.PolicyPublished.Policy,
		PolicyPublishedSubdomainPolicy:         r.PolicyPublished.SubdomainPolicy,
		PolicyPublishedPercentage:              r.PolicyPublished.Percentage,
		PolicyPublishedFailureReportingOptions: r.PolicyPublished.FailureReportingOptions,
	}
}

// Converts a ReportModel to a parsers.Report
func ModelToReport(r *ReportModel, recs []*parsers.Record) *parsers.Report {
	records := make([]parsers.Record, len(recs))
	for idx, rec := range recs {
		records[idx] = *rec
	}

	return &parsers.Report{
		Version: r.Version,
		Records: records,
		ReportMetadata: parsers.ReportMetadata{
			OrgName:          r.ReportMetadataOrgName,
			Email:            r.ReportMetadataEmail,
			ExtraContactInfo: r.ReportMetadataExtraContactInfo,
			ReportID:         r.ReportID,
			DateRange: parsers.DateRange{
				Begin: r.ReportDateRangeBegin.Unix(),
				End:   r.ReportDateRangeEnd.Unix(),
			},
		},
		PolicyPublished: parsers.PolicyPublished{
			Domain:                  r.PolicyPublishedDomain,
			AlignmentModeDKIM:       r.PolicyPublishedAlignmentModeDKIM,
			AlignmentModeSPF:        r.PolicyPublishedAlignmentModeSPF,
			Policy:                  r.PolicyPublishedPolicy,
			SubdomainPolicy:         r.PolicyPublishedSubdomainPolicy,
			Percentage:              r.PolicyPublishedPercentage,
			FailureReportingOptions: r.PolicyPublishedFailureReportingOptions,
		},
	}
}
