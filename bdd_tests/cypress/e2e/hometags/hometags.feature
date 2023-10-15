Feature: Browse list of tags
  As a Meme Lord,
  I want to see a list of all available tags on the home page,
  So that I can see what kinds of memes are in my vault.
  Scenario: Meme Lord sees list of all tags on the home page.
    Given Tags exist in the database
    When I go to the home page
    Then I should see all the available tags