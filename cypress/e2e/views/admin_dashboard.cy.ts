import { faker } from "@faker-js/faker";

describe("Admin Dashboard", () => {
  beforeEach(() => {
    cy.visit("/users");
  });

  context("Mobile View", () => {
    beforeEach(() => {
      cy.viewport("iphone-6");
    });

    it("should open and close user delete modal", () => {
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
    });

    it("should delete user on confirm", () => {
      cy.get("[data-cy=user-row]")
        .first()
        .find("[data-cy=user-email]")
        .invoke("text")
        .then((email) => {
          expect(email).to.not.be.empty;
          expect(email).to.not.be.undefined;
          cy.get("[data-cy=open-actions]").first().click();
          cy.get("[data-cy=actions]").should("be.visible");
          cy.get("[data-cy=action-delete-user]:visible").click();
          cy.contains(
            "h2",
            "Bist du sicher, dass du den Benutzer löschen möchtest?",
          ).should("be.visible");
          cy.contains("button", "Ja").click();
          cy.contains(
            "h2",
            "Bist du sicher, dass du den Benutzer löschen möchtest?",
          ).should("not.exist");
          cy.get("[data-cy=user-table]").contains("email").should("not.exist");
        });
    });

    it("should only allow valid input in add user modal", () => {
      cy.get("[data-cy=open-add-user-overlay]").click();
      cy.get("[data-cy=add-user-overlay]").should("be.visible");

      const fillFieldAndCheckValidity = (
        fieldSelector: string,
        value: string,
      ) => {
        cy.get(fieldSelector + ":invalid").should("have.length", 1);
        cy.get(fieldSelector).type(value);
        cy.get(fieldSelector + ":invalid").should("have.length", 0);
      };

      fillFieldAndCheckValidity(
        "[data-cy=firstname]",
        faker.person.firstName(),
      );
      fillFieldAndCheckValidity("[data-cy=lastname]", faker.person.lastName());

      const email = faker.internet.email();
      cy.get("[data-cy=email]").type(email.split("@")[0]);
      cy.get("[data-cy=email]:invalid").should("have.length", 1);
      cy.get("[data-cy=email]").clear();
      cy.get("[data-cy=email]").type(email);
      cy.get("[data-cy=email]:invalid").should("have.length", 0);

      fillFieldAndCheckValidity(
        "[data-cy=password]",
        faker.internet.password(),
      );

      fillFieldAndCheckValidity(
        "[data-cy=dateOfBirth]",
        faker.date.past().toISOString().split("T")[0],
      );
      cy.get("[data-cy=role]").select("admin");

      cy.get("[data-cy=add-user-overlay-submit]:visible").click();
      cy.get("[data-cy=add-user-overlay]").should("not.be.visible");

      cy.get("[data-cy=user-row]").contains(email).should("be.visible");
    });

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

    it("should open and close user delete modal", () => {
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
    });

    it("should delete user on confirm", () => {
      cy.get("[data-cy=user-row]")
        .first()
        .find("[data-cy=user-email]")
        .invoke("text")
        .then((email) => {
          expect(email).to.not.be.empty;
          expect(email).to.not.be.undefined;
          cy.get("[data-cy=open-actions]").first().click();
          cy.get("[data-cy=actions]").should("be.visible");
          cy.get("[data-cy=action-delete-user]:visible").click();
          cy.contains(
            "h2",
            "Bist du sicher, dass du den Benutzer löschen möchtest?",
          ).should("be.visible");
          cy.contains("button", "Ja").click();
          cy.contains(
            "h2",
            "Bist du sicher, dass du den Benutzer löschen möchtest?",
          ).should("not.exist");
          cy.get("[data-cy=user-table]").contains("email").should("not.exist");
        });
    });

    it("should only allow valid input in add user modal", () => {
      cy.get("[data-cy=open-add-user-overlay]").click();
      cy.get("[data-cy=add-user-overlay]").should("be.visible");

      const fillFieldAndCheckValidity = (
        fieldSelector: string,
        value: string,
      ) => {
        cy.get(fieldSelector + ":invalid").should("have.length", 1);
        cy.get(fieldSelector).type(value);
        cy.get(fieldSelector + ":invalid").should("have.length", 0);
      };

      fillFieldAndCheckValidity(
        "[data-cy=firstname]",
        faker.person.firstName(),
      );
      fillFieldAndCheckValidity("[data-cy=lastname]", faker.person.lastName());

      const email = faker.internet.email();
      cy.get("[data-cy=email]").type(email.split("@")[0]);
      cy.get("[data-cy=email]:invalid").should("have.length", 1);
      cy.get("[data-cy=email]").clear();
      cy.get("[data-cy=email]").type(email);
      cy.get("[data-cy=email]:invalid").should("have.length", 0);

      fillFieldAndCheckValidity(
        "[data-cy=password]",
        faker.internet.password(),
      );

      fillFieldAndCheckValidity(
        "[data-cy=dateOfBirth]",
        faker.date.past().toISOString().split("T")[0],
      );
      cy.get("[data-cy=role]").select("admin");

      cy.get("[data-cy=add-user-overlay-submit]:visible").click();
      cy.get("[data-cy=add-user-overlay]").should("not.be.visible");

      cy.get("[data-cy=user-row]").contains(email).should("be.visible");
    });
  });
});
