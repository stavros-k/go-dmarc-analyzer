package main

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/inputs"
	"github.com/stavros-k/go-dmarc-analyzer/internal/server"
)

var directories = []string{"reports"}

const processFileAtBoot = false
const processFileInterval = time.Second * 30

func main() {
	// Create a new storage
	store, err := database.NewSqliteStorage("dmarc.db")
	if err != nil {
		panic(err)
	}
	// Migrate the database
	store.Migrate()

	inputers := []inputs.Inputer{}

	// Create file inputer(s)
	for _, dir := range directories {
		p, err := inputs.NewFileInput(dir, store)

		if err != nil {
			log.Errorf("Failed to create provider for directory %s: %s", dir, err)
			continue
		}
		inputers = append(inputers, p)
	}

	// Start processing
	for _, p := range inputers {
		switch p.(type) {
		case *inputs.FileInput:
			if processFileAtBoot {
				go p.ProcessAll()
			}
			go p.Watch(processFileInterval)
		}
	}

	// go func() {
	// 	records, err := store.FindRecords()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	recLen := len(records)
	// 	for idx, record := range records {
	// 		log.Infof("[%d/%d] Creating address %s", idx+1, recLen, record.Row.SourceIP)
	// 		store.CreateAddress(&types.Address{IP: record.Row.SourceIP})
	// 	}
	// }()

	s := server.NewAPIServer("localhost", 8080, store)
	s.RegisterRoutesAndStart()
}
