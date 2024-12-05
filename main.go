package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gofr.dev/pkg/gofr"

	impHandler "github.com/zopdev/zop-cli/handler/cloud/import"
	impStore "github.com/zopdev/zop-cli/store/cloud/import"
)

func main() {
	app := gofr.NewCMD()

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
	importHandler := impHandler.New(accountStore)

	app.SubCommand("cloud import", importHandler.Import)

	app.Run()
}
