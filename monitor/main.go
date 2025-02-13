package main

import (
	// "fmt"

	"database/sql"
	"log"

	"github.com/marcorentap/hallucinet/config"
	"github.com/marcorentap/hallucinet/types"
	"github.com/marcorentap/hallucinet/watcher"
	"github.com/marcorentap/hallucinet/web"
	"github.com/marcorentap/hallucinet/writer"
)

func main() {
	config := config.NewHallucinetConfig()

	db, err := sql.Open("sqlite3", config.SqlitePath)
	if err != nil {
		log.Fatalf("Cannot open hallucinet database %v: %v\n", config.SqlitePath, err)
	}
	defer db.Close()
	writer.InitializeDB(db)

	hctx := types.HallucinetContext{
		Config:    config,
		EventChan: make(chan types.HallucinetEvent),
		DB:        db,
	}

	go watcher.WatchDockerEvents(hctx)
	go writer.UpdateDatabase(hctx)
	go writer.UpdateHosts(hctx)
	web.Serve(hctx)
}
