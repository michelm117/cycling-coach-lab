import users from "../../fixtures/users.json";

describe("App close", () => {
  beforeEach(() => {
    const { email, password } = users.createdAdmin;
    cy.login(email, password);
  });

  it("should clear the database and close the server", () => {
    cy.visit("/settings");

    cy.get("button[data-cy=reset]").click();
  });
});
