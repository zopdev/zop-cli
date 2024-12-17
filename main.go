package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/service"
	_ "modernc.org/sqlite"

	applicationHandler "zop.dev/cli/zop/application/handler"
	applicationSvc "zop.dev/cli/zop/application/service"
	impHandler "zop.dev/cli/zop/cloud/handler"
	impService "zop.dev/cli/zop/cloud/service/gcp"
	listSvc "zop.dev/cli/zop/cloud/service/list"
	impStore "zop.dev/cli/zop/cloud/store/gcp"
	envHandler "zop.dev/cli/zop/environment/handler"
	envService "zop.dev/cli/zop/environment/service"
)

const (
	//nolint:gosec //This is a tokenURL for google oauth2
	tokenURL = "https://oauth2.googleapis.com"
)

func main() {
	app := gofr.NewCMD()

	app.AddHTTPService(impService.ZopAPIService, app.Config.Get("ZOP_API_URL"))
	app.AddHTTPService(impService.GcloudService, tokenURL,
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

	accStore := impStore.New(db)
	accSvc := impService.New(accStore)
	lSvc := listSvc.New()
	h := impHandler.New(accSvc, lSvc)

	app.SubCommand("cloud import", h.Import)
	app.SubCommand("cloud list", h.List)

	appSvc := applicationSvc.New()
	appH := applicationHandler.New(appSvc)

	app.SubCommand("application add", appH.Add)
	app.SubCommand("application list", appH.List)

	envSvc := envService.New(appSvc)
	envH := envHandler.New(envSvc)

	app.SubCommand("environment add", envH.Add)

	app.Run()
}
