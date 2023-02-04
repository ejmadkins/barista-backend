package main

import (
	"context"
	"log"
	"os"
	"testing"
)

type QAPayload struct {
	ProjectID string `json:"project_id"`
	CI        string `json:"ci,omitempty"`
}

// This will actually init the app, as we need the team to be registered to test something and log the score.
func setupTest(ctx context.Context) {
	initConfig(ctx)
	initBond()
	intro(ctx)
}

func Test_QAChapter(t *testing.T) {
	//Setup the context once per test suite
	ctx := context.Background()
	setupTest(ctx)

	//Check where the test is being run from
	CI := "false"
	if os.Getenv("CI") != "" {
		CI = os.Getenv("CI")
	}

	ctx = context.WithValue(ctx, "ci", CI)

	tests := []struct {
		name    string
		payload QAPayload
		want    string
	}{
		{
			name: "bonus",
			payload: QAPayload{
				ProjectID: cfg.ProjectID,
				CI:        ctx.Value("ci").(string),
			},
			want: "true",
		},
	}

	//Run all the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("Starting test: %s", tt.name)
			log.Printf("POST sendJson with payload = %v", tt.payload)
			_, err := sendJson(ctx, "/v1/qa", tt.payload)
			if ctx.Value("ci").(string) != tt.want {
				log.Printf("Hint: Run the tests in a pipeline to get bonus points!")
			}
			if err != nil {
				t.Errorf("sendJson error = %v, expected %v", err, tt.want)
				return
			}
		})
	}
}
