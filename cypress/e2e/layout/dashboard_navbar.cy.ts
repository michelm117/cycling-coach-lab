import users from "../../fixtures/users.json";

describe("Dashboard Navbar", () => {
  it("should log out", () => {
    Cypress.session.clearAllSavedSessions();

    cy.login(users.createdAdmin.email, users.createdAdmin.password);
    cy.visit("/");
    cy.url().should("not.include", "/login");

    cy.get('[data-cy="navbar"]').should("exist");
    cy.get('[data-cy="open-navbar-action-menu"]').click();
    cy.get('[data-cy="navbar-action-menu"]').should("be.visible");
    cy.get('[data-cy="navbar-action-menu"]').should("be.visible");
    cy.get('[data-cy="logout"]').click();

    cy.url().should("include", "/login");

    cy.visit("/");
    cy.url().should("include", "/login");
  });
});
