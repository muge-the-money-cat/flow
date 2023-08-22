Feature: Subtotal HTTP API
  Scenario: GET non-existent Subtotal
    Given a Subtotal HTTP API server is up
    When we GET a Subtotal by name "Assets"
    Then we should see HTTP response status 404
