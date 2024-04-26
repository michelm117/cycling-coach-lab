const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    baseUrl: "http://localhost:8080",
    supportFile: "cypress/support/e2e.{js,jsx,ts,tsx}",
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    specPattern: [
      "cypress/e2e/setup/start.cy.ts",
      "cypress/e2e/pages/*.cy.ts",
      "cypress/e2e/api/*.cy.ts",
      "cypress/e2e/base-layout/*.cy.ts",
      "cypress/e2e/setup/stop.cy.ts",
    ],
  },
});
