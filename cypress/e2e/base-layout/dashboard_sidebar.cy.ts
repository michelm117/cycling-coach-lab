import users from "../../fixtures/users.json";

describe("Dashboard Sidebar", () => {
  const testSidebarMobile = () => {
    beforeEach(() => {
      cy.viewport("iphone-6");
      const { email, password } = users.createdAdmin;
      cy.login(email, password);
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
  };

  const testSidebarDesktop = () => {
    beforeEach(() => {
      cy.viewport("macbook-15");
      const { email, password } = users.createdAdmin;
      cy.login(email, password);
      cy.visit("/");
    });

    it("should be open by default", () => {
      cy.get('[data-cy="sidebar"]').should("be.visible");
    });

    it("should not be able to close", () => {
      cy.get('[data-cy="sidebar-open"]').should("not.be.visible");
    });
  };

  const testSidebarNavigation = () => {
    beforeEach(() => {
      cy.viewport("macbook-15");
      const { email, password } = users.createdAdmin;
      cy.login(email, password);
      cy.visit("/");
    });

    it("should navigate to the home page", () => {
      cy.get('[data-cy="sidebar-user-management"]').click();
      cy.url().should("eq", Cypress.config().baseUrl + "/users");
    });
  };

  describe("on mobile", testSidebarMobile);

  describe("on desktop", testSidebarDesktop);

  describe("navigation", testSidebarNavigation);
});
