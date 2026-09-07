const { defineConfig } = require('cypress')

module.exports = defineConfig({
  e2e: {
    // Override with CYPRESS_BASE_URL to point at a different environment
    // (e.g. a deployed instance), same idea as the Go E2E suite's API_BASE_URL.
    baseUrl: 'http://localhost:5173',
    supportFile: 'cypress/support/e2e.js',
  },
})
