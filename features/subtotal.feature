Feature: Subtotal
  Scenario: GET non-existent Subtotal
    Given a Flow HTTP API v1 server is up
    When we GET a Subtotal by name "Assets"
    Then we should see HTTP response status 404

  Scenario: POST Subtotal with no parent
    Given a Flow HTTP API v1 server is up
    When we POST a Subtotal with name "Balance Sheet" and no parent
    Then we should see HTTP response status 201

  Scenario: POST Subtotal with no parent and then GET
    Given a Flow HTTP API v1 server is up
    When we POST a Subtotal with name "Profit & Loss" and no parent
    And we GET a Subtotal by name "Profit & Loss"
    Then we should see HTTP response status 200
    And we should see a Subtotal with name "Profit & Loss" and no parent
