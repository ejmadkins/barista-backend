/*
Cymbal Coffee Backend Service

This handles the communication with Cymbal Coffee's bond service.

There is no need to understand this code or how it works to complete the tasks.

You can look through it though if you like!
*/
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Store configuration globally
var cfg config

const (
	defaultPort       = "8080"
	defaultCollection = "bond"
)

type config struct {
	Port      string
	ProjectID string
}

type AppInstance struct {
	ProjectID string `json:"project_id"`
	Team      string `json:"team"`
}

func init() {
	ctx := context.Background()
	initConfig(ctx)
	initBond()
	intro(ctx)
	DDDInit()
}

// Init all our config and everything else
func initConfig(ctx context.Context) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Obtain Project ID from metadata server unless specified
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if projectID == "" {
		projectID = os.Getenv("DEVSHELL_PROJECT_ID")
	}
	if projectID == "" {
		log.Println("Fetching Project ID from metadata server")
		metadataURL := "http://metadata.google.internal/computeMetadata/v1/project/project-id"
		// Send a GET request to the metadata server
		// Do this to make it super-simple for CEs to deploy
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, metadataURL, nil)
		if err != nil {
			log.Printf("Warning - could not retrieve project ID from metadata server")
			log.Fatalln(err)
		}
		req.Header.Set("Metadata-Flavor", "Google")
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Warning - could not retrieve project ID from metadata server")
			log.Fatalln(err)
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Warning - could not retrieve project ID from metadata server")
			log.Fatalln(err)
		}
		projectID = string(b)
	}
	if projectID == "" {
		log.Fatalf("Expected PROJECT_ID environment variable to be set")
	}

	log.Printf("Running in project: %v\n", projectID)

	cfg = config{
		Port:      port,
		ProjectID: projectID,
	}
}

func intro(ctx context.Context) {
	log.Printf("Registering with bond service at %v\n", bondCfg.BondURL)
	ai := AppInstance{
		ProjectID: cfg.ProjectID,
	}
	body, err := sendJson(ctx, "/v1/intro", ai)
	if err != nil {
		if body != nil {
			log.Printf("Response from Bond: %v\n", string(body))
		}
		log.Println("Your project must be registered. Please register here: go/techday-reg")
		log.Fatalf("Could not register with Bond Service: %v\n", err)
	}
	log.Printf("Bond Service Replied: %v\n", string(body))
}

func main() {
	// Give curious CEs bonus points if they uncomment this line
	// getBonusPoints()

	log.Printf("Starting backend service on project %v\n", cfg.ProjectID)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/", defaultHandler)

	// Eventful Day Story
	r.Route("/eventful_day", eventfulDayRouter)

	// Data-Driven Decaf
	r.Route("/data_driven_decaf", DDDRouter)

	// Start HTTP server.
	log.Printf("Listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from %v!\n", cfg.ProjectID)
}
