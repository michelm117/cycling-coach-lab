import users from "../../fixtures/users.json";

describe("Settings", () => {
  const getThemes = () => {
    return cy
      .get("[data-cy=theme-switcher]")
      .invoke("text")
      .then((themeText) => {
        const themes = themeText.split(" ").slice(0, -1);
        return themes.map((theme) => theme.toLowerCase());
      });
  };

  beforeEach(() => {
    const { email, password } = users.createdAdmin;
    cy.login(email, password);
    cy.visit("/settings");
  });

  it("should switch to themes", () => {
    getThemes().then((themes) => {
      themes.forEach((theme) => {
        // check if a theme is already selected
        cy.get('[data-cy="theme-switcher"] input[type="radio"]').should(
          "be.checked",
        );

        // select a different theme
        cy.get(`input[type="radio"][value="${theme}"]`).click();
        cy.get("[data-cy=save-theme]").click();

        // check if the theme has been saved
        cy.get("html").invoke("attr", "data-theme").should("equal", theme);
      });
    });
  });
});
