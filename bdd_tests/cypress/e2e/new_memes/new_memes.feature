Feature: Scan for new memes
  As a Meme Lord
  I want to store fresh memes into my vault
  So that I can tag and find them later
  Scenario: Meme Lord dropped some fresh memes and wants to load them into the system
    Given There are new images in the full image path
    And I am on the new memes page
     When I click the Scan button
     Then I should see a list of thumbnails for the new memes

