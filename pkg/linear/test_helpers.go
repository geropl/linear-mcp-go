package linear

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// NewTestClient creates a LinearClient for testing
// If record is true, it will record HTTP interactions
// If record is false, it will replay recorded interactions
func NewTestClient(t *testing.T, cassetteName string, record bool) (*LinearClient, func()) {
	if record {
		// Ensure API key is set when recording
		if os.Getenv("LINEAR_API_KEY") == "" {
			t.Fatal("LINEAR_API_KEY environment variable is required for recording")
		}
	}

	wipeAuthorizationHook := func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	}

	wipeChangingMetadataHook := func(i *cassette.Interaction) error {
		delete(i.Response.Headers, "Cf-Ray")
		delete(i.Response.Headers, "Date")

		for k := range i.Response.Headers {
			if strings.HasPrefix(strings.ToLower(k), "x-") {
				delete(i.Response.Headers, k)
			}
		}
		i.Response.Duration = 0
		return nil
	}

	// Create the recorder with appropriate mode
	options := []recorder.Option{
		// don't record authorization header in cassettes
		recorder.WithHook(wipeAuthorizationHook, recorder.AfterCaptureHook),
		recorder.WithHook(wipeChangingMetadataHook, recorder.AfterCaptureHook),
		recorder.WithMatcher(cassette.NewDefaultMatcher(cassette.WithIgnoreAuthorization())),
	}
	if record {
		options = append(options, recorder.WithMode(recorder.ModeRecordOnly))
	} else {
		options = append(options, recorder.WithMode(recorder.ModeReplayOnly))
	}

	r, err := recorder.New("../../testdata/fixtures/"+cassetteName, options...)
	if err != nil {
		t.Fatalf("Failed to create recorder: %v", err)
	}

	// Create a Linear client that uses the recorder's HTTP client
	apiKey := os.Getenv("LINEAR_API_KEY")
	client := &LinearClient{
		apiKey:      apiKey,
		httpClient:  r.GetDefaultClient(),
		rateLimiter: NewRateLimiter(1400),
	}

	// Return the client and a cleanup function
	cleanup := func() {
		r.Stop()
	}

	return client, cleanup
}
