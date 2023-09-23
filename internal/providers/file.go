package providers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/parsers"
)

type reportChannel struct {
	report *parsers.Report
	file   string
}

type FileProvider struct {
	ReportsPath          string
	FailedReportsPath    string
	ProcessedReportsPath string
	store                database.Storage
}

func NewFileProvider(reportsPath, failedReportsPath, processedReportsPath string, store database.Storage) *FileProvider {
	for _, path := range []string{reportsPath, failedReportsPath, processedReportsPath} {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				panic(err)
			}
		}
	}

	return &FileProvider{
		ReportsPath:          strings.TrimSuffix(reportsPath, "/"),
		FailedReportsPath:    strings.TrimSuffix(failedReportsPath, "/"),
		ProcessedReportsPath: strings.TrimSuffix(processedReportsPath, "/"),
		store:                store,
	}
}

func (f *FileProvider) ProcessAll() {
	files, err := filepath.Glob(f.ReportsPath + "/*.xml")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		f.Store(file)
	}
}

func (f *FileProvider) Store(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
		log.Errorf("Failed to read file %s: %s", file, err)
	}

	report, err := parsers.NewReport(data)
	if err != nil {
		os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
		log.Errorf("Failed to parse report %s: %s", file, err)
		return
	}

	log.Infof("Saving report %s", report.ReportMetadata.ReportID)

	if err := f.store.CreateReport(report); err != nil {
		log.Errorf("Failed to save report %s: %s", report.ReportMetadata.ReportID, err)
		os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
		return
	}

	log.Infof("Saved report %s", report.ReportMetadata.ReportID)
	os.Rename(file, f.ProcessedReportsPath+"/"+filepath.Base(file))
}
