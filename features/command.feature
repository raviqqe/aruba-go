Feature: Command

  Scenario: Run a command
    When I run `echo hello world`

  Scenario: Run a command successfully
    When I successfully run `echo hello world`

  Scenario: Check an exit status
    When I run `true`
    Then the exit status should be 0

  Scenario: Check an exit status of 1
    When I run `false`
    Then the exit status should be 1

  Scenario: Check a non-zero exit status
    When I run `false`
    Then the exit status should not be 0

  Scenario: Run a command interactively
    When I run `echo`
    Then the exit status should be 0
