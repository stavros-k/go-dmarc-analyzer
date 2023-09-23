package main

import (
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/providers"
	"github.com/stavros-k/go-dmarc-analyzer/internal/server"
)

func main() {
	// Create a new storage
	store, err := database.NewSqliteStorage("dmarc.db")
	if err != nil {
		panic(err)
	}
	// Migrate the database
	store.Migrate()

	// Create reports provider
	fileProvider := providers.NewFileProvider("reports", "reports/failed", "reports/processed", store)
	// Process reports
	go fileProvider.ProcessAll()

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
