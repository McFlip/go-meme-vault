import { When, Then, Given } from "@badeball/cypress-cucumber-preprocessor";

Given("I am on the home page", () => {
  cy.visit("localhost:8080/")
})

When("I scroll through the list of tags", () => {
  console.log("TODO: implement scroll into view")
})

Then("I should see all the available tags", () => {
  cy.get("#tag-list").should("contain.text", "test tag")
})