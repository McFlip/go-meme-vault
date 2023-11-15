Feature: The search bar allows you to search tags and select a tag
  As a Meme Lord
  I want to search for and select a starting tag
  So that I can see all memes with that tag
  Scenario: Meme Lord searches tags and selects one
    Given I am on the home page
    When I search for cats
    And I click on the cats tag
    Then I should go to the memes page
    And I should see my tag has been selected