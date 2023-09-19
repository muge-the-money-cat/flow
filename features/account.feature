Feature: Account
  Background:
    Given an Account endpoint is available
    And a Subtotal endpoint is available

  Scenario: GET non-existent Account
    When we GET an Account by name "Cash on Hand"
    Then we should see HTTP response status 404

  Scenario: POST Account
    Given we POST a Subtotal with name "Bank" and no parent
    When we POST an Account with name "Current Account" and Subtotal "Bank"
    Then we should see HTTP response status 204

  Scenario: POST Account with non-existent Subtotal
    When we POST an Account with name "Capital" and Subtotal "Equity"
    Then we should see HTTP response status 404

  Scenario: POST Account and then GET
    Given we POST a Subtotal with name "Inventory" and no parent
    When we POST an Account with name "Gold" and Subtotal "Inventory"
    And we GET an Account by name "Gold"
    Then we should see HTTP response status 200
    And we should see an Account with name "Gold" and Subtotal "Inventory"
