Feature: Chart (of Accounts)
  Scenario: GET Chart
    Given Chart endpoint is available
    And Subtotal endpoint is available
    And Account endpoint is available
    And we POST Subtotal "Inventory" with no parent
    And we POST Subtotal "Inputs" with parent "Inventory"
    And we POST Account "Raw Materials" with Subtotal "Inputs"
    And we POST Account "Parts" with Subtotal "Inputs"
    And we POST Account "Supplies" with Subtotal "Inputs"
    And we POST Account "Work in Progress" with Subtotal "Inventory"
    And we POST Account "Finished Goods" with Subtotal "Inventory"
    And we POST Subtotal "Merchandise" with parent "Inventory"
    When we GET Chart based on Subtotal "Inventory"
    Then we should see HTTP response status 200
    And we should see Chart with edge "Inventory" -> "Inputs"
    And we should see Chart with edge "Inputs" -> "Raw Materials"
    And we should see Chart with edge "Inputs" -> "Parts"
    And we should see Chart with edge "Inputs" -> "Supplies"
    And we should see Chart with edge "Inventory" -> "Work in Progress"
    And we should see Chart with edge "Inventory" -> "Finished Goods"
    And we should see Chart with edge "Inventory" -> "Merchandise"
