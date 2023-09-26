Feature: Subtotal
  Background:
    Given Subtotal endpoint is available

  Scenario: GET non-existent Subtotal
    When we GET Subtotal "Assets"
    Then we should see HTTP response status 404

  Scenario: POST Subtotal
    When we POST Subtotal "Balance Sheet" with no parent
    Then we should see HTTP response status 204

  Scenario: POST Subtotal with no parent and then GET
    When we POST Subtotal "Profit & Loss" with no parent
    And we GET Subtotal "Profit & Loss"
    Then we should see HTTP response status 200
    And we should see Subtotal "Profit & Loss" with no parent

  Scenario: POST Subtotal with parent and then GET
    Given we POST Subtotal "Current Assets" with no parent
    When we POST Subtotal "Cash" with parent "Current Assets"
    And we GET Subtotal "Cash"
    Then we should see HTTP response status 200
    And we should see Subtotal "Cash" with parent "Current Assets"

  Scenario: POST Subtotal with same name as existing
    Given we POST Subtotal "Liabilities" with no parent
    When we POST Subtotal "Liabilities" with no parent
    Then we should see HTTP response status 409

  Scenario: POST Subtotal with non-existent parent
    When we POST Subtotal "Land" with parent "Assets"
    Then we should see HTTP response status 404

  Scenario: PATCH Subtotal
    Given we POST Subtotal "Discounts" with no parent
    When we PATCH Subtotal "Discounts" with new name "Sales Discounts"
    Then we should see HTTP response status 204

  Scenario: PATCH non-existent Subtotal
    When we PATCH Subtotal "Assets" with new parent "Balance Sheet"
    Then we should see HTTP response status 404

  Scenario: PATCH Subtotal with same name as existing
    Given we POST Subtotal "Machinery" with no parent
    And we POST Subtotal "Plant" with no parent
    When we PATCH Subtotal "Plant" with new name "Machinery"
    Then we should see HTTP response status 409

  Scenario: PATCH Subtotal with new name and then GET by new name
    Given we POST Subtotal "Depreciation" with no parent
    When we PATCH Subtotal "Depreciation" with new name "Amortisation"
    And we GET Subtotal "Amortisation"
    Then we should see HTTP response status 200
    And we should see Subtotal "Amortisation" with no parent

  Scenario: PATCH Subtotal with new name and then GET by old name
    Given we POST Subtotal "Provisions" with no parent
    When we PATCH Subtotal "Provisions" with new name "Reserves"
    And we GET Subtotal "Provisions"
    Then we should see HTTP response status 404

  Scenario: PATCH Subtotal with new parent and then GET
    Given we POST Subtotal "Sales" with no parent
    And we POST Subtotal "Discounts" with parent "Sales"
    And we POST Subtotal "Revenues" with no parent
    When we PATCH Subtotal "Discounts" with new parent "Revenues"
    And we GET Subtotal "Discounts"
    Then we should see HTTP response status 200
    And we should see Subtotal "Discounts" with parent "Revenues"

  Scenario: PATCH Subtotal with non-existent parent
    Given we POST Subtotal "Gold" with no parent
    When we PATCH Subtotal "Gold" with new parent "Assets"
    Then we should see HTTP response status 404

  Scenario: DELETE Subtotal
    Given we POST Subtotal "Loans Payable" with no parent
    When we DELETE Subtotal "Loans Payable"
    Then we should see HTTP response status 200
    And we should see Subtotal "Loans Payable" with no parent

  Scenario: DELETE non-existent Subtotal
    When we DELETE Subtotal "Assets"
    Then we should see HTTP response status 404

  Scenario: DELETE Subtotal and then GET
    Given we POST Subtotal "Taxes Payable" with no parent
    When we DELETE Subtotal "Taxes Payable"
    And we GET Subtotal "Taxes Payable"
    Then we should see HTTP response status 404

  Scenario: DELETE Subtotal with existing child
    Given we POST Subtotal "Payables" with no parent
    And we POST Subtotal "Notes Payable" with parent "Payables"
    When we DELETE Subtotal "Payables"
    Then we should see HTTP response status 409

  Scenario: DELETE Subtotal with existing Account
    Given Account endpoint is available
    And we POST Subtotal "Shares" with no parent
    And we POST Account "Preferred Shares" with Subtotal "Shares"
    When we DELETE Subtotal "Shares"
    Then we should see HTTP response status 409
