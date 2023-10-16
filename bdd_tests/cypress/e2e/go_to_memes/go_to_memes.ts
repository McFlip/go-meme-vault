import { When, Then, Given } from "@badeball/cypress-cucumber-preprocessor"
import '@testing-library/cypress/add-commands'

Given("I am on the home page", () => {
  var tag = {name: "Doom"}
  cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
  // cy.request({method: 'POST', url: 'http://localhost:8080/api/testhooks/tags', body: tag, retryOnNetworkFailure: false})
  cy.visit("localhost:8080/")
})

When("I select a starting tag", () => {
  cy.findByText("Doom").click()
})

Then("I should go to the memes page", () => {
  cy.location('pathname').should('equal', '/memes')
})

Then("I should see my tag has been selected", () => {
  cy.findByText("Doom").should('exist')
})