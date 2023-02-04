/*
This is the code for the Eventful Day story

Listens for a Pub/Sub and/or Eventarc message and sends it to Bond.
*/

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Subset of what we get from EventArc
type EventarcPayload struct {
	Kind    string `json:"kind,omitempty"`
	Name    string `json:"name,omitempty"`
	Project string `json:"project,omitempty"`
}

// Subset of what we get from PubSub
type PubSubPayload struct {
	Message struct {
		Attributes struct {
			EventType string `json:"eventType,omitempty"`
			ObjectID  string `json:"objectId,omitempty"`
		} `json:"attributes,omitempty"`
	} `json:"message,omitempty"`
	Subscription string `json:"subscription,omitempty"`
}

// Handle an incoming event (Pub/Sub or Eventarc)
func EventfulDayHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Eventful Day Task: Event received")

	var eventarcPayload EventarcPayload
	var pubSubPayload PubSubPayload

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Eventful Day Task: Error: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// Attempt to read as PubSub and Eventarc
	err = json.Unmarshal(body, &eventarcPayload)
	if err != nil {
		log.Printf("Eventful Day Task: Error: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &pubSubPayload)
	if err != nil {
		log.Printf("Eventful Day Task: Error: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if it's a pubsub message first
	if eventarcPayload.Kind == "" && pubSubPayload.Message.Attributes.EventType == "OBJECT_FINALIZE" {
		// Copy fields from Pubsub to eventarc struct to process later
		eventarcPayload.Kind = "storage#object"
		eventarcPayload.Name = pubSubPayload.Message.Attributes.ObjectID
		log.Println("Eventful Day Task: Got event from Pub/Sub")
	} else {
		log.Println("Eventful Day Task: Got event from Eventarc")
	}

	// Ensure it's valid
	if eventarcPayload.Kind != "storage#object" {
		log.Printf("Eventful Day Task: Error: invalid kind: %v (expecting storage#object)\n", eventarcPayload.Kind)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// Check the name field is populated - this is used for verification later
	if eventarcPayload.Name == "" {
		log.Printf("Eventful Day Task: Error: missing Name in payload: %v\n", eventarcPayload)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Add the project name to the payload
	eventarcPayload.Project = cfg.ProjectID

	log.Printf("Eventful Day Task: Storage object \"%v\" in project %v\n", eventarcPayload.Name, eventarcPayload.Project)

	// Send payload to bond service for verification
	log.Println("Eventful Day Task: Verifying event payload with bond service")

	// Verify with Bond Service
	res, err := sendJson(r.Context(), "/v1/eventful_day/verify", eventarcPayload)
	if err != nil {
		log.Printf("Eventful Day Task: Error: %v\n", err)
		http.Error(w, "Error validating event", http.StatusInternalServerError)
		return
	}
	log.Printf("Response: %v\n", res)
}

// Chi router to handle incoming POST
func eventfulDayRouter(r chi.Router) {
	r.Post("/", EventfulDayHandler)
}
