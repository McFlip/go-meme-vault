import { When, Then, Given, Before } from "@badeball/cypress-cucumber-preprocessor"
import '@testing-library/cypress/add-commands'

Before(() => {
  cy.request('DELETE', 'http://localhost:8080/api/testhooks/nuke')
})

Given("I am on the home page", () => {
  var tag = {name: "cats"}
  cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
  cy.visit("localhost:8080/")
})

When("I search for cats", () => {
  cy.get("#searchbar").findByLabelText("Search").type("cat")
})

When("I click on the cats tag", () => {
  cy.get("#searchbar").findByText("cats").click()
})

Then("I should go to the memes page", () => {
  cy.location('pathname').should('equal', '/memes')
})

Then("I should see my tag has been selected", () => {
  cy.get("#main").findByText("cats").should('exist')
})