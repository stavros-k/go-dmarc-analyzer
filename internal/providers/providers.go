package providers

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/parsers"
)

type Provider interface {
	ProcessAll()
	Watcher(interval time.Duration)
}

// StoreReport takes a byte slice of a report and stores it in the database
func (f *FileProvider) StoreReport(data []byte) error {
	report, err := parsers.NewReport(data)
	if err != nil {
		return err
	}

	log.Infof("Saving report %s", report.ReportMetadata.ReportID)

	if err := f.store.CreateReport(report); err != nil {
		log.Errorf("Failed to save report %s: %s", report.ReportMetadata.ReportID, err)
		return err
	}

	log.Infof("Saved report %s", report.ReportMetadata.ReportID)
	return nil
}
