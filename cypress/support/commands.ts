// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
export {};

Cypress.Commands.add("login", (email, password) => {
  cy.session([email, password], () => {
    cy.log("Logging in");
    cy.request({
      method: "POST",
      url: "/auth/login",
      form: true,
      body: {
        email: email,
        password: password,
      },
    });
  });
});

Cypress.Commands.add(
  "addUser",
  (firstname, lastname, email, dateOfBirth, password, role) => {
    cy.request({
      method: "POST",
      url: "/users",
      form: true,
      body: {
        firstname,
        lastname,
        email,
        dateOfBirth,
        password,
        confirmPassword: password,
        role,
      },
    }).then((response) => {
      expect(response.status).to.eq(200);
    });
  },
);

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Custom command to login
       * @example cy.login('email', 'password')
       */
      login(email: string, password: string): void;
      addUser(
        firstName: string,
        lastName: string,
        email: string,
        dateOfBirth: string,
        password: string,
        role: string,
      ): void;
    }
  }
}
