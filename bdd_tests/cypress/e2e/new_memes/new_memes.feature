Feature: Scan for new memes
  As a Meme Lord
  I want to store fresh memes into my vault
  So that I can tag and find them later

  Scenario: Meme Lord dropped some fresh memes and wants to load them into the system
    Given There are new images in the full image path
    And I am on the new memes page
    When I click the Scan button
    Then I should see a list of thumbnails for the new memes

  Scenario: Meme Lord only wants to see the fresh memes found after a scan
    Given All memes in the full image path are already in the DB
    When I click the Scan button
    Then I should not see any memes listed

  Scenario: Meme Lord wants to see a large version of one of the fresh memes
    Given There are new images in the full image path
    And The new memes have not been loaded in the DB yet
    And I am on the new memes page
    When I click the Scan button
    And I select the first thumbnail
    Then I should see a modal with the full image

  Scenario: Meme Lord wants to create a new tag for a new meme
    Given There are new images in the full image path
    And The new memes have not been loaded in the DB yet
    And I am on the new memes page
    When I click the Scan button
    And I select the first thumbnail
    And I submit a new tag
    Then I should see the new tag in the list of tags for this meme

  Scenario: Meme Lord wants to tag a new meme with an existing tag
    Given There are new images in the full image path
    And The new memes have not been loaded in the DB yet
    And I am on the new memes page
    And A tag exists in the DB
    When I click the Scan button
    And I select the first thumbnail
    And I select the existing tag
    Then I should see the new tag in the list of tags for this meme

  Scenario: Meme Lord changes their mind after adding a tag and removes it
    Given There are new images in the full image path
    And The new memes have not been loaded in the DB yet
    And I am on the new memes page
    When I click the Scan button
    And I select the first thumbnail
    And I submit a new tag
    And I delete the tag
    Then I should not see any tags for this meme