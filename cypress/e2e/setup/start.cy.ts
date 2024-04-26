import users from "../../fixtures/users.json";

describe("App Setup", () => {
  it("should setup a new application", () => {
    cy.visit("/setup");

    const { firstname, lastname, email, dateOfBirth, password } =
      users.createdAdmin;
    cy.get("input[name=firstname]").type(firstname);
    cy.get("input[name=lastname]").type(lastname);
    cy.get("input[name=email]").type(email);
    cy.get("input[name=password]").type(password);
    cy.get("input[name=dateOfBirth]").type(dateOfBirth);

    cy.get("button[type=submit]").click();
    cy.url().should("contain", "login");
  });
});
