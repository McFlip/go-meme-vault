import { When, Then, Given, Before } from "@badeball/cypress-cucumber-preprocessor"
import '@testing-library/cypress/add-commands'

Before(() => {
  cy.request('DELETE', 'http://localhost:8080/api/testhooks/nuke')
})

Given("There are new images in the full image path", () => {
  cy.request("http://localhost:8080/public/img/full/doom_eternal.jpg").its('status').should('equal', 200)
})

Given("I am on the new memes page", () => {
  cy.visit("http://localhost:8080/")
  cy.get('menu').findByRole('link', {name: 'New Memes'}).click()
})

When("I click the Scan button", () => {
  cy.findByRole('button', {name: 'Scan'}).click()
})

Then("I should see a list of thumbnails for the new memes", () => {
  cy.get('main').findAllByRole('img').each(($el, index, $list) => {
    cy.wrap($el).should('have.attr', 'src').and('match', /public\/img\/tn\//)
  }).then(($list) => {
    expect($list).to.have.length(4)
  })
})

Given("All memes in the full image path are already in the DB", () => {
  cy.visit("http://localhost:8080/")
  cy.get('menu').findByRole('link', {name: 'New Memes'}).click()
  cy.findByRole('button', {name: 'Scan'}).click()
})

Then("I should not see any memes listed", () => {
  cy.get('#memes').findAllByRole('img').should('not.exist')
})