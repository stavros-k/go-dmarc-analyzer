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
	ReportChannel        chan *reportChannel
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
		ReportChannel:        make(chan *reportChannel, 10),
		store:                store,
	}
}

func (f *FileProvider) Process() {
	files, err := filepath.Glob(f.ReportsPath + "/*.xml")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
			log.Errorf("Failed to read file %s: %s", file, err)
		}
		report, err := parsers.NewReport(data)
		if err != nil {
			os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
			log.Errorf("Failed to parse report %s: %s", file, err)
			continue
		}
		f.ReportChannel <- &reportChannel{
			report: report,
			file:   file,
		}
	}
}

func (f *FileProvider) Save() {
	for data := range f.ReportChannel {
		log.Infof("Saving report %s", data.report.ReportMetadata.ReportID)

		if err := f.store.CreateReport(data.report); err != nil {
			log.Errorf("Failed to save report %s: %s", data.report.ReportMetadata.ReportID, err)
			os.Rename(data.file, f.FailedReportsPath+"/"+filepath.Base(data.file))
			continue
		}
		log.Infof("Saved report %s", data.report.ReportMetadata.ReportID)
		os.Rename(data.file, f.ProcessedReportsPath+"/"+filepath.Base(data.file))
	}
}
