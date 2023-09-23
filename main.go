package main

import (
	"time"

	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/providers"
	"github.com/stavros-k/go-dmarc-analyzer/internal/server"
)

var directories = []string{"reports"}

func main() {
	// Create a new storage
	store, err := database.NewSqliteStorage("dmarc.db")
	if err != nil {
		panic(err)
	}
	// Migrate the database
	store.Migrate()

	// Create reports provider
	for _, dir := range directories {
		p := providers.NewFileProvider(dir, store)
		go p.ProcessAll()
		go p.Watcher(time.Second * 30)
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
