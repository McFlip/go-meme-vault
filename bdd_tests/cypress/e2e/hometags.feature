Feature: Browse list of tags
  As a Meme Lord,
  I want to see a list of all available tags on the home page,
  So that I can see what kinds of memes are in my vault.
  Scenario: Meme Lord sees list of all tags on the home page.
    Given I am on the home page
    When I scroll through the list of tags
    Then I should see all the available tags