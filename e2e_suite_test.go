// Package e2e_test holds black-box end-to-end specs for the Todo app
// backend: they only ever talk to it over real HTTP, at the address given by
// API_BASE_URL, the same way a real client would. See README.md for how to
// run these.
package e2e_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"e2e/apiclient"
)

const defaultBaseURL = "http://localhost:8080"

var api *apiclient.Client

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo App E2E Suite")
}

var _ = BeforeSuite(func(ctx context.Context) {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	api = apiclient.New(baseURL)

	By("waiting for the backend to report ready at " + baseURL)
	Eventually(func(ctx context.Context) (int, error) {
		resp, err := api.Readyz(ctx)
		if err != nil {
			return 0, err
		}
		return resp.StatusCode, nil
	}).WithContext(ctx).WithTimeout(30*time.Second).WithPolling(500*time.Millisecond).Should(
		Equal(200),
		"backend at %s never became ready — is it running? See README.md.", baseURL,
	)
})
