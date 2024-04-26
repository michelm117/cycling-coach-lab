describe("Utils handler routes", () => {
  it("version route should return version", () => {
    cy.request("/version").then((response) => {
      const version = Cypress.env("APP_VERSION");
      if (!version) {
        throw new Error("Environment variable 'APP_VERSION' is not set");
      }
      expect(response.status).to.eq(200);
      expect(response.body).to.contain(version);
    });
  });

  it("health route should return 200 and text", () => {
    cy.request("/health").then((response) => {
      expect(response.status).to.eq(200);
      expect(response.body).to.contains("Service is healthy!");
    });
  });
});
