Feature: Subtotal API
  Scenario: POST Subtotal with no parent
    Given there is a Subtotal API
    When we POST a Subtotal with name "Balance Sheet" and no parent
    Then there should be a Subtotal with name "Balance Sheet" and no parent
