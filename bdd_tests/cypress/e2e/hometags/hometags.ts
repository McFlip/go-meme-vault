import { When, Then, Given, Before } from "@badeball/cypress-cucumber-preprocessor"
import '@testing-library/cypress/add-commands'

Before(() => {
  cy.request('DELETE', 'http://localhost:8080/api/testhooks/nuke')
})

Given("Tags exist in the database", () => {
  cy.fixture('tags.json').then((tags) => {
    tags.forEach(tag => {
      cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
    });
  })
})

When("I go to the home page", () => {
  cy.visit("localhost:8080/")
})

// When("I scroll through the list of tags", () => {
//   cy.findByRole('link', {name: 'third tag'}).scrollIntoView()
// })

Then("I should see all the available tags", () => {
  cy.get('main').findAllByRole('listitem').should('have.length.at.least', 3)
})