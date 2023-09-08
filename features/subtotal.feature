Feature: Subtotal
  Scenario: GET non-existent Subtotal
    Given a Subtotal endpoint is available
    When we GET a Subtotal by name "Assets"
    Then we should see HTTP response status 404

  Scenario: POST Subtotal with no parent
    Given a Subtotal endpoint is available
    When we POST a Subtotal with name "Balance Sheet" and no parent
    Then we should see HTTP response status 201

  Scenario: POST Subtotal with no parent and then GET
    Given a Subtotal endpoint is available
    When we POST a Subtotal with name "Profit & Loss" and no parent
    And we GET a Subtotal by name "Profit & Loss"
    Then we should see HTTP response status 200
    And we should see a Subtotal with name "Profit & Loss" and no parent

  Scenario: POST Subtotal and then POST another Subtotal with same name
    Given a Subtotal endpoint is available
    When we POST a Subtotal with name "Liabilities" and no parent
    And we POST a Subtotal with name "Liabilities" and no parent
    Then we should see HTTP response status 409
