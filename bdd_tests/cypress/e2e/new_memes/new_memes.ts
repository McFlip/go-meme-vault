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

Given("The new memes have not been loaded in the DB yet", () => {
  cy.request('DELETE', 'http://localhost:8080/api/testhooks/nuke')
})

When("I select the first thumbnail", () => {
  cy.get('#memes').findAllByRole('img').first().click()
})

Then("I should see a modal with the full image", () => {
  cy.get('#modal').findByRole('img').should('exist')
})

When("I submit a new tag", () => {
  cy.get("#modal").findByRole('searchbox').type('first tag')
  cy.get("#modal").findByRole('button', {name: 'Create'}).click()
})

Then("I should see the new tag in the list of tags for this meme", () => {
  cy.get("#tags").findByText('first tag').should('exist')
})

Given("A tag exists in the DB", () => {
    cy.fixture('tags.json').then((tags) => {
    tags.forEach(tag => {
      cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
    });
  })
})

When("I select the existing tag", () => {
  cy.get("#modal").findByRole('searchbox').type('first')
  cy.get("#modal").findByRole('button', {name: 'first tag'}).click()
})