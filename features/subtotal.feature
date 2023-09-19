Feature: Subtotal
  Scenario: GET non-existent Subtotal
    Given a Subtotal endpoint is available
    When we GET a Subtotal by name "Assets"
    Then we should see HTTP response status 404

  Scenario: POST Subtotal
    Given a Subtotal endpoint is available
    When we POST a Subtotal with name "Balance Sheet" and no parent
    Then we should see HTTP response status 204

  Scenario: POST Subtotal with no parent and then GET
    Given a Subtotal endpoint is available
    When we POST a Subtotal with name "Profit & Loss" and no parent
    And we GET a Subtotal by name "Profit & Loss"
    Then we should see HTTP response status 200
    And we should see a Subtotal with name "Profit & Loss" and no parent

  Scenario: POST Subtotal with parent and then GET
    Given a Subtotal endpoint is available
    When we POST a Subtotal with name "Current Assets" and no parent
    And we POST a Subtotal with name "Cash" and parent "Current Assets"
    And we GET a Subtotal by name "Cash"
    Then we should see HTTP response status 200
    And we should see a Subtotal with name "Cash" and parent "Current Assets"

  Scenario: POST Subtotal and then POST another Subtotal with same name
    Given a Subtotal endpoint is available
    When we POST a Subtotal with name "Liabilities" and no parent
    And we POST a Subtotal with name "Liabilities" and no parent
    Then we should see HTTP response status 409

  Scenario: PATCH Subtotal
    Given a Subtotal endpoint is available
    And we POST a Subtotal with name "Discounts" and no parent
    When we PATCH a Subtotal named "Discounts" with new name "Sales Discounts"
    Then we should see HTTP response status 204

  Scenario: PATCH non-existent Subtotal
    Given a Subtotal endpoint is available
    When we PATCH a Subtotal named "Assets" with new parent "Balance Sheet"
    Then we should see HTTP response status 404

  Scenario: PATCH Subtotal with new name and then GET by new name
    Given a Subtotal endpoint is available
    And we POST a Subtotal with name "Depreciation" and no parent
    When we PATCH a Subtotal named "Depreciation" with new name "Amortisation"
    And we GET a Subtotal by name "Amortisation"
    Then we should see HTTP response status 200
    And we should see a Subtotal with name "Amortisation" and no parent

  Scenario: PATCH Subtotal with new name and then GET by old name
    Given a Subtotal endpoint is available
    And we POST a Subtotal with name "Provisions" and no parent
    When we PATCH a Subtotal named "Provisions" with new name "Reserves"
    And we GET a Subtotal by name "Provisions"
    Then we should see HTTP response status 404

  Scenario: PATCH Subtotal with new parent and then GET
    Given a Subtotal endpoint is available
    And we POST a Subtotal with name "Sales" and no parent
    And we POST a Subtotal with name "Discounts" and parent "Sales"
    And we POST a Subtotal with name "Revenues" and no parent
    When we PATCH a Subtotal named "Discounts" with new parent "Revenues"
    And we GET a Subtotal by name "Discounts"
    Then we should see HTTP response status 200
    And we should see a Subtotal with name "Discounts" and parent "Revenues"

  Scenario: DELETE Subtotal
    Given a Subtotal endpoint is available
    And we POST a Subtotal with name "Loans Payable" and no parent
    When we DELETE a Subtotal named "Loans Payable"
    Then we should see HTTP response status 200
    And we should see a Subtotal with name "Loans Payable" and no parent

  Scenario: DELETE non-existent Subtotal
    Given a Subtotal endpoint is available
    When we DELETE a Subtotal named "Assets"
    Then we should see HTTP response status 404

  Scenario: DELETE Subtotal and then GET
    Given a Subtotal endpoint is available
    And we POST a Subtotal with name "Taxes Payable" and no parent
    When we DELETE a Subtotal named "Taxes Payable"
    And we GET a Subtotal by name "Taxes Payable"
    Then we should see HTTP response status 404

  Scenario: DELETE Subtotal with existing child
    Given a Subtotal endpoint is available
    And we POST a Subtotal with name "Payables" and no parent
    When we POST a Subtotal with name "Notes Payable" and parent "Payables"
    And we DELETE a Subtotal named "Payables"
    Then we should see HTTP response status 409
