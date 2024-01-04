import { When, Then, Given, Before } from "@badeball/cypress-cucumber-preprocessor"
import '@testing-library/cypress/add-commands'

Before(() => {
  cy.request('DELETE', 'http://localhost:8080/api/testhooks/nuke')
})

Given("I am on the home page", () => {
  cy.visit("localhost:8080/")
})

Given("There are four memes all with the tag four", () => {
  cy.visit("http://localhost:8080/memes/new");
  cy.findByRole("button", { name: "Scan" }).click();

  const tag = {name: "4"}
  cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
  for (let i = 1; i < 5; i++) {
    cy.request('POST', `http://localhost:8080/api/testhooks/memes/${i}/addtag/1`)
  }
})

Given("3 of those memes have the tag '3'", () => {
  const tag = {name: "3"}
  cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
  for (let i = 2; i < 5; i++) {
    cy.request('POST', `http://localhost:8080/api/testhooks/memes/${i}/addtag/2`)
  }
})

Given("2 of those memes have the tag '2'", () => {
  const tag = {name: "2"}
  cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
  for (let i = 3; i < 5; i++) {
    cy.request('POST', `http://localhost:8080/api/testhooks/memes/${i}/addtag/3`)
  }
})

Given("1 of those memes has the tag '1'", () => {
  const tag = {name: "1"}
  cy.request('POST', 'http://localhost:8080/api/testhooks/tags', tag)
  cy.request('POST', 'http://localhost:8080/api/testhooks/memes/4/addtag/4')
})

When("I select the 4 tag", () => {
  cy.findByText("4").click()
})

When("I select the 3 tag", () => {
  cy.findByText("3").click()
})

When("I select the 2 tag", () => {
  cy.findByText("2").click()
})

When("I select the 1 tag", () => {
  cy.findByText("1").click()
})

Then("I should only see 1 meme", () => {
  cy.findByRole("img")
})
