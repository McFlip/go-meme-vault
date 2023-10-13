import { When, Then, Given } from "@badeball/cypress-cucumber-preprocessor";
import '@testing-library/cypress/add-commands'

Given("I am on the home page", () => {
  cy.visit("localhost:8080/")
})

When("I scroll through the list of tags", () => {
  cy.findByRole('link', {name: 'third tag'}).scrollIntoView()
})

Then("I should see all the available tags", () => {
  cy.get('main').findAllByRole('listitem').should('have.length', 3)
})