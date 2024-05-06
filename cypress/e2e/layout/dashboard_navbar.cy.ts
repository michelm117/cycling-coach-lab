import users from "../../fixtures/users.json";

describe("Dashboard Navbar", () => {
  const openNavbarActionMenu = () => {
    cy.get('[data-cy="navbar"]').should("exist");
    cy.get('[data-cy="open-navbar-action-menu"]').click();
    cy.get('[data-cy="navbar-action-menu"]').should("be.visible");
  };

  beforeEach(() => {
    cy.viewport("macbook-15");
    const { email, password } = users.createdAdmin;
    cy.login(email, password);
    cy.visit("/");
    cy.url().should("not.include", "/login");
  });

  it("should open and close profile drop down", () => {
    openNavbarActionMenu();

    // close by clicking outside
    cy.get('[data-cy="navbar"]').click();
    cy.get('[data-cy="navbar-action-menu"]').should("not.be.visible");
  });

  it("should navigate to the settings page", () => {
    openNavbarActionMenu();

    cy.get('[data-cy="settings"]').click();
    cy.url().should("eq", Cypress.config().baseUrl + "/settings");
  });

  it("should log out", () => {
    Cypress.session.clearAllSavedSessions();
    openNavbarActionMenu();

    cy.get('[data-cy="logout"]').click();
    cy.url().should("include", "/login");

    cy.visit("/");
    cy.url().should("include", "/login");
  });
});
