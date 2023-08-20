Feature: Subtotal API
  Scenario: POST Subtotal with no parent
    Given there is a Subtotal API
    When we POST a Subtotal with name "Balance Sheet" and no parent
    Then there should be a Subtotal with name "Balance Sheet" and no parent

  Scenario: POST Subtotal with parent
    Given there is a Subtotal API
    And there is a Subtotal with name "Assets"
    When we POST a Subtotal with name "Inventory" and parent "Assets"
    Then there should be a Subtotal with name "Inventory" and parent "Assets"
