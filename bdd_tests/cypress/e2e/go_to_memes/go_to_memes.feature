Feature: A tag link on the home page takes you to the memes page
  As a Meme Lord
  I want to select a starting tag
  So that I can see all memes with that tag
  Scenario: Meme Lord selects a tag from the home page
    Given I am on the home page
    When I select a starting tag
    Then I should go to the memes page
    And I should see my tag has been selected