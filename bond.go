/*
This handles all communications to and from the Bond Service.

Bond service is expecting you, so if you can figure out how to hack it to get more points then it might be worth it.

Your call of course.
*/
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const defaultBondURL = "https://bond-service-l5xebjflvq-ew.a.run.app"

var bondCfg bondConfig

type bondConfig struct {
	BondURL string
}

func initBond() {
	url := os.Getenv("BOND_SERVICE_URL")
	if url == "" {
		url = defaultBondURL
	}

	bondCfg = bondConfig{
		BondURL: url,
	}

}

// Sends a JSON request as a POST body to bond and returns the raw bytes from the response
func sendJson(ctx context.Context, endpoint string, body any) (b []byte, err error) {
	// Marshall
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return b, err
	}
	url := bondCfg.BondURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return b, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return b, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return b, fmt.Errorf("expected 200 response instead got %v", res.StatusCode)
	}

	b, err = io.ReadAll(res.Body)
	if err != nil {
		return b, err
	}
	return b, nil
}

type BonusPayload struct {
	ProjectID string `json:"project_id"`
	Points    int64  `json:"points"`
}

// TODO: There is a way to get even more points here.
// I wonder if anyone will bother?
func getBonusPoints() {
	log.Println("Attempting to get bonus points")
	p := BonusPayload{
		ProjectID: cfg.ProjectID,
	}

	ep := "/v1/misc/bonus_points"
	b, err := sendJson(context.Background(), ep, p)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	log.Printf("Success! %v\n", string(b))
}
