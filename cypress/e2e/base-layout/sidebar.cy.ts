describe("Sidebar", () => {
  describe("on mobile", () => {
    beforeEach(() => {
      cy.viewport("iphone-6");
      console.log("it worked")
      cy.login("user1@example.com", "1234")
      cy.visit("/");
    });

    it("should be closed by default", () => {
      cy.get('[data-cy="sidebar"]').should("not.be.visible");
    });

    it("should open and close when clicking the button", () => {
      cy.get('[data-cy="sidebar-open"]').should("be.visible");
      cy.get('[data-cy="sidebar-open"]').click();
      cy.get('[data-cy="sidebar"]').should("be.visible");
      cy.get('[data-cy="sidebar-overlay"]').click();
      cy.get('[data-cy="sidebar"]').should("not.be.visible");
    });
  });

  describe("on desktop", () => {
    beforeEach(() => {
      cy.viewport("macbook-15");
      cy.login("user1@example.com", "1234")
      cy.visit("/");
    });

    it("should be open by default", () => {
      cy.get('[data-cy="sidebar"]').should("be.visible");
    });

    it("should not be able to close", () => {
      cy.get('[data-cy="sidebar-open"]').should("not.be.visible");
    });
  });

  describe("navigation", () => {
    beforeEach(() => {
      cy.viewport("macbook-15");
      cy.login("user1@example.com", "1234")
      cy.visit("/");
    });

    it("should navigate to the home page", () => {
      cy.get('[data-cy="sidebar-user-management"]').click();
      cy.url().should("eq", Cypress.config().baseUrl + "/users");
    });
  });
});
