Feature: Account
  Background:
    Given Account endpoint is available
    And Subtotal endpoint is available
    And we POST Subtotal "Expenses" with no parent
    And we POST Subtotal "Income" with no parent

  Scenario: GET non-existent Account
    When we GET Account "Advertising"
    Then we should see HTTP response status 404

  Scenario: POST Account
    When we POST Account "Entertainment" with Subtotal "Expenses"
    Then we should see HTTP response status 204

  Scenario: POST Account with non-existent Subtotal
    When we POST Account "Gold" with Subtotal "Assets"
    Then we should see HTTP response status 404

  Scenario: POST Account and then GET
    When we POST Account "Insurance" with Subtotal "Expenses"
    And we GET Account "Insurance"
    Then we should see HTTP response status 200
    And we should see Account "Insurance" with Subtotal "Expenses"

  Scenario: POST Account with same name as existing
    Given we POST Account "Interest" with Subtotal "Expenses"
    When we POST Account "Interest" with Subtotal "Income"
    Then we should see HTTP response status 409

  Scenario: PATCH Account
    Given we POST Account "Telephone" with Subtotal "Expenses"
    When we PATCH Account "Telephone" with new name "Internet"
    Then we should see HTTP response status 204

  Scenario: PATCH non-existent Account
    When we PATCH Account "Stock" with new name "Shares"
    Then we should see HTTP response status 404

  Scenario: PATCH Account with same name as existing
    Given we POST Account "Salaries" with Subtotal "Expenses"
    And we POST Account "Wages" with Subtotal "Expenses"
    When we PATCH Account "Wages" with new name "Salaries"
    Then we should see HTTP response status 409

  Scenario: PATCH Account with new name and then GET by new name
    Given we POST Account "Printing" with Subtotal "Expenses"
    When we PATCH Account "Printing" with new name "Stationery"
    And we GET Account "Stationery"
    Then we should see HTTP response status 200
    And we should see Account "Stationery" with Subtotal "Expenses"

  Scenario: PATCH Account with new name and then GET by old name
    Given we POST Account "Audit" with Subtotal "Expenses"
    When we PATCH Account "Audit" with new name "Professional Services"
    And we GET Account "Audit"
    Then we should see HTTP response status 404

  Scenario: PATCH Account with new Subtotal and then GET
    Given we POST Account "Rent" with Subtotal "Expenses"
    When we PATCH Account "Rent" with new Subtotal "Income"
    And we GET Account "Rent"
    Then we should see HTTP response status 200
    And we should see Account "Rent" with Subtotal "Income"

  Scenario: PATCH Account with non-existent Subtotal
    Given we POST Account "Prepaid Rent" with Subtotal "Expenses"
    When we PATCH Account "Prepaid Rent" with new Subtotal "Assets"
    Then we should see HTTP response status 404

  Scenario: DELETE Account
    Given we POST Account "Postage" with Subtotal "Expenses"
    When we DELETE Account "Postage"
    Then we should see HTTP response status 200
    And we should see Account "Postage" with Subtotal "Expenses"

  Scenario: DELETE non-existent Account
    When we DELETE Account "Licenses"
    Then we should see HTTP response status 404

  Scenario: DELETE Account and then GET
    Given we POST Account "Maintenance" with Subtotal "Expenses"
    When we DELETE Account "Maintenance"
    And we GET Account "Maintenance"
    Then we should see HTTP response status 404

  # TODO
  # Scenario: DELETE Account with existing Journal Entry Line
