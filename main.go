package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gofr.dev/pkg/gofr"

	impHandler "zop/handler/cloud/import"
	impService "zop/service/cloud/import/gcp"
	impStore "zop/store/cloud/import/gcp"
)

func main() {
	app := gofr.NewCMD()

	app.AddHTTPService("api-service", app.Config.Get("ZOP_API_URL"))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get the user's home directory: %v", err)
	}

	// Build the path to the credentials.db file
	dbPath := filepath.Join(homeDir, ".config", "gcloud", "credentials.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	accountStore := impStore.New(db)
	accountSvc := impService.New(accountStore)
	importHandler := impHandler.New(accountSvc)

	app.SubCommand("cloud import", importHandler.Import)

	app.Run()
}
