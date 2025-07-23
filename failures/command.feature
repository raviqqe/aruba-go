Feature: Command

  Scenario: Run a command successfully
    When I successfully run `false`

  Scenario: Check an exit status
    When I run `true`
    Then the exit status should not be 0

  Scenario: Check an exit status of 1
    When I run `false`
    Then the exit status should not be 1

  Scenario: Check an exit status of 0
    When I run `true`
    Then the exit status should be 1

  Scenario: Run a command interactively
    When I run `echo` interactively
    Then the exit status should not be 0
