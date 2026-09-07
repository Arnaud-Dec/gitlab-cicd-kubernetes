package e2e_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health endpoints", func() {
	It("reports liveness on /healthz", func(ctx context.Context) {
		resp, err := api.Healthz(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("reports readiness on /readyz, meaning the database is reachable", func(ctx context.Context) {
		resp, err := api.Readyz(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})
