package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/service"

	impHandler "zop.dev/cli/zop/cloud/handler"
	impService "zop.dev/cli/zop/cloud/service/gcp"
	impStore "zop.dev/cli/zop/cloud/store/gcp"
)

const (
	//nolint:gosec //This is a tokenURL for google oauth2
	tokenURL = "https://oauth2.googleapis.com"
)

func main() {
	app := gofr.NewCMD()

	app.AddHTTPService("api-service", app.Config.Get("ZOP_API_URL"))
	app.AddHTTPService("gcloud-service", tokenURL,
		&service.DefaultHeaders{Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"}})

	homeDir, err := os.UserHomeDir()
	if err != nil {
		app.Logger().Fatalf("Failed to get the user's home directory: %v", err)
	}

	// Build the path to the credentials.db file of gcloud
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
