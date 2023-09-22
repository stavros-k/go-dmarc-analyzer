package main

import (
	"encoding/xml"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := gorm.Open(sqlite.Open("dmarc.db"), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}
	store := database.NewSqliteStorage(db)
	store.Migrate()

	// s := server.NewAPIServer("localhost", 8080, store)

	files, err := os.ReadDir("reports")
	if err != nil {
		panic(err)
	}

	waitch := make(chan struct{}, len(files))
	for _, file := range files {
		go func(file os.DirEntry, waitch chan struct{}) {
			log.Infof("Processing file %s", file.Name())
			f, err := os.ReadFile("reports/" + file.Name())
			if err != nil {
				panic(err)
			}

			report := &types.Report{}
			if err := xml.Unmarshal(f, &report); err != nil {
				log.Errorf("Failed to unmarshal report: %s", err)
			}

			store.CreateReport(report)
			waitch <- struct{}{}
		}(file, waitch)
	}

	for i := 0; i < len(files); i++ {
		<-waitch
	}

	// s.RegisterRoutesAndStart()
}
