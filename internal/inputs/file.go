package inputs

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/parsers"
)

type FileInput struct {
	ReportsPath          string
	FailedReportsPath    string
	ProcessedReportsPath string
	store                database.Storage
	mutexProcess         sync.Mutex
}

// NewFileInput creates a new FileInput
func NewFileInput(reportsPath string, store database.Storage) (*FileInput, error) {
	for _, path := range []string{reportsPath, reportsPath + "/failed", reportsPath + "/processed"} {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				log.Errorf("Failed to create directory %s: %s", path, err)
				return nil, err
			}
		}
	}

	path := strings.TrimSuffix(reportsPath, "/")

	return &FileInput{
		ReportsPath:          path,
		FailedReportsPath:    path + "/failed",
		ProcessedReportsPath: path + "/processed",
		store:                store,
		mutexProcess:         sync.Mutex{},
	}, nil
}

// Watcher watches the reports directory for new files
// and processes them on a given interval
func (f *FileInput) Watch(interval time.Duration) {
	for {
		f.ProcessAll()
		time.Sleep(interval)
	}
}

// StoreReport takes a byte slice of a report and stores it in the database
func (f *FileInput) StoreReport(data []byte) error {
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

// ProcessAll processes all the reports in the reports directory
func (f *FileInput) ProcessAll() {
	// We lock the process mutex to avoid multiple calls to ProcessAll
	// Multiple calls to ProcessAll can happen if the watcher interval is
	// less than the time it takes to process all the reports
	// Or the initial call to ProcessAll that runs on startup
	// is still running when the watcher interval is reached
	f.mutexProcess.Lock()
	defer f.mutexProcess.Unlock()

	files, err := filepath.Glob(f.ReportsPath + "/*.xml")
	if err != nil {
		log.Errorf("Failed to glob reports path %s: %s", f.ReportsPath, err)
	}

	for _, file := range files {
		if err := f.Process(file); err != nil {
			log.Errorf("Failed to process file %s: %s", file, err)
		}
	}
}

// Process processes a single report file
func (f *FileInput) Process(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
		log.Errorf("Failed to read file %s: %s", file, err)
		return err
	}
	if err := f.StoreReport(data); err != nil {
		log.Errorf("Failed to store file %s: %s", file, err)
		os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
		return err
	}

	os.Rename(file, f.ProcessedReportsPath+"/"+filepath.Base(file))
	return nil
}
