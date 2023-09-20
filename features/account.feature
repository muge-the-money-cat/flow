Feature: Account
  Background:
    Given an Account endpoint is available
    And a Subtotal endpoint is available
    And we POST a Subtotal with name "Expenses" and no parent
    And we POST a Subtotal with name "Income" and no parent

  Scenario: GET non-existent Account
    When we GET an Account by name "Advertising"
    Then we should see HTTP response status 404

  Scenario: POST Account
    When we POST an Account with name "Entertainment" and Subtotal "Expenses"
    Then we should see HTTP response status 204

  Scenario: POST Account with non-existent Subtotal
    When we POST an Account with name "Gold" and Subtotal "Assets"
    Then we should see HTTP response status 404

  Scenario: POST Account and then GET
    When we POST an Account with name "Insurance" and Subtotal "Expenses"
    And we GET an Account by name "Insurance"
    Then we should see HTTP response status 200
    And we should see an Account with name "Insurance" and Subtotal "Expenses"

  Scenario: POST Account with same name as existing
    Given we POST an Account with name "Interest" and Subtotal "Expenses"
    When we POST an Account with name "Interest" and Subtotal "Income"
    Then we should see HTTP response status 409

  Scenario: PATCH Account
    Given we POST an Account with name "Telephone" and Subtotal "Expenses"
    When we PATCH an Account named "Telephone" with new name "Internet"
    Then we should see HTTP response status 204

  Scenario: PATCH non-existent Account
    When we PATCH an Account named "Stock" with new name "Shares"
    Then we should see HTTP response status 404

  Scenario: PATCH Account with same name as existing
    Given we POST an Account with name "Salaries" and Subtotal "Expenses"
    And we POST an Account with name "Wages" and Subtotal "Expenses"
    When we PATCH an Account named "Wages" with new name "Salaries"
    Then we should see HTTP response status 409

  Scenario: PATCH Account with new name and then GET by new name
    Given we POST an Account with name "Printing" and Subtotal "Expenses"
    When we PATCH an Account named "Printing" with new name "Stationery"
    And we GET an Account by name "Stationery"
    Then we should see HTTP response status 200
    And we should see an Account with name "Stationery" and Subtotal "Expenses"

  Scenario: PATCH Account with new name and then GET by old name
    Given we POST an Account with name "Audit" and Subtotal "Expenses"
    When we PATCH an Account named "Audit" with new name "Professional Services"
    And we GET an Account by name "Audit"
    Then we should see HTTP response status 404

  Scenario: PATCH Account with new Subtotal and then GET
    Given we POST an Account with name "Rent" and Subtotal "Expenses"
    When we PATCH an Account named "Rent" with new Subtotal "Income"
    And we GET an Account by name "Rent"
    Then we should see HTTP response status 200
    And we should see an Account with name "Rent" and Subtotal "Income"

  Scenario: PATCH Account with non-existent Subtotal
    Given we POST an Account with name "Prepaid Rent" and Subtotal "Expenses"
    When we PATCH an Account named "Prepaid Rent" with new Subtotal "Assets"
    Then we should see HTTP response status 404

  Scenario: DELETE Account
    Given we POST an Account with name "Postage" and Subtotal "Expenses"
    When we DELETE an Account named "Postage"
    Then we should see HTTP response status 200
    And we should see an Account with name "Postage" and Subtotal "Expenses"

  Scenario: DELETE non-existent Account
    When we DELETE an Account named "Licenses"
    Then we should see HTTP response status 404

  Scenario: DELETE Account and then GET
    Given we POST an Account with name "Maintenance" and Subtotal "Expenses"
    When we DELETE an Account named "Maintenance"
    And we GET an Account by name "Maintenance"
    Then we should see HTTP response status 404

  # TODO
  # Scenario: DELETE Account with existing Journal Entry Line
