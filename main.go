package main

import (
	"encoding/xml"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/server"
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

	s := server.NewAPIServer("localhost", 8080, store)

	_, err = os.Stat("reports")
	if os.IsNotExist(err) {
		err = os.Mkdir("processed", 0755)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		files, err := os.ReadDir("reports")
		if err != nil {
			panic(err)
		}

		filesLen := len(files)
		for idx, file := range files {
			log.Infof("[%d/%d] Processing file %s", idx+1, filesLen, file.Name())
			f, err := os.ReadFile("reports/" + file.Name())
			if err != nil {
				panic(err)
			}

			report := &types.Report{}
			if err := xml.Unmarshal(f, &report); err != nil {
				log.Errorf("Failed to unmarshal report: %s", err)
			}

			store.CreateReport(report)

			err = os.Rename("reports/"+file.Name(), "processed/"+file.Name())
			if err != nil {
				log.Errorf("Failed to move file %s to processed: %s", file.Name(), err)
			}
		}
	}()

	go func() {
		records, err := store.FindRecords()
		if err != nil {
			panic(err)
		}

		recLen := len(records)
		for idx, record := range records {
			log.Infof("[%d/%d] Creating address %s", idx+1, recLen, record.Row.SourceIP)
			store.CreateAddress(&types.Address{IP: record.Row.SourceIP})
		}
	}()

	s.RegisterRoutesAndStart()
}
