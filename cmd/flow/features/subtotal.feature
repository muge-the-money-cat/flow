Feature: Subtotal
  Scenario: Create Subtotal with no parent
    When we create Subtotal "Profit & Loss" with no parent
    Then we should see message "Subtotal successfully created"
    And we should see Subtotal "Profit & Loss" with no parent

  Scenario: Create Subtotal with parent
    Given we create Subtotal "Current Assets" with no parent
    When we create Subtotal "Cash" with parent "Current Assets"
    Then we should see message "Subtotal successfully created"
    And we should see Subtotal "Cash" with parent "Current Assets"

  Scenario: Create Subtotal with same name as existing
    Given we create Subtotal "Balance Sheet" with no parent
    And we create Subtotal "Liabilities" with parent "Balance Sheet"
    When we create Subtotal "Liabilities" with no parent
    Then we should see error "Subtotal with same name exists"
    And we should see Subtotal "Liabilities" with parent "Balance Sheet"

  Scenario: Create Subtotal with non-existent parent
    When we create Subtotal "Land" with parent "Assets"
    Then we should see error "Parent Subtotal does not exist"

  #Scenario: PATCH Subtotal
  #  Given we POST Subtotal "Discounts" with no parent
  #  When we PATCH Subtotal "Discounts" with new name "Sales Discounts"
  #  Then we should see HTTP response status 204

  #Scenario: PATCH non-existent Subtotal
  #  When we PATCH Subtotal "Assets" with new parent "Balance Sheet"
  #  Then we should see HTTP response status 404

  #Scenario: PATCH Subtotal with same name as existing
  #  Given we POST Subtotal "Machinery" with no parent
  #  And we POST Subtotal "Plant" with no parent
  #  When we PATCH Subtotal "Plant" with new name "Machinery"
  #  Then we should see HTTP response status 409

  #Scenario: PATCH Subtotal with new name and then GET by new name
  #  Given we POST Subtotal "Depreciation" with no parent
  #  When we PATCH Subtotal "Depreciation" with new name "Amortisation"
  #  And we GET Subtotal "Amortisation"
  #  Then we should see HTTP response status 200
  #  And we should see Subtotal "Amortisation" with no parent

  #Scenario: PATCH Subtotal with new name and then GET by old name
  #  Given we POST Subtotal "Provisions" with no parent
  #  When we PATCH Subtotal "Provisions" with new name "Reserves"
  #  And we GET Subtotal "Provisions"
  #  Then we should see HTTP response status 404

  #Scenario: PATCH Subtotal with new parent and then GET
  #  Given we POST Subtotal "Sales" with no parent
  #  And we POST Subtotal "Discounts" with parent "Sales"
  #  And we POST Subtotal "Revenues" with no parent
  #  When we PATCH Subtotal "Discounts" with new parent "Revenues"
  #  And we GET Subtotal "Discounts"
  #  Then we should see HTTP response status 200
  #  And we should see Subtotal "Discounts" with parent "Revenues"

  #Scenario: PATCH Subtotal with non-existent parent
  #  Given we POST Subtotal "Gold" with no parent
  #  When we PATCH Subtotal "Gold" with new parent "Assets"
  #  Then we should see HTTP response status 404

  Scenario: Delete Subtotal
    Given we create Subtotal "Loans Payable" with no parent
    When we delete Subtotal "Loans Payable"
    Then we should see message "Subtotal successfully deleted"
    And we should see Subtotal "Loans Payable" with no parent

  #Scenario: DELETE non-existent Subtotal
  #  When we DELETE Subtotal "Assets"
  #  Then we should see HTTP response status 404

  #Scenario: DELETE Subtotal and then GET
  #  Given we POST Subtotal "Taxes Payable" with no parent
  #  When we DELETE Subtotal "Taxes Payable"
  #  And we GET Subtotal "Taxes Payable"
  #  Then we should see HTTP response status 404

  #Scenario: DELETE Subtotal with existing child
  #  Given we POST Subtotal "Payables" with no parent
  #  And we POST Subtotal "Notes Payable" with parent "Payables"
  #  When we DELETE Subtotal "Payables"
  #  Then we should see HTTP response status 409

  #Scenario: DELETE Subtotal with existing Account
  #  Given Account endpoint is available
  #  And we POST Subtotal "Shares" with no parent
  #  And we POST Account "Preferred Shares" with Subtotal "Shares"
  #  When we DELETE Subtotal "Shares"
  #  Then we should see HTTP response status 409
