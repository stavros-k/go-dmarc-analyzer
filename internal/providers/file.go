package providers

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
)

type FileProvider struct {
	ReportsPath          string
	FailedReportsPath    string
	ProcessedReportsPath string
	store                database.Storage
	mutexProcess         sync.Mutex
}

func NewFileProvider(reportsPath string, store database.Storage) *FileProvider {
	for _, path := range []string{reportsPath, reportsPath + "/failed", reportsPath + "/processed"} {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				panic(err)
			}
		}
	}

	path := strings.TrimSuffix(reportsPath, "/")

	return &FileProvider{
		ReportsPath:          path,
		FailedReportsPath:    path + "/failed",
		ProcessedReportsPath: path + "/processed",
		store:                store,
		mutexProcess:         sync.Mutex{},
	}
}

func (f *FileProvider) ProcessAll() {
	// Lock the process mutex to avoid initial process and watcher to run at the same time
	f.mutexProcess.Lock()
	defer f.mutexProcess.Unlock()

	files, err := filepath.Glob(f.ReportsPath + "/*.xml")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
			log.Errorf("Failed to read file %s: %s", file, err)
			continue
		}
		if err := f.StoreReport(data); err != nil {
			log.Errorf("Failed to store file %s: %s", file, err)
			os.Rename(file, f.FailedReportsPath+"/"+filepath.Base(file))
			continue
		}

		os.Rename(file, f.ProcessedReportsPath+"/"+filepath.Base(file))
	}
}

func (f *FileProvider) Watcher(interval time.Duration) {
	for {
		f.ProcessAll()
		time.Sleep(interval)
	}
}
