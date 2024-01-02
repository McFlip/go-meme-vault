Feature: Navigate through memes like drilling down into a tree
  As a Meme Lord
  I want to narrow down search results by selecting additional search tags
  So that I can find a specific meme
  Scenario: Meme Lord selects a tag from the home page and keeps selecting tags
    Given There are four memes all with the tag four
    And 3 of those memes have the tag '3'
    And 2 of those memes have the tag '2'
    And 1 of those memes has the tag '1'
    And I am on the home page
    When I select the 4 tag
    And I select the 3 tag
    And I select the 2 tag
    And I select the 1 tag
    Then I should only see 1 meme
