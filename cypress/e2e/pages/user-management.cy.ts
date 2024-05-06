import { faker } from "@faker-js/faker";
import users from "../../fixtures/users.json";

describe("User Management", () => {
  const openAndCloseDeleteModal = () => {
    cy.get("[data-cy=open-actions]").first().click();
    cy.get("[data-cy=actions]").should("be.visible");
    cy.get("[data-cy=action-delete-user]:visible").click();
    cy.contains(
      "h2",
      "Bist du sicher, dass du den Benutzer löschen möchtest?",
    ).should("be.visible");
    cy.contains("button", "Abbrechen").click();
    cy.contains(
      "h2",
      "Bist du sicher, dass du den Benutzer löschen möchtest?",
    ).should("not.exist");
  };

  const deleteUser = (user: any) => {
    cy.addUser(
      user.firstname,
      user.lastname,
      user.email,
      user.dateOfBirth,
      user.password,
      user.role,
    );
    cy.visit("/users");
    cy.get('tr[data-cy="user-row"]')
      .contains('div[data-cy="user-email"]', user.email)
      .closest("tr")
      .within(() => {
        cy.get('[data-cy="open-actions"]').click();
        cy.get('[data-cy="actions"]').should("be.visible");
        cy.get('[data-cy="action-delete-user"]').click();
      })
      .then(() => {
        cy.contains(
          "h2",
          "Bist du sicher, dass du den Benutzer löschen möchtest?",
        ).should("be.visible");
        cy.contains("button", "Ja").click();
        cy.contains(
          "h2",
          "Bist du sicher, dass du den Benutzer löschen möchtest?",
        ).should("not.exist");
        cy.get("[data-cy=user-table]").contains(user.email).should("not.exist");
      });
  };

  const fillFieldAndCheckValidity = (fieldSelector: string, value: string) => {
    cy.get(fieldSelector + ":invalid").should("have.length", 1);
    cy.get(fieldSelector).type(value);
    cy.get(fieldSelector + ":invalid").should("have.length", 0);
  };

  const testValidInputInAddUserModal = () => {
    cy.get("[data-cy=open-add-user-overlay]").click();
    cy.get("[data-cy=add-user-overlay]").should("be.visible");

    fillFieldAndCheckValidity("[data-cy=firstname]", faker.person.firstName());
    fillFieldAndCheckValidity("[data-cy=lastname]", faker.person.lastName());

    const email = faker.internet.email().toLowerCase();
    cy.get("[data-cy=email]").type(email.split("@")[0]);
    cy.get("[data-cy=email]:invalid").should("have.length", 1);
    cy.get("[data-cy=email]").clear();
    cy.get("[data-cy=email]").type(email);
    cy.get("[data-cy=email]:invalid").should("have.length", 0);

    const password = faker.internet.password();
    fillFieldAndCheckValidity("[data-cy=password]", password);
    fillFieldAndCheckValidity("[data-cy=confirmPassword]", password);

    fillFieldAndCheckValidity(
      "[data-cy=dateOfBirth]",
      faker.date.past().toISOString().split("T")[0],
    );
    cy.get("[data-cy=role]").select("admin");

    cy.get("[data-cy=add-user-overlay-submit]:visible").click();
    cy.get("[data-cy=add-user-overlay]").should("not.be.visible");
    cy.get("[data-cy=user-table]").contains(email).should("exist");
  };

  beforeEach(() => {
    cy.login(users.createdAdmin.email, users.createdAdmin.password);
    cy.visit("/users");
  });

  context("Mobile View", () => {
    beforeEach(() => {
      cy.viewport("iphone-6");
    });

    it("should open and close user delete modal", openAndCloseDeleteModal);

    it("should delete user on confirm", () => {
      const dateOfBirth = faker.date.past().toISOString().split("T")[0];
      const user = {
        firstname: faker.person.firstName(),
        lastname: faker.person.lastName(),
        email: faker.internet.email().toLowerCase(),
        dateOfBirth,
        password: faker.internet.password(),
        role: "admin",
      };
      deleteUser(user);
    });

    it(
      "should only allow valid input in add user modal",
      testValidInputInAddUserModal,
    );

    describe("should open and close add user modal", () => {
      it("using the close button", () => {
        cy.get("[data-cy=open-add-user-overlay]").click();
        cy.get("[data-cy=add-user-overlay]").should("be.visible");
        cy.get("[data-cy=add-user-overlay-close]:visible").click();
        cy.get("[data-cy=add-user-overlay]").should("not.be.visible");
      });

      it("using the esc button", () => {
        cy.get("[data-cy=open-add-user-overlay]").click();
        cy.get("[data-cy=add-user-overlay]").should("be.visible");
        cy.get("[data-cy=add-user-overlay]").type("{esc}");
        cy.get("[data-cy=add-user-overlay]").should("not.be.visible");
      });
    });
  });

  context("Desktop View", () => {
    beforeEach(() => {
      cy.viewport(1280, 720);
    });

    it("should open and close user delete modal", openAndCloseDeleteModal);

    it("should delete user on confirm", () => {
      const dateOfBirth = faker.date.past().toISOString().split("T")[0];
      const user = {
        firstname: faker.person.firstName(),
        lastname: faker.person.lastName(),
        email: faker.internet.email().toLowerCase(),
        dateOfBirth,
        password: faker.internet.password(),
        role: "admin",
      };
      deleteUser(user);
    });

    it(
      "should only allow valid input in add user modal",
      testValidInputInAddUserModal,
    );
  });
});
